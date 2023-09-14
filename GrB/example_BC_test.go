package GrB_test

import (
	"fmt"
	"github.com/intel/forGraphBLASGo/GrB"
	"testing"
)

/*
 * Given a boolean n x n adjacency matrix A and a source vertex s,
 * compute the BC-metric vector delta.
 */
func BetweennessCentrality(A GrB.Matrix[bool], s GrB.Index) (delta GrB.Vector[float32], err error) {
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

	// sigma[d, k] = #shortest paths to node k at level d
	sigma, err := GrB.MatrixNew[int](n, n)
	GrB.OK(err)
	defer func() {
		GrB.OK(sigma.Free())
	}()

	// path counts
	q, err := GrB.VectorNew[int](n)
	GrB.OK(err)
	defer func() {
		GrB.OK(q.Free())
	}()
	GrB.OK(q.SetElement(1, s))

	// shortest path counts so far
	p, err := q.Dup()
	GrB.OK(err)
	defer func() {
		GrB.OK(p.Free())
	}()

	// get the first set of out neighbors
	GrB.OK(q.VxM(p.AsMask(), nil, GrB.PlusTimesSemiring[int](), q, GrB.MatrixView[int, bool](A), GrB.DescRC))

	/*
	 * BFS phase
	 */

	// BFS level number
	d := 0
	// sum == 0 when BFS phase is complete
	for sum := 1; sum > 0; {
		// sigma[d,:] = q
		GrB.OK(sigma.RowAssign(nil, nil, q, d, GrB.All(n), nil))
		// accum path counts on this level
		GrB.OK(p.EWiseAddBinaryOp(nil, nil, GrB.Plus[int](), p, q, nil))
		// q = #paths to nodes reachable from current level
		GrB.OK(q.VxM(p.AsMask(), nil, GrB.PlusTimesSemiring[int](), q, GrB.MatrixView[int, bool](A), GrB.DescRC))
		// sum path counts at this level
		sum, err = q.Reduce(GrB.PlusMonoid[int](), nil)
		d++
	}

	/*
	 * BC computation phase
	 * (t1, t2, t3, t4) are temporary vectors
	 */
	t1, err := GrB.VectorNew[float32](n)
	GrB.OK(err)
	defer func() {
		GrB.OK(t1.Free())
	}()
	t2, err := GrB.VectorNew[float32](n)
	GrB.OK(err)
	defer func() {
		GrB.OK(t2.Free())
	}()
	t3, err := GrB.VectorNew[float32](n)
	GrB.OK(err)
	defer func() {
		GrB.OK(t3.Free())
	}()
	t4, err := GrB.VectorNew[float32](n)
	GrB.OK(err)
	defer func() {
		GrB.OK(t4.Free())
	}()

	for i := d - 1; i > 0; i-- {
		// t1 = 1+delta
		GrB.OK(t1.AssignConstant(nil, nil, 1, GrB.All(n), nil))
		GrB.OK(t1.EWiseAddBinaryOp(nil, nil, GrB.Plus[float32](), t1, delta, nil))
		// t2 = sigma[i,:]
		GrB.OK(t2.ColExtract(nil, nil, GrB.MatrixView[float32, int](sigma), GrB.All(n), i, GrB.DescT0))
		// t2 = (1+delta)/sigma[i,:]
		GrB.OK(t2.EWiseMultBinaryOp(nil, nil, GrB.Div[float32](), t1, t2, nil))
		// add contributions made by successors of a node
		GrB.OK(t3.MxV(nil, nil, GrB.PlusTimesSemiring[float32](), GrB.MatrixView[float32, bool](A), t2, nil))
		// t4 = sigma[i-1,:]
		GrB.OK(t4.ColExtract(nil, nil, GrB.MatrixView[float32, int](sigma), GrB.All(n), i-1, GrB.DescT0))
		// t4 = sigma[i-1,:]*t3
		GrB.OK(t4.EWiseMultBinaryOp(nil, nil, GrB.Times[float32](), t4, t3, nil))
		// accumulate into delta
		GrB.OK(delta.EWiseAddBinaryOp(nil, nil, GrB.Plus[float32](), delta, t4, nil))
	}

	return delta, nil
}

func Example_betweennessCentrality() {
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

	delta, err := BetweennessCentrality(A, 0)
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
	// [1 2 3 4 6 7 10 11]
	// [2.5 2.5 2.5 2.5 0.5 0.5 0.5 0.5]
}
