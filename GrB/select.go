package GrB

// #include "GraphBLAS.h"
import "C"
import "unsafe"

// VectorSelect applies a select operator (an index unary operator) to the elements of a vector to determine
// whether or not to keep them.
//
// To pass a [Scalar] object instead of a non-opaque variable, use [VectorSelectScalar].
//
// Parameters:
//
//   - w (INOUT): An existing GraphBLAS vector. On input, the vector provides values
//     that may be accumulated with the result of the select operation. On output,
//     this vector holds the results of the operation.
//
//   - mask (IN): An optional "write" [VectorMask].
//
//   - accum (IN): An optional binary operator used for accumulating entries into existing w
//     entries. If assignment rather than accumulation is desired, nil should be specified.
//
//   - op (IN): An index unary operator applied to each element stored in the input vector u. It is a function
//     of the stored element's value, its location index, and a user supplied scalar value val.
//
//   - u (IN): The GraphBLAS vector whose elements are passed to the index unary operator.
//
//   - val (IN): An additional scalar value that is passed to the index unary operator.
//
//   - desc (IN): An optional operation [Descriptor]. If a default descriptor is desired,
//     nil should be specified.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func VectorSelect[D, T any](
	w Vector[D],
	mask *Vector[bool],
	accum *BinaryOp[D, D, D],
	op IndexUnaryOp[bool, D, T],
	u Vector[D],
	val T,
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	var info Info
	switch x := any(val).(type) {
	case bool:
		info = Info(C.GrB_Vector_select_BOOL(w.grb, cmask, caccum, op.grb, u.grb, C.bool(x), cdesc))
	case int:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Vector_select_INT32(w.grb, cmask, caccum, op.grb, u.grb, C.int32_t(x), cdesc))
		} else {
			info = Info(C.GrB_Vector_select_INT64(w.grb, cmask, caccum, op.grb, u.grb, C.int64_t(x), cdesc))
		}
	case int8:
		info = Info(C.GrB_Vector_select_INT8(w.grb, cmask, caccum, op.grb, u.grb, C.int8_t(x), cdesc))
	case int16:
		info = Info(C.GrB_Vector_select_INT16(w.grb, cmask, caccum, op.grb, u.grb, C.int16_t(x), cdesc))
	case int32:
		info = Info(C.GrB_Vector_select_INT32(w.grb, cmask, caccum, op.grb, u.grb, C.int32_t(x), cdesc))
	case int64:
		info = Info(C.GrB_Vector_select_INT64(w.grb, cmask, caccum, op.grb, u.grb, C.int64_t(x), cdesc))
	case uint:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Vector_select_UINT32(w.grb, cmask, caccum, op.grb, u.grb, C.uint32_t(x), cdesc))
		} else {
			info = Info(C.GrB_Vector_select_UINT64(w.grb, cmask, caccum, op.grb, u.grb, C.uint64_t(x), cdesc))
		}
	case uint8:
		info = Info(C.GrB_Vector_select_UINT8(w.grb, cmask, caccum, op.grb, u.grb, C.uint8_t(x), cdesc))
	case uint16:
		info = Info(C.GrB_Vector_select_UINT16(w.grb, cmask, caccum, op.grb, u.grb, C.uint16_t(x), cdesc))
	case uint32:
		info = Info(C.GrB_Vector_select_UINT32(w.grb, cmask, caccum, op.grb, u.grb, C.uint32_t(x), cdesc))
	case uint64:
		info = Info(C.GrB_Vector_select_UINT64(w.grb, cmask, caccum, op.grb, u.grb, C.uint64_t(x), cdesc))
	case float32:
		info = Info(C.GrB_Vector_select_FP32(w.grb, cmask, caccum, op.grb, u.grb, C.float(x), cdesc))
	case float64:
		info = Info(C.GrB_Vector_select_FP64(w.grb, cmask, caccum, op.grb, u.grb, C.double(x), cdesc))
	case complex64:
		info = Info(C.GxB_Vector_select_FC32(w.grb, cmask, caccum, op.grb, u.grb, C.complexfloat(x), cdesc))
	case complex128:
		info = Info(C.GxB_Vector_select_FC64(w.grb, cmask, caccum, op.grb, u.grb, C.complexdouble(x), cdesc))
	default:
		info = Info(C.GrB_Vector_select_UDT(w.grb, cmask, caccum, op.grb, u.grb, unsafe.Pointer(&val), cdesc))
	}
	if info == success {
		return nil
	}
	return makeError(info)
}

