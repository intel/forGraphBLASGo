package GrB

// #include "GraphBLAS.h"
import "C"
import "unsafe"

// VectorAssign assigns values from one GraphBLAS vector to a subset of a vector as specified by a set of indices.
// The size of the input vector is the same size as the index slice provided.
//
// Parameters:
//
//   - w (INOUT): An existing GraphBLAS vector. On input, the vector provides values
//     that may be accumulated with the result of the assign operation. On output,
//     this vector holds the results of the operation.
//
//   - mask (IN): An optional "write" [VectorMask].
//
//   - accum (IN): An optional binary operator used for accumulating entries into existing w
//     entries. If assignment rather than accumulation is desired, nil should be specified.
//
//   - u (IN): The GraphBLAS vector whose contents are assigned to a subset of w.
//
//   - indices (IN): The ordered set (slice) of indices corresponding to the locations in w
//     that are to be assigned. If all elements of w are to be assigned in order from 0 to nindices − 1,
//     then [All](nindices) should be specified. Regardless of execution mode and return value, this slice
//     may be manipulated by the caller after this operation returns without affecting any deferred
//     computations for this operation. If this slice contains duplicate values, it implies an assignment
//     of more than one value to the same location which leads to undefined results. len(indices) must be
//     equal to size(u).
//
//   - desc (IN): An optional operation [Descriptor]. If a default descriptor is desired,
//     nil should be specified.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [IndexOutOfBounds], [InvalidObject], [OutOfMemory], [Panic]
func VectorAssign[D any](
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
	info := Info(C.GrB_Vector_assign(w.grb, cmask, caccum, u.grb, cindices, cnindices, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixAssign assigns values from one GraphBLAS matrix to a subset of a matrix as specified by a set of indices.
// The dimensions of the input matrix are the same size as the row and column index slices provided.
//
// Parameters:
//
//   - c (INOUT): An existing GraphBLAS matrix. On input, the matrix provides values
//     that may be accumulated with the result of the assign operation. On output,
//     this matrix holds the results of the operation.
//
//   - mask (IN): An optional "write" [MatrixMask].
//
//   - accum (IN): An optional binary operator used for accumulating entries into existing c
//     entries. If assignment rather than accumulation is desired, nil should be specified.
//
//   - a (IN): The GraphBLAS matrix whose contents are assigned to a subset of c.
//
//   - rowIndices (IN): The ordered set (slice) of indices corresponding to the rows of c
//     that are assigned. If all rows of c are to be assigned in order from 0 to nrows − 1,
//     then [All](nrows) should be specified. Regardless of execution mode and return value, this slice
//     may be manipulated by the caller after this operation returns without affecting any deferred
//     computations for this operation. If this slice contains duplicate values, it implies an assignment
//     of more than one value to the same location which leads to undefined results. len(rowIndices)
//     must be equal to nrows(a) if a is not transposed, or equal to ncols(a) if a is transposed.
//
//   - colIndices (IN): The ordered set (slice) of indices corresponding to the columns of c
//     that are assigned. If all columns of c are to be assigned in order from 0 to ncols − 1,
//     then [All](ncols) should be specified. Regardless of execution mode and return value, this slice
//     may be manipulated by the caller after this operation returns without affecting any deferred
//     computations for this operation. If this slice contains duplicate values, it implies an assignment
//     of more than one value to the same location which leads to undefined results. len(colIndices)
//     must be equal to ncols(a) if a is not transposed, or equal to nrows(a) if a is transposed.
//
//   - desc (IN): An optional operation [Descriptor]. If a default descriptor is desired,
//     nil should be specified.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [IndexOutOfBounds], [InvalidObject], [OutOfMemory], [Panic]
func MatrixAssign[D any](
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
	info := Info(C.GrB_Matrix_assign(c.grb, cmask, caccum, a.grb, crowindices, cnrows, ccolindices, cncols, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixColAssign assigns the contents of a vector to a subset of elements in one column of a matrix.
// Note that since the output cannot be transposed, [MatrixRowAssign] is also provided to
// assign to a row of a matrix.
//
// Parameters:
//
//   - c (INOUT): An existing GraphBLAS matrix. On input, the matrix provides values
//     that may be accumulated with the result of the assign operation. On output,
//     this matrix holds the results of the operation.
//
//   - mask (IN): An optional "write" [MatrixMask].
//
//   - accum (IN): An optional binary operator used for accumulating entries into existing c
//     entries. If assignment rather than accumulation is desired, nil should be specified.
//
//   - u (IN): The GraphBLAS vector whose contents are assigned to (a subset of) a column of c.
//
//   - rowIndices (IN): The ordered set (slice) of indices corresponding to the rows of c
//     that are assigned. If all rows of c are to be assigned in order from 0 to nrows − 1,
//     then [All](nrows) should be specified. Regardless of execution mode and return value, this slice
//     may be manipulated by the caller after this operation returns without affecting any deferred
//     computations for this operation. If this slice contains duplicate values, it implies an assignment
//     of more than one value to the same location which leads to undefined results. len(rowIndices) must be
//     equal to size(u).
//
//   - colIndex (IN): The index of the column in C to assign. Must be in the range [0, ncols(c)).
//
//   - desc (IN): An optional operation [Descriptor]. If a default descriptor is desired,
//     nil should be specified.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [InvalidIndex], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [IndexOutOfBounds], [InvalidObject], [OutOfMemory], [Panic]
func MatrixColAssign[D any](
	c Matrix[D],
	mask *Vector[bool],
	accum *BinaryOp[D, D, D],
	u Vector[D],
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
	info := Info(C.GrB_Col_assign(c.grb, cmask, caccum, u.grb, crowindices, cnrows, C.GrB_Index(colIndex), cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixRowAssign assigns the contents of a vector to a subset of elements in one row of a matrix.
// Note that since the output cannot be transposed, [MatrixColAssign] is also provided to
// assign to a column of a matrix.
//
// Parameters:
//
//   - c (INOUT): An existing GraphBLAS matrix. On input, the matrix provides values
//     that may be accumulated with the result of the assign operation. On output,
//     this matrix holds the results of the operation.
//
//   - mask (IN): An optional "write" [MatrixMask].
//
//   - accum (IN): An optional binary operator used for accumulating entries into existing c
//     entries. If assignment rather than accumulation is desired, nil should be specified.
//
//   - u (IN): The GraphBLAS vector whose contents are assigned to (a subset of) a row of c.
//
//   - rowIndex (IN): The index of the row in C to assign. Must be in the range [0, nrows(c)).
//
//   - colIndices (IN): The ordered set (slice) of indices corresponding to the columns of c
//     that are assigned. If all columns of c are to be assigned in order from 0 to nrows − 1,
//     then [All](ncols) should be specified. Regardless of execution mode and return value, this slice
//     may be manipulated by the caller after this operation returns without affecting any deferred
//     computations for this operation. If this slice contains duplicate values, it implies an assignment
//     of more than one value to the same location which leads to undefined results. len(colIndices) must
//     be equal to size(u).
//
//   - desc (IN): An optional operation [Descriptor]. If a default descriptor is desired,
//     nil should be specified.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [InvalidIndex], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [IndexOutOfBounds], [InvalidObject], [OutOfMemory], [Panic]
func MatrixRowAssign[D any](
	c Matrix[D],
	mask *Vector[bool],
	accum *BinaryOp[D, D, D],
	u Vector[D],
	rowIndex int,
	colIndices []int,
	desc *Descriptor,
) error {
	if rowIndex < 0 {
		return makeError(InvalidIndex)
	}
	ccolindices, cncols, err := cIndices(colIndices)
	if err != nil {
		return err
	}
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	info := Info(C.GrB_Row_assign(c.grb, cmask, caccum, u.grb, C.GrB_Index(rowIndex), ccolindices, cncols, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// VectorAssignConstant assigns the same value to a subset of a vector as specified by a set of indices.
// With the use of [All], the entire destination vector can be filled with the constant.
//
// To pass a [Scalar] object instead of a non-opaque variable, use [VectorAssignScalar].
//
// Parameters:
//
//   - w (INOUT): An existing GraphBLAS vector. On input, the vector provides values
//     that may be accumulated with the result of the assign operation. On output,
//     this vector holds the results of the operation.
//
//   - mask (IN): An optional "write" [VectorMask].
//
//   - accum (IN): An optional binary operator used for accumulating entries into existing w
//     entries. If assignment rather than accumulation is desired, nil should be specified.
//
//   - val (IN): Scalar value to assign to (a subset of) w.
//
//   - indices (IN): The ordered set (slice) of indices corresponding to the locations in w
//     that are to be assigned. If all elements of w are to be assigned in order from 0 to nindices − 1,
//     then [All](nindices) should be specified. Regardless of execution mode and return value, this slice
//     may be manipulated by the caller after this operation returns without affecting any deferred
//     computations for this operation. In this function, the specific order of the values in the slice
//     has no effect on the result. Unlike other assign functions, if there are duplicated values in this
//     slice the result is still defined. len(indices) must be in the range [0, size(w)]. If len(indices)
//     is zero, the operation becomes a NO-OP.
//
//   - desc (IN): An optional operation [Descriptor]. If a default descriptor is desired,
//     nil should be specified.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [IndexOutOfBounds], [InvalidObject], [OutOfMemory], [Panic]
func VectorAssignConstant[D any](
	w Vector[D],
	mask *Vector[bool],
	accum *BinaryOp[D, D, D],
	val D,
	indices []int,
	desc *Descriptor,
) error {
	cindices, cnindices, err := cIndices(indices)
	if err != nil {
		return err
	}
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	var info Info
	switch x := any(val).(type) {
	case bool:
		info = Info(C.GrB_Vector_assign_BOOL(w.grb, cmask, caccum, C.bool(x), cindices, cnindices, cdesc))
	case int:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Vector_assign_INT32(w.grb, cmask, caccum, C.int32_t(x), cindices, cnindices, cdesc))
		} else {
			info = Info(C.GrB_Vector_assign_INT64(w.grb, cmask, caccum, C.int64_t(x), cindices, cnindices, cdesc))
		}
	case int8:
		info = Info(C.GrB_Vector_assign_INT8(w.grb, cmask, caccum, C.int8_t(x), cindices, cnindices, cdesc))
	case int16:
		info = Info(C.GrB_Vector_assign_INT16(w.grb, cmask, caccum, C.int16_t(x), cindices, cnindices, cdesc))
	case int32:
		info = Info(C.GrB_Vector_assign_INT32(w.grb, cmask, caccum, C.int32_t(x), cindices, cnindices, cdesc))
	case int64:
		info = Info(C.GrB_Vector_assign_INT64(w.grb, cmask, caccum, C.int64_t(x), cindices, cnindices, cdesc))
	case uint:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Vector_assign_UINT32(w.grb, cmask, caccum, C.uint32_t(x), cindices, cnindices, cdesc))
		} else {
			info = Info(C.GrB_Vector_assign_UINT64(w.grb, cmask, caccum, C.uint64_t(x), cindices, cnindices, cdesc))
		}
	case uint8:
		info = Info(C.GrB_Vector_assign_UINT8(w.grb, cmask, caccum, C.uint8_t(x), cindices, cnindices, cdesc))
	case uint16:
		info = Info(C.GrB_Vector_assign_UINT16(w.grb, cmask, caccum, C.uint16_t(x), cindices, cnindices, cdesc))
	case uint32:
		info = Info(C.GrB_Vector_assign_UINT32(w.grb, cmask, caccum, C.uint32_t(x), cindices, cnindices, cdesc))
	case uint64:
		info = Info(C.GrB_Vector_assign_UINT64(w.grb, cmask, caccum, C.uint64_t(x), cindices, cnindices, cdesc))
	case float32:
		info = Info(C.GrB_Vector_assign_FP32(w.grb, cmask, caccum, C.float(x), cindices, cnindices, cdesc))
	case float64:
		info = Info(C.GrB_Vector_assign_FP64(w.grb, cmask, caccum, C.double(x), cindices, cnindices, cdesc))
	case complex64:
		info = Info(C.GxB_Vector_assign_FC32(w.grb, cmask, caccum, C.complexfloat(x), cindices, cnindices, cdesc))
	case complex128:
		info = Info(C.GxB_Vector_assign_FC64(w.grb, cmask, caccum, C.complexdouble(x), cindices, cnindices, cdesc))
	default:
		info = Info(C.GrB_Vector_assign_UDT(w.grb, cmask, caccum, unsafe.Pointer(&val), cindices, cnindices, cdesc))
	}
	if info == success {
		return nil
	}
	return makeError(info)
}

// VectorAssignScalar is like [VectorAssignConstant], except that the scalar value is passed as a [Scalar]
// object. It may be empty.
func VectorAssignScalar[D any](
	w Vector[D],
	mask *Vector[bool],
	accum *BinaryOp[D, D, D],
	val Scalar[D],
	indices []int,
	desc *Descriptor,
) error {
	cindices, cnindices, err := cIndices(indices)
	if err != nil {
		return err
	}
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	info := Info(C.GrB_Vector_assign_Scalar(w.grb, cmask, caccum, val.grb, cindices, cnindices, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixAssignConstant assigns the same value to a subset of a matrix as specified by a set of indices.
// With the use of [All], the entire destination matrix can be filled with the constant.
//
// To pass a [Scalar] object instead of a non-opaque variable, use [MatrixAssignScalar].
//
// Parameters:
//
//   - c (INOUT): An existing GraphBLAS matrix. On input, the matrix provides values
//     that may be accumulated with the result of the assign operation. On output,
//     this matrix holds the results of the operation.
//
//   - mask (IN): An optional "write" [MatrixMask].
//
//   - accum (IN): An optional binary operator used for accumulating entries into existing c
//     entries. If assignment rather than accumulation is desired, nil should be specified.
//
//   - val (IN): Scalar value to assign to (a subset of) c.
//
//   - rowIndices (IN): The ordered set (slice) of indices corresponding to the rows of c
//     that are to be assigned. If all rows of c are to be assigned in order from 0 to nrows − 1,
//     then [All](nrows) should be specified. Regardless of execution mode and return value, this slice
//     may be manipulated by the caller after this operation returns without affecting any deferred
//     computations for this operation. In this function, the specific order of the values in the slice
//     has no effect on the result. Unlike other assign functions, if there are duplicated values in this
//     slice the result is still defined. len(rowIndices) must be in the range [0, nrows(c)].
//     If len(rowIndices) is zero, the operation becomes a NO-OP.
//
//   - colIndices (IN): The ordered set (slice) of indices corresponding to the columns of c
//     that are to be assigned. If all columns of c are to be assigned in order from 0 to ncols − 1,
//     then [All](ncols) should be specified. Regardless of execution mode and return value, this slice
//     may be manipulated by the caller after this operation returns without affecting any deferred
//     computations for this operation. In this function, the specific order of the values in the slice
//     has no effect on the result. Unlike other assign functions, if there are duplicated values in this
//     slice the result is still defined. len(colIndices) must be in the range [0, ncols(c)].
//     If len(colIndices) is zero, the operation becomes a NO-OP.
//
//   - desc (IN): An optional operation [Descriptor]. If a default descriptor is desired,
//     nil should be specified.
//
// GraphBLAS API errors that may be returned:
//   - [DimensionMismatch], [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [IndexOutOfBounds], [InvalidObject], [OutOfMemory], [Panic]
func MatrixAssignConstant[D any](
	c Matrix[D],
	mask *Matrix[bool],
	accum *BinaryOp[D, D, D],
	val D,
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
	var info Info
	switch x := any(val).(type) {
	case bool:
		info = Info(C.GrB_Matrix_assign_BOOL(c.grb, cmask, caccum, C.bool(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	case int:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Matrix_assign_INT32(c.grb, cmask, caccum, C.int32_t(x), crowindices, cnrows, ccolindices, cncols, cdesc))
		} else {
			info = Info(C.GrB_Matrix_assign_INT64(c.grb, cmask, caccum, C.int64_t(x), crowindices, cnrows, ccolindices, cncols, cdesc))
		}
	case int8:
		info = Info(C.GrB_Matrix_assign_INT8(c.grb, cmask, caccum, C.int8_t(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	case int16:
		info = Info(C.GrB_Matrix_assign_INT16(c.grb, cmask, caccum, C.int16_t(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	case int32:
		info = Info(C.GrB_Matrix_assign_INT32(c.grb, cmask, caccum, C.int32_t(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	case int64:
		info = Info(C.GrB_Matrix_assign_INT64(c.grb, cmask, caccum, C.int64_t(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	case uint:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Matrix_assign_UINT32(c.grb, cmask, caccum, C.uint32_t(x), crowindices, cnrows, ccolindices, cncols, cdesc))
		} else {
			info = Info(C.GrB_Matrix_assign_UINT64(c.grb, cmask, caccum, C.uint64_t(x), crowindices, cnrows, ccolindices, cncols, cdesc))
		}
	case uint8:
		info = Info(C.GrB_Matrix_assign_UINT8(c.grb, cmask, caccum, C.uint8_t(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	case uint16:
		info = Info(C.GrB_Matrix_assign_UINT16(c.grb, cmask, caccum, C.uint16_t(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	case uint32:
		info = Info(C.GrB_Matrix_assign_UINT32(c.grb, cmask, caccum, C.uint32_t(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	case uint64:
		info = Info(C.GrB_Matrix_assign_UINT64(c.grb, cmask, caccum, C.uint64_t(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	case float32:
		info = Info(C.GrB_Matrix_assign_FP32(c.grb, cmask, caccum, C.float(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	case float64:
		info = Info(C.GrB_Matrix_assign_FP64(c.grb, cmask, caccum, C.double(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	case complex64:
		info = Info(C.GxB_Matrix_assign_FC32(c.grb, cmask, caccum, C.complexfloat(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	case complex128:
		info = Info(C.GxB_Matrix_assign_FC64(c.grb, cmask, caccum, C.complexdouble(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	default:
		info = Info(C.GrB_Matrix_assign_UDT(c.grb, cmask, caccum, unsafe.Pointer(&val), crowindices, cnrows, ccolindices, cncols, cdesc))
	}
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixAssignScalar is like [MatrixAssignConstant], except that the scalar value is passed as a [Scalar]
// object. It may be empty.
func MatrixAssignScalar[D any](
	c Matrix[D],
	mask *Matrix[bool],
	accum *BinaryOp[D, D, D],
	val Scalar[D],
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
	info := Info(C.GrB_Matrix_assign_Scalar(c.grb, cmask, caccum, val.grb, crowindices, cnrows, ccolindices, cncols, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}
