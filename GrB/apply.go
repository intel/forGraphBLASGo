package GrB

// #include "GraphBLAS.h"
import "C"
import "unsafe"

// VectorApply computes the transformation of the values of the elements of a vector
// using a unary function.
//
// Parameters:
//
//   - w (INOUT): An existing GraphBLAS vector. On input, the vector provides values
//     that may be accumulated with the result of the apply operation. On output,
//     this vector holds the results of the operation.
//
//   - mask (IN): An optional "write" [VectorMask].
//
//   - accum (IN): An optional binary operator used for accumulating entries into existing w
//     entries. If assignment rather than accumulation is desired, nil should be specified.
//
//   - op (IN): A unary operator applied to each element of input vector u.
//
//   - u (IN): The GraphBLAS vector to which the unary function is applied.
//
//   - desc (IN): An optional operation [Descriptor]. If a default descriptor is desired,
//     nil should be specified.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func VectorApply[Dw, Du any](
	w Vector[Dw],
	mask *Vector[bool],
	accum *BinaryOp[Dw, Dw, Dw],
	op UnaryOp[Dw, Du],
	u Vector[Du],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	info := Info(C.GrB_Vector_apply(w.grb, cmask, caccum, op.grb, u.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixApply computes the transformation of the values of the elements of a matrix
// using a unary function.
//
// Parameters:
//
//   - c (INOUT): An existing GraphBLAS matrix. On input, the matrix provides values
//     that may be accumulated with the result of the apply operation. On output,
//     this matrix holds the results of the operation.
//
//   - mask (IN): An optional "write" [MatrixMask].
//
//   - accum (IN): An optional binary operator used for accumulating entries into existing c
//     entries. If assignment rather than accumulation is desired, nil should be specified.
//
//   - op (IN): A unary operator applied to each element of input matrix a.
//
//   - a (IN): The GraphBLAS matrix to which the unary function is applied.
//
//   - desc (IN): An optional operation [Descriptor]. If a default descriptor is desired,
//     nil should be specified.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func MatrixApply[DC, DA any](
	c Matrix[DC],
	mask *Matrix[bool],
	accum *BinaryOp[DC, DC, DC],
	op UnaryOp[DC, DA],
	a Matrix[DA],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADM(mask, accum, desc)
	info := Info(C.GrB_Matrix_apply(c.grb, cmask, caccum, op.grb, a.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// VectorApplyBinaryOp1st computes the transformation of the values of the stored elements of a vector using a binary
// operator and a scalar value. The specified scalar value is passed as the first argument to the binary operator and
// stored elements of the vector are passed as the second argument. The scalar is passed as a non-opaque variable.
//
// To pass a [Scalar] object instead of a non-opaque variable, use [VectorApplyBinaryOp1stScalar].
//
// To pass the stored elements of the vector as the first argument to the binary operator, and the specified scalar
// value as the second argument, use [VectorApplyBinaryOp2nd] or [VectorApplyBinaryOp2ndScalar].
//
// Parameters:
//
//   - w (INOUT): An existing GraphBLAS vector. On input, the vector provides values
//     that may be accumulated with the result of the apply operation. On output,
//     this vector holds the results of the operation.
//
//   - mask (IN): An optional "write" [VectorMask].
//
//   - accum (IN): An optional binary operator used for accumulating entries into existing w
//     entries. If assignment rather than accumulation is desired, nil should be specified.
//
//   - op (IN): A binary operator applied to the scalar value val and each element of input vector u.
//
//   - val (IN): Scalar value that is passed to the binary operator as the left-hand (first) argument.
//
//   - u (IN): The GraphBLAS vector whose elements are passed to the binary operator as the right-hand
//     (second) argument.
//
//   - desc (IN): An optional operation [Descriptor]. If a default descriptor is desired,
//     nil should be specified.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func VectorApplyBinaryOp1st[Dw, D, Du any](
	w Vector[Dw],
	mask *Vector[bool],
	accum *BinaryOp[Dw, Dw, Dw],
	op BinaryOp[Dw, D, Du],
	val D,
	u Vector[Du],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	var info Info
	switch x := any(val).(type) {
	case bool:
		info = Info(C.GrB_Vector_apply_BinaryOp1st_BOOL(w.grb, cmask, caccum, op.grb, C.bool(x), u.grb, cdesc))
	case int:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Vector_apply_BinaryOp1st_INT32(w.grb, cmask, caccum, op.grb, C.int32_t(x), u.grb, cdesc))
		} else {
			info = Info(C.GrB_Vector_apply_BinaryOp1st_INT64(w.grb, cmask, caccum, op.grb, C.int64_t(x), u.grb, cdesc))
		}
	case int8:
		info = Info(C.GrB_Vector_apply_BinaryOp1st_INT8(w.grb, cmask, caccum, op.grb, C.int8_t(x), u.grb, cdesc))
	case int16:
		info = Info(C.GrB_Vector_apply_BinaryOp1st_INT16(w.grb, cmask, caccum, op.grb, C.int16_t(x), u.grb, cdesc))
	case int32:
		info = Info(C.GrB_Vector_apply_BinaryOp1st_INT32(w.grb, cmask, caccum, op.grb, C.int32_t(x), u.grb, cdesc))
	case int64:
		info = Info(C.GrB_Vector_apply_BinaryOp1st_INT64(w.grb, cmask, caccum, op.grb, C.int64_t(x), u.grb, cdesc))
	case uint:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Vector_apply_BinaryOp1st_UINT32(w.grb, cmask, caccum, op.grb, C.uint32_t(x), u.grb, cdesc))
		} else {
			info = Info(C.GrB_Vector_apply_BinaryOp1st_UINT64(w.grb, cmask, caccum, op.grb, C.uint64_t(x), u.grb, cdesc))
		}
	case uint8:
		info = Info(C.GrB_Vector_apply_BinaryOp1st_UINT8(w.grb, cmask, caccum, op.grb, C.uint8_t(x), u.grb, cdesc))
	case uint16:
		info = Info(C.GrB_Vector_apply_BinaryOp1st_UINT16(w.grb, cmask, caccum, op.grb, C.uint16_t(x), u.grb, cdesc))
	case uint32:
		info = Info(C.GrB_Vector_apply_BinaryOp1st_UINT32(w.grb, cmask, caccum, op.grb, C.uint32_t(x), u.grb, cdesc))
	case uint64:
		info = Info(C.GrB_Vector_apply_BinaryOp1st_UINT64(w.grb, cmask, caccum, op.grb, C.uint64_t(x), u.grb, cdesc))
	case float32:
		info = Info(C.GrB_Vector_apply_BinaryOp1st_FP32(w.grb, cmask, caccum, op.grb, C.float(x), u.grb, cdesc))
	case float64:
		info = Info(C.GrB_Vector_apply_BinaryOp1st_FP64(w.grb, cmask, caccum, op.grb, C.double(x), u.grb, cdesc))
	case complex64:
		info = Info(C.GxB_Vector_apply_BinaryOp1st_FC32(w.grb, cmask, caccum, op.grb, C.complexfloat(x), u.grb, cdesc))
	case complex128:
		info = Info(C.GxB_Vector_apply_BinaryOp1st_FC64(w.grb, cmask, caccum, op.grb, C.complexdouble(x), u.grb, cdesc))
	default:
		info = Info(C.GrB_Vector_apply_BinaryOp1st_UDT(w.grb, cmask, caccum, op.grb, unsafe.Pointer(&val), u.grb, cdesc))
	}
	if info == success {
		return nil
	}
	return makeError(info)
}

