package GrB

// #include "GraphBLAS.h"
import "C"

// MxM multiplies a matrix with another matrix on a semiring.
//
// Parameters:
//
//   - c (INOUT): An existing GraphBLAS matrix. On input, the matrix provides values
//     that may be accumulated with the result of the matrix product. On output,
//     this matrix holds the results of the operation.
//
//   - mask (IN): An optional "write" [MatrixMask].
//
//   - accum (IN): An optional binary operator used for accumulating entries into existing c
//     entries. If assignment rather than accumulation is desired, nil should be specified.
//
//   - op (IN): The semiring used in the matrix-matrix multiply.
//
//   - a (IN): The GraphBLAS matrix holding the values for the left-hand matrix in the
//     multiplication.
//
//   - b (IN): The GraphBLAS matrix holding the values for the right-hand matrix in the
//     multiplication.
//
//   - desc (IN): An optional operation [Descriptor]. If a default descriptor is desired,
//     nil should be specified.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func MxM[DC, DA, DB any](
	c Matrix[DC],
	mask *Matrix[bool],
	accum *BinaryOp[DC, DC, DC],
	op Semiring[DC, DA, DB],
	a Matrix[DA],
	b Matrix[DB],
	desc *Descriptor) error {
	cmask, caccum, cdesc := processMADM(mask, accum, desc)
	info := Info(C.GrB_mxm(c.grb, cmask, caccum, op.grb, a.grb, b.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MxM is the method variant of [MxM].
func (matrix Matrix[D]) MxM(
	mask *Matrix[bool],
	accum *BinaryOp[D, D, D],
	op Semiring[D, D, D],
	a Matrix[D],
	b Matrix[D],
	desc *Descriptor) error {
	return MxM(matrix, mask, accum, op, a, b, desc)
}

// VxM multiplies a (row) vector with a matrix on a semiring. The result is a vector.
//
// Parameters:
//
//   - w (INOUT): An existing GraphBLAS vector. On input, the vector provides values
//     that may be accumulated with the result of the vector-matrix product. On output,
//     this vector holds the results of the operation.
//
//   - mask (IN): An optional "write" [VectorMask].
//
//   - accum (IN): An optional binary operator used for accumulating entries into existing w
//     entries. If assignment rather than accumulation is desired, nil should be specified.
//
//   - op (IN): The semiring used in the vector-matrix multiply.
//
//   - u (IN): The GraphBLAS vector holding the values for the left-hand vector in the
//     multiplication.
//
//   - a (IN): The GraphBLAS matrix holding the values for the right-hand matrix in the
//     multiplication.
//
//   - desc (IN): An optional operation [Descriptor]. If a default descriptor is desired,
//     nil should be specified. Set the [Tran] descriptor for [Inp1] to use the transpose
//     of the input matrix a.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func VxM[Dw, Du, DA any](
	w Vector[Dw],
	mask *Vector[bool],
	accum *BinaryOp[Dw, Dw, Dw],
	op Semiring[Dw, Du, DA],
	u Vector[Du],
	a Matrix[DA],
	desc *Descriptor) error {
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	info := Info(C.GrB_vxm(w.grb, cmask, caccum, op.grb, u.grb, a.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// VxM is the method variant of [VxM].
func (vector Vector[D]) VxM(
	mask *Vector[bool],
	accum *BinaryOp[D, D, D],
	op Semiring[D, D, D],
	u Vector[D],
	a Matrix[D],
	desc *Descriptor) error {
	return VxM(vector, mask, accum, op, u, a, desc)
}

// MxV multiplies a matrix by a vector on a semiring. The result is a vector.
//
// Parameters:
//
//   - w (INOUT): An existing GraphBLAS vector. On input, the vector provides values
//     that may be accumulated with the result of the matrix-vector product. On output,
//     this vector holds the results of the operation.
//
//   - mask (IN): An optional "write" [VectorMask].
//
//   - accum (IN): An optional binary operator used for accumulating entries into existing w
//     entries. If assignment rather than accumulation is desired, nil should be specified.
//
//   - op (IN): The semiring used in the vector-matrix multiply.
//
//   - a (IN): The GraphBLAS matrix holding the values for the left-hand matrix in the
//     multiplication.
//
//   - u (IN): The GraphBLAS vector holding the values for the right-hand vector in the
//     multiplication.
//
//   - desc (IN): An optional operation [Descriptor]. If a default descriptor is desired,
//     nil should be specified. Set the [Tran] descriptor for [Inp0] to use the transpose
//     of the input matrix a.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func MxV[Dw, DA, Du any](
	w Vector[Dw],
	mask *Vector[bool],
	accum *BinaryOp[Dw, Dw, Dw],
	op Semiring[Dw, DA, Du],
	a Matrix[DA],
	u Vector[Du],
	desc *Descriptor) error {
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	info := Info(C.GrB_mxv(w.grb, cmask, caccum, op.grb, a.grb, u.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MxV is the method variant of [MxV].
func (vector Vector[D]) MxV(
	mask *Vector[bool],
	accum *BinaryOp[D, D, D],
	op Semiring[D, D, D],
	a Matrix[D],
	u Vector[D],
	desc *Descriptor) error {
	return MxV(vector, mask, accum, op, a, u, desc)
}
