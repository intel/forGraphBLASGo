package GrB

// #include "GraphBLAS.h"
import "C"

// VectorEWiseAddSemiring performs element-wise (general) addition on the elements of two vectors, producing a third
// vector as a result.
//
// The addition operator is the [Semiring.Add] operation of the provided [Semiring]. To use a [Monoid] instead of
// a [Semiring], use [VectorEWiseAddMonoid]. To use a [BinaryOp] instead of a Semiring, use [VectorEWiseAddBinaryOp].
//
// Parameters:
//
//   - w (INOUT): An existing GraphBLAS vector. On input, the vector provides values
//     that may be accumulated with the result of the element-wise operation. On output,
//     this vector holds the results of the operation.
//
//   - mask (IN): An optional "write" [VectorMask].
//
//   - accum (IN): An optional binary operator used for accumulating entries into existing w
//     entries. If assignment rather than accumulation is desired, nil should be specified.
//
//   - op (IN): The semiring used in the element-wise "sum" operation. The multiplicative binary
//     operator and additive identity of the semiring are ignored.
//
//   - u (IN): The GraphBLAS vector holding the values for the left-hand vector in the operation.
//
//   - v (IN): The GraphBLAS vector holding the values for the right-hand vector in the operation.
//
//   - desc (IN): An optional operation [Descriptor]. If a default descriptor is desired,
//     nil should be specified.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func VectorEWiseAddSemiring[Dw, Du, Dv any](
	w Vector[Dw],
	mask *Vector[bool],
	accum *BinaryOp[Dw, Dw, Dw],
	op Semiring[Dw, Du, Dv],
	u Vector[Du],
	v Vector[Dv],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	info := Info(C.GrB_Vector_eWiseAdd_Semiring(w.grb, cmask, caccum, op.grb, u.grb, v.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// VectorEWiseAddMonoid is like [VectorEWiseAddSemiring], except that a [Monoid] is used instead of a [Semiring]
// to specify the binary operator op. The identity element of the monoid is ignored.
func VectorEWiseAddMonoid[D any](
	w Vector[D],
	mask *Vector[bool],
	accum *BinaryOp[D, D, D],
	op Monoid[D],
	u Vector[D],
	v Vector[D],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	info := Info(C.GrB_Vector_eWiseAdd_Monoid(w.grb, cmask, caccum, op.grb, u.grb, v.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// VectorEWiseAddBinaryOp is like [VectorEWiseAddSemiring], except that a [BinaryOp] is used instead of a [Semiring]
// to specify the binary operator op.
func VectorEWiseAddBinaryOp[Dw, Du, Dv any](
	w Vector[Dw],
	mask *Vector[bool],
	accum *BinaryOp[Dw, Dw, Dw],
	op BinaryOp[Dw, Du, Dv],
	u Vector[Du],
	v Vector[Dv],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	info := Info(C.GrB_Vector_eWiseAdd_BinaryOp(w.grb, cmask, caccum, op.grb, u.grb, v.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixEWiseAddSemiring performs element-wise (general) addition on the elements of two matrices, producing a third
// matrix as a result.
//
// The addition operator is the [Semiring.Add] operation of the provided [Semiring]. To use a [Monoid] instead of
// a [Semiring], use [MatrixEWiseAddMonoid]. To use a [BinaryOp] instead of a Semiring, use [MatrixEWiseAddBinaryOp].
//
// Parameters:
//
//   - c (INOUT): An existing GraphBLAS matrix. On input, the matrix provides values
//     that may be accumulated with the result of the element-wise operation. On output,
//     this matrix holds the results of the operation.
//
//   - mask (IN): An optional "write" [MatrixMask].
//
//   - accum (IN): An optional binary operator used for accumulating entries into existing c
//     entries. If assignment rather than accumulation is desired, nil should be specified.
//
//   - op (IN): The semiring used in the element-wise "sum" operation. The multiplicative binary
//     operator and additive identity of the semiring are ignored.
//
//   - a (IN): The GraphBLAS matrix holding the values for the left-hand matrix in the operation.
//
//   - b (IN): The GraphBLAS matrix holding the values for the right-hand matrix in the operation.
//
//   - desc (IN): An optional operation [Descriptor]. If a default descriptor is desired,
//     nil should be specified.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func MatrixEWiseAddSemiring[DC, DA, DB any](
	c Matrix[DC],
	mask *Matrix[bool],
	accum *BinaryOp[DC, DC, DC],
	op Semiring[DC, DA, DB],
	a Matrix[DA],
	b Matrix[DB],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADM(mask, accum, desc)
	info := Info(C.GrB_Matrix_eWiseAdd_Semiring(c.grb, cmask, caccum, op.grb, a.grb, b.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixEWiseAddMonoid is like [MatrixEWiseAddSemiring], except that a [Monoid] is used instead of a [Semiring]
// to specify the binary operator op. The identity element of the monoid is ignored.
func MatrixEWiseAddMonoid[D any](
	c Matrix[D],
	mask *Matrix[bool],
	accum *BinaryOp[D, D, D],
	op Monoid[D],
	a Matrix[D],
	b Matrix[D],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADM(mask, accum, desc)
	info := Info(C.GrB_Matrix_eWiseAdd_Monoid(c.grb, cmask, caccum, op.grb, a.grb, b.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixEWiseAddBinaryOp is like [MatrixEWiseAddSemiring], except that a [BinaryOp] is used instead of a [Semiring]
// to specify the binary operator op.
func MatrixEWiseAddBinaryOp[DC, DA, DB any](
	c Matrix[DC],
	mask *Matrix[bool],
	accum *BinaryOp[DC, DC, DC],
	op BinaryOp[DC, DA, DB],
	a Matrix[DA],
	b Matrix[DB],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADM(mask, accum, desc)
	info := Info(C.GrB_Matrix_eWiseAdd_BinaryOp(c.grb, cmask, caccum, op.grb, a.grb, b.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}