// VectorApplyBinaryOp1stScalar is like [VectorApplyBinaryOp1st], except that the scalar value is passed as a [Scalar]
// object. It must not be empty.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [EmptyObject], [InvalidObject], [OutOfMemory], [Panic]
func VectorApplyBinaryOp1stScalar[Dw, D, Du any](
	w Vector[Dw],
	mask *Vector[bool],
	accum *BinaryOp[Dw, Dw, Dw],
	op BinaryOp[Dw, D, Du],
	val Scalar[D],
	u Vector[Du],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	info := Info(C.GrB_Vector_apply_BinaryOp1st_Scalar(w.grb, cmask, caccum, op.grb, val.grb, u.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// VectorApplyBinaryOp2nd is like [VectorApplyBinaryOp1st], except that the stored elements of the vector are
// passed as the first argument to the binary operator and the specified scalar value is passed as the second argument.
func VectorApplyBinaryOp2nd[Dw, Du, D any](
	w Vector[Dw],
	mask *Vector[bool],
	accum *BinaryOp[Dw, Dw, Dw],
	op BinaryOp[Dw, Du, D],
	u Vector[Du],
	val D,
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	var info Info
	switch x := any(val).(type) {
	case bool:
		info = Info(C.GrB_Vector_apply_BinaryOp2nd_BOOL(w.grb, cmask, caccum, op.grb, u.grb, C.bool(x), cdesc))
	case int:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Vector_apply_BinaryOp2nd_INT32(w.grb, cmask, caccum, op.grb, u.grb, C.int32_t(x), cdesc))
		} else {
			info = Info(C.GrB_Vector_apply_BinaryOp2nd_INT64(w.grb, cmask, caccum, op.grb, u.grb, C.int64_t(x), cdesc))
		}
	case int8:
		info = Info(C.GrB_Vector_apply_BinaryOp2nd_INT8(w.grb, cmask, caccum, op.grb, u.grb, C.int8_t(x), cdesc))
	case int16:
		info = Info(C.GrB_Vector_apply_BinaryOp2nd_INT16(w.grb, cmask, caccum, op.grb, u.grb, C.int16_t(x), cdesc))
	case int32:
		info = Info(C.GrB_Vector_apply_BinaryOp2nd_INT32(w.grb, cmask, caccum, op.grb, u.grb, C.int32_t(x), cdesc))
	case int64:
		info = Info(C.GrB_Vector_apply_BinaryOp2nd_INT64(w.grb, cmask, caccum, op.grb, u.grb, C.int64_t(x), cdesc))
	case uint:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Vector_apply_BinaryOp2nd_UINT32(w.grb, cmask, caccum, op.grb, u.grb, C.uint32_t(x), cdesc))
		} else {
			info = Info(C.GrB_Vector_apply_BinaryOp2nd_UINT64(w.grb, cmask, caccum, op.grb, u.grb, C.uint64_t(x), cdesc))
		}
	case uint8:
		info = Info(C.GrB_Vector_apply_BinaryOp2nd_UINT8(w.grb, cmask, caccum, op.grb, u.grb, C.uint8_t(x), cdesc))
	case uint16:
		info = Info(C.GrB_Vector_apply_BinaryOp2nd_UINT16(w.grb, cmask, caccum, op.grb, u.grb, C.uint16_t(x), cdesc))
	case uint32:
		info = Info(C.GrB_Vector_apply_BinaryOp2nd_UINT32(w.grb, cmask, caccum, op.grb, u.grb, C.uint32_t(x), cdesc))
	case uint64:
		info = Info(C.GrB_Vector_apply_BinaryOp2nd_UINT64(w.grb, cmask, caccum, op.grb, u.grb, C.uint64_t(x), cdesc))
	case float32:
		info = Info(C.GrB_Vector_apply_BinaryOp2nd_FP32(w.grb, cmask, caccum, op.grb, u.grb, C.float(x), cdesc))
	case float64:
		info = Info(C.GrB_Vector_apply_BinaryOp2nd_FP64(w.grb, cmask, caccum, op.grb, u.grb, C.double(x), cdesc))
	case complex64:
		info = Info(C.GxB_Vector_apply_BinaryOp2nd_FC32(w.grb, cmask, caccum, op.grb, u.grb, C.complexfloat(x), cdesc))
	case complex128:
		info = Info(C.GxB_Vector_apply_BinaryOp2nd_FC64(w.grb, cmask, caccum, op.grb, u.grb, C.complexdouble(x), cdesc))
	default:
		info = Info(C.GrB_Vector_apply_BinaryOp2nd_UDT(w.grb, cmask, caccum, op.grb, u.grb, unsafe.Pointer(&val), cdesc))
	}
	if info == success {
		return nil
	}
	return makeError(info)
}

