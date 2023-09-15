package GrB

// #include "GraphBLAS.h"
import "C"

// VectorEWiseMultSemiring performs element-wise (general) multiplication on the intersection of the elements of two
// vectors, producing a third vector as a result.
//
// The multiplication operator is the [Semiring.Multiply] operation of the provided [Semiring]. To use a [Monoid]
// instead of a [Semiring], use [VectorEWiseMultMonoid]. To use a [BinaryOp] instead of a Semiring, use
// [VectorEWiseMultBinaryOp].
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
//   - op (IN): The semiring used in the element-wise "product" operation. The additive monoid
//     of the semiring is ignored.
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
func VectorEWiseMultSemiring[Dw, Du, Dv any](
	w Vector[Dw],
	mask *Vector[bool],
	accum *BinaryOp[Dw, Dw, Dw],
	op Semiring[Dw, Du, Dv],
	u Vector[Du],
	v Vector[Dv],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	info := Info(C.GrB_Vector_eWiseMult_Semiring(w.grb, cmask, caccum, op.grb, u.grb, v.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// VectorEWiseMultMonoid is like [VectorEWiseMultSemiring], except that a [Monoid] is used instead of a [Semiring]
// to specify the binary operator op. The identity element of the monoid is ignored.
func VectorEWiseMultMonoid[D any](
	w Vector[D],
	mask *Vector[bool],
	accum *BinaryOp[D, D, D],
	op Monoid[D],
	u Vector[D],
	v Vector[D],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	info := Info(C.GrB_Vector_eWiseMult_Monoid(w.grb, cmask, caccum, op.grb, u.grb, v.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// VectorEWiseMultBinaryOp is like [VectorEWiseMultSemiring], except that a [BinaryOp] is used instead of a [Semiring]
// to specify the binary operator op.
func VectorEWiseMultBinaryOp[Dw, Du, Dv any](
	w Vector[Dw],
	mask *Vector[bool],
	accum *BinaryOp[Dw, Dw, Dw],
	op BinaryOp[Dw, Du, Dv],
	u Vector[Du],
	v Vector[Dv],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	info := Info(C.GrB_Vector_eWiseMult_BinaryOp(w.grb, cmask, caccum, op.grb, u.grb, v.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixEWiseMultSemiring performs element-wise (general) multiplication on the intersection of the elements of two
// matrices, producing a third matrix as a result.
//
// The multiplication operator is the [Semiring.Multiply] operation of the provided [Semiring]. To use a [Monoid]
// instead of a [Semiring], use [MatrixEWiseMultMonoid]. To use a [BinaryOp] instead of a Semiring, use
// [MatrixEWiseMultBinaryOp].
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
//   - op (IN): The semiring used in the element-wise "product" operation. The additive monoid
//     of the semiring is ignored.
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
func MatrixEWiseMultSemiring[DC, DA, DB any](
	c Matrix[DC],
	mask *Matrix[bool],
	accum *BinaryOp[DC, DC, DC],
	op Semiring[DC, DA, DB],
	a Matrix[DA],
	b Matrix[DB],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADM(mask, accum, desc)
	info := Info(C.GrB_Matrix_eWiseMult_Semiring(c.grb, cmask, caccum, op.grb, a.grb, b.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixEWiseMultMonoid is like [MatrixEWiseMultSemiring], except that a [Monoid] is used instead of a [Semiring]
// to specify the binary operator op. The identity element of the monoid is ignored.
func MatrixEWiseMultMonoid[D any](
	c Matrix[D],
	mask *Matrix[bool],
	accum *BinaryOp[D, D, D],
	op Monoid[D],
	a Matrix[D],
	b Matrix[D],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADM(mask, accum, desc)
	info := Info(C.GrB_Matrix_eWiseMult_Monoid(c.grb, cmask, caccum, op.grb, a.grb, b.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixEWiseMultBinaryOp is like [MatrixEWiseMultSemiring], except that a [BinaryOp] is used instead of a [Semiring]
// to specify the binary operator op.
func MatrixEWiseMultBinaryOp[DC, DA, DB any](
	c Matrix[DC],
	mask *Matrix[bool],
	accum *BinaryOp[DC, DC, DC],
	op BinaryOp[DC, DA, DB],
	a Matrix[DA],
	b Matrix[DB],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADM(mask, accum, desc)
	info := Info(C.GrB_Matrix_eWiseMult_BinaryOp(c.grb, cmask, caccum, op.grb, a.grb, b.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}
