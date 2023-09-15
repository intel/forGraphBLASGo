package GrB_test

import (
	"fmt"
	"github.com/intel/forGraphBLASGo/GrB"
	"testing"
)

/*
 * Compute partial BC metric for a subset of source vertices, s, in graph A
 */
func BatchedBetweennessCentrality(A GrB.Matrix[bool], s []GrB.Index) (delta GrB.Vector[float32], err error) {
	defer GrB.CheckErrors(&err)

	n, err := A.Nrows()
	GrB.OK(err)

	delta, err = GrB.VectorNew[float32](n)
	GrB.OK(err)
	defer func() {
		if err != nil {
			_ = delta.Free()
		}
	}()

	// index and value arrays needed to build numsp
	iLens := make([]GrB.Index, len(s))
	ones := make([]int, len(s))
	for i := range s {
		iLens[i] = i
		ones[i] = 1
	}

	// numsp: structure holds the number of shortest paths for each node and starting vertex
	// discovered so far. Initialized to source vertices: numsp[s[i], i] = 1, i = [0, len(s)]
	numsp, err := GrB.MatrixNew[int](n, len(s))
	GrB.OK(err)
	defer func() {
		GrB.OK(numsp.Free())
	}()
	dup := GrB.Plus[int]()
	GrB.OK(numsp.Build(s, iLens, ones, &dup))
	iLens, ones = nil, nil

	// frontier: Holds the current frontier where values are path counts.
	// Initialized to out vertices of each source node in s.
	frontier, err := GrB.MatrixNew[int](n, len(s))
	GrB.OK(err)
	defer func() {
		GrB.OK(frontier.Free())
	}()
	GrB.OK(GrB.MatrixExtract(frontier, numsp.AsMask(), nil, GrB.MatrixView[int, bool](A), GrB.All(n), s, GrB.DescRCT0))

	// sigma: stores frontier information for each level of BFS phase. The memory
	// for an entry in sigmas is only allocated within the for-loop if needed.
	var sigmas []GrB.Matrix[bool]

	// nvals == 0 when BFS phase is complete
	for nvals := 1; nvals > 0; {
		// sigmas[level](:,s) = frontier from source vertex s for the current level
		sigma, err := GrB.MatrixNew[bool](n, len(s))
		GrB.OK(err)
		defer func() {
			GrB.OK(sigma.Free())
		}()
		sigmas = append(sigmas, sigma)

		// sigmas[level](:,:) = bool(frontier)
		GrB.OK(GrB.MatrixApply(sigma, nil, nil, GrB.Identity[bool](), GrB.MatrixView[bool, int](frontier), nil))
		// numsp += frontier (accum path counts)
		GrB.OK(GrB.MatrixEWiseAddBinaryOp(numsp, nil, nil, GrB.Plus[int](), numsp, frontier, nil))
		// f<!numsp> = A' +.* f (update frontier)
		GrB.OK(GrB.MxM(frontier, numsp.AsMask(), nil, GrB.PlusTimesSemiring[int](), GrB.MatrixView[int, bool](A), frontier, GrB.DescRCT0))
		// number of nodes in frontier at this level
		nvals, err = frontier.Nvals()
		GrB.OK(err)
	}

	// nspinv: the inverse of the number of shortest paths for each node and starting vertex.
	nspinv, err := GrB.MatrixNew[float32](n, len(s))
	GrB.OK(err)
	defer func() {
		GrB.OK(nspinv.Free())
	}()
	// nspinv = 1/numsp
	GrB.OK(GrB.MatrixApply(nspinv, nil, nil, GrB.Minv[float32](), GrB.MatrixView[float32, int](numsp), nil))

	// bcu: BC updates for each vertex for each starting vertex in s
	bcu, err := GrB.MatrixNew[float32](n, len(s))
	GrB.OK(err)
	defer func() {
		GrB.OK(bcu.Free())
	}()
	// filled with 1 to avoid sparsity issues
	GrB.OK(GrB.MatrixAssignConstant(bcu, nil, nil, 1, GrB.All(n), GrB.All(len(s)), nil))

	// temporary workspace matrix
	w, err := GrB.MatrixNew[float32](n, len(s))
	GrB.OK(err)
	defer func() {
		GrB.OK(w.Free())
	}()

	plusFloat32 := GrB.Plus[float32]()

	// Tally phase (backward sweep)
	for i := len(sigmas) - 1; i > 0; i-- {
		// w<sigmas[i]> = (1 ./ nsp) .* bcu
		GrB.OK(GrB.MatrixEWiseMultBinaryOp(w, sigmas[i].AsMask(), nil, GrB.Times[float32](), bcu, nspinv, GrB.DescR))

		// add contributions by successors and mask with that BFS level's frontier
		// w<sigmas[i-1]> = (A +.* w)
		GrB.OK(GrB.MxM(w, sigmas[i-1].AsMask(), nil, GrB.PlusTimesSemiring[float32](), GrB.MatrixView[float32, bool](A), w, GrB.DescR))
		// bcu += w .* numsp
		GrB.OK(GrB.MatrixEWiseMultBinaryOp(bcu, nil, &plusFloat32, GrB.Times[float32](), w, GrB.MatrixView[float32, int](numsp), nil))
	}

	// row reduce bcu and subtract "len(s)" from every entry to account
	// for 1 extra value per bcu row element.
	GrB.OK(GrB.MatrixReduceBinaryOp(delta, nil, nil, GrB.Plus[float32](), bcu, nil))
	GrB.OK(GrB.VectorApplyBinaryOp2nd(delta, nil, nil, GrB.Minus[float32](), delta, float32(len(s)), nil))

	return delta, nil
}

func Example_batchedBetweennessCentrality() {
	OK := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	if !testing.Testing() {
		// When run by "go test", this initialization of
		// GraphBLAS is done elsewhere in TestMain.
		OK(GrB.Init(GrB.NonBlocking))
		defer func() {
			OK(GrB.Finalize())
		}()
	}

	A, err := GrB.MatrixNew[bool](15, 15)
	OK(err)

	OK(A.SetElement(true, 0, 1))
	OK(A.SetElement(true, 0, 2))
	OK(A.SetElement(true, 0, 3))
	OK(A.SetElement(true, 0, 4))
	OK(A.SetElement(true, 1, 5))
	OK(A.SetElement(true, 1, 6))
	OK(A.SetElement(true, 2, 7))
	OK(A.SetElement(true, 2, 8))
	OK(A.SetElement(true, 3, 9))
	OK(A.SetElement(true, 3, 10))
	OK(A.SetElement(true, 4, 11))
	OK(A.SetElement(true, 4, 12))
	OK(A.SetElement(true, 6, 13))
	OK(A.SetElement(true, 7, 13))
	OK(A.SetElement(true, 10, 14))
	OK(A.SetElement(true, 11, 14))

	delta, err := BatchedBetweennessCentrality(A, []int{0, 1, 2, 3, 4})
	OK(err)
	defer func() {
		OK(delta.Free())
	}()

	var indices []int
	var values []float32
	OK(delta.ExtractTuples(&indices, &values))
	fmt.Println(indices)
	fmt.Println(values)
	// Output:
	// [0 1 2 3 4 5 6 7 8 9 10 11 12 13 14]
	// [0 2.5 2.5 2.5 2.5 0 1.5 1.5 0 0 1.5 1.5 0 0 0]
}
