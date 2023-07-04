package GrB

// #include "GraphBLAS.h"
import "C"
import "unsafe"

// VectorReduceMonoidValue reduces all stored values into a single scalar.
//
// As an exceptional case, the forGraphBLASGo version of this GraphBLAS function
// does not provide a way to accumulate the result with an already existing value.
//
// To reduce the stored values into a [Scalar] object instead of a non-opaque variable,
// and/or provide an accumulation function, use [VectorReduceMonoidScalar].
//
// Parameters:
//
//   - op (IN): The monoid used in the reduction operation. The operator must be
//     commutative and associative; otherwise, the outcome of the operation is
//     undefined.
//
//   - u (IN): The GraphBLAS vector on which the reduction will be performed.
//
//   - desc (IN): Currently unused.
//
// GraphBLAS API errors that may be returned:
//   - [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func VectorReduceMonoidValue[D any](
	op Monoid[D],
	u Vector[D],
	desc *Descriptor,
) (val D, err error) {
	caccum := C.GrB_BinaryOp(C.GrB_NULL)
	cdesc := processDescriptor(desc)
	var info Info
	switch x := any(&val).(type) {
	case *bool:
		var cx C.bool
		info = Info(C.GrB_Vector_reduce_BOOL(&cx, caccum, op.grb, u.grb, cdesc))
		if info == success {
			*x = bool(cx)
			return
		}
	case *int:
		if unsafe.Sizeof(0) == 4 {
			var cx C.int32_t
			info = Info(C.GrB_Vector_reduce_INT32(&cx, caccum, op.grb, u.grb, cdesc))
			if info == success {
				*x = int(cx)
				return
			}
		} else {
			var cx C.int64_t
			info = Info(C.GrB_Vector_reduce_INT64(&cx, caccum, op.grb, u.grb, cdesc))
			if info == success {
				*x = int(cx)
				return
			}
		}
	case *int8:
		var cx C.int8_t
		info = Info(C.GrB_Vector_reduce_INT8(&cx, caccum, op.grb, u.grb, cdesc))
		if info == success {
			*x = int8(cx)
			return
		}
	case *int16:
		var cx C.int16_t
		info = Info(C.GrB_Vector_reduce_INT16(&cx, caccum, op.grb, u.grb, cdesc))
		if info == success {
			*x = int16(cx)
			return
		}
	case *int32:
		var cx C.int32_t
		info = Info(C.GrB_Vector_reduce_INT32(&cx, caccum, op.grb, u.grb, cdesc))
		if info == success {
			*x = int32(cx)
			return
		}
	case *int64:
		var cx C.int64_t
		info = Info(C.GrB_Vector_reduce_INT64(&cx, caccum, op.grb, u.grb, cdesc))
		if info == success {
			*x = int64(cx)
			return
		}
	case *uint:
		if unsafe.Sizeof(0) == 4 {
			var cx C.uint32_t
			info = Info(C.GrB_Vector_reduce_UINT32(&cx, caccum, op.grb, u.grb, cdesc))
			if info == success {
				*x = uint(cx)
				return
			}
		} else {
			var cx C.uint64_t
			info = Info(C.GrB_Vector_reduce_UINT64(&cx, caccum, op.grb, u.grb, cdesc))
			if info == success {
				*x = uint(cx)
				return
			}
		}
	case *uint8:
		var cx C.uint8_t
		info = Info(C.GrB_Vector_reduce_UINT8(&cx, caccum, op.grb, u.grb, cdesc))
		if info == success {
			*x = uint8(cx)
			return
		}
	case *uint16:
		var cx C.uint16_t
		info = Info(C.GrB_Vector_reduce_UINT16(&cx, caccum, op.grb, u.grb, cdesc))
		if info == success {
			*x = uint16(cx)
			return
		}
	case *uint32:
		var cx C.uint32_t
		info = Info(C.GrB_Vector_reduce_UINT32(&cx, caccum, op.grb, u.grb, cdesc))
		if info == success {
			*x = uint32(cx)
			return
		}
	case *uint64:
		var cx C.uint64_t
		info = Info(C.GrB_Vector_reduce_UINT64(&cx, caccum, op.grb, u.grb, cdesc))
		if info == success {
			*x = uint64(cx)
			return
		}
	case *float32:
		var cx C.float
		info = Info(C.GrB_Vector_reduce_FP32(&cx, caccum, op.grb, u.grb, cdesc))
		if info == success {
			*x = float32(cx)
			return
		}
	case *float64:
		var cx C.double
		info = Info(C.GrB_Vector_reduce_FP64(&cx, caccum, op.grb, u.grb, cdesc))
		if info == success {
			*x = float64(cx)
			return
		}
	case *complex64:
		var cx C.complexfloat
		info = Info(C.GxB_Vector_reduce_FC32(&cx, caccum, op.grb, u.grb, cdesc))
		if info == success {
			*x = complex64(cx)
			return
		}
	case *complex128:
		var cx C.complexdouble
		info = Info(C.GxB_Vector_reduce_FC64(&cx, caccum, op.grb, u.grb, cdesc))
		if info == success {
			*x = complex128(cx)
			return
		}
	default:
		info = Info(C.GrB_Vector_reduce_UDT(unsafe.Pointer(&val), caccum, op.grb, u.grb, cdesc))
		if info == success {
			return
		}
	}
	err = makeError(info)
	return
}

