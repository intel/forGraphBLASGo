package GrB

// #include "GraphBLAS.h"
import "C"

// Index is the index type used for accessing elements of vectors and matrices.
//
// Unlike in the GraphBLAS C API, this is Go's default signed int type, not C's
// unsigned uint64_t type. This fits much better with the rest of Go due to int
// being the default type for integer constants in Go.
//
// The forGraphBLASGo API does not use the Index type, but directly uses int instead.
//
// The range of valid values for a variable of type Index is [0, [IndexMax]].
//
// Sets of indices are represented as slices of int ([]int]). Likewise, sets of
// scalar values are represented as slices of the corresponding type, for example
// []float64. Some GraphBLAS operations (for example, [MatrixAssign]) include an
// input parameter with the type of an index slice. This input index slice selects
// a subset of elements from a GraphBLAS vector or matrix object to be used in the
// operation. In these cases, the function [All] can be used to indicate that all
// indices of the associated GraphBLAS vector or matrix object up to a certain index
// should be used.
//
// User-defined index slices should not include negative indices. Functions [All],
// [Range], [Stride] and [Backwards] may create slices with negative indices to
// mark special cases that are properly understood by forGraphBLASGo functions.
// The meaning of these markers may silently change in future versions of the
// forGraphBLASGo API.
type Index = int

// IndexMax is the permissible maximum value for [Index].
const IndexMax = 1<<60 - 1

// All is used in various GraphBLAS functions to indicate that all indices in the
// range 0 <= index < size are to be used. This representation uses significantly
// less memory than enumerating the corresponding indices explicitly, and is also
// handled more efficiently by the underlying GraphBLAS implementation in C.
func All(size int) []int {
	if size < 0 {
		panic(IndexOutOfBounds.Error())
	}
	if size == 0 {
		return nil
	}
	return []int{-size}
}

// Range can be used everywhere [All] can be used as well.
// Range indicates that all indices in the range begin <= index < end are
// to be used. This representation uses significantly less memory than enumerating the
// corresponding indices explicitly, and is also handled more efficiently by the
// underlying SuiteSparse:GraphBLAS implementation.
//
// Range is a SuiteSparse:GraphBLAS extension.
func Range(begin, end int) []int {
	if begin < 0 || end < 0 || begin > end {
		panic(IndexOutOfBounds.Error())
	}
	return []int{-1, begin, end}
}

// Stride can be used everywhere [All] can be used as well.
// Stride indicates that indices in the range begin <= index < end are
// to be used, starting at begin and incremented by inc.
// For example, Stride(3, 10, 2) corresponds to the slice []int{3, 5, 7, 9}.
// This representation uses significantly less memory than enumerating the
// corresponding indices explicitly, and is also handled more efficiently by the
// underlying SuiteSparse:GraphBLAS implementation.
//
// Stride is a SuiteSparse:GraphBLAS extension.
func Stride(begin, end, inc int) []int {
	if begin < 0 || end < 0 || inc < 0 || begin > end {
		panic(IndexOutOfBounds.Error())
	}
	return []int{-1, begin, end, inc}
}

// Backwards is like [Stride], except that begin > end, and inc is used
// as a decrement. For example, Stride(10, 3, 2) corresponds to the slice
// []int{10, 8, 6, 4}. The resulting indices are in the range begin >= index > end.
// This representation uses significantly less memory than enumerating the
// corresponding indices explicitly, and is also handled more efficiently by the
// underlying SuiteSparse:GraphBLAS implementation.
//
// Backwards is a SuiteSparse:GraphBLAS extension.
func Backwards(begin, end, inc int) []int {
	if begin < 0 || end < 0 || inc < 0 || begin < end {
		panic(IndexOutOfBounds.Error())
	}
	return []int{-2, begin, end, inc}
}

func cIndices(indices []int) (cindices *C.GrB_Index, cnindices C.GrB_Index, e error) {
	switch len(indices) {
	case 0:
		return (*C.GrB_Index)(nil), 0, nil
	case 1:
		if sz := indices[0]; sz < 0 {
			return C.GrB_ALL, C.GrB_Index(-sz), nil
		}
	case 3:
		if indices[0] < 0 {
			return grbIndices(indices[1:]), C.GxB_RANGE, nil
		}
	case 4:
		if indices[0] < 0 {
			switch indices[0] {
			case -1:
				return grbIndices(indices[1:]), C.GxB_STRIDE, nil
			case -2:
				return grbIndices(indices[1:]), C.GxB_BACKWARDS, nil
			default:
				panic("Ambiguous indices kind.")
			}
		}
	}
	for _, index := range indices {
		if index < 0 {
			e = makeError(IndexOutOfBounds)
			return
		}
	}
	return grbIndices(indices), C.GrB_Index(len(indices)), nil
}
