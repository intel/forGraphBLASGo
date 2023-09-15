package GrB_test

import (
	"fmt"
	"github.com/intel/forGraphBLASGo/GrB"
	"testing"
)

/*
 * Given an n x n boolean adjacency matrix, A, of an undirected graph, computes
 * the number of triangles in the graph.
 */
func TriangleCount[T GrB.Predefined](A GrB.Matrix[T]) (count int, err error) {
	defer GrB.CheckErrors(&err)

	n, err := A.Nrows()
	GrB.OK(err)

	// L: NxN, lower-triangular, bool
	L, err := GrB.MatrixNew[bool](n, n)
	GrB.OK(err)
	defer func() {
		GrB.OK(L.Free())
	}()
	GrB.OK(GrB.MatrixSelect(L, nil, nil, GrB.Tril[bool](), GrB.MatrixView[bool, T](A), 0, nil))
	Lint := GrB.MatrixView[int, bool](L)

	C, err := GrB.MatrixNew[int](n, n)
	GrB.OK(err)
	defer func() {
		GrB.OK(C.Free())
	}()

	// C<L> = L +.* L
	GrB.OK(GrB.MxM(C, &L, nil, GrB.PlusTimesSemiring[int](), Lint, Lint, nil))

	// 1-norm of C
	return GrB.MatrixReduce(GrB.PlusMonoid[int](), C, nil)
}

func Example_triangleCount() {
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
	OK(A.SetElement(true, 1, 3))
	OK(A.SetElement(true, 1, 5))
	OK(A.SetElement(true, 1, 6))
	OK(A.SetElement(true, 2, 4))
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

	OK(GrB.MatrixEWiseAddBinaryOp(A, nil, nil, GrB.LorBool, A, A, GrB.DescT1))

	count, err := TriangleCount(A)
	OK(err)
	fmt.Println(count)
	// Output:
	// 2
}