// VectorApplyBinaryOp2ndScalar is like [VectorApplyBinaryOp2nd], except that the scalar value is passed as a
// [Scalar] object. It must not be empty.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [EmptyObject], [InvalidObject], [OutOfMemory], [Panic]
func VectorApplyBinaryOp2ndScalar[Dw, Du, D any](
	w Vector[Dw],
	mask *Vector[bool],
	accum *BinaryOp[Dw, Dw, Dw],
	op BinaryOp[Dw, Du, D],
	u Vector[Du],
	val Scalar[D],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	info := Info(C.GrB_Vector_apply_BinaryOp2nd_Scalar(w.grb, cmask, caccum, op.grb, u.grb, val.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixApplyBinaryOp1st computes the transformation of the values of the stored elements of a matrix using a binary
// operator and a scalar value. The specified scalar value is passed as the first argument to the binary operator and
// stored elements of the matrix are passed as the second argument. The scalar is passed as a non-opaque variable.
//
// To pass a [Scalar] object instead of a non-opaque variable, use [MatrixApplyBinaryOp1stScalar].
//
// To pass the stored elements of the matrix as the first argument to the binary operator and the specified scalar
// value as the second argument, use [MatrixApplyBinaryOp2nd] or [MatrixApplyBinaryOp2ndScalar].
//
// Parameters:
//
//   - c (INOUT): An existing GraphBLAS matrix. On input, the matrix provides values
//     that may be accumulated with the result of the apply operation. On output,
//     this matrix holds the results of the operation.
//
//   - mask (IN): An optional "write" [MatrixMask].
//
//   - accum (IN): An optional binary operator used for accumulating entries into existing c
//     entries. If assignment rather than accumulation is desired, nil should be specified.
//
//   - op (IN): A binary operator applied to the scalar value val and each element of input matrix a.
//
//   - val (IN): Scalar value that is passed to the binary operator as the left-hand (first) argument.
//
//   - a (IN): The GraphBLAS matrix whose elements are passed to the binary operator as the right-hand
//     (second) argument.
//
//   - desc (IN): An optional operation [Descriptor]. If a default descriptor is desired,
//     nil should be specified.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func MatrixApplyBinaryOp1st[DC, D, DA any](
	c Matrix[DC],
	mask *Matrix[bool],
	accum *BinaryOp[DC, DC, DC],
	op BinaryOp[DC, D, DA],
	val D,
	a Matrix[DA],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADM(mask, accum, desc)
	var info Info
	switch x := any(val).(type) {
	case bool:
		info = Info(C.GrB_Matrix_apply_BinaryOp1st_BOOL(c.grb, cmask, caccum, op.grb, C.bool(x), a.grb, cdesc))
	case int:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Matrix_apply_BinaryOp1st_INT32(c.grb, cmask, caccum, op.grb, C.int32_t(x), a.grb, cdesc))
		} else {
			info = Info(C.GrB_Matrix_apply_BinaryOp1st_INT64(c.grb, cmask, caccum, op.grb, C.int64_t(x), a.grb, cdesc))
		}
	case int8:
		info = Info(C.GrB_Matrix_apply_BinaryOp1st_INT8(c.grb, cmask, caccum, op.grb, C.int8_t(x), a.grb, cdesc))
	case int16:
		info = Info(C.GrB_Matrix_apply_BinaryOp1st_INT16(c.grb, cmask, caccum, op.grb, C.int16_t(x), a.grb, cdesc))
	case int32:
		info = Info(C.GrB_Matrix_apply_BinaryOp1st_INT32(c.grb, cmask, caccum, op.grb, C.int32_t(x), a.grb, cdesc))
	case int64:
		info = Info(C.GrB_Matrix_apply_BinaryOp1st_INT64(c.grb, cmask, caccum, op.grb, C.int64_t(x), a.grb, cdesc))
	case uint:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Matrix_apply_BinaryOp1st_UINT32(c.grb, cmask, caccum, op.grb, C.uint32_t(x), a.grb, cdesc))
		} else {
			info = Info(C.GrB_Matrix_apply_BinaryOp1st_UINT64(c.grb, cmask, caccum, op.grb, C.uint64_t(x), a.grb, cdesc))
		}
	case uint8:
		info = Info(C.GrB_Matrix_apply_BinaryOp1st_UINT8(c.grb, cmask, caccum, op.grb, C.uint8_t(x), a.grb, cdesc))
	case uint16:
		info = Info(C.GrB_Matrix_apply_BinaryOp1st_UINT16(c.grb, cmask, caccum, op.grb, C.uint16_t(x), a.grb, cdesc))
	case uint32:
		info = Info(C.GrB_Matrix_apply_BinaryOp1st_UINT32(c.grb, cmask, caccum, op.grb, C.uint32_t(x), a.grb, cdesc))
	case uint64:
		info = Info(C.GrB_Matrix_apply_BinaryOp1st_UINT64(c.grb, cmask, caccum, op.grb, C.uint64_t(x), a.grb, cdesc))
	case float32:
		info = Info(C.GrB_Matrix_apply_BinaryOp1st_FP32(c.grb, cmask, caccum, op.grb, C.float(x), a.grb, cdesc))
	case float64:
		info = Info(C.GrB_Matrix_apply_BinaryOp1st_FP64(c.grb, cmask, caccum, op.grb, C.double(x), a.grb, cdesc))
	case complex64:
		info = Info(C.GxB_Matrix_apply_BinaryOp1st_FC32(c.grb, cmask, caccum, op.grb, C.complexfloat(x), a.grb, cdesc))
	case complex128:
		info = Info(C.GxB_Matrix_apply_BinaryOp1st_FC64(c.grb, cmask, caccum, op.grb, C.complexdouble(x), a.grb, cdesc))
	default:
		info = Info(C.GrB_Matrix_apply_BinaryOp1st_UDT(c.grb, cmask, caccum, op.grb, unsafe.Pointer(&val), a.grb, cdesc))
	}
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixApplyBinaryOp1stScalar is like [MatrixApplyBinaryOp1st], except that the scalar value is passed as a
// [Scalar] object. It must not be empty.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [EmptyObject], [InvalidObject], [OutOfMemory], [Panic]
func MatrixApplyBinaryOp1stScalar[DC, D, DA any](
	c Matrix[DC],
	mask *Matrix[bool],
	accum *BinaryOp[DC, DC, DC],
	op BinaryOp[DC, D, DA],
	val Scalar[D],
	a Matrix[DA],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADM(mask, accum, desc)
	info := Info(C.GrB_Matrix_apply_BinaryOp1st_Scalar(c.grb, cmask, caccum, op.grb, val.grb, a.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixApplyBinaryOp2nd is like [MatrixApplyBinaryOp1st], except that the stored elements of the matrix are passed
// as the first argument to the binary operator and the specified scalar value is passed as the second argument.
func MatrixApplyBinaryOp2nd[DC, DA, D any](
	c Matrix[DC],
	mask *Matrix[bool],
	accum *BinaryOp[DC, DC, DC],
	op BinaryOp[DC, DA, D],
	a Matrix[DA],
	val D,
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADM(mask, accum, desc)
	var info Info
	switch x := any(val).(type) {
	case bool:
		info = Info(C.GrB_Matrix_apply_BinaryOp2nd_BOOL(c.grb, cmask, caccum, op.grb, a.grb, C.bool(x), cdesc))
	case int:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Matrix_apply_BinaryOp2nd_INT32(c.grb, cmask, caccum, op.grb, a.grb, C.int32_t(x), cdesc))
		} else {
			info = Info(C.GrB_Matrix_apply_BinaryOp2nd_INT64(c.grb, cmask, caccum, op.grb, a.grb, C.int64_t(x), cdesc))
		}
	case int8:
		info = Info(C.GrB_Matrix_apply_BinaryOp2nd_INT8(c.grb, cmask, caccum, op.grb, a.grb, C.int8_t(x), cdesc))
	case int16:
		info = Info(C.GrB_Matrix_apply_BinaryOp2nd_INT16(c.grb, cmask, caccum, op.grb, a.grb, C.int16_t(x), cdesc))
	case int32:
		info = Info(C.GrB_Matrix_apply_BinaryOp2nd_INT32(c.grb, cmask, caccum, op.grb, a.grb, C.int32_t(x), cdesc))
	case int64:
		info = Info(C.GrB_Matrix_apply_BinaryOp2nd_INT64(c.grb, cmask, caccum, op.grb, a.grb, C.int64_t(x), cdesc))
	case uint:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Matrix_apply_BinaryOp2nd_UINT32(c.grb, cmask, caccum, op.grb, a.grb, C.uint32_t(x), cdesc))
		} else {
			info = Info(C.GrB_Matrix_apply_BinaryOp2nd_UINT64(c.grb, cmask, caccum, op.grb, a.grb, C.uint64_t(x), cdesc))
		}
	case uint8:
		info = Info(C.GrB_Matrix_apply_BinaryOp2nd_UINT8(c.grb, cmask, caccum, op.grb, a.grb, C.uint8_t(x), cdesc))
	case uint16:
		info = Info(C.GrB_Matrix_apply_BinaryOp2nd_UINT16(c.grb, cmask, caccum, op.grb, a.grb, C.uint16_t(x), cdesc))
	case uint32:
		info = Info(C.GrB_Matrix_apply_BinaryOp2nd_UINT32(c.grb, cmask, caccum, op.grb, a.grb, C.uint32_t(x), cdesc))
	case uint64:
		info = Info(C.GrB_Matrix_apply_BinaryOp2nd_UINT64(c.grb, cmask, caccum, op.grb, a.grb, C.uint64_t(x), cdesc))
	case float32:
		info = Info(C.GrB_Matrix_apply_BinaryOp2nd_FP32(c.grb, cmask, caccum, op.grb, a.grb, C.float(x), cdesc))
	case float64:
		info = Info(C.GrB_Matrix_apply_BinaryOp2nd_FP64(c.grb, cmask, caccum, op.grb, a.grb, C.double(x), cdesc))
	case complex64:
		info = Info(C.GxB_Matrix_apply_BinaryOp2nd_FC32(c.grb, cmask, caccum, op.grb, a.grb, C.complexfloat(x), cdesc))
	case complex128:
		info = Info(C.GxB_Matrix_apply_BinaryOp2nd_FC64(c.grb, cmask, caccum, op.grb, a.grb, C.complexdouble(x), cdesc))
	default:
		info = Info(C.GrB_Matrix_apply_BinaryOp2nd_UDT(c.grb, cmask, caccum, op.grb, a.grb, unsafe.Pointer(&val), cdesc))
	}
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixApplyBinaryOp2ndScalar is like [MatrixApplyBinaryOp2nd], except that the scalar value is passed as a
// [Scalar] object. It must not be empty.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [EmptyObject], [InvalidObject], [OutOfMemory], [Panic]
func MatrixApplyBinaryOp2ndScalar[DC, DA, D any](
	c Matrix[DC],
	mask *Matrix[bool],
	accum *BinaryOp[DC, DC, DC],
	op BinaryOp[DC, DA, D],
	a Matrix[DA],
	val Scalar[D],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADM(mask, accum, desc)
	info := Info(C.GrB_Matrix_apply_BinaryOp2nd_Scalar(c.grb, cmask, caccum, op.grb, a.grb, val.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// VectorApplyIndexOp computes the transformation of the values of the stored elements of a vector using an index unary
// operator that is a function of the stored value, its location indices, and an user provided scalar
// value. The scalar is passed as a non-opaque variable.
//
// To pass a [Scalar] object instead of a non-opaque variable, use [VectorApplyIndexOpScalar].
//
// Parameters:
//
//   - w (INOUT): An existing GraphBLAS vector. On input, the vector provides values
//     that may be accumulated with the result of the apply operation. On output,
//     this vector holds the results of the operation.
//
//   - mask (IN): An optional "write" [VectorMask].
//
//   - accum (IN): An optional binary operator used for accumulating entries into existing w
//     entries. If assignment rather than accumulation is desired, nil should be specified.
//
//   - op (IN): An index unary operator applied to each element stored in the input vector u.
//     It is a function of the stored element's value, its location index, and a user supplied
//     scalar value val.
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
func VectorApplyIndexOp[Dw, Du, D any](
	w Vector[Dw],
	mask *Vector[bool],
	accum *BinaryOp[Dw, Dw, Dw],
	op IndexUnaryOp[Dw, Du, D],
	u Vector[Du],
	val D,
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	var info Info
	switch x := any(val).(type) {
	case bool:
		info = Info(C.GrB_Vector_apply_IndexOp_BOOL(w.grb, cmask, caccum, op.grb, u.grb, C.bool(x), cdesc))
	case int:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Vector_apply_IndexOp_INT32(w.grb, cmask, caccum, op.grb, u.grb, C.int32_t(x), cdesc))
		} else {
			info = Info(C.GrB_Vector_apply_IndexOp_INT64(w.grb, cmask, caccum, op.grb, u.grb, C.int64_t(x), cdesc))
		}
	case int8:
		info = Info(C.GrB_Vector_apply_IndexOp_INT8(w.grb, cmask, caccum, op.grb, u.grb, C.int8_t(x), cdesc))
	case int16:
		info = Info(C.GrB_Vector_apply_IndexOp_INT16(w.grb, cmask, caccum, op.grb, u.grb, C.int16_t(x), cdesc))
	case int32:
		info = Info(C.GrB_Vector_apply_IndexOp_INT32(w.grb, cmask, caccum, op.grb, u.grb, C.int32_t(x), cdesc))
	case int64:
		info = Info(C.GrB_Vector_apply_IndexOp_INT64(w.grb, cmask, caccum, op.grb, u.grb, C.int64_t(x), cdesc))
	case uint:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Vector_apply_IndexOp_UINT32(w.grb, cmask, caccum, op.grb, u.grb, C.uint32_t(x), cdesc))
		} else {
			info = Info(C.GrB_Vector_apply_IndexOp_UINT64(w.grb, cmask, caccum, op.grb, u.grb, C.uint64_t(x), cdesc))
		}
	case uint8:
		info = Info(C.GrB_Vector_apply_IndexOp_UINT8(w.grb, cmask, caccum, op.grb, u.grb, C.uint8_t(x), cdesc))
	case uint16:
		info = Info(C.GrB_Vector_apply_IndexOp_UINT16(w.grb, cmask, caccum, op.grb, u.grb, C.uint16_t(x), cdesc))
	case uint32:
		info = Info(C.GrB_Vector_apply_IndexOp_UINT32(w.grb, cmask, caccum, op.grb, u.grb, C.uint32_t(x), cdesc))
	case uint64:
		info = Info(C.GrB_Vector_apply_IndexOp_UINT64(w.grb, cmask, caccum, op.grb, u.grb, C.uint64_t(x), cdesc))
	case float32:
		info = Info(C.GrB_Vector_apply_IndexOp_FP32(w.grb, cmask, caccum, op.grb, u.grb, C.float(x), cdesc))
	case float64:
		info = Info(C.GrB_Vector_apply_IndexOp_FP64(w.grb, cmask, caccum, op.grb, u.grb, C.double(x), cdesc))
	case complex64:
		info = Info(C.GxB_Vector_apply_IndexOp_FC32(w.grb, cmask, caccum, op.grb, u.grb, C.complexfloat(x), cdesc))
	case complex128:
		info = Info(C.GxB_Vector_apply_IndexOp_FC64(w.grb, cmask, caccum, op.grb, u.grb, C.complexdouble(x), cdesc))
	default:
		info = Info(C.GrB_Vector_apply_IndexOp_UDT(w.grb, cmask, caccum, op.grb, u.grb, unsafe.Pointer(&val), cdesc))
	}
	if info == success {
		return nil
	}
	return makeError(info)
}

