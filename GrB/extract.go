package GrB

// #include "GraphBLAS.h"
import "C"

// VectorExtract extracts a sub-vector from a larger vector as specified by a set of indices. The result is a vector
// whose size is equal to the number of indices.
//
// Parameters:
//
//   - w (INOUT): An existing GraphBLAS vector. On input, the vector provides values
//     that may be accumulated with the result of the extract operation. On output,
//     this vector holds the results of the operation.
//
//   - mask (IN): An optional "write" [VectorMask].
//
//   - accum (IN): An optional binary operator used for accumulating entries into existing w
//     entries. If assignment rather than accumulation is desired, nil should be specified.
//
//   - u (IN): The GraphBLAS vector from which the subset is extracted.
//
//   - indices (IN): The ordered set (slice) of indices corresponding to the locations of elements
//     from u that are extracted. If all elements of u are to be extracted in order from 0 to nindices − 1,
//     then [All](nindices) should be specified. Regardless of execution mode and return value, this slice
//     may be manipulated by the caller after this operation returns without affecting any deferred
//     computations for this operation.
//
//   - desc (IN): An optional operation [Descriptor]. If a default descriptor is desired,
//     nil should be specified.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [IndexOutOfBounds], [InvalidObject], [OutOfMemory], [Panic]
func VectorExtract[D any](
	w Vector[D],
	mask *Vector[bool],
	accum *BinaryOp[D, D, D],
	u Vector[D],
	indices []int,
	desc *Descriptor,
) error {
	cindices, cnindices, err := cIndices(indices)
	if err != nil {
		return err
	}
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	info := Info(C.GrB_Vector_extract(w.grb, cmask, caccum, u.grb, cindices, cnindices, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Extract is the method variant of [VectorExtract].
func (vector Vector[D]) Extract(
	mask *Vector[bool],
	accum *BinaryOp[D, D, D],
	u Vector[D],
	indices []int,
	desc *Descriptor,
) error {
	return VectorExtract(vector, mask, accum, u, indices, desc)
}

// MatrixExtract extracts a sub-matrix from a larger matrix as specified by a set of row indices and a set of column
// indices. The result is a matrix whose size is equal to the size of the sets of indices.
//
// Parameters:
//
//   - c (INOUT): An existing GraphBLAS matrix. On input, the matrix provides values
//     that may be accumulated with the result of the extract operation. On output,
//     this matrix holds the results of the operation.
//
//   - mask (IN): An optional "write" [MatrixMask].
//
//   - accum (IN): An optional binary operator used for accumulating entries into existing c
//     entries. If assignment rather than accumulation is desired, nil should be specified.
//
//   - a (IN): The GraphBLAS matrix from which the subset is extracted.
//
//   - rowIndices (IN): The ordered set (slice) of indices corresponding to the rows of a from which
//     elements are extracted. If elements of all rows are to be extracted in order,
//     then [All](nrows) should be specified. Regardless of execution mode and return value, this slice
//     may be manipulated by the caller after this operation returns without affecting any deferred
//     computations for this operation. len(rowIndices) must be equal to nrows(c).
//
//   - colIndices (IN): The ordered set (slice) of indices corresponding to the columns of a from which
//     elements are extracted. If elements of all columns are to be extracted in order,
//     then [All](ncols) should be specified. Regardless of execution mode and return value, this slice
//     may be manipulated by the caller after this operation returns without affecting any deferred
//     computations for this operation. len(colIndices) must be equal to ncols(c).
//
//   - desc (IN): An optional operation [Descriptor]. If a default descriptor is desired,
//     nil should be specified.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [IndexOutOfBounds], [InvalidObject], [OutOfMemory], [Panic]
func MatrixExtract[D any](
	c Matrix[D],
	mask *Matrix[bool],
	accum *BinaryOp[D, D, D],
	a Matrix[D],
	rowIndices, colIndices []int,
	desc *Descriptor,
) error {
	crowindices, cnrows, err := cIndices(rowIndices)
	if err != nil {
		return err
	}
	ccolindices, cncols, err := cIndices(colIndices)
	if err != nil {
		return err
	}
	cmask, caccum, cdesc := processMADM(mask, accum, desc)
	info := Info(C.GrB_Matrix_extract(c.grb, cmask, caccum, a.grb, crowindices, cnrows, ccolindices, cncols, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Extract is the method variant of [MatrixExtract].
func (matrix Matrix[D]) Extract(
	mask *Matrix[bool],
	accum *BinaryOp[D, D, D],
	a Matrix[D],
	rowIndices, colIndices []int,
	desc *Descriptor,
) error {
	return MatrixExtract(matrix, mask, accum, a, rowIndices, colIndices, desc)
}

// MatrixColExtract extracts elements from one column of a matrix into a vector. Note that with the transpose
// descriptor for the source matrix, elements of an arbitrary row of the matrix can be extracted with this function
// as well.
//
// Parameters:
//
//   - w (INOUT): An existing GraphBLAS vector. On input, the vector provides values
//     that may be accumulated with the result of the extract operation. On output,
//     this vector holds the results of the operation.
//
//   - mask (IN): An optional "write" [VectorMask].
//
//   - accum (IN): An optional binary operator used for accumulating entries into existing w
//     entries. If assignment rather than accumulation is desired, nil should be specified.
//
//   - a (IN): The GraphBLAS matrix from which the column subset is extracted.
//
//   - rowIndices (IN): The ordered set (slice) of indices corresponding to the locations
//     within the specified column of a from which elements are extracted. If elements of all rows in a are to be
//     extracted in order, then [All](nrows) should be specified. Regardless of execution mode and return value, this
//     slice may be manipulated by the caller after this operation returns without affecting any deferred
//     computations for this operation. len(rowIndices) must be equal to size(w).
//
//   - colIndex (IN): The index of the column of a from which to extract values. It must be in the range [0, ncols(a)).
//
//   - desc (IN): An optional operation [Descriptor]. If a default descriptor is desired,
//     nil should be specified.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [InvalidIndex], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [IndexOutOfBounds], [InvalidObject], [OutOfMemory], [Panic]
func MatrixColExtract[D any](
	w Vector[D],
	mask *Vector[bool],
	accum *BinaryOp[D, D, D],
	a Matrix[D],
	rowIndices []int,
	colIndex int,
	desc *Descriptor,
) error {
	crowindices, cnrows, err := cIndices(rowIndices)
	if err != nil {
		return err
	}
	if colIndex < 0 {
		return makeError(InvalidIndex)
	}
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	info := Info(C.GrB_Col_extract(w.grb, cmask, caccum, a.grb, crowindices, cnrows, C.GrB_Index(colIndex), cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// ColExtract is the method variant of [MatrixColExtract].
func (vector Vector[D]) ColExtract(
	mask *Vector[bool],
	accum *BinaryOp[D, D, D],
	a Matrix[D],
	rowIndices []int,
	colIndex int,
	desc *Descriptor,
) error {
	return MatrixColExtract(vector, mask, accum, a, rowIndices, colIndex, desc)
}
