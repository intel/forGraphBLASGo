package GrB

// #include "GraphBLAS.h"
import "C"
import "unsafe"

// VectorSubassign is the same as [VectorAssign], except that the mask
// is restricted to the passed indices, and if the [Replace] descriptor
// is set for [Outp], then entries outside of the passed indices
// are not affected.
//
// VectorSubassign is a SuiteSparse:GraphBLAS extension.
func VectorSubassign[D any](
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
	info := Info(C.GxB_Vector_subassign(w.grb, cmask, caccum, u.grb, cindices, cnindices, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixSubassign is the same as [MatrixAssign], except that the mask
// is restricted to the passed indices, and if the [Replace] descriptor
// is set for [Outp], then entries outside of the passed indices
// are not affected.
//
// MatrixSubassign is a SuiteSparse:GraphBLAS extension.
func MatrixSubassign[D any](
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
	info := Info(C.GxB_Matrix_subassign(c.grb, cmask, caccum, a.grb, crowindices, cnrows, ccolindices, cncols, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixColSubassign is the same as [MatrixColAssign], except that the mask
// is restricted to the passed indices, and if the [Replace] descriptor
// is set for [Outp], then entries outside of the passed indices
// are not affected.
//
// MatrixColSubassign is a SuiteSparse:GraphBLAS extension.
func MatrixColSubassign[D any](
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
	info := Info(C.GxB_Col_subassign(c.grb, cmask, caccum, u.grb, crowindices, cnrows, C.GrB_Index(colIndex), cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixRowSubassign is the same as [MatrixRowAssign], except that the mask
// is restricted to the passed indices, and if the [Replace] descriptor
// is set for [Outp], then entries outside of the passed indices
// are not affected.
//
// MatrixRowSubassign is a SuiteSparse:GraphBLAS extension.
func MatrixRowSubassign[D any](
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
	info := Info(C.GxB_Row_subassign(c.grb, cmask, caccum, u.grb, C.GrB_Index(rowIndex), ccolindices, cncols, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// VectorSubassignConstant is the same as [VectorAssignConstant], except that the mask
// is restricted to the passed indices, and if the [Replace] descriptor
// is set for [Outp], then entries outside of the passed indices
// are not affected.
//
// VectorSubassignConstant is a SuiteSparse:GraphBLAS extension.
func VectorSubassignConstant[D any](
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
		info = Info(C.GxB_Vector_subassign_BOOL(w.grb, cmask, caccum, C.bool(x), cindices, cnindices, cdesc))
	case int:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GxB_Vector_subassign_INT32(w.grb, cmask, caccum, C.int32_t(x), cindices, cnindices, cdesc))
		} else {
			info = Info(C.GxB_Vector_subassign_INT64(w.grb, cmask, caccum, C.int64_t(x), cindices, cnindices, cdesc))
		}
	case int8:
		info = Info(C.GxB_Vector_subassign_INT8(w.grb, cmask, caccum, C.int8_t(x), cindices, cnindices, cdesc))
	case int16:
		info = Info(C.GxB_Vector_subassign_INT16(w.grb, cmask, caccum, C.int16_t(x), cindices, cnindices, cdesc))
	case int32:
		info = Info(C.GxB_Vector_subassign_INT32(w.grb, cmask, caccum, C.int32_t(x), cindices, cnindices, cdesc))
	case int64:
		info = Info(C.GxB_Vector_subassign_INT64(w.grb, cmask, caccum, C.int64_t(x), cindices, cnindices, cdesc))
	case uint:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GxB_Vector_subassign_UINT32(w.grb, cmask, caccum, C.uint32_t(x), cindices, cnindices, cdesc))
		} else {
			info = Info(C.GxB_Vector_subassign_UINT64(w.grb, cmask, caccum, C.uint64_t(x), cindices, cnindices, cdesc))
		}
	case uint8:
		info = Info(C.GxB_Vector_subassign_UINT8(w.grb, cmask, caccum, C.uint8_t(x), cindices, cnindices, cdesc))
	case uint16:
		info = Info(C.GxB_Vector_subassign_UINT16(w.grb, cmask, caccum, C.uint16_t(x), cindices, cnindices, cdesc))
	case uint32:
		info = Info(C.GxB_Vector_subassign_UINT32(w.grb, cmask, caccum, C.uint32_t(x), cindices, cnindices, cdesc))
	case uint64:
		info = Info(C.GxB_Vector_subassign_UINT64(w.grb, cmask, caccum, C.uint64_t(x), cindices, cnindices, cdesc))
	case float32:
		info = Info(C.GxB_Vector_subassign_FP32(w.grb, cmask, caccum, C.float(x), cindices, cnindices, cdesc))
	case float64:
		info = Info(C.GxB_Vector_subassign_FP64(w.grb, cmask, caccum, C.double(x), cindices, cnindices, cdesc))
	case complex64:
		info = Info(C.GxB_Vector_subassign_FC32(w.grb, cmask, caccum, C.complexfloat(x), cindices, cnindices, cdesc))
	case complex128:
		info = Info(C.GxB_Vector_subassign_FC64(w.grb, cmask, caccum, C.complexdouble(x), cindices, cnindices, cdesc))
	default:
		info = Info(C.GxB_Vector_subassign_UDT(w.grb, cmask, caccum, unsafe.Pointer(&val), cindices, cnindices, cdesc))
	}
	if info == success {
		return nil
	}
	return makeError(info)
}

// VectorSubassignScalar is the same as [VectorAssignScalar], except that the mask
// is restricted to the passed indices, and if the [Replace] descriptor
// is set for [Outp], then entries outside of the passed indices
// are not affected.
//
// VectorSubassignScalar is a SuiteSparse:GraphBLAS extension.
func VectorSubassignScalar[D any](
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
	info := Info(C.GxB_Vector_subassign_Scalar(w.grb, cmask, caccum, val.grb, cindices, cnindices, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixSubassignConstant is the same as [MatrixAssignConstant], except that the mask
// is restricted to the passed indices, and if the [Replace] descriptor
// is set for [Outp], then entries outside of the passed indices
// are not affected.
//
// MatrixSubassignConstant is a SuiteSparse:GraphBLAS extension.
func MatrixSubassignConstant[D any](
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
		info = Info(C.GxB_Matrix_subassign_BOOL(c.grb, cmask, caccum, C.bool(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	case int:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GxB_Matrix_subassign_INT32(c.grb, cmask, caccum, C.int32_t(x), crowindices, cnrows, ccolindices, cncols, cdesc))
		} else {
			info = Info(C.GxB_Matrix_subassign_INT64(c.grb, cmask, caccum, C.int64_t(x), crowindices, cnrows, ccolindices, cncols, cdesc))
		}
	case int8:
		info = Info(C.GxB_Matrix_subassign_INT8(c.grb, cmask, caccum, C.int8_t(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	case int16:
		info = Info(C.GxB_Matrix_subassign_INT16(c.grb, cmask, caccum, C.int16_t(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	case int32:
		info = Info(C.GxB_Matrix_subassign_INT32(c.grb, cmask, caccum, C.int32_t(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	case int64:
		info = Info(C.GxB_Matrix_subassign_INT64(c.grb, cmask, caccum, C.int64_t(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	case uint:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GxB_Matrix_subassign_UINT32(c.grb, cmask, caccum, C.uint32_t(x), crowindices, cnrows, ccolindices, cncols, cdesc))
		} else {
			info = Info(C.GxB_Matrix_subassign_UINT64(c.grb, cmask, caccum, C.uint64_t(x), crowindices, cnrows, ccolindices, cncols, cdesc))
		}
	case uint8:
		info = Info(C.GxB_Matrix_subassign_UINT8(c.grb, cmask, caccum, C.uint8_t(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	case uint16:
		info = Info(C.GxB_Matrix_subassign_UINT16(c.grb, cmask, caccum, C.uint16_t(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	case uint32:
		info = Info(C.GxB_Matrix_subassign_UINT32(c.grb, cmask, caccum, C.uint32_t(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	case uint64:
		info = Info(C.GxB_Matrix_subassign_UINT64(c.grb, cmask, caccum, C.uint64_t(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	case float32:
		info = Info(C.GxB_Matrix_subassign_FP32(c.grb, cmask, caccum, C.float(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	case float64:
		info = Info(C.GxB_Matrix_subassign_FP64(c.grb, cmask, caccum, C.double(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	case complex64:
		info = Info(C.GxB_Matrix_subassign_FC32(c.grb, cmask, caccum, C.complexfloat(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	case complex128:
		info = Info(C.GxB_Matrix_subassign_FC64(c.grb, cmask, caccum, C.complexdouble(x), crowindices, cnrows, ccolindices, cncols, cdesc))
	default:
		info = Info(C.GxB_Matrix_subassign_UDT(c.grb, cmask, caccum, unsafe.Pointer(&val), crowindices, cnrows, ccolindices, cncols, cdesc))
	}
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixSubassignScalar is the same as [MatrixSubassignScalar], except that the mask
// is restricted to the passed indices, and if the [Replace] descriptor
// is set for [Outp], then entries outside of the passed indices
// are not affected.
//
// MatrixSubassignScalar is a SuiteSparse:GraphBLAS extension.
func MatrixSubassignScalar[D any](
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
	info := Info(C.GxB_Matrix_subassign_Scalar(c.grb, cmask, caccum, val.grb, crowindices, cnrows, ccolindices, cncols, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}
