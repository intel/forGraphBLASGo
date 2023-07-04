package GrB

// #include "GraphBLAS.h"
import "C"
import (
	"unsafe"
)

// A Matrix is defined by a domain D, its number of rows M > 0, its number of columnns N > 0,
// and a set up tuples (i, j, a(i, j)), where 0 <= i < M, 0 <= j < N, and a(i, j) ∈ D. A
// particular pair of values i, j can occur at most once in a.
type Matrix[D any] struct {
	grb C.GrB_Matrix
}

// A MatrixMask can be used to optionally control which results from a GraphBLAS operation
// are stored into an output matrix. The mask dimensions must match those of the
// output. If the [Structure] descriptor is not set for the mask, the domain
// of the mask matrix must be of type bool, any of the [Predefined] "built-in" types,
// or any of the [Complex] "built-in" types. Use [Matrix.AsMask] to convert
// to the required parameter type. If the default mask is desired (i.e., a mask that is all
// true with the dimensions of the output matrix), nil should be specified.
//
// The forGraphBLASGo API does not use the MatrixMask type, but directly uses *Matrix[bool] instead.
type MatrixMask = *Matrix[bool]

// MatrixView returns a view on the given matrix (with domain From) using a different domain To.
//
// In the GraphBLAS specification for the C programming language, collections (scalars, vectors and matrices) of
// [Predefined] domains can be arbitrarily intermixed. In SuiteSparse:GraphBLAS, this extends to collections of [Complex]
// domains. When entries of collections are accessed expecting a particular domain (type), then
// the entry values are typecast using the rules of the C programming language. (Collections of
// user-defined domains are not compatible with any other collections in this way.)
//
// In Go, generally only identical types are compatible with each other, and conversions are
// not implicit. To get around this restriction, [ScalarView], [VectorView] and MatrixView can be used to view a
// collection using a different domain. These functions do not perform any conversion themselves, but are essentially
// NO-OPs.
//
// MatrixView is a forGraphBLASGo extension.
func MatrixView[To, From Predefined | Complex](matrix Matrix[From]) (view Matrix[To]) {
	view.grb = matrix.grb
	return
}

// AsMask returns a view on the given matrix using the domain bool.
//
// In GraphBLAS, whenever a mask is required as an input parameter for a GraphBLAS operation,
// a matrix of any domain can be passed, and depending on whether [Structure] is set or not in the
// [Descriptor] passed to that operation, the only requirement is that the domain is compatible
// with bool. In the C programming language, this holds for any of the [Predefined] domains.
// In SuiteSparse:GraphBLAS, this extends to any of the [Complex] domains.
//
// In Go, generally only identical types are compatible with each other, and conversions are
// not implicit. To get around this restriction, AsMask can be used to view a matrix as a bool
// mask. AsMask does not perform any conversion itself, but is essentially a NO-OP.
//
// AsMask is a forGraphBLASGo extension.
func (matrix Matrix[D]) AsMask() *Matrix[bool] {
	return &Matrix[bool]{matrix.grb}
}

// Type returns the actual [Type] object representing the domain of the given matrix.
// This is not necessarily the [Type] object corresponding to D, if Type is called
// on a [MatrixView] of a matrix of some other domain.
//
// Type might return false as a second return value if the domain is not a [Predefined]
// or [Complex] domain, or if the type has not been registered with [TypeNew] or
// [NamedTypeNew].
//
// Type is a forGraphBLASGo extension. It can be used in place of GxB_Matrix_type_name
// and GxB_Type_from_name, which are SuiteSparse:GraphBLAS extensions.
func (matrix Matrix[D]) Type() (typ Type, ok bool, err error) {
	var ctypename [C.GxB_MAX_NAME_LEN]C.char
	info := Info(C.GxB_Matrix_type_name(&ctypename[0], matrix.grb))
	if info != success {
		err = makeError(info)
		return
	}
	var grb C.GrB_Type
	info = Info(C.GxB_Type_from_name(&grb, &ctypename[0]))
	if info != success {
		err = makeError(info)
	}
	typ, ok = goType[grb]
	return
}

