package GrB_test

import (
	"fmt"
	"github.com/intel/forGraphBLASGo/GrB"
	"testing"
)

// Assign a random number to each element scaled by the inverse of the node's degree.
// This will increase the probability that low degree nodes are selected and larger
// sets are selected.
const (
	setRandomName = "setRandom"
	setRandomDef  = `void setRandom(void *out, const void *in) {
       uint32_t degree = *(uint32_t*)in;
       *(float*)out = (0.0001f + random()/(1. + 2.*degree)); // add 1 to prevent divied by zero
    }`
)

/*
 * A variant of Luby's randomized algorithm [Luby 1985].
 *
 * Given a numeric n x n adjacency matrix A of an unweighted and undirected graph (where
 * the value true represents an edge), compute a maximal set of independent vertices and
 * return it in a boolean n-vector, 'iset' where set[i] == true implies vertex i is a member
 * of the set.
 */
func MaximalIndependentSet(A GrB.Matrix[bool]) (iset GrB.Vector[bool], err error) {
	defer GrB.CheckErrors(&err)

	n, err := A.Nrows()
	GrB.OK(err)

	// Initialize independent set vector
	iset, err = GrB.VectorNew[bool](n)
	GrB.OK(err)
	defer func() {
		if err != nil {
			_ = iset.Free()
		}
	}()

	// holds random probabilities for each node
	prob, err := GrB.VectorNew[float32](n)
	GrB.OK(err)
	defer func() {
		GrB.OK(prob.Free())
	}()
	// holds value of max neighbor probability
	neighborMax, err := GrB.VectorNew[float32](n)
	GrB.OK(err)
	defer func() {
		GrB.OK(neighborMax.Free())
	}()
	// holds set of new members to iset
	newMembers, err := GrB.VectorNew[bool](n)
	GrB.OK(err)
	defer func() {
		GrB.OK(newMembers.Free())
	}()
	// holds set of new neighbors to new iset members
	newNeighbors, err := GrB.VectorNew[bool](n)
	GrB.OK(err)
	defer func() {
		GrB.OK(newNeighbors.Free())
	}()
	// candidate members to iset
	candidates, err := GrB.VectorNew[bool](n)
	GrB.OK(err)
	defer func() {
		GrB.OK(candidates.Free())
	}()

	setRandom, err := GrB.NamedUnaryOpNew[float32, uint32](nil, setRandomName, setRandomDef)
	GrB.OK(err)
	defer func() {
		GrB.OK(setRandom.Free())
	}()

	// compute the degree of each vertex
	degrees, err := GrB.VectorNew[float64](n)
	GrB.OK(err)
	defer func() {
		GrB.OK(degrees.Free())
	}()
	GrB.OK(degrees.MatrixReduceBinaryOp(nil, nil, GrB.Plus[float64](), GrB.MatrixView[float64, bool](A), nil))

	// Isolated vertices are not candidates: candidates[degrees != 0] = true
	GrB.OK(candidates.AssignConstant(degrees.AsMask(), nil, true, GrB.All(n), nil))

	// add all singletons to iset: iset[degree == 0] = 1
	GrB.OK(iset.AssignConstant(degrees.AsMask(), nil, true, GrB.All(n), nil))

	// Iterate while there are candidates to check.
	nvals, err := candidates.Nvals()
	GrB.OK(err)
	for nvals > 0 {
		// compute a random probability scaled by inverse of degree
		GrB.OK(GrB.VectorApply(prob, &candidates, nil, setRandom, GrB.VectorView[uint32, float64](degrees), GrB.DescR))

		// compute the max probability of all neighbors
		GrB.OK(neighborMax.MxV(&candidates, nil, GrB.MaxSecondSemiring[float32](), GrB.MatrixView[float32, bool](A), prob, GrB.DescR))

		// select vertex if its probability is larger than all its active neighbors,
		// and apply a "masked no-op" to remove stored falses
		GrB.OK(GrB.VectorEWiseAddBinaryOp(newMembers, nil, nil, GrB.Gt[float32](), prob, neighborMax, nil))
		GrB.OK(newMembers.Apply(&newMembers, nil, GrB.Identity[bool](), newMembers, GrB.DescR))

		// add new members to independent set
		GrB.OK(iset.EWiseAddBinaryOp(nil, nil, GrB.LorBool, iset, newMembers, nil))

		// remove new members from set of candidates c = c & !new
		GrB.OK(candidates.EWiseMultBinaryOp(&newMembers, nil, GrB.LandBool, candidates, candidates, GrB.DescRC))

		nvals, err = candidates.Nvals()
		GrB.OK(err)
		if nvals == 0 {
			break
		}

		// Neighbors of new members can also be removed from candidates
		GrB.OK(newNeighbors.MxV(&candidates, nil, GrB.LorLandSemiringBool, A, newMembers, nil))
		GrB.OK(candidates.EWiseMultBinaryOp(&newNeighbors, nil, GrB.LandBool, candidates, candidates, GrB.DescRC))

		nvals, err = candidates.Nvals()
		GrB.OK(err)
	}

	return iset, nil
}

func Example_maximalIndependentSet() {
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

	iset, err := MaximalIndependentSet(A)
	OK(err)
	defer func() {
		OK(iset.Free())
	}()

	var indices []int
	OK(iset.ExtractTuples(&indices, nil))
	fmt.Println(indices)
	// Output:
	// [0 1 2 3 4 6 7 10 11]
}
