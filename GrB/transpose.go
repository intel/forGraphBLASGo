package GrB

// #include "GraphBLAS.h"
import "C"

// Transpose computes a new matrix that is the transpose of the source matrix.
//
// Parameters:
//
//   - c (INOUT): An existing GraphBLAS matrix. On input, the matrix provides values
//     that may be accumulated with the result of the transpose operation. On output,
//     this matrix holds the results of the operation.
//
//   - mask (IN): An optional "write" [MatrixMask].
//
//   - accum (IN): An optional binary operator used for accumulating entries into existing c
//     entries. If assignment rather than accumulation is desired, nil should be specified.
//
//   - a (IN): The GraphBLAS matrix on which transposition will be performed.
//
//   - desc (IN): An optional operation [Descriptor]. If a default descriptor is desired,
//     nil should be specified.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func Transpose[D any](
	c Matrix[D],
	mask *Matrix[bool],
	accum *BinaryOp[D, D, D],
	a Matrix[D],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADM(mask, accum, desc)
	info := Info(C.GrB_transpose(c.grb, cmask, caccum, a.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Transpose is the method variant of [Transpose].
func (matrix Matrix[D]) Transpose(
	mask *Matrix[bool],
	accum *BinaryOp[D, D, D],
	a Matrix[D],
	desc *Descriptor,
) error {
	return Transpose(matrix, mask, accum, a, desc)
}