// MatrixNew creates a new matrix with specified domain and dimensions.
//
// Parameters:
//
//   - D: The type corresponding to the domain of the matrix being created.
//     Can be one of the [Predefined] or [Complex] types, or an existing
//     user-defined GraphBLAS type.
//
//   - nrows (IN): The number of rows of the matrix being created.
//
//   - ncols (IN): The number of columns of the matrix being created.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func MatrixNew[D any](nrows, ncols int) (matrix Matrix[D], err error) {
	if nrows < 0 || ncols < 0 {
		err = makeError(InvalidValue)
	}
	var d D
	dt, ok := grbType[TypeOf(d)]
	if !ok {
		err = makeError(UninitializedObject)
		return
	}
	info := Info(C.GrB_Matrix_new(&matrix.grb, dt, C.GrB_Index(nrows), C.GrB_Index(ncols)))
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// Dup creates a new matrix with the same domain, dimensions, and contents as another matrix.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (matrix Matrix[D]) Dup() (dup Matrix[D], err error) {
	info := Info(C.GrB_Matrix_dup(&dup.grb, matrix.grb))
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// Resize changes the dimensions of an existing matrix.
//
// Parameters:
//
//   - nrows (IN): The new number of rows of the matrix. It can be smaller or larger than the current number of rows.
//
//   - ncolums (IN): The new number of colunms of the matrix. It can be smaller or larger than the current number of columns.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (matrix Matrix[D]) Resize(nrows, ncols int) error {
	if nrows < 0 || ncols < 0 {
		return makeError(InvalidValue)
	}
	info := Info(C.GrB_Matrix_resize(matrix.grb, C.GrB_Index(nrows), C.GrB_Index(ncols)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Clear removes all elements (tuples) from the matrix.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (matrix Matrix[D]) Clear() error {
	info := Info(C.GrB_Matrix_clear(matrix.grb))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Nrows retrieves the number of rows in a matrix.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [Panic]
func (matrix Matrix[D]) Nrows() (nrows int, err error) {
	var cnrows C.GrB_Index
	info := Info(C.GrB_Matrix_nrows(&cnrows, matrix.grb))
	if info == success {
		return int(cnrows), nil
	}
	err = makeError(info)
	return
}

// Ncols retrieves the number of columns in a matrix.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [Panic]
func (matrix Matrix[D]) Ncols() (ncols int, err error) {
	var cncols C.GrB_Index
	info := Info(C.GrB_Matrix_ncols(&cncols, matrix.grb))
	if info == success {
		return int(cncols), nil
	}
	err = makeError(info)
	return
}

// Size retrieves the number of rows and columns in a matrix.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [Panic]
//
// Size is a forGraphBLASGo extension.
func (matrix Matrix[D]) Size() (nrows, ncols int, err error) {
	if nrows, err = matrix.Nrows(); err != nil {
		return
	}
	ncols, err = matrix.Ncols()
	return
}

// Nvals retrieves the number of stored elements (tuples) in a matrix.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (matrix Matrix[D]) Nvals() (nvals int, err error) {
	var cnvals C.GrB_Index
	info := Info(C.GrB_Matrix_nvals(&cnvals, matrix.grb))
	if info == success {
		return int(cnvals), nil
	}
	err = makeError(info)
	return
}

// Build stores elements from tuples in a matrix.
//
// Parameters:
//
//   - rowIndices: A slice of row indices.
//
//   - colIndices: A slice of column indices.
//
//   - values: A slice of scalars of type D.
//
//   - dup: An associative and commutative binary operator to apply when duplicate
//     values for the same location are present in the input slices. All three domains
//     of dup must be D. If dup is nil, then duplicate locations will result in an [InvalidValue] error.
//
// GraphBLAS API errors that may be returned:
//   - [DomainMismatch], [InvalidValue], [SliceMismatch], [OutputNotEmpty], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [IndexOutOfBounds], [InvalidObject], [OutOfMemory], [Panic]
func (matrix Matrix[D]) Build(rowIndices, colIndices []int, values []D, dup *BinaryOp[D, D, D]) error {
	if len(rowIndices) != len(colIndices) || len(colIndices) != len(values) {
		return makeError(SliceMismatch)
	}
	for _, index := range rowIndices {
		if index < 0 {
			return makeError(IndexOutOfBounds)
		}
	}
	for _, index := range colIndices {
		if index < 0 {
			return makeError(IndexOutOfBounds)
		}
	}
	var cdup C.GrB_BinaryOp
	if dup == nil {
		cdup = C.GrB_BinaryOp(C.GrB_NULL)
	} else {
		cdup = dup.grb
	}
	var info Info
	switch vals := any(values).(type) {
	case []bool:
		info = Info(C.GrB_Matrix_build_BOOL(
			matrix.grb, grbIndices(rowIndices), grbIndices(colIndices),
			cSlice[C.bool, bool](vals),
			C.GrB_Index(len(rowIndices)), cdup,
		))
	case []int:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Matrix_build_INT32(
				matrix.grb, grbIndices(rowIndices), grbIndices(colIndices),
				cSlice[C.int32_t, int](vals),
				C.GrB_Index(len(rowIndices)), cdup,
			))
		} else {
			info = Info(C.GrB_Matrix_build_INT64(
				matrix.grb, grbIndices(rowIndices), grbIndices(colIndices),
				cSlice[C.int64_t, int](vals),
				C.GrB_Index(len(rowIndices)), cdup,
			))
		}
	case []int8:
		info = Info(C.GrB_Matrix_build_INT8(
			matrix.grb, grbIndices(rowIndices), grbIndices(colIndices),
			cSlice[C.int8_t, int8](vals),
			C.GrB_Index(len(rowIndices)), cdup,
		))
	case []int16:
		info = Info(C.GrB_Matrix_build_INT16(
			matrix.grb, grbIndices(rowIndices), grbIndices(colIndices),
			cSlice[C.int16_t, int16](vals),
			C.GrB_Index(len(rowIndices)), cdup,
		))
	case []int32:
		info = Info(C.GrB_Matrix_build_INT32(
			matrix.grb, grbIndices(rowIndices), grbIndices(colIndices),
			cSlice[C.int32_t, int32](vals),
			C.GrB_Index(len(rowIndices)), cdup,
		))
	case []int64:
		info = Info(C.GrB_Matrix_build_INT64(
			matrix.grb, grbIndices(rowIndices), grbIndices(colIndices),
			cSlice[C.int64_t, int64](vals),
			C.GrB_Index(len(rowIndices)), cdup,
		))
	case []uint:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Matrix_build_UINT32(
				matrix.grb, grbIndices(rowIndices), grbIndices(colIndices),
				cSlice[C.uint32_t, uint](vals),
				C.GrB_Index(len(rowIndices)), cdup,
			))
		} else {
			info = Info(C.GrB_Matrix_build_UINT64(
				matrix.grb, grbIndices(rowIndices), grbIndices(colIndices),
				cSlice[C.uint64_t, uint](vals),
				C.GrB_Index(len(rowIndices)), cdup,
			))
		}
	case []uint8:
		info = Info(C.GrB_Matrix_build_UINT8(
			matrix.grb, grbIndices(rowIndices), grbIndices(colIndices),
			cSlice[C.uint8_t, uint8](vals),
			C.GrB_Index(len(rowIndices)), cdup,
		))
	case []uint16:
		info = Info(C.GrB_Matrix_build_UINT16(
			matrix.grb, grbIndices(rowIndices), grbIndices(colIndices),
			cSlice[C.uint16_t, uint16](vals),
			C.GrB_Index(len(rowIndices)), cdup,
		))
	case []uint32:
		info = Info(C.GrB_Matrix_build_UINT32(
			matrix.grb, grbIndices(rowIndices), grbIndices(colIndices),
			cSlice[C.uint32_t, uint32](vals),
			C.GrB_Index(len(rowIndices)), cdup,
		))
	case []uint64:
		info = Info(C.GrB_Matrix_build_UINT64(
			matrix.grb, grbIndices(rowIndices), grbIndices(colIndices),
			cSlice[C.uint64_t, uint64](vals),
			C.GrB_Index(len(rowIndices)), cdup,
		))
	case []float32:
		info = Info(C.GrB_Matrix_build_FP32(
			matrix.grb, grbIndices(rowIndices), grbIndices(colIndices),
			cSlice[C.float, float32](vals),
			C.GrB_Index(len(rowIndices)), cdup,
		))
	case []float64:
		info = Info(C.GrB_Matrix_build_FP64(
			matrix.grb, grbIndices(rowIndices), grbIndices(colIndices),
			cSlice[C.double, float64](vals),
			C.GrB_Index(len(rowIndices)), cdup,
		))
	case []complex64:
		info = Info(C.GxB_Matrix_build_FC32(
			matrix.grb, grbIndices(rowIndices), grbIndices(colIndices),
			cSlice[C.complexfloat, complex64](vals),
			C.GrB_Index(len(rowIndices)), cdup,
		))
	case []complex128:
		info = Info(C.GxB_Matrix_build_FC64(
			matrix.grb, grbIndices(rowIndices), grbIndices(colIndices),
			cSlice[C.complexdouble, complex128](vals),
			C.GrB_Index(len(rowIndices)), cdup,
		))
	default:
		info = Info(C.GrB_Matrix_build_UDT(
			matrix.grb, grbIndices(rowIndices), grbIndices(colIndices),
			unsafe.Pointer(unsafe.SliceData(values)),
			C.GrB_Index(len(rowIndices)), cdup,
		))
	}
	if info == success {
		return nil
	}
	return makeError(info)
}

// BuildScalar is like [Matrix.Build], except that the scalar is the value of all the tuples.
//
// Unlike [Matrix.Build], there is no dup operator to handle duplicate entries. Instead, any
// duplicates are silently ignored.
//
// BuildScalar is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) BuildScalar(rowIndices, colIndices []int, scalar Scalar[D]) error {
	if len(rowIndices) != len(colIndices) {
		return makeError(SliceMismatch)
	}
	for _, index := range rowIndices {
		if index < 0 {
			return makeError(InvalidIndex)
		}
	}
	for _, index := range colIndices {
		if index < 0 {
			return makeError(InvalidIndex)
		}
	}
	info := Info(C.GxB_Matrix_build_Scalar(
		matrix.grb, grbIndices(rowIndices), grbIndices(colIndices),
		scalar.grb, C.GrB_Index(len(rowIndices)),
	))
	if info == success {
		return nil
	}
	return makeError(info)
}