// VectorSelectScalar is like [VectorSelect], except that the scalar value is passed as a [Scalar]
// object. It must not be empty.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [EmptyObject], [InvalidObject], [OutOfMemory], [Panic]
func VectorSelectScalar[D, T any](
	w Vector[D],
	mask *Vector[bool],
	accum *BinaryOp[D, D, D],
	op IndexUnaryOp[bool, D, T],
	u Vector[D],
	val Scalar[T],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	info := Info(C.GrB_Vector_select_Scalar(w.grb, cmask, caccum, op.grb, u.grb, val.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixSelect applies a select operator (an index unary operator to the elements of a matrix to determine
// whether or not to keep them.
//
// To pass a [Scalar] object instead of a non-opaque variable, use [MatrixSelectScalar].
//
// Parameters:
//
//   - c (INOUT): An existing GraphBLAS matrix. On input, the matrix provides values
//     that may be accumulated with the result of the select operation. On output,
//     this matrix holds the results of the operation.
//
//   - mask (IN): An optional "write" [MatrixMask].
//
//   - accum (IN): An optional binary operator used for accumulating entries into existing c
//     entries. If assignment rather than accumulation is desired, nil should be specified.
//
//   - op (IN): An index unary operator applied to each element stored in the input vector u. It is a function
//     of the stored element's value, its row and column indices, and a user supplied scalar value val.
//
//   - a (IN): The GraphBLAS matrix whose elements are passed to the index unary operator.
//
//   - val (IN): An additional scalar value that is passed to the index unary operator.
//
//   - desc (IN): An optional operation [Descriptor]. If a default descriptor is desired,
//     nil should be specified.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func MatrixSelect[D, T any](
	c Matrix[D],
	mask *Matrix[bool],
	accum *BinaryOp[D, D, D],
	op IndexUnaryOp[bool, D, T],
	a Matrix[D],
	val T,
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADM(mask, accum, desc)
	var info Info
	switch x := any(val).(type) {
	case bool:
		info = Info(C.GrB_Matrix_select_BOOL(c.grb, cmask, caccum, op.grb, a.grb, C.bool(x), cdesc))
	case int:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Matrix_select_INT32(c.grb, cmask, caccum, op.grb, a.grb, C.int32_t(x), cdesc))
		} else {
			info = Info(C.GrB_Matrix_select_INT64(c.grb, cmask, caccum, op.grb, a.grb, C.int64_t(x), cdesc))
		}
	case int8:
		info = Info(C.GrB_Matrix_select_INT8(c.grb, cmask, caccum, op.grb, a.grb, C.int8_t(x), cdesc))
	case int16:
		info = Info(C.GrB_Matrix_select_INT16(c.grb, cmask, caccum, op.grb, a.grb, C.int16_t(x), cdesc))
	case int32:
		info = Info(C.GrB_Matrix_select_INT32(c.grb, cmask, caccum, op.grb, a.grb, C.int32_t(x), cdesc))
	case int64:
		info = Info(C.GrB_Matrix_select_INT64(c.grb, cmask, caccum, op.grb, a.grb, C.int64_t(x), cdesc))
	case uint:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Matrix_select_UINT32(c.grb, cmask, caccum, op.grb, a.grb, C.uint32_t(x), cdesc))
		} else {
			info = Info(C.GrB_Matrix_select_UINT64(c.grb, cmask, caccum, op.grb, a.grb, C.uint64_t(x), cdesc))
		}
	case uint8:
		info = Info(C.GrB_Matrix_select_UINT8(c.grb, cmask, caccum, op.grb, a.grb, C.uint8_t(x), cdesc))
	case uint16:
		info = Info(C.GrB_Matrix_select_UINT16(c.grb, cmask, caccum, op.grb, a.grb, C.uint16_t(x), cdesc))
	case uint32:
		info = Info(C.GrB_Matrix_select_UINT32(c.grb, cmask, caccum, op.grb, a.grb, C.uint32_t(x), cdesc))
	case uint64:
		info = Info(C.GrB_Matrix_select_UINT64(c.grb, cmask, caccum, op.grb, a.grb, C.uint64_t(x), cdesc))
	case float32:
		info = Info(C.GrB_Matrix_select_FP32(c.grb, cmask, caccum, op.grb, a.grb, C.float(x), cdesc))
	case float64:
		info = Info(C.GrB_Matrix_select_FP64(c.grb, cmask, caccum, op.grb, a.grb, C.double(x), cdesc))
	case complex64:
		info = Info(C.GxB_Matrix_select_FC32(c.grb, cmask, caccum, op.grb, a.grb, C.complexfloat(x), cdesc))
	case complex128:
		info = Info(C.GxB_Matrix_select_FC64(c.grb, cmask, caccum, op.grb, a.grb, C.complexdouble(x), cdesc))
	default:
		info = Info(C.GrB_Matrix_select_UDT(c.grb, cmask, caccum, op.grb, a.grb, unsafe.Pointer(&val), cdesc))
	}
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixSelectScalar is like [MatrixSelect], except that the scalar value is passed as a [Scalar]
// object. It must not be empty.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [EmptyObject], [InvalidObject], [OutOfMemory], [Panic]
func MatrixSelectScalar[D, T any](
	c Matrix[D],
	mask *Matrix[bool],
	accum *BinaryOp[D, D, D],
	op IndexUnaryOp[bool, D, T],
	a Matrix[D],
	val Scalar[T],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADM(mask, accum, desc)
	info := Info(C.GrB_Matrix_select_Scalar(c.grb, cmask, caccum, op.grb, a.grb, val.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}