// Reduce is the method variant of [VectorReduceMonoidValue].
func (vector Vector[D]) Reduce(op Monoid[D], desc *Descriptor) (val D, err error) {
	return VectorReduceMonoidValue(op, vector, desc)
}

// VectorReduceMonoidScalar reduces all stored values into a single scalar.
//
// To reduce the stored values into a non-opaque variable, use [VectorReduceMonoidValue].
//
// To use a [BinaryOp] instead of a [Monoid], use [VectorReduceBinaryOpScalar].
//
// Parameters:
//
//   - s (INOUT): Scalar to store the final reduced value into. On input, the scalar provides
//     a value that may be accumulated (optionally) with the result of the reduction
//     operation. On output, this scalar holds the result of the operation.
//
//   - accum (IN): An optional binary operator used for accumulating entries into an existing
//     scalar value. If assignment rather than accumulation is desired, nil should be specified.
//
//   - op (IN): The monoid used in the reduction operation. The operator must be
//     commutative and associative; otherwise, the outcome of the operation is
//     undefined.
//
//   - u (IN): The GraphBLAS vector on which the reduction will be performed.
//
//   - desc (IN): Currently unused.
//
// GraphBLAS API errors that may be returned:
//   - [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func VectorReduceMonoidScalar[D any](
	s Scalar[D],
	accum *BinaryOp[D, D, D],
	op Monoid[D],
	u Vector[D],
	desc *Descriptor,
) error {
	caccum, cdesc := processAD(accum, desc)
	info := Info(C.GrB_Vector_reduce_Monoid_Scalar(s.grb, caccum, op.grb, u.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// VectorReduceMonoid is the method variant of [VectorReduceMonoidScalar].
func (scalar Scalar[D]) VectorReduceMonoid(
	accum *BinaryOp[D, D, D],
	op Monoid[D],
	u Vector[D],
	desc *Descriptor,
) error {
	return VectorReduceMonoidScalar(scalar, accum, op, u, desc)
}

// VectorReduceBinaryOpScalar is like [VectorReduceMonoidScalar], except that a [BinaryOp] is used instead of a [Monoid]
// to specify the reduction operator.
//
// SuiteSparse:GraphBLAS supports this function only, if op is a built-in binary operator,
// and corresponds to a built-in monoid. For other binary operators, including user-defined
// ones, [NotImplemented] is returned.
func VectorReduceBinaryOpScalar[D any](
	s Scalar[D],
	accum *BinaryOp[D, D, D],
	op BinaryOp[D, D, D],
	u Vector[D],
	desc *Descriptor,
) error {
	caccum, cdesc := processAD(accum, desc)
	info := Info(C.GrB_Vector_reduce_BinaryOp_Scalar(s.grb, caccum, op.grb, u.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// VectorReduceBinaryOp is the method variant of [VectorReduceBinaryOpScalar].
func (scalar Scalar[D]) VectorReduceBinaryOp(
	accum *BinaryOp[D, D, D],
	op BinaryOp[D, D, D],
	u Vector[D],
	desc *Descriptor,
) error {
	return VectorReduceBinaryOpScalar(scalar, accum, op, u, desc)
}

// MatrixReduceMonoidValue reduces all stored values into a single scalar.
//
// As an exceptional case, the forGraphBLASGo version of this GraphBLAS function
// does not provide a way to accumulate the result with an already existing value.
//
// To reduce the stored values into a [Scalar] object instead of a non-opaque variable,
// and/or provide an accumulation function, use [MatrixReduceMonoidScalar].
//
// Parameters:
//
//   - op (IN): The monoid used in the reduction operation. The operator must be
//     commutative and associative; otherwise, the outcome of the operation is
//     undefined.
//
//   - a (IN): The GraphBLAS matrix on which the reduction will be performed.
//
//   - desc (IN): Currently unused.
//
// GraphBLAS API errors that may be returned:
//   - [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func MatrixReduceMonoidValue[D any](
	op Monoid[D],
	a Matrix[D],
	desc *Descriptor,
) (val D, err error) {
	caccum := C.GrB_BinaryOp(C.GrB_NULL)
	cdesc := processDescriptor(desc)
	var info Info
	switch x := any(&val).(type) {
	case *bool:
		var cx C.bool
		info = Info(C.GrB_Matrix_reduce_BOOL(&cx, caccum, op.grb, a.grb, cdesc))
		if info == success {
			*x = bool(cx)
			return
		}
	case *int:
		if unsafe.Sizeof(0) == 4 {
			var cx C.int32_t
			info = Info(C.GrB_Matrix_reduce_INT32(&cx, caccum, op.grb, a.grb, cdesc))
			if info == success {
				*x = int(cx)
				return
			}
		} else {
			var cx C.int64_t
			info = Info(C.GrB_Matrix_reduce_INT64(&cx, caccum, op.grb, a.grb, cdesc))
			if info == success {
				*x = int(cx)
				return
			}
		}
	case *int8:
		var cx C.int8_t
		info = Info(C.GrB_Matrix_reduce_INT8(&cx, caccum, op.grb, a.grb, cdesc))
		if info == success {
			*x = int8(cx)
			return
		}
	case *int16:
		var cx C.int16_t
		info = Info(C.GrB_Matrix_reduce_INT16(&cx, caccum, op.grb, a.grb, cdesc))
		if info == success {
			*x = int16(cx)
			return
		}
	case *int32:
		var cx C.int32_t
		info = Info(C.GrB_Matrix_reduce_INT32(&cx, caccum, op.grb, a.grb, cdesc))
		if info == success {
			*x = int32(cx)
			return
		}
	case *int64:
		var cx C.int64_t
		info = Info(C.GrB_Matrix_reduce_INT64(&cx, caccum, op.grb, a.grb, cdesc))
		if info == success {
			*x = int64(cx)
			return
		}
	case *uint:
		if unsafe.Sizeof(0) == 4 {
			var cx C.uint32_t
			info = Info(C.GrB_Matrix_reduce_UINT32(&cx, caccum, op.grb, a.grb, cdesc))
			if info == success {
				*x = uint(cx)
				return
			}
		} else {
			var cx C.uint64_t
			info = Info(C.GrB_Matrix_reduce_UINT64(&cx, caccum, op.grb, a.grb, cdesc))
			if info == success {
				*x = uint(cx)
				return
			}
		}
	case *uint8:
		var cx C.uint8_t
		info = Info(C.GrB_Matrix_reduce_UINT8(&cx, caccum, op.grb, a.grb, cdesc))
		if info == success {
			*x = uint8(cx)
			return
		}
	case *uint16:
		var cx C.uint16_t
		info = Info(C.GrB_Matrix_reduce_UINT16(&cx, caccum, op.grb, a.grb, cdesc))
		if info == success {
			*x = uint16(cx)
			return
		}
	case *uint32:
		var cx C.uint32_t
		info = Info(C.GrB_Matrix_reduce_UINT32(&cx, caccum, op.grb, a.grb, cdesc))
		if info == success {
			*x = uint32(cx)
			return
		}
	case *uint64:
		var cx C.uint64_t
		info = Info(C.GrB_Matrix_reduce_UINT64(&cx, caccum, op.grb, a.grb, cdesc))
		if info == success {
			*x = uint64(cx)
			return
		}
	case *float32:
		var cx C.float
		info = Info(C.GrB_Matrix_reduce_FP32(&cx, caccum, op.grb, a.grb, cdesc))
		if info == success {
			*x = float32(cx)
			return
		}
	case *float64:
		var cx C.double
		info = Info(C.GrB_Matrix_reduce_FP64(&cx, caccum, op.grb, a.grb, cdesc))
		if info == success {
			*x = float64(cx)
			return
		}
	case *complex64:
		var cx C.complexfloat
		info = Info(C.GxB_Matrix_reduce_FC32(&cx, caccum, op.grb, a.grb, cdesc))
		if info == success {
			*x = complex64(cx)
			return
		}
	case *complex128:
		var cx C.complexdouble
		info = Info(C.GxB_Matrix_reduce_FC64(&cx, caccum, op.grb, a.grb, cdesc))
		if info == success {
			*x = complex128(cx)
			return
		}
	default:
		info = Info(C.GrB_Matrix_reduce_UDT(unsafe.Pointer(&val), caccum, op.grb, a.grb, cdesc))
		if info == success {
			return
		}
	}
	err = makeError(info)
	return
}

// Reduce is the method variant of [MatrixReduceMonoidValue].
func (matrix Matrix[D]) Reduce(op Monoid[D], desc *Descriptor) (val D, err error) {
	return MatrixReduceMonoidValue(op, matrix, desc)
}

// MatrixReduceMonoidScalar reduces all stored values into a single scalar.
//
// To reduce the stored values into a non-opaque variable, use [MatrixReduceMonoidValue].
//
// To use a [BinaryOp] instead of a [Monoid], use [MatrixReduceBinaryOpScalar].
//
// Parameters:
//
//   - s (INOUT): Scalar to store the final reduced value into. On input, the scalar provides
//     a value that may be accumulated (optionally) with the result of the reduction
//     operation. On output, this scalar holds the result of the operation.
//
//   - accum (IN): An optional binary operator used for accumulating entries into an existing
//     scalar value. If assignment rather than accumulation is desired, nil should be specified.
//
//   - op (IN): The monoid used in the reduction operation. The operator must be
//     commutative and associative; otherwise, the outcome of the operation is
//     undefined.
//
//   - a (IN): The GraphBLAS matrix on which the reduction will be performed.
//
//   - desc (IN): Currently unused.
//
// GraphBLAS API errors that may be returned:
//   - [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func MatrixReduceMonoidScalar[D any](
	s Scalar[D],
	accum *BinaryOp[D, D, D],
	op Monoid[D],
	a Matrix[D],
	desc *Descriptor,
) error {
	caccum, cdesc := processAD(accum, desc)
	info := Info(C.GrB_Matrix_reduce_Monoid_Scalar(s.grb, caccum, op.grb, a.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixReduceMonoid is the method variant of [MatrixReduceMonoidScalar].
func (scalar Scalar[D]) MatrixReduceMonoid(
	accum *BinaryOp[D, D, D],
	op Monoid[D],
	a Matrix[D],
	desc *Descriptor,
) error {
	return MatrixReduceMonoidScalar(scalar, accum, op, a, desc)
}

// MatrixReduceBinaryOpScalar is like [MatrixReduceMonoidScalar], except that a [BinaryOp] is used instead of a [Monoid]
// to specify the reduction operator.
//
// SuiteSparse:GraphBLAS supports this function only, if op is a built-in binary operator,
// and corresponds to a built-in monoid. For other binary operators, including user-defined
// ones, [NotImplemented] is returned.
func MatrixReduceBinaryOpScalar[D any](
	s Scalar[D],
	accum *BinaryOp[D, D, D],
	op BinaryOp[D, D, D],
	a Matrix[D],
	desc *Descriptor,
) error {
	caccum, cdesc := processAD(accum, desc)
	info := Info(C.GrB_Matrix_reduce_BinaryOp_Scalar(s.grb, caccum, op.grb, a.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixReduceBinaryOp is the method variant of [MatrixReduceBinaryOpScalar].
func (scalar Scalar[D]) MatrixReduceBinaryOp(
	accum *BinaryOp[D, D, D],
	op BinaryOp[D, D, D],
	a Matrix[D],
	desc *Descriptor,
) error {
	return MatrixReduceBinaryOpScalar(scalar, accum, op, a, desc)
}

// MatrixReduceMonoidVector performs a reduction across rows of a matrix to produce
// a vector. If reduction down columns is desired, the input matrix should be transposed
// using the descriptor.
//
// To use a [BinaryOp] instead of a [Monoid], use [MatrixReduceBinaryOpVector].
//
// Parameters:
//
//   - w (INOUT): An existing GraphBLAS vector. On input, the scalar provides
//     values that may be accumulated with the result of the reduction
//     operation. On output, this vector holds the result of the operation.
//
//   - mask (IN): An optional "write" [VectorMask].
//
//   - accum (IN): An optional binary operator used for accumulating entries into existing w
//     entries. If assignment rather than accumulation is desired, nil should be specified.
//
//   - op (IN): The monoid used in the reduction operation. The operator must be
//     commutative and associative; otherwise, the outcome of the operation is
//     undefined.
//
//   - a (IN): The GraphBLAS matrix on which the reduction will be performed.
//
//   - desc (IN): An optional operation [Descriptor]. If a default descriptor is desired,
//     nil should be specified.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func MatrixReduceMonoidVector[D any](
	w Vector[D],
	mask *Vector[bool],
	accum *BinaryOp[D, D, D],
	op Monoid[D],
	a Matrix[D],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	info := Info(C.GrB_Matrix_reduce_Monoid(w.grb, cmask, caccum, op.grb, a.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixReduceMonoid is the method variant of [MatrixReduceMonoidVector].
func (vector Vector[D]) MatrixReduceMonoid(
	mask *Vector[bool],
	accum *BinaryOp[D, D, D],
	op Monoid[D],
	a Matrix[D],
	desc *Descriptor,
) error {
	return MatrixReduceMonoidVector(vector, mask, accum, op, a, desc)
}

// MatrixReduceBinaryOpVector is like [MatrixReduceMonoidVector], except that a [BinaryOp] is used instead of a [Monoid]
// to specify the reduction operator.
//
// SuiteSparse:GraphBLAS supports this function only, if op is a built-in binary operator,
// and corresponds to a built-in monoid. For other binary operators, including user-defined
// ones, [NotImplemented] is returned.
func MatrixReduceBinaryOpVector[D any](
	w Vector[D],
	mask *Vector[bool],
	accum *BinaryOp[D, D, D],
	op BinaryOp[D, D, D],
	a Matrix[D],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	info := Info(C.GrB_Matrix_reduce_BinaryOp(w.grb, cmask, caccum, op.grb, a.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixReduceBinaryOp is the method variant of [MatrixReduceBinaryOpVector].
func (vector Vector[D]) MatrixReduceBinaryOp(
	mask *Vector[bool],
	accum *BinaryOp[D, D, D],
	op BinaryOp[D, D, D],
	a Matrix[D],
	desc *Descriptor,
) error {
	return MatrixReduceBinaryOpVector(vector, mask, accum, op, a, desc)
}
