package GrB

// #include "GraphBLAS.h"
import "C"

// KroneckerSemiring computes the Kronecker product of two matrices. The result is a matrix.
//
// The multiplication operator is the [Semiring.Multiply] operation of the provided [Semiring]. To use a [Monoid]
// instead of a [Semiring], use [KroneckerMonoid]. To use a [BinaryOp] instead of a Semiring, use
// [KroneckerBinaryOp].
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
//   - op (IN): The semiring used in the "product" operation. The additive monoid of the semiring
//     is ignored.
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
func KroneckerSemiring[DC, DA, DB any](
	c Matrix[DC],
	mask *Matrix[bool],
	accum *BinaryOp[DC, DC, DC],
	op Semiring[DC, DA, DB],
	a Matrix[DA],
	b Matrix[DB],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADM(mask, accum, desc)
	info := Info(C.GrB_Matrix_kronecker_Semiring(c.grb, cmask, caccum, op.grb, a.grb, b.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// KroneckerMonoid is like [KroneckerSemiring], except that a [Monoid] is used instead of a [Semiring]
// to specify the binary operator op. The identity element of the monoid is ignored.
func KroneckerMonoid[D any](
	c Matrix[D],
	mask *Matrix[bool],
	accum *BinaryOp[D, D, D],
	op Monoid[D],
	a Matrix[D],
	b Matrix[D],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADM(mask, accum, desc)
	info := Info(C.GrB_Matrix_kronecker_Monoid(c.grb, cmask, caccum, op.grb, a.grb, b.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// KroneckerBinaryOp is like [KroneckerSemiring], except that a [BinaryOp] is used instead of a [Semiring]
// to specify the binary operator op.
func KroneckerBinaryOp[DC, DA, DB any](
	c Matrix[DC],
	mask *Matrix[bool],
	accum *BinaryOp[DC, DC, DC],
	op BinaryOp[DC, DA, DB],
	a Matrix[DA],
	b Matrix[DB],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADM(mask, accum, desc)
	info := Info(C.GrB_Matrix_kronecker_BinaryOp(c.grb, cmask, caccum, op.grb, a.grb, b.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}