// SetElement sets one element of a matrix to a given value.
//
// To pass a [Scalar] object instead of a non-opaque variable, use [Matrix.SetElementScalar].
//
// Parameters:
//
//   - val (IN): Scalar to assign.
//
//   - rowIndex (IN): Row index of element to be assigned.
//
//   - colIndex (IN): Column index of element to be assigned.
//
// GraphBLAS API errors that may be returned:
//   - [DomainMismatch], [InvalidIndex], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (matrix Matrix[D]) SetElement(val D, rowIndex, colIndex int) error {
	if rowIndex < 0 || colIndex < 0 {
		return makeError(InvalidIndex)
	}
	var info Info
	switch value := any(val).(type) {
	case bool:
		info = Info(C.GrB_Matrix_setElement_BOOL(matrix.grb, C.bool(value), C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
	case int:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Matrix_setElement_INT32(matrix.grb, C.int32_t(value), C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
		} else {
			info = Info(C.GrB_Matrix_setElement_INT64(matrix.grb, C.int64_t(value), C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
		}
	case int8:
		info = Info(C.GrB_Matrix_setElement_INT8(matrix.grb, C.int8_t(value), C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
	case int16:
		info = Info(C.GrB_Matrix_setElement_INT16(matrix.grb, C.int16_t(value), C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
	case int32:
		info = Info(C.GrB_Matrix_setElement_INT32(matrix.grb, C.int32_t(value), C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
	case int64:
		info = Info(C.GrB_Matrix_setElement_INT64(matrix.grb, C.int64_t(value), C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
	case uint:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Matrix_setElement_UINT32(matrix.grb, C.uint32_t(value), C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
		} else {
			info = Info(C.GrB_Matrix_setElement_UINT64(matrix.grb, C.uint64_t(value), C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
		}
	case uint8:
		info = Info(C.GrB_Matrix_setElement_UINT8(matrix.grb, C.uint8_t(value), C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
	case uint16:
		info = Info(C.GrB_Matrix_setElement_UINT16(matrix.grb, C.uint16_t(value), C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
	case uint32:
		info = Info(C.GrB_Matrix_setElement_UINT32(matrix.grb, C.uint32_t(value), C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
	case uint64:
		info = Info(C.GrB_Matrix_setElement_UINT64(matrix.grb, C.uint64_t(value), C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
	case float32:
		info = Info(C.GrB_Matrix_setElement_FP32(matrix.grb, C.float(value), C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
	case float64:
		info = Info(C.GrB_Matrix_setElement_FP64(matrix.grb, C.double(value), C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
	case complex64:
		info = Info(C.GxB_Matrix_setElement_FC32(matrix.grb, C.complexfloat(value), C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
	case complex128:
		info = Info(C.GxB_Matrix_setElement_FC64(matrix.grb, C.complexdouble(value), C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
	default:
		info = Info(C.GrB_Matrix_setElement_UDT(matrix.grb, unsafe.Pointer(&val), C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
	}
	if info == success {
		return nil
	}
	return makeError(info)
}

// SetElementScalar is like [Matrix.SetElement], except that the scalar value is passed as a [Scalar]
// object. It may be empty.
func (matrix Matrix[D]) SetElementScalar(val Scalar[D], rowIndex, colIndex int) error {
	if rowIndex < 0 || colIndex < 0 {
		return makeError(InvalidIndex)
	}
	info := Info(C.GrB_Matrix_setElement_Scalar(matrix.grb, val.grb, C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// RemoveElement removes (annihilates) one stored element from a matrix.
//
// Parameters:
//
//   - rowIndex (IN): Row index of element to be removed.
//
//   - colIndex (IN): Column index of element to be removed.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidIndex], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (matrix Matrix[D]) RemoveElement(rowIndex, colIndex int) error {
	if rowIndex < 0 || colIndex < 0 {
		return makeError(InvalidIndex)
	}
	info := Info(C.GrB_Matrix_removeElement(matrix.grb, C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// ExtractElement extracts one element of a matrix.
//
// When there is no stored value at the specified location, ExtractElement returns
// ok == false. Otherwise, it returns ok == true.
//
// To store the element in a [Scalar] object instead of returning a non-opaque value,
// use [Matrix.ExtractElementScalar].
//
// Parameters:
//
//   - rowIndex (IN): Row index of element to be assigned.
//
//   - colIndex (IN): Column index of element to be assigned.
//
// GraphBLAS API errors that may be returned:
//   - [DomainMismatch], [InvalidIndex], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (matrix Matrix[D]) ExtractElement(rowIndex, colIndex int) (result D, ok bool, err error) {
	if rowIndex < 0 || colIndex < 0 {
		err = makeError(InvalidIndex)
		return
	}
	var info Info
	switch res := any(&result).(type) {
	case *bool:
		var cresult C.bool
		info = Info(C.GrB_Matrix_extractElement_BOOL(&cresult, matrix.grb, C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
		if info == success {
			*res = bool(cresult)
			ok = true
			return
		}
	case *int:
		if unsafe.Sizeof(0) == 4 {
			var cresult C.int32_t
			info = Info(C.GrB_Matrix_extractElement_INT32(&cresult, matrix.grb, C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
			if info == success {
				*res = int(cresult)
				ok = true
				return
			}
		} else {
			var cresult C.int64_t
			info = Info(C.GrB_Matrix_extractElement_INT64(&cresult, matrix.grb, C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
			if info == success {
				*res = int(cresult)
				ok = true
				return
			}
		}
	case *int8:
		var cresult C.int8_t
		info = Info(C.GrB_Matrix_extractElement_INT8(&cresult, matrix.grb, C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
		if info == success {
			*res = int8(cresult)
			ok = true
			return
		}
	case *int16:
		var cresult C.int16_t
		info = Info(C.GrB_Matrix_extractElement_INT16(&cresult, matrix.grb, C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
		if info == success {
			*res = int16(cresult)
			ok = true
			return
		}
	case *int32:
		var cresult C.int32_t
		info = Info(C.GrB_Matrix_extractElement_INT32(&cresult, matrix.grb, C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
		if info == success {
			*res = int32(cresult)
			ok = true
			return
		}
	case *int64:
		var cresult C.int64_t
		info = Info(C.GrB_Matrix_extractElement_INT64(&cresult, matrix.grb, C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
		if info == success {
			*res = int64(cresult)
			ok = true
			return
		}
	case *uint:
		if unsafe.Sizeof(0) == 4 {
			var cresult C.uint32_t
			info = Info(C.GrB_Matrix_extractElement_UINT32(&cresult, matrix.grb, C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
			if info == success {
				*res = uint(cresult)
				ok = true
				return
			}
		} else {
			var cresult C.uint64_t
			info = Info(C.GrB_Matrix_extractElement_UINT64(&cresult, matrix.grb, C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
			if info == success {
				*res = uint(cresult)
				ok = true
				return
			}
		}
	case *uint8:
		var cresult C.uint8_t
		info = Info(C.GrB_Matrix_extractElement_UINT8(&cresult, matrix.grb, C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
		if info == success {
			*res = uint8(cresult)
			ok = true
			return
		}
	case *uint16:
		var cresult C.uint16_t
		info = Info(C.GrB_Matrix_extractElement_UINT16(&cresult, matrix.grb, C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
		if info == success {
			*res = uint16(cresult)
			ok = true
			return
		}
	case *uint32:
		var cresult C.uint32_t
		info = Info(C.GrB_Matrix_extractElement_UINT32(&cresult, matrix.grb, C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
		if info == success {
			*res = uint32(cresult)
			ok = true
			return
		}
	case *uint64:
		var cresult C.uint64_t
		info = Info(C.GrB_Matrix_extractElement_UINT64(&cresult, matrix.grb, C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
		if info == success {
			*res = uint64(cresult)
			ok = true
			return
		}
	case *float32:
		var cresult C.float
		info = Info(C.GrB_Matrix_extractElement_FP32(&cresult, matrix.grb, C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
		if info == success {
			*res = float32(cresult)
			ok = true
			return
		}
	case *float64:
		var cresult C.double
		info = Info(C.GrB_Matrix_extractElement_FP64(&cresult, matrix.grb, C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
		if info == success {
			*res = float64(cresult)
			ok = true
			return
		}
	case *complex64:
		var cresult C.complexfloat
		info = Info(C.GxB_Matrix_extractElement_FC32(&cresult, matrix.grb, C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
		if info == success {
			*res = complex64(cresult)
			ok = true
			return
		}
	case *complex128:
		var cresult C.complexdouble
		info = Info(C.GxB_Matrix_extractElement_FC64(&cresult, matrix.grb, C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
		if info == success {
			*res = complex128(cresult)
			ok = true
			return
		}
	default:
		info = Info(C.GrB_Matrix_extractElement_UDT(unsafe.Pointer(&result), matrix.grb, C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
		if info == success {
			ok = true
			return
		}
	}
	if info == noValue {
		return
	}
	err = makeError(info)
	return
}

// ExtractElementScalar is like [Matrix.ExtractElement], except that the element is stored in a [Scalar]
// object.
//
// When there is no stored value at the specified location, the result becomes empty.
func (matrix Matrix[D]) ExtractElementScalar(result Scalar[D], rowIndex, colIndex int) error {
	if rowIndex < 0 || colIndex < 0 {
		return makeError(InvalidIndex)
	}
	info := Info(C.GrB_Matrix_extractElement_Scalar(result.grb, matrix.grb, C.GrB_Index(rowIndex), C.GrB_Index(colIndex)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// IsStoredElement determines whether there is a stored value at the specified
// location or not.
//
// Parameters:
//
//   - rowIndex (IN): Row index of element to be assigned.
//
//   - colIndex (IN): Column index of element to be assigned.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidIndex], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
//
// IsStoredElement is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) IsStoredElement(rowIndex, colIndex int) (ok bool, err error) {
	if rowIndex < 0 || colIndex < 0 {
		err = makeError(InvalidIndex)
		return
	}
	switch info := Info(C.GxB_Matrix_isStoredElement(matrix.grb, C.GrB_Index(rowIndex), C.GrB_Index(colIndex))); info {
	case success:
		return true, nil
	case noValue:
		return false, nil
	default:
		err = makeError(info)
		return
	}
}

// ExtractTuples extracts the contents of a GraphBLAS matrix into non-opaque slices,
// by appending the row indices, column indices, and values to the slices
// passed to this function (by using Go's built-in append function).
//
// Parameters:
//
//   - rowIndices (INOUT): Pointer to a slice of indices. If nil, ExtractTuples does not
//     produces the row indices of the matrix.
//
//   - colIndices (INOUT): Pointer to a slice of indices. If nil, ExtractTuples does not
//     produces the column indices of the matrix.
//
//   - values (INOUT): Pointer to a slice of indices. If nil, ExtractTuples does not
//     produces the values of the matrix.
//
// It is valid to pass pointers to nil slices, and ExtractTuples then produces the
// corresponding indices or values.
//
// GraphBLAS API errors that may be returned:
//   - [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (matrix Matrix[D]) ExtractTuples(rowIndices, colIndices *[]int, values *[]D) error {
	nvals, err := matrix.Nvals()
	if err != nil {
		return err
	}
	targetRowIndices, finalizeTargetRowIndices := growIndices(rowIndices, nvals)
	targetColIndices, finalizeTargetColIndices := growIndices(colIndices, nvals)
	targetValues := growslice(values, nvals)
	var info Info
	cnvals := C.GrB_Index(nvals)
	switch vals := any(targetValues).(type) {
	case []bool:
		info = Info(C.GrB_Matrix_extractTuples_BOOL(
			targetRowIndices, targetColIndices,
			cSlice[C.bool, bool](vals),
			&cnvals, matrix.grb,
		))
	case []int:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Matrix_extractTuples_INT32(
				targetRowIndices, targetColIndices,
				cSlice[C.int32_t, int](vals),
				&cnvals, matrix.grb,
			))
		} else {
			info = Info(C.GrB_Matrix_extractTuples_INT64(
				targetRowIndices, targetColIndices,
				cSlice[C.int64_t, int](vals),
				&cnvals, matrix.grb,
			))
		}
	case []int8:
		info = Info(C.GrB_Matrix_extractTuples_INT8(
			targetRowIndices, targetColIndices,
			cSlice[C.int8_t, int8](vals),
			&cnvals, matrix.grb,
		))
	case []int16:
		info = Info(C.GrB_Matrix_extractTuples_INT16(
			targetRowIndices, targetColIndices,
			cSlice[C.int16_t, int16](vals),
			&cnvals, matrix.grb,
		))
	case []int32:
		info = Info(C.GrB_Matrix_extractTuples_INT32(
			targetRowIndices, targetColIndices,
			cSlice[C.int32_t, int32](vals),
			&cnvals, matrix.grb,
		))
	case []int64:
		info = Info(C.GrB_Matrix_extractTuples_INT64(
			targetRowIndices, targetColIndices,
			cSlice[C.int64_t, int64](vals),
			&cnvals, matrix.grb,
		))
	case []uint:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Matrix_extractTuples_UINT32(
				targetRowIndices, targetColIndices,
				cSlice[C.uint32_t, uint](vals),
				&cnvals, matrix.grb,
			))
		} else {
			info = Info(C.GrB_Matrix_extractTuples_UINT64(
				targetRowIndices, targetColIndices,
				cSlice[C.uint64_t, uint](vals),
				&cnvals, matrix.grb,
			))
		}
	case []uint8:
		info = Info(C.GrB_Matrix_extractTuples_UINT8(
			targetRowIndices, targetColIndices,
			cSlice[C.uint8_t, uint8](vals),
			&cnvals, matrix.grb,
		))
	case []uint16:
		info = Info(C.GrB_Matrix_extractTuples_UINT16(
			targetRowIndices, targetColIndices,
			cSlice[C.uint16_t, uint16](vals),
			&cnvals, matrix.grb,
		))
	case []uint32:
		info = Info(C.GrB_Matrix_extractTuples_UINT32(
			targetRowIndices, targetColIndices,
			cSlice[C.uint32_t, uint32](vals),
			&cnvals, matrix.grb,
		))
	case []uint64:
		info = Info(C.GrB_Matrix_extractTuples_UINT64(
			targetRowIndices, targetColIndices,
			cSlice[C.uint64_t, uint64](vals),
			&cnvals, matrix.grb,
		))
	case []float32:
		info = Info(C.GrB_Matrix_extractTuples_FP32(
			targetRowIndices, targetColIndices,
			cSlice[C.float, float32](vals),
			&cnvals, matrix.grb,
		))
	case []float64:
		info = Info(C.GrB_Matrix_extractTuples_FP64(
			targetRowIndices, targetColIndices,
			cSlice[C.double, float64](vals),
			&cnvals, matrix.grb,
		))
	case []complex64:
		info = Info(C.GxB_Matrix_extractTuples_FC32(
			targetRowIndices, targetColIndices,
			cSlice[C.complexfloat, complex64](vals),
			&cnvals, matrix.grb,
		))
	case []complex128:
		info = Info(C.GxB_Matrix_extractTuples_FC64(
			targetRowIndices, targetColIndices,
			cSlice[C.complexdouble, complex128](vals),
			&cnvals, matrix.grb,
		))
	default:
		info = Info(C.GrB_Matrix_extractTuples_UDT(
			targetRowIndices, targetColIndices,
			unsafe.Pointer(unsafe.SliceData(targetValues)),
			&cnvals, matrix.grb,
		))
	}
	if info == success {
		if nvals != int(cnvals) {
			return makeError(InvalidObject)
		}
		finalizeTargetRowIndices()
		finalizeTargetColIndices()
		return nil
	}
	return makeError(info)
}

// Reshape changes the size of a matrix, taking its entries either column-wise or row-wise. If the
// matrix on input is nrows-by-ncols, and the requested dimensions on output are nrowsNew-by-ncolsNew,
// then the condition nrows*ncols == nrowsNew*ncolsNew must hold. The matrix is modified in-place.
//
// To create a new matrix, use [Matrix.ReshapeDup] instead.
//
// Parameters:
//
//   - byCol (IN): true if reshape by column, false if by row
//
//   - nrowsNew (IN): new number of rows
//
//   - ncolsNew (IN): new number of columns
//
// Reshape is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) Reshape(byCol bool, nrowsNew, ncolsNew int, desc *Descriptor) error {
	if nrowsNew < 0 || ncolsNew < 0 {
		return makeError(InvalidValue)
	}
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Matrix_reshape(matrix.grb, C.bool(byCol), C.GrB_Index(nrowsNew), C.GrB_Index(ncolsNew), cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// ReshapeDup is identical to [Matrix.Reshape], except that it creates a new output matrix
// instead of modified the matrix in-place.
//
// ReshapeDup is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) ReshapeDup(byCol bool, nrowsNew, ncolsNew int, desc *Descriptor) (dup Matrix[D], err error) {
	if nrowsNew < 0 || ncolsNew < 0 {
		err = makeError(InvalidValue)
		return
	}
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Matrix_reshapeDup(&dup.grb, matrix.grb, C.bool(byCol), C.GrB_Index(nrowsNew), C.GrB_Index(ncolsNew), cdesc))
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// ExportHint provides a hint as to which storage format might be most efficient
// for exporting the matrix with [Matrix.Export].
//
// If the implementation does not have a preferred format, it may return ok == false.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (matrix Matrix[D]) ExportHint() (format Format, ok bool, err error) {
	var cformat C.GrB_Format
	info := Info(C.GrB_Matrix_exportHint(&cformat, matrix.grb))
	switch info {
	case success:
		return Format(cformat), true, nil
	case noValue:
		return
	}
	err = makeError(info)
	return
}

// Export exports a GraphBLAS matrix to a pre-defined format.
//
// Parameters:
//
//   - format (IN): A value indicating the [Format] in which the matrix will be exported.
//
// Return Values:
//
//   - indptr: A slice that will hold row or column offsets, or row indices, depending
//     on the value of format.
//
//   - indices: A slice that will hold row or column indices of the elements in values,
//     depending on the value of format.
//
//   - values: A slice that will hold stored values.
//
// GraphBLAS API errors that may be returned:
//   - [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (matrix Matrix[D]) Export(format Format) (indptr, indices []int, values []D, err error) {
	var nindptr, nindices, nvalues C.GrB_Index
	info := Info(C.GrB_Matrix_exportSize(&nindptr, &nindices, &nvalues, C.GrB_Format(format), matrix.grb))
	if info != success {
		return nil, nil, nil, makeError(info)
	}
	cindptr := make([]C.GrB_Index, nindptr)
	cindices := make([]C.GrB_Index, nindices)
	values = make([]D, nvalues)
	switch vals := any(values).(type) {
	case []bool:
		info = Info(C.GrB_Matrix_export_BOOL(
			unsafe.SliceData(cindptr), unsafe.SliceData(cindices),
			cSlice[C.bool, bool](vals),
			&nindptr, &nindices, &nvalues, C.GrB_Format(format), matrix.grb,
		))
	case []int:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Matrix_export_INT32(
				unsafe.SliceData(cindptr), unsafe.SliceData(cindices),
				cSlice[C.int32_t, int](vals),
				&nindptr, &nindices, &nvalues, C.GrB_Format(format), matrix.grb,
			))
		} else {
			info = Info(C.GrB_Matrix_export_INT64(
				unsafe.SliceData(cindptr), unsafe.SliceData(cindices),
				cSlice[C.int64_t, int](vals),
				&nindptr, &nindices, &nvalues, C.GrB_Format(format), matrix.grb,
			))
		}
	case []int8:
		info = Info(C.GrB_Matrix_export_INT8(
			unsafe.SliceData(cindptr), unsafe.SliceData(cindices),
			cSlice[C.int8_t, int8](vals),
			&nindptr, &nindices, &nvalues, C.GrB_Format(format), matrix.grb,
		))
	case []int16:
		info = Info(C.GrB_Matrix_export_INT16(
			unsafe.SliceData(cindptr), unsafe.SliceData(cindices),
			cSlice[C.int16_t, int16](vals),
			&nindptr, &nindices, &nvalues, C.GrB_Format(format), matrix.grb,
		))
	case []int32:
		info = Info(C.GrB_Matrix_export_INT32(
			unsafe.SliceData(cindptr), unsafe.SliceData(cindices),
			cSlice[C.int32_t, int32](vals),
			&nindptr, &nindices, &nvalues, C.GrB_Format(format), matrix.grb,
		))
	case []int64:
		info = Info(C.GrB_Matrix_export_INT64(
			unsafe.SliceData(cindptr), unsafe.SliceData(cindices),
			cSlice[C.int64_t, int64](vals),
			&nindptr, &nindices, &nvalues, C.GrB_Format(format), matrix.grb,
		))
	case []uint:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Matrix_export_UINT32(
				unsafe.SliceData(cindptr), unsafe.SliceData(cindices),
				cSlice[C.uint32_t, uint](vals),
				&nindptr, &nindices, &nvalues, C.GrB_Format(format), matrix.grb,
			))
		} else {
			info = Info(C.GrB_Matrix_export_UINT64(
				unsafe.SliceData(cindptr), unsafe.SliceData(cindices),
				cSlice[C.uint64_t, uint](vals),
				&nindptr, &nindices, &nvalues, C.GrB_Format(format), matrix.grb,
			))
		}
	case []uint8:
		info = Info(C.GrB_Matrix_export_UINT8(
			unsafe.SliceData(cindptr), unsafe.SliceData(cindices),
			cSlice[C.uint8_t, uint8](vals),
			&nindptr, &nindices, &nvalues, C.GrB_Format(format), matrix.grb,
		))
	case []uint16:
		info = Info(C.GrB_Matrix_export_UINT16(
			unsafe.SliceData(cindptr), unsafe.SliceData(cindices),
			cSlice[C.uint16_t, uint16](vals),
			&nindptr, &nindices, &nvalues, C.GrB_Format(format), matrix.grb,
		))
	case []uint32:
		info = Info(C.GrB_Matrix_export_UINT32(
			unsafe.SliceData(cindptr), unsafe.SliceData(cindices),
			cSlice[C.uint32_t, uint32](vals),
			&nindptr, &nindices, &nvalues, C.GrB_Format(format), matrix.grb,
		))
	case []uint64:
		info = Info(C.GrB_Matrix_export_UINT64(
			unsafe.SliceData(cindptr), unsafe.SliceData(cindices),
			cSlice[C.uint64_t, uint64](vals),
			&nindptr, &nindices, &nvalues, C.GrB_Format(format), matrix.grb,
		))
	case []float32:
		info = Info(C.GrB_Matrix_export_FP32(
			unsafe.SliceData(cindptr), unsafe.SliceData(cindices),
			cSlice[C.float, float32](vals),
			&nindptr, &nindices, &nvalues, C.GrB_Format(format), matrix.grb,
		))
	case []float64:
		info = Info(C.GrB_Matrix_export_FP64(
			unsafe.SliceData(cindptr), unsafe.SliceData(cindices),
			cSlice[C.double, float64](vals),
			&nindptr, &nindices, &nvalues, C.GrB_Format(format), matrix.grb,
		))
	case []complex64:
		info = Info(C.GxB_Matrix_export_FC32(
			unsafe.SliceData(cindptr), unsafe.SliceData(cindices),
			cSlice[C.complexfloat, complex64](vals),
			&nindptr, &nindices, &nvalues, C.GrB_Format(format), matrix.grb,
		))
	case []complex128:
		info = Info(C.GxB_Matrix_export_FC64(
			unsafe.SliceData(cindptr), unsafe.SliceData(cindices),
			cSlice[C.complexdouble, complex128](vals),
			&nindptr, &nindices, &nvalues, C.GrB_Format(format), matrix.grb,
		))
	default:
		info = Info(C.GrB_Matrix_export_UDT(
			unsafe.SliceData(cindptr), unsafe.SliceData(cindices),
			unsafe.Pointer(unsafe.SliceData(values)),
			&nindptr, &nindices, &nvalues, C.GrB_Format(format), matrix.grb,
		))
	}
	if info == success {
		if int(nindptr) != len(cindptr) || int(nindices) != len(cindices) || int(nvalues) != len(values) {
			return nil, nil, nil, makeError(InvalidObject)
		}
		return goIndices(cindptr), goIndices(cindices), values, nil
	}
	return nil, nil, nil, makeError(info)
}

// MatrixImport imports a matrix into a GraphBLAS object.
//
// Parameters:
//
//   - indptr (IN): A slice of row or column offsets, or row indices, depending on the
//     value of format.
//
//   - indices (IN): A slice of row or column indices of the elements in values, depending
//     on the value of format.
//
//   - values (IN): A slice of values.
//
//   - format (IN): A value indicating the [Format] of the matrix being imported.
//
// GraphBLAS API errors that may be returned:
//   - [DomainMismatch], [InvalidValue], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [IndexOutOfBounds], [OutOfMemory], [Panic]
func MatrixImport[D any](
	nrows, ncols int,
	indptr, indices []int,
	values []D,
	format Format,
) (a Matrix[D], err error) {
	var d D
	dt, ok := grbType[TypeOf(d)]
	if !ok {
		err = makeError(UninitializedObject)
		return
	}
	if nrows < 0 || ncols < 0 {
		err = makeError(InvalidValue)
		return
	}
	for _, index := range indptr {
		if index < 0 {
			err = makeError(IndexOutOfBounds)
			return
		}
	}
	for _, index := range indices {
		if index < 0 {
			err = makeError(IndexOutOfBounds)
			return
		}
	}
	var info Info
	switch vals := any(values).(type) {
	case []bool:
		info = Info(C.GrB_Matrix_import_BOOL(
			&a.grb, dt, C.GrB_Index(nrows), C.GrB_Index(ncols),
			grbIndices(indptr),
			grbIndices(indices),
			cSlice[C.bool, bool](vals),
			C.GrB_Index(len(indptr)),
			C.GrB_Index(len(indices)),
			C.GrB_Index(len(values)),
			C.GrB_Format(format),
		))
	case []int:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Matrix_import_INT32(
				&a.grb, dt, C.GrB_Index(nrows), C.GrB_Index(ncols),
				grbIndices(indptr),
				grbIndices(indices),
				cSlice[C.int32_t, int](vals),
				C.GrB_Index(len(indptr)),
				C.GrB_Index(len(indices)),
				C.GrB_Index(len(values)),
				C.GrB_Format(format),
			))
		} else {
			info = Info(C.GrB_Matrix_import_INT64(
				&a.grb, dt, C.GrB_Index(nrows), C.GrB_Index(ncols),
				grbIndices(indptr),
				grbIndices(indices),
				cSlice[C.int64_t, int](vals),
				C.GrB_Index(len(indptr)),
				C.GrB_Index(len(indices)),
				C.GrB_Index(len(values)),
				C.GrB_Format(format),
			))
		}
	case []int8:
		info = Info(C.GrB_Matrix_import_INT8(
			&a.grb, dt, C.GrB_Index(nrows), C.GrB_Index(ncols),
			grbIndices(indptr),
			grbIndices(indices),
			cSlice[C.int8_t, int8](vals),
			C.GrB_Index(len(indptr)),
			C.GrB_Index(len(indices)),
			C.GrB_Index(len(values)),
			C.GrB_Format(format),
		))
	case []int16:
		info = Info(C.GrB_Matrix_import_INT16(
			&a.grb, dt, C.GrB_Index(nrows), C.GrB_Index(ncols),
			grbIndices(indptr),
			grbIndices(indices),
			cSlice[C.int16_t, int16](vals),
			C.GrB_Index(len(indptr)),
			C.GrB_Index(len(indices)),
			C.GrB_Index(len(values)),
			C.GrB_Format(format),
		))
	case []int32:
		info = Info(C.GrB_Matrix_import_INT32(
			&a.grb, dt, C.GrB_Index(nrows), C.GrB_Index(ncols),
			grbIndices(indptr),
			grbIndices(indices),
			cSlice[C.int32_t, int32](vals),
			C.GrB_Index(len(indptr)),
			C.GrB_Index(len(indices)),
			C.GrB_Index(len(values)),
			C.GrB_Format(format),
		))
	case []int64:
		info = Info(C.GrB_Matrix_import_INT64(
			&a.grb, dt, C.GrB_Index(nrows), C.GrB_Index(ncols),
			grbIndices(indptr),
			grbIndices(indices),
			cSlice[C.int64_t, int64](vals),
			C.GrB_Index(len(indptr)),
			C.GrB_Index(len(indices)),
			C.GrB_Index(len(values)),
			C.GrB_Format(format),
		))
	case []uint:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Matrix_import_UINT32(
				&a.grb, dt, C.GrB_Index(nrows), C.GrB_Index(ncols),
				grbIndices(indptr),
				grbIndices(indices),
				cSlice[C.uint32_t, uint](vals),
				C.GrB_Index(len(indptr)),
				C.GrB_Index(len(indices)),
				C.GrB_Index(len(values)),
				C.GrB_Format(format),
			))
		} else {
			info = Info(C.GrB_Matrix_import_UINT64(
				&a.grb, dt, C.GrB_Index(nrows), C.GrB_Index(ncols),
				grbIndices(indptr),
				grbIndices(indices),
				cSlice[C.uint64_t, uint](vals),
				C.GrB_Index(len(indptr)),
				C.GrB_Index(len(indices)),
				C.GrB_Index(len(values)),
				C.GrB_Format(format),
			))
		}
	case []uint8:
		info = Info(C.GrB_Matrix_import_UINT8(
			&a.grb, dt, C.GrB_Index(nrows), C.GrB_Index(ncols),
			grbIndices(indptr),
			grbIndices(indices),
			cSlice[C.uint8_t, uint8](vals),
			C.GrB_Index(len(indptr)),
			C.GrB_Index(len(indices)),
			C.GrB_Index(len(values)),
			C.GrB_Format(format),
		))
	case []uint16:
		info = Info(C.GrB_Matrix_import_UINT16(
			&a.grb, dt, C.GrB_Index(nrows), C.GrB_Index(ncols),
			grbIndices(indptr),
			grbIndices(indices),
			cSlice[C.uint16_t, uint16](vals),
			C.GrB_Index(len(indptr)),
			C.GrB_Index(len(indices)),
			C.GrB_Index(len(values)),
			C.GrB_Format(format),
		))
	case []uint32:
		info = Info(C.GrB_Matrix_import_UINT32(
			&a.grb, dt, C.GrB_Index(nrows), C.GrB_Index(ncols),
			grbIndices(indptr),
			grbIndices(indices),
			cSlice[C.uint32_t, uint32](vals),
			C.GrB_Index(len(indptr)),
			C.GrB_Index(len(indices)),
			C.GrB_Index(len(values)),
			C.GrB_Format(format),
		))
	case []uint64:
		info = Info(C.GrB_Matrix_import_UINT64(
			&a.grb, dt, C.GrB_Index(nrows), C.GrB_Index(ncols),
			grbIndices(indptr),
			grbIndices(indices),
			cSlice[C.uint64_t, uint64](vals),
			C.GrB_Index(len(indptr)),
			C.GrB_Index(len(indices)),
			C.GrB_Index(len(values)),
			C.GrB_Format(format),
		))
	case []float32:
		info = Info(C.GrB_Matrix_import_FP32(
			&a.grb, dt, C.GrB_Index(nrows), C.GrB_Index(ncols),
			grbIndices(indptr),
			grbIndices(indices),
			cSlice[C.float, float32](vals),
			C.GrB_Index(len(indptr)),
			C.GrB_Index(len(indices)),
			C.GrB_Index(len(values)),
			C.GrB_Format(format),
		))
	case []float64:
		info = Info(C.GrB_Matrix_import_FP64(
			&a.grb, dt, C.GrB_Index(nrows), C.GrB_Index(ncols),
			grbIndices(indptr),
			grbIndices(indices),
			cSlice[C.double, float64](vals),
			C.GrB_Index(len(indptr)),
			C.GrB_Index(len(indices)),
			C.GrB_Index(len(values)),
			C.GrB_Format(format),
		))
	case []complex64:
		info = Info(C.GxB_Matrix_import_FC32(
			&a.grb, dt, C.GrB_Index(nrows), C.GrB_Index(ncols),
			grbIndices(indptr),
			grbIndices(indices),
			cSlice[C.complexfloat, complex64](vals),
			C.GrB_Index(len(indptr)),
			C.GrB_Index(len(indices)),
			C.GrB_Index(len(values)),
			C.GrB_Format(format),
		))
	case []complex128:
		info = Info(C.GxB_Matrix_import_FC64(
			&a.grb, dt, C.GrB_Index(nrows), C.GrB_Index(ncols),
			grbIndices(indptr),
			grbIndices(indices),
			cSlice[C.complexdouble, complex128](vals),
			C.GrB_Index(len(indptr)),
			C.GrB_Index(len(indices)),
			C.GrB_Index(len(values)),
			C.GrB_Format(format),
		))
	default:
		info = Info(C.GrB_Matrix_import_UDT(
			&a.grb, dt, C.GrB_Index(nrows), C.GrB_Index(ncols),
			grbIndices(indptr),
			grbIndices(indices),
			unsafe.Pointer(unsafe.SliceData(values)),
			C.GrB_Index(len(indptr)),
			C.GrB_Index(len(indices)),
			C.GrB_Index(len(values)),
			C.GrB_Format(format),
		))
	}
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// SerializeSize computes the buffer size (in bytes) necessary to serialize the matrix using [Matrix.Serialize].
//
// GraphBLAS execution errors that may cause a panic:
//   - [OutOfMemory], [Panic]
func (matrix Matrix[D]) SerializeSize() (size int, err error) {
	var csize C.GrB_Index
	info := Info(C.GrB_Matrix_serializeSize(&csize, matrix.grb))
	if info == success {
		return int(csize), nil
	}
	err = makeError(info)
	return
}

// Serialize a GraphBLAS matrix object into an opaque slice of bytes.
// Serialize returns the number of bytes written to data.
//
// Parameters:
//
//   - data (INOUT): A preallocated buffer where the serialized matrix will be written.
//
// GraphBLAS API errors that may be returned:
//   - [InsufficientSpace], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (matrix Matrix[D]) Serialize(data []byte) (size int, err error) {
	csize := C.GrB_Index(len(data))
	info := Info(C.GrB_Matrix_serialize(unsafe.Pointer(unsafe.SliceData(data)), &csize, matrix.grb))
	if info == success {
		return int(csize), nil
	}
	err = makeError(info)
	return
}

// MatrixDeserialize constructs a new GraphBLAS matrix from a serialized object.
//
// Parameters:
//
//   - data (IN): A slice that holds a GraphBLAS matrix created with [Matrix.Serialize].
//
// GraphBLAS API errors that may be returned:
//   - [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func MatrixDeserialize[D any](data []byte) (a Matrix[D], err error) {
	var d D
	dt, ok := grbType[TypeOf(d)]
	if !ok {
		err = makeError(UninitializedObject)
		return
	}
	info := Info(C.GrB_Matrix_deserialize(&a.grb, dt, unsafe.Pointer(unsafe.SliceData(data)), C.GrB_Index(len(data))))
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// RowIteratorNew creates a row iterator and attaches it to the matrix.
//
// GraphBLAS API errors that may be returned:
//   - [NotImplemented]: The matrix cannot be iterated by row.
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
//
// RowIteratorNew is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) RowIteratorNew(desc *Descriptor) (it RowIterator[D], err error) {
	info := Info(C.GxB_Iterator_new(&it.grb))
	if info != success {
		err = makeError(info)
		return
	}
	cdesc := processDescriptor(desc)
	info = Info(C.GxB_rowIterator_attach(it.grb, matrix.grb, cdesc))
	if info == success {
		it.init()
		return
	}
	err = makeError(info)
	return
}

// ColIteratorNew creates a column iterator and attaches it to the matrix.
//
// GraphBLAS API errors that may be returned:
//   - [NotImplemented]: The matrix cannot be iterated by column.
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
//
// ColIteratorNew is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) ColIteratorNew(desc *Descriptor) (it ColIterator[D], err error) {
	info := Info(C.GxB_Iterator_new(&it.grb))
	if info != success {
		err = makeError(info)
		return
	}
	cdesc := processDescriptor(desc)
	info = Info(C.GxB_colIterator_attach(it.grb, matrix.grb, cdesc))
	if info == success {
		it.init()
		return
	}
	err = makeError(info)
	return
}

// IteratorNew creates a entry iterator and attaches it to the matrix.
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
//
// IteratorNew is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) IteratorNew(desc *Descriptor) (it EntryIterator[D], err error) {
	info := Info(C.GxB_Iterator_new(&it.grb))
	if info != success {
		err = makeError(info)
		return
	}
	cdesc := processDescriptor(desc)
	info = Info(C.GxB_Matrix_Iterator_attach(it.grb, matrix.grb, cdesc))
	if info == success {
		it.init()
		return
	}
	err = makeError(info)
	return
}

// Sort all the rows or all the columns of a matrix. Each row or column is sorted separately.
//
// Parameters:
//
//   - into (OUT): Contains the matrix of sorted values. If nil, this output is not produced.
//
//   - p (OUT): Contains the permutations of the sorted values. If nil, this output is not produced.
//
//   - op (IN): The comparator operation.
//
//   - desc (IN): If the [Tran] descriptor is set for [Inp0], then the columns of the matrix are sorted
//     instead of the rows.
//
// Sort is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) Sort(
	into *Matrix[D],
	p *Matrix[int],
	op BinaryOp[bool, D, D],
	desc *Descriptor,
) error {
	var cinto, cp C.GrB_Matrix
	if into == nil {
		cinto = C.GrB_Matrix(C.NULL)
	} else {
		cinto = into.grb
	}
	if p == nil {
		cp = C.GrB_Matrix(C.NULL)
	} else {
		cp = p.grb
	}
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Matrix_sort(cinto, cp, op.grb, matrix.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MemoryUsage returns the memory space required for a matrix, in bytes.
//
// MemoryUsage is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) MemoryUsage() (size int, err error) {
	var csize C.size_t
	info := Info(C.GxB_Matrix_memoryUsage(&csize, matrix.grb))
	if info == success {
		return int(csize), nil
	}
	err = makeError(info)
	return
}

// Iso returns true if the matrix is iso-valued, false otherwise.
//
// Iso is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) Iso() (iso bool, err error) {
	var ciso C.bool
	info := Info(C.GxB_Matrix_iso(&ciso, matrix.grb))
	if info == success {
		return bool(ciso), nil
	}
	err = makeError(info)
	return
}

// Concat concatenates a slice of matrices (tiles) into a single matrix.
//
// Concat is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) Concat(tiles []Matrix[D], m, n int, desc *Descriptor) error {
	if m <= 0 || n <= 0 {
		return makeError(InvalidValue)
	}
	if len(tiles) != m*n {
		return makeError(DimensionMismatch)
	}
	cdesc := processDescriptor(desc)
	ctiles := make([]C.GrB_Matrix, len(tiles))
	for i, tile := range tiles {
		ctiles[i] = tile.grb
	}
	info := Info(C.GxB_Matrix_concat(matrix.grb, unsafe.SliceData(ctiles), C.GrB_Index(m), C.GrB_Index(n), cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Split  a single input matrix into a 2D slice of tiles.
//
// Split is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) Split(tileNrows, tileNcols []int, desc *Descriptor) (tiles []Matrix[D], err error) {
	cdesc := processDescriptor(desc)
	m := len(tileNrows)
	n := len(tileNcols)
	ctiles := make([]C.GrB_Matrix, m*n)
	info := Info(C.GxB_Matrix_split(
		unsafe.SliceData(ctiles), C.GrB_Index(m), C.GrB_Index(n),
		grbIndices(tileNrows), grbIndices(tileNcols),
		matrix.grb, cdesc,
	))
	if info == success {
		tiles = make([]Matrix[D], m*n)
		for i, tile := range ctiles {
			tiles[i].grb = tile
		}
		return
	}
	return nil, makeError(info)
}

// BuildDiag is identical to [Vector.Diag], except for the extra [Descriptor]
// parameter, and this function is not a constructor. The matrix must already
// exist on input, of the correct size.
//
// BuildDiag is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) BuildDiag(v Vector[D], k int, desc *Descriptor) error {
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Matrix_diag(matrix.grb, v.grb, C.int64_t(int64(k)), cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// SetHyperSwitch determines how the matrix is converted between the hypersparse and
// non-hypersparse formats.
//
// Parameters:
//
//   - hyperSwitch (IN): A value between 0 and 1. To force a matrix to always be non-hypersparse,
//     use [NeverHyper]. To force a matrix to always stay hypersparse, use [AlwaysHyper].
//
// SetHyperSwitch is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) SetHyperSwitch(hyperSwitch float64) error {
	info := Info(C.GxB_Matrix_Option_set_FP64(matrix.grb, C.GxB_HYPER_SWITCH, C.double(hyperSwitch)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// GetHyperSwitch retrieves the current switch to hypersparse. See [Matrix.SetHyperSwitch].
//
// GetHyperSwitch is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) GetHyperSwitch() (hyperSwitch float64, err error) {
	var cHyperSwitch C.double
	info := Info(C.GxB_Matrix_Option_get_FP64(matrix.grb, C.GxB_HYPER_SWITCH, &cHyperSwitch))
	if info == success {
		return float64(cHyperSwitch), nil
	}
	err = makeError(info)
	return
}

// SetBitmapSwitch determines how the matrix is converted to the bitmap format.
//
// Parameters:
//
//   - bitmapSwitch (IN): A value between 0 and 1.
//
// SetBitmapSwitch is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) SetBitmapSwitch(bitmapSwitch float64) error {
	info := Info(C.GxB_Matrix_Option_set_FP64(matrix.grb, C.GxB_BITMAP_SWITCH, C.double(bitmapSwitch)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// GetBitmapSwitch retrieves the current switch to bitmap. See [Matrix.SetBitmapSwitch].
//
// GetBitmapSwitch is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) GetBitmapSwitch() (bitmapSwitch float64, err error) {
	var cBitmapSwitch C.double
	info := Info(C.GxB_Matrix_Option_get_FP64(matrix.grb, C.GxB_BITMAP_SWITCH, &cBitmapSwitch))
	if info == success {
		return float64(cBitmapSwitch), nil
	}
	err = makeError(info)
	return
}

// SetLayout sets the [Layout] (GxB_FORMAT) of the matrix.
//
// SetLayout is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) SetLayout(format Layout) error {
	info := Info(C.GxB_Matrix_Option_set_INT32(matrix.grb, C.GxB_FORMAT, C.int32_t(format)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// GetLayout retrieves the [Layout] (GxB_FORMAT) of the matrix.
//
// GetLayout a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) GetLayout() (format Layout, err error) {
	var cformat C.int32_t
	info := Info(C.GxB_Matrix_Option_get_INT32(matrix.grb, C.GxB_FORMAT, &cformat))
	if info == success {
		return Layout(cformat), nil
	}
	err = makeError(info)
	return
}

// SetSparsityControl determines the valid [Sparsity] format(s) for the matrix.
//
// SetSparsityControl is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) SetSparsityControl(sparsity Sparsity) error {
	info := Info(C.GxB_Matrix_Option_set_INT32(matrix.grb, C.GxB_SPARSITY_CONTROL, C.int32_t(sparsity)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// GetSparsityControl retrieves the valid [Sparsity] format(s) of the matrix.
//
// GetSparsityControl is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) GetSparsityControl() (sparsity Sparsity, err error) {
	var csparsity C.int32_t
	info := Info(C.GxB_Matrix_Option_get_INT32(matrix.grb, C.GxB_SPARSITY_CONTROL, &csparsity))
	if info == success {
		return Sparsity(csparsity), nil
	}
	err = makeError(info)
	return
}

// GetSparsityStatus retrieves the current [Sparsity] format of the matrix.
//
// GetSparsityStatus is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) GetSparsityStatus() (status Sparsity, err error) {
	var cstatus C.int32_t
	info := Info(C.GxB_Matrix_Option_get_INT32(matrix.grb, C.GxB_SPARSITY_STATUS, &cstatus))
	if info == success {
		return Sparsity(cstatus), nil
	}
	err = makeError(info)
	return
}

// Valid returns true if matrix has been created by a successful call to [Matrix.Dup], [MatrixDeserialize],
// [MatrixImport], [MatrixNew], or [Vector.Diag].
//
// Valid is a forGraphBLASGo extension. It is used in place of comparing against GrB_INVALID_HANDLE.
func (matrix Matrix[D]) Valid() bool {
	return matrix.grb != C.GrB_Matrix(C.GrB_INVALID_HANDLE)
}

// Free destroys a previously created [Matrix] and releases any resources associated with
// it. Calling Free on an object that is not [Matrix.Valid]() is legal.
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func (matrix *Matrix[D]) Free() error {
	info := Info(C.GrB_Matrix_free(&matrix.grb))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Wait until function calls in a sequence put the matrix into a state of completion or
// materialization.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [IndexOutOfBounds], [OutOfMemory], [Panic]
func (matrix Matrix[D]) Wait(mode WaitMode) error {
	info := Info(C.GrB_Matrix_wait(matrix.grb, C.GrB_WaitMode(mode)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Err returns an error message about any errors encountered during the processing associated with
// the matrix.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func (matrix Matrix[D]) Err() (string, error) {
	var cerror *C.char
	info := Info(C.GrB_Matrix_error(&cerror, matrix.grb))
	if info == success {
		return C.GoString(cerror), nil
	}
	return "", makeError(info)
}

// Print the contents of the matrix to stdout.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue]: The underlying print routine returned an I/O error.
//   - [NullPointer]: matrix is a nil pointer.
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [Panic]
//
// Print is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) Print(name string, pr PrintLevel) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	info := Info(C.GxB_Matrix_fprint(matrix.grb, cname, C.GxB_Print_Level(pr), (*C.FILE)(C.NULL)))
	if info == success {
		return nil
	}
	return makeError(info)
}
