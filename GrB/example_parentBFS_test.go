package GrB_test

import (
	"fmt"
	"github.com/intel/forGraphBLASGo/GrB"
	"testing"
)

/*
 * Given a binary n x n adjacency matrix A and a source vertex s, performs a BFS
 * traversal of the graph and sets parents[i] to the index of vertex i's parent.
 * The parent of the root vertex, s, will be set to itself (parents[s] = s). If
 * vertex i is not reachable from s, parents[i] will not contain a stored value.
 */
func ParentBreadthFirstSearch(A GrB.Matrix[int], s GrB.Index) (parents GrB.Vector[int], err error) {
	defer GrB.CheckErrors(&err)

	n, err := A.Nrows()
	GrB.OK(err)

	parents, err = GrB.VectorNew[int](n)
	GrB.OK(err)
	defer func() {
		if err != nil {
			_ = parents.Free()
		}
	}()
	GrB.OK(parents.SetElement(s, s))

	wavefront, err := GrB.VectorNew[int](n)
	GrB.OK(err)
	defer func() {
		GrB.OK(wavefront.Free())
	}()
	GrB.OK(wavefront.SetElement(1, s))

	/*
	 * BFS traversal and label the vertices.
	 */

	plusInt := GrB.Plus[int]()

	for nvals := 1; nvals > 0; {
		// convert all stored values in wavefront to their 0-based index
		GrB.OK(wavefront.ApplyIndexOp(nil, nil, GrB.RowIndex[int, int](), wavefront, 0, nil))

		// "First" because left-multiplying wavefront rows. Masking out the parent
		// list ensures wavefront values do not overwrite parents already stored.
		GrB.OK(wavefront.VxM(parents.AsMask(), nil, GrB.MinFirstSemiring[int](), wavefront, A, GrB.DescRSC))

		// Don't need to mask here since we did it in VxM. Merges new parents in
		// current wavefrontwith existing parents: parents += wavefront
		GrB.OK(parents.Apply(nil, &plusInt, GrB.Identity[int](), wavefront, nil))

		nvals, err = wavefront.Nvals()
		GrB.OK(err)
	}

	return parents, nil
}

func Example_parentBreadthFirstSearch() {
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

	A, err := GrB.MatrixNew[int](15, 15)
	OK(err)

	OK(A.SetElement(1, 0, 1))
	OK(A.SetElement(1, 0, 2))
	OK(A.SetElement(1, 0, 3))
	OK(A.SetElement(1, 0, 4))
	OK(A.SetElement(1, 1, 5))
	OK(A.SetElement(1, 1, 6))
	OK(A.SetElement(1, 2, 7))
	OK(A.SetElement(1, 2, 8))
	OK(A.SetElement(1, 3, 9))
	OK(A.SetElement(1, 3, 10))
	OK(A.SetElement(1, 4, 11))
	OK(A.SetElement(1, 4, 12))
	OK(A.SetElement(1, 6, 13))
	OK(A.SetElement(1, 7, 13))
	OK(A.SetElement(1, 10, 14))
	OK(A.SetElement(1, 11, 14))

	v, err := ParentBreadthFirstSearch(A, 0)
	OK(err)
	defer func() {
		OK(v.Free())
	}()

	var indices, values []int
	OK(v.ExtractTuples(&indices, &values))
	fmt.Println(indices)
	fmt.Println(values)
	// Output:
	// [0 1 2 3 4 5 6 7 8 9 10 11 12 13 14]
	// [0 0 0 0 0 1 1 2 2 3 3 4 4 6 10]
}