// VectorApplyIndexOpScalar is like [VectorApplyIndexOp], except that the scalar value is passed as a
// [Scalar] object. It must not be empty.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [EmptyObject], [InvalidObject], [OutOfMemory], [Panic]
func VectorApplyIndexOpScalar[Dw, Du, D any](
	w Vector[Dw],
	mask *Vector[bool],
	accum *BinaryOp[Dw, Dw, Dw],
	op IndexUnaryOp[Dw, Du, D],
	u Vector[Du],
	val Scalar[D],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	info := Info(C.GrB_Vector_apply_IndexOp_Scalar(w.grb, cmask, caccum, op.grb, u.grb, val.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixApplyIndexOp computes the transformation of the values of the stored elements of a matrix using an index unary
// operator that is a function of the stored value, its location indices, and an user provided scalar
// value. The scalar is passed as a non-opaque variable.
//
// To pass a [Scalar] object instead of a non-opaque variable, use [MatrixApplyIndexOpScalar].
//
// Parameters:
//
//   - c (INOUT): An existing GraphBLAS matrix. On input, the matrix provides values
//     that may be accumulated with the result of the apply operation. On output,
//     this matrix holds the results of the operation.
//
//   - mask (IN): An optional "write" [MatrixMask].
//
//   - accum (IN): An optional binary operator used for accumulating entries into existing c
//     entries. If assignment rather than accumulation is desired, nil should be specified.
//
//   - op (IN): An index unary operator applied to each element stored in the input matrix c.
//     It is a function of the stored element's value, its row and column indices, and a user supplied
//     scalar value val.
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
func MatrixApplyIndexOp[DC, DA, D any](
	c Matrix[DC],
	mask *Matrix[bool],
	accum *BinaryOp[DC, DC, DC],
	op IndexUnaryOp[DC, DA, D],
	a Matrix[DA],
	val D,
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADM(mask, accum, desc)
	var info Info
	switch x := any(val).(type) {
	case bool:
		info = Info(C.GrB_Matrix_apply_IndexOp_BOOL(c.grb, cmask, caccum, op.grb, a.grb, C.bool(x), cdesc))
	case int:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Matrix_apply_IndexOp_INT32(c.grb, cmask, caccum, op.grb, a.grb, C.int32_t(x), cdesc))
		} else {
			info = Info(C.GrB_Matrix_apply_IndexOp_INT64(c.grb, cmask, caccum, op.grb, a.grb, C.int64_t(x), cdesc))
		}
	case int8:
		info = Info(C.GrB_Matrix_apply_IndexOp_INT8(c.grb, cmask, caccum, op.grb, a.grb, C.int8_t(x), cdesc))
	case int16:
		info = Info(C.GrB_Matrix_apply_IndexOp_INT16(c.grb, cmask, caccum, op.grb, a.grb, C.int16_t(x), cdesc))
	case int32:
		info = Info(C.GrB_Matrix_apply_IndexOp_INT32(c.grb, cmask, caccum, op.grb, a.grb, C.int32_t(x), cdesc))
	case int64:
		info = Info(C.GrB_Matrix_apply_IndexOp_INT64(c.grb, cmask, caccum, op.grb, a.grb, C.int64_t(x), cdesc))
	case uint:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Matrix_apply_IndexOp_UINT32(c.grb, cmask, caccum, op.grb, a.grb, C.uint32_t(x), cdesc))
		} else {
			info = Info(C.GrB_Matrix_apply_IndexOp_UINT64(c.grb, cmask, caccum, op.grb, a.grb, C.uint64_t(x), cdesc))
		}
	case uint8:
		info = Info(C.GrB_Matrix_apply_IndexOp_UINT8(c.grb, cmask, caccum, op.grb, a.grb, C.uint8_t(x), cdesc))
	case uint16:
		info = Info(C.GrB_Matrix_apply_IndexOp_UINT16(c.grb, cmask, caccum, op.grb, a.grb, C.uint16_t(x), cdesc))
	case uint32:
		info = Info(C.GrB_Matrix_apply_IndexOp_UINT32(c.grb, cmask, caccum, op.grb, a.grb, C.uint32_t(x), cdesc))
	case uint64:
		info = Info(C.GrB_Matrix_apply_IndexOp_UINT64(c.grb, cmask, caccum, op.grb, a.grb, C.uint64_t(x), cdesc))
	case float32:
		info = Info(C.GrB_Matrix_apply_IndexOp_FP32(c.grb, cmask, caccum, op.grb, a.grb, C.float(x), cdesc))
	case float64:
		info = Info(C.GrB_Matrix_apply_IndexOp_FP64(c.grb, cmask, caccum, op.grb, a.grb, C.double(x), cdesc))
	case complex64:
		info = Info(C.GxB_Matrix_apply_IndexOp_FC32(c.grb, cmask, caccum, op.grb, a.grb, C.complexfloat(x), cdesc))
	case complex128:
		info = Info(C.GxB_Matrix_apply_IndexOp_FC64(c.grb, cmask, caccum, op.grb, a.grb, C.complexdouble(x), cdesc))
	default:
		info = Info(C.GrB_Matrix_apply_IndexOp_UDT(c.grb, cmask, caccum, op.grb, a.grb, unsafe.Pointer(&val), cdesc))
	}
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixApplyIndexOpScalar is like [MatrixApplyIndexOp], except that the scalar value is passed as a
// [Scalar] object. It must not be empty.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [EmptyObject], [InvalidObject], [OutOfMemory], [Panic]
func MatrixApplyIndexOpScalar[DC, DA, D any](
	c Matrix[DC],
	mask *Matrix[bool],
	accum *BinaryOp[DC, DC, DC],
	op IndexUnaryOp[DC, DA, D],
	a Matrix[DA],
	val Scalar[D],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADM(mask, accum, desc)
	info := Info(C.GrB_Matrix_apply_IndexOp_Scalar(c.grb, cmask, caccum, op.grb, a.grb, val.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}
