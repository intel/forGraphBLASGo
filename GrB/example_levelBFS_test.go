package GrB_test

import (
	"fmt"
	"github.com/intel/forGraphBLASGo/GrB"
	"testing"
)

/*
 * Given a boolean n x n adjacency matrix A and a source vertex s, performs a BFS traversal
 * of the graph and sets v[i] to the level in which vertex i is visited (v[i] = 1).
 * If i is not reachable from s, then v[i] = 0.
 */
func LevelBreadthFirstSearch(A GrB.Matrix[bool], s GrB.Index) (v GrB.Vector[int], err error) {
	defer GrB.CheckErrors(&err)

	n, err := A.Nrows()
	GrB.OK(err)

	v, err = GrB.VectorNew[int](n)
	GrB.OK(err)
	defer func() {
		if err != nil {
			_ = v.Free()
		}
	}()

	// vertices visited in each level
	q, err := GrB.VectorNew[bool](n)
	GrB.OK(err)
	defer func() {
		// q vector no longer needed
		GrB.OK(q.Free())
	}()
	// q[s] = true, false everywhere else
	GrB.OK(q.SetElement(true, s))

	/*
	 * BFS traversal and label the vertices.
	 */

	// d = level in BFS traversal
	d := 0

	// succ == true when some successor found
	for succ := true; succ; {
		// next level (start with 1)
		d++
		// v[q] = d
		GrB.OK(GrB.VectorAssignConstant(v, &q, nil, d, GrB.All(n), nil))
		// q [!v] = q ||.&& A ; finds all the unvisited successors from current q
		GrB.OK(GrB.VxM(q, v.AsMask(), nil, GrB.LorLandSemiringBool, q, A, GrB.DescRC))

		// succ = ||(q)
		succ, err = GrB.VectorReduce(GrB.LorMonoidBool, q, nil)
		GrB.OK(err)
	}

	return v, nil
}

func Example_levelBreadthFirstSearch() {
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

	v, err := LevelBreadthFirstSearch(A, 0)
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
	// [1 2 2 2 2 3 3 3 3 3 3 3 3 4 4]
}
