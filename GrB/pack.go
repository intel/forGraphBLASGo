package GrB

// #include "GraphBLAS.h"
import "C"
import (
	"unsafe"
)

func realType(typ Type) Type {
	if typ == Int || typ == Uint {
		if unsafe.Sizeof(0) == 4 {
			if typ == Int {
				return Int32
			}
			return Uint32
		} else {
			if typ == Int {
				return Int64
			}
			return Uint64
		}
	}
	return typ
}

func checkType[D any](typ Type, ok bool, err error) error {
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	var d D
	if realType(TypeOf(d)) != realType(typ) {
		return makeError(DomainMismatch)
	}
	return nil
}

// Vector CSC

// PackCSC packs a vector from two user slices in CSC format.
//
// In the resulting vector, the CSC format is a sparse vector with [ByCol] [Layout].
// The vector must exist on input with the right type and size. No type casting is done,
// so the domain D must correctly reflect [Vector.Type]().
//
// Parameters:
//
//   - vi (INOUT): A [SystemSlice] of indices of the corresponding values in vx, without duplicates.
//     This is not checked, so the result is undefined if this is not the case.
//
//   - vx (INOUT): A [SystemSlice] of values.
//
//   - iso (IN): If true, the resulting vector is iso.
//
//   - nvals (IN): The number of entries in the resulting vector.
//
//   - jumbled (IN): If false, the indices in vi must appear in sorted order. This is not
//     checked, so the result is undefined if this is not the case.
//
// On successful return, vi and vx are empty, to indicate that the user application no longer
// owns them. They have instead been moved to the resulting vector. If not successful, vi and vx
// are not modified.
//
// PackCSC is a SuiteSparse:GraphBLAS extension.
func (vector Vector[D]) PackCSC(vi *SystemSlice[int], vx *SystemSlice[D], iso bool, nvals int, jumbled bool, desc *Descriptor) error {
	if err := checkType[D](vector.Type()); err != nil {
		return err
	}
	uvi, viCopied := makeGrBIndexSlice(vi)
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Vector_pack_CSC(
		vector.grb,
		(**C.GrB_Index)(unsafe.Pointer(&uvi.ptr)), &vx.ptr,
		C.GrB_Index(uvi.size), C.GrB_Index(vx.size),
		C.bool(iso), C.GrB_Index(nvals), C.bool(jumbled), cdesc,
	))
	if info != success {
		if viCopied {
			uvi.Free()
		}
		return makeError(info)
	}
	vi.Free()
	vx.size = 0
	return nil
}

// UnpackCSC unpacks a vector to user slices in CSC format.
//
// No type casting is done, so the domain D must correctly reflect [Vector.Type]().
//
// Parameters:
//
//   - allowJumbled (IN): If false, the indices in vi appear in ascending order. If true,
//     the indices may appear in any order.
//
// Return Values:
//
//   - vi: A [SystemSlice] of indices of the corresponding values in vx.
//
//   - vx: A [SystemSlice] of values.
//
//   - iso: If true, the input vector was iso.
//
//   - nvals: The number of entries in the resulting slices.
//
//   - jumbled: If false, the indices in vi appear in sorted order.
//
// On successful return, the input vector has no entries anymore, and the user application now
// owns the resulting slices. If not successful, the input vector is not modified.
//
// UnpackCSC is a SuiteSparse:GraphBLAS extension.
func (vector Vector[D]) UnpackCSC(allowJumbled bool, desc *Descriptor) (vi SystemSlice[int], vx SystemSlice[D], iso bool, nvals int, jumbled bool, err error) {
	if err = checkType[D](vector.Type()); err != nil {
		return
	}
	var cvi *C.GrB_Index
	var cvx unsafe.Pointer
	var cviSize, cvxSize, cnvals C.GrB_Index
	var ciso C.bool
	cjumbled := C.bool(false)
	cdesc := processDescriptor(desc)
	var info Info
	if allowJumbled {
		info = Info(C.GxB_Vector_unpack_CSC(vector.grb, &cvi, &cvx, &cviSize, &cvxSize, &ciso, &cnvals, &cjumbled, cdesc))
	} else {
		info = Info(C.GxB_Vector_unpack_CSC(vector.grb, &cvi, &cvx, &cviSize, &cvxSize, &ciso, &cnvals, nil, cdesc))
	}
	if info != success {
		err = makeError(info)
		return
	}
	uvi := AsSystemSlice[uint64](unsafe.Pointer(cvi), int(cviSize))
	vi, uviCopied := makeGoIndexSlice(&uvi)
	if uviCopied {
		uvi.Free()
	}
	vx = AsSystemSlice[D](cvx, int(cvxSize))
	iso = bool(ciso)
	nvals = int(cnvals)
	jumbled = bool(cjumbled)
	return
}

// Vector Bitmap

// PackBitmap packs a vector from two user slices in bitmap format.
//
// The vector must exist on input with the right type and size. No type casting is done,
// so the domain D must correctly reflect [Vector.Type]().
//
// Parameters:
//
//   - vb (INOUT): A [SystemSlice] that indicates which indices are present: If vb is true
//     at a given index, then the entry at that index is present with value given by vx at
//     the same index. If vb is false at a given index, the entry at that index is not present,
//     and the value given by vx at that index is ignored.
//
//   - vx (INOUT): A [SystemSlice] of values.
//
//   - iso (IN): If true, the resulting vector is iso.
//
//   - nvals (IN): The number of entries in the resulting vector.
//
// On successful return, vb and vx are empty, to indicate that the user application no longer
// owns them. They have instead been moved to the resulting vector. If not successful, vb and vx
// are not modified.
//
// PackBitmap is a SuiteSparse:GraphBLAS extension.
func (vector Vector[D]) PackBitmap(vb *SystemSlice[bool], vx *SystemSlice[D], iso bool, nvals int, desc *Descriptor) error {
	if err := checkType[D](vector.Type()); err != nil {
		return err
	}
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Vector_pack_Bitmap(
		vector.grb,
		(**C.int8_t)(unsafe.Pointer(&vb.ptr)), &vx.ptr,
		C.GrB_Index(vb.size), C.GrB_Index(vx.size),
		C.bool(iso), C.GrB_Index(nvals), cdesc,
	))
	if info != success {
		return makeError(info)
	}
	vb.size = 0
	vx.size = 0
	return nil
}

// UnpackBitmap unpacks a vector to user slices in bitmap format.
//
// No type casting is done, so the domain D must correctly reflect [Vector.Type]().
//
// Return Values:
//
//   - vb: A [SystemSlice] that indicates which indices are present: If vb is true
//     at a given index, then the entry at that index is present with value given by vx at
//     the same index. If vb is false at a given index, the entry at that index is not present,
//     and the value given by vx at that index is ignored.
//
//   - vx: A [SystemSlice] of values.
//
//   - iso: If true, the input vector was iso.
//
//   - nvals: The number of entries in the resulting slices.
//
// On successful return, the input vector has no entries anymore, and the user application now
// owns the resulting slices. If not successful, the input vector is not modified.
//
// UnpackBitmap is a SuiteSparse:GraphBLAS extension.
func (vector Vector[D]) UnpackBitmap(desc *Descriptor) (vb SystemSlice[bool], vx SystemSlice[D], iso bool, nvals int, err error) {
	if err = checkType[D](vector.Type()); err != nil {
		return
	}
	var cvb *C.int8_t
	var cvx unsafe.Pointer
	var cvbSize, cvxSize, cnvals C.GrB_Index
	var ciso C.bool
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Vector_unpack_Bitmap(vector.grb, &cvb, &cvx, &cvbSize, &cvxSize, &ciso, &cnvals, cdesc))
	if info != success {
		err = makeError(info)
		return
	}
	vb = AsSystemSlice[bool](unsafe.Pointer(cvb), int(cvbSize))
	vx = AsSystemSlice[D](cvx, int(cvxSize))
	iso = bool(ciso)
	nvals = int(cnvals)
	return
}

// Vector Full

// PackFull packs a vector from a user slice in full format.
//
// The vector must exist on input with the right type and size. No type casting is done,
// so the domain D must correctly reflect [Vector.Type]().
//
// Parameters:
//
//   - vx (INOUT): A [SystemSlice] of values. All entries with index < nvals are present.
//
//   - iso (IN): If true, the resulting vector is iso.
//
//   - nvals (IN): The number of entries in the resulting vector.
//
// On successful return, vx is empty, to indicate that the user application no longer
// owns it. It has instead been moved to the resulting vector. If not successful, vx
// is not modified.
//
// PackFull is a SuiteSparse:GraphBLAS extension.
func (vector Vector[D]) PackFull(vx *SystemSlice[D], iso bool, desc *Descriptor) error {
	if err := checkType[D](vector.Type()); err != nil {
		return err
	}
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Vector_pack_Full(vector.grb, &vx.ptr, C.GrB_Index(vx.size), C.bool(iso), cdesc))
	if info != success {
		return makeError(info)
	}
	vx.size = 0
	return nil
}

// UnpackFull unpacks a vector to a user slice in full format.
//
// No type casting is done, so the domain D must correctly reflect [Vector.Type]().
//
// Return Values:
//
//   - vx: A [SystemSlice] of values. All entries with index < size of vector are present.
//
//   - iso: If true, the input vector was iso.
//
//   - nvals: The number of entries in the resulting slices.
//
// On successful return, the input vector has no entries anymore, and the user application now
// owns the resulting slice. If not successful, the input vector is not modified.
//
// UnpackFull is a SuiteSparse:GraphBLAS extension.
func (vector Vector[D]) UnpackFull(desc *Descriptor) (vx SystemSlice[D], iso bool, err error) {
	if err = checkType[D](vector.Type()); err != nil {
		return
	}
	var cvx unsafe.Pointer
	var cvxSize C.GrB_Index
	var ciso C.bool
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Vector_unpack_Full(vector.grb, &cvx, &cvxSize, &ciso, cdesc))
	if info != success {
		err = makeError(info)
		return
	}
	vx = AsSystemSlice[D](cvx, int(cvxSize))
	iso = bool(ciso)
	return
}

// Matrix CSR

// PackCSR packs a matrix from three user slices in CSR format.
//
// In the resulting matrix, the CSR format is a sparse matrix with [ByRow] [Layout].
// The matrix must exist on input with the right type and dimensions. No type casting is done,
// so the domain D must correctly reflect [Matrix.Type]().
//
// Parameters:
//
//   - ap (INOUT): A [SystemSlice] of integers that defines where the column indices and
//     values appear in aj and ax, for each row. The number of entries in row i is given
//     by the expression ap[i+1] - ap[i]. ap must therefore have a length of at least
//     nrows + 1.
//
//   - aj (INOUT): A [SystemSlice] of column indices of the corresponding values in ax in each row,
//     without duplicates per row. This is not checked, so the result is undefined if this is not the case.
//
//   - ax (INOUT): A [SystemSlice] of values.
//
//   - iso (IN): If true, the resulting matrix is iso.
//
//   - nvals (IN): The number of entries in the resulting matrix.
//
//   - jumbled (IN): If false, the indices in aj must appear in sorted order within each row.
//     This is not checked, so the result is undefined if this is not the case.
//
// On successful return, ap, aj and ax are empty, to indicate that the user application no longer
// owns them. They have instead been moved to the resulting matrix. If not successful, ap, aj and ax
// are not modified.
//
// PackCSR is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) PackCSR(ap, aj *SystemSlice[int], ax *SystemSlice[D], iso, jumbled bool, desc *Descriptor) error {
	if err := checkType[D](matrix.Type()); err != nil {
		return err
	}
	uap, apCopied := makeGrBIndexSlice(ap)
	uaj, ajCopied := makeGrBIndexSlice(aj)
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Matrix_pack_CSR(
		matrix.grb,
		(**C.GrB_Index)(unsafe.Pointer(&uap.ptr)), (**C.GrB_Index)(unsafe.Pointer(&uaj.ptr)), &ax.ptr,
		C.GrB_Index(uap.size), C.GrB_Index(uaj.size), C.GrB_Index(ax.size),
		C.bool(iso), C.bool(jumbled), cdesc,
	))
	if info != success {
		if apCopied {
			uap.Free()
		}
		if ajCopied {
			uaj.Free()
		}
		return makeError(info)
	}
	ap.Free()
	aj.Free()
	ax.size = 0
	return nil
}

// UnpackCSR unpacks a matrix to user slices in CSR format.
//
// No type casting is done, so the domain D must correctly reflect [Vector.Type]().
//
// Parameters:
//
//   - allowJumbled (IN): If false, the indices in aj appear in ascending order within each row.
//     If true, the indices may appear in any order within each row.
//
// Return Values:
//
//   - ap: A [SystemSlice] of integers that defines where the column indices and
//     values appear in aj and ax, for each row. The number of entries in row i is given
//     by the expression ap[i+1] - ap[i]. ap therefore has a length of at least
//     nrows + 1.
//
//   - aj: A [SystemSlice] of column indices of the corresponding values in ax in each row,
//     without duplicates per row.
//
//   - ax: A [SystemSlice] of values.
//
//   - iso: If true, the input matrix was iso.
//
//   - nvals: The number of entries in the resulting slices.
//
//   - jumbled: If false, the indices in aj appear in sorted order within each row.
//
// On successful return, the input matrix has no entries anymore, and the user application now
// owns the resulting slices. If not successful, the input matrix is not modified.
//
// UnpackCSR is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) UnpackCSR(allowJumbled bool, desc *Descriptor) (ap, aj SystemSlice[int], ax SystemSlice[D], iso, jumbled bool, err error) {
	if err = checkType[D](matrix.Type()); err != nil {
		return
	}
	var cap, caj *C.GrB_Index
	var cax unsafe.Pointer
	var capSize, cajSize, caxSize C.GrB_Index
	var ciso C.bool
	cjumbled := C.bool(false)
	cdesc := processDescriptor(desc)
	var info Info
	if allowJumbled {
		info = Info(C.GxB_Matrix_unpack_CSR(matrix.grb, &cap, &caj, &cax, &capSize, &cajSize, &caxSize, &ciso, &cjumbled, cdesc))
	} else {
		info = Info(C.GxB_Matrix_unpack_CSR(matrix.grb, &cap, &caj, &cax, &capSize, &cajSize, &caxSize, &ciso, nil, cdesc))
	}
	if info != success {
		err = makeError(info)
		return
	}
	uap := AsSystemSlice[uint64](unsafe.Pointer(cap), int(capSize))
	ap, uapCopied := makeGoIndexSlice(&uap)
	if uapCopied {
		uap.Free()
	}
	uaj := AsSystemSlice[uint64](unsafe.Pointer(caj), int(cajSize))
	aj, uajCopied := makeGoIndexSlice(&uaj)
	if uajCopied {
		uaj.Free()
	}
	ax = AsSystemSlice[D](cax, int(caxSize))
	iso = bool(ciso)
	jumbled = bool(cjumbled)
	return
}

// Matrix CSC

// PackCSC packs a matrix from three user slices in CSC format.
//
// In the resulting matrix, the CSC format is a sparse matrix with [ByCol] [Layout].
// The matrix must exist on input with the right type and dimensions. No type casting is done,
// so the domain D must correctly reflect [Matrix.Type]().
//
// Parameters:
//
//   - ap (INOUT): A [SystemSlice] of integers that defines where the row indices and
//     values appear in ai and ax, for each column. The number of entries in column i is given
//     by the expression ap[i+1] - ap[i]. ap must therefore have a length of at least
//     ncols + 1.
//
//   - ai (INOUT): A [SystemSlice] of row indices of the corresponding values in ax in each column,
//     without duplicates per column. This is not checked, so the result is undefined if this is not the case.
//
//   - ax (INOUT): A [SystemSlice] of values.
//
//   - iso (IN): If true, the resulting matrix is iso.
//
//   - nvals (IN): The number of entries in the resulting matrix.
//
//   - jumbled (IN): If false, the indices in ai must appear in sorted order within each column.
//     This is not checked, so the result is undefined if this is not the case.
//
// On successful return, ap, ai and ax are empty, to indicate that the user application no longer
// owns them. They have instead been moved to the resulting matrix. If not successful, ap, ai and ax
// are not modified.
//
// PackCSC is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) PackCSC(ap, ai *SystemSlice[int], ax *SystemSlice[D], iso, jumbled bool, desc *Descriptor) error {
	if err := checkType[D](matrix.Type()); err != nil {
		return err
	}
	uap, apCopied := makeGrBIndexSlice(ap)
	uai, aiCopied := makeGrBIndexSlice(ai)
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Matrix_pack_CSC(
		matrix.grb,
		(**C.GrB_Index)(unsafe.Pointer(&uap.ptr)), (**C.GrB_Index)(unsafe.Pointer(&uai.ptr)), &ax.ptr,
		C.GrB_Index(uap.size), C.GrB_Index(uai.size), C.GrB_Index(ax.size),
		C.bool(iso), C.bool(jumbled), cdesc,
	))
	if info != success {
		if apCopied {
			uap.Free()
		}
		if aiCopied {
			uai.Free()
		}
		return makeError(info)
	}
	ap.Free()
	ai.Free()
	ax.size = 0
	return nil
}

// UnpackCSC unpacks a matrix to user slices in CSC format.
//
// No type casting is done, so the domain D must correctly reflect [Vector.Type]().
//
// Parameters:
//
//   - allowJumbled (IN): If false, the indices in ai appear in ascending order within each column.
//     If true, the indices may appear in any order within each column.
//
// Return Values:
//
//   - ap: A [SystemSlice] of integers that defines where the row indices and
//     values appear in ai and ax, for each column. The number of entries in column i is given
//     by the expression ap[i+1] - ap[i]. ap therefore has a length of at least
//     ncols + 1.
//
//   - ai: A [SystemSlice] of row indices of the corresponding values in ax in each column,
//     without duplicates per column.
//
//   - ax: A [SystemSlice] of values.
//
//   - iso: If true, the input matrix was iso.
//
//   - nvals: The number of entries in the resulting slices.
//
//   - jumbled: If false, the indices in ai appear in sorted order within each column.
//
// On successful return, the input matrix has no entries anymore, and the user application now
// owns the resulting slices. If not successful, the input matrix is not modified.
//
// UnpackCSC is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) UnpackCSC(allowJumbled bool, desc *Descriptor) (ap, ai SystemSlice[int], ax SystemSlice[D], iso, jumbled bool, err error) {
	if err = checkType[D](matrix.Type()); err != nil {
		return
	}
	var cap, cai *C.GrB_Index
	var cax unsafe.Pointer
	var capSize, caiSize, caxSize C.GrB_Index
	var ciso C.bool
	cjumbled := C.bool(false)
	cdesc := processDescriptor(desc)
	var info Info
	if allowJumbled {
		info = Info(C.GxB_Matrix_unpack_CSC(matrix.grb, &cap, &cai, &cax, &capSize, &caiSize, &caxSize, &ciso, &cjumbled, cdesc))
	} else {
		info = Info(C.GxB_Matrix_unpack_CSC(matrix.grb, &cap, &cai, &cax, &capSize, &caiSize, &caxSize, &ciso, nil, cdesc))
	}
	if info != success {
		err = makeError(info)
		return
	}
	uap := AsSystemSlice[uint64](unsafe.Pointer(cap), int(capSize))
	ap, uapCopied := makeGoIndexSlice(&uap)
	if uapCopied {
		uap.Free()
	}
	uai := AsSystemSlice[uint64](unsafe.Pointer(cai), int(caiSize))
	ai, uaiCopied := makeGoIndexSlice(&uai)
	if uaiCopied {
		uai.Free()
	}
	ax = AsSystemSlice[D](cax, int(caxSize))
	iso = bool(ciso)
	jumbled = bool(cjumbled)
	return
}

// Matrix HyperCSR

// PackHyperCSR packs a matrix from four user slices in hypersparse CSR format.
//
// In the resulting matrix, the hypersparse CSR format is a sparse matrix with [ByRow] [Layout].
// The matrix must exist on input with the right type and dimensions. No type casting is done,
// so the domain D must correctly reflect [Matrix.Type]().
//
// Parameters:
//
//   - ap (INOUT): A [SystemSlice] of integers that defines where the column indices and
//     values appear in aj and ax, for each present row. The number of entries in row ah[i] is given
//     by the expression ap[i+1] - ap[i].
//
//   - ah (INOUT): A [SystemSlice] of integers that defines which rows (row indices) are present. Rows
//     that are not present are considered completely empty.
//
//   - aj (INOUT): A [SystemSlice] of column indices of the corresponding values in ax in present rows,
//     without duplicates per row. This is not checked, so the result is undefined if this is not the case.
//
//   - ax (INOUT): A [SystemSlice] of values.
//
//   - iso (IN): If true, the resulting matrix is iso.
//
//   - nvec (IN): The number of rows that appear in ah.
//
//   - jumbled (IN): If false, the indices in aj must appear in sorted order within each row.
//     This is not checked, so the result is undefined if this is not the case.
//
// On successful return, ap, ah, aj and ax are empty, to indicate that the user application no longer
// owns them. They have instead been moved to the resulting matrix. If not successful, ap, ah, aj and ax
// are not modified.
//
// PackHyperCSR is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) PackHyperCSR(ap, ah, aj *SystemSlice[int], ax *SystemSlice[D], iso bool, nvec int, jumbled bool, desc *Descriptor) error {
	if err := checkType[D](matrix.Type()); err != nil {
		return err
	}
	uap, apCopied := makeGrBIndexSlice(ap)
	uah, ahCopied := makeGrBIndexSlice(ah)
	uaj, ajCopied := makeGrBIndexSlice(aj)
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Matrix_pack_HyperCSR(
		matrix.grb,
		(**C.GrB_Index)(unsafe.Pointer(&uap.ptr)), (**C.GrB_Index)(unsafe.Pointer(&uah.ptr)),
		(**C.GrB_Index)(unsafe.Pointer(&uaj.ptr)), &ax.ptr,
		C.GrB_Index(uap.size), C.GrB_Index(uah.size), C.GrB_Index(uaj.size), C.GrB_Index(ax.size),
		C.bool(iso), C.GrB_Index(nvec), C.bool(jumbled), cdesc,
	))
	if info != success {
		if apCopied {
			uap.Free()
		}
		if ahCopied {
			uah.Free()
		}
		if ajCopied {
			uaj.Free()
		}
		return makeError(info)
	}
	ap.Free()
	ah.Free()
	aj.Free()
	ax.size = 0
	return nil
}

// UnpackHyperCSR unpacks a matrix to user slices in hypersparse CSR format.
//
// No type casting is done, so the domain D must correctly reflect [Vector.Type]().
//
// Parameters:
//
//   - allowJumbled (IN): If false, the indices in aj appear in ascending order within each row.
//     If true, the indices may appear in any order within each row.
//
// Return Values:
//
//   - ap: A [SystemSlice] of integers that defines where the column indices and
//     values appear in aj and ax, for each present row. The number of entries in row ah[i] is given
//     by the expression ap[i+1] - ap[i].
//
//   - ah: A [SystemSlice] of integers that defines which rows (row indices) are present. Rows
//     that are not present are considered completely empty.
//
//   - aj: A [SystemSlice] of column indices of the corresponding values in ax in present rows,
//     without duplicates per row.
//
//   - ax: A [SystemSlice] of values.
//
//   - iso: If true, the input matrix was iso.
//
//   - nvec: The number of rows that appear in ah.
//
//   - jumbled: If false, the indices in aj appear in sorted order within each row.
//
// On successful return, the input matrix has no entries anymore, and the user application now
// owns the resulting slices. If not successful, the input matrix is not modified.
//
// UnpackHyperCSR is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) UnpackHyperCSR(allowJumbled bool, desc *Descriptor) (ap, ah, aj SystemSlice[int], ax SystemSlice[D], iso bool, nvec int, jumbled bool, err error) {
	if err = checkType[D](matrix.Type()); err != nil {
		return
	}
	var cap, cah, caj *C.GrB_Index
	var cax unsafe.Pointer
	var capSize, cahSize, cajSize, caxSize, cnvec C.GrB_Index
	var ciso C.bool
	cjumbled := C.bool(false)
	cdesc := processDescriptor(desc)
	var info Info
	if allowJumbled {
		info = Info(C.GxB_Matrix_unpack_HyperCSR(matrix.grb, &cap, &cah, &caj, &cax, &capSize, &cahSize, &cajSize, &caxSize, &ciso, &cnvec, &cjumbled, cdesc))
	} else {
		info = Info(C.GxB_Matrix_unpack_HyperCSR(matrix.grb, &cap, &cah, &caj, &cax, &capSize, &cahSize, &cajSize, &caxSize, &ciso, &cnvec, nil, cdesc))
	}
	if info != success {
		err = makeError(info)
		return
	}
	uap := AsSystemSlice[uint64](unsafe.Pointer(cap), int(capSize))
	ap, uapCopied := makeGoIndexSlice(&uap)
	if uapCopied {
		uap.Free()
	}
	uah := AsSystemSlice[uint64](unsafe.Pointer(cah), int(cahSize))
	ah, uahCopied := makeGoIndexSlice(&uah)
	if uahCopied {
		uah.Free()
	}
	uaj := AsSystemSlice[uint64](unsafe.Pointer(caj), int(cajSize))
	aj, uajCopied := makeGoIndexSlice(&uaj)
	if uajCopied {
		uaj.Free()
	}
	ax = AsSystemSlice[D](cax, int(caxSize))
	iso = bool(ciso)
	nvec = int(cnvec)
	jumbled = bool(cjumbled)
	return
}

// Matrix HyperCSC

// PackHyperCSC packs a matrix from four user slices in hypersparse CSC format.
//
// In the resulting matrix, the hypersparse CSC format is a sparse matrix with [ByCol] [Layout].
// The matrix must exist on input with the right type and dimensions. No type casting is done,
// so the domain D must correctly reflect [Matrix.Type]().
//
// Parameters:
//
//   - ap (INOUT): A [SystemSlice] of integers that defines where the row indices and
//     values appear in ai and ax, for each present column. The number of entries in column ah[i] is given
//     by the expression ap[i+1] - ap[i].
//
//   - ah (INOUT): A [SystemSlice] of integers that defines which columns (column indices) are present. Columns
//     that are not present are considered completely empty.
//
//   - ai (INOUT): A [SystemSlice] of row indices of the corresponding values in ax in present columns,
//     without duplicates per columns. This is not checked, so the result is undefined if this is not the case.
//
//   - ax (INOUT): A [SystemSlice] of values.
//
//   - iso (IN): If true, the resulting matrix is iso.
//
//   - nvec (IN): The number of columns that appear in ah.
//
//   - jumbled (IN): If false, the indices in aj must appear in sorted order within each column.
//     This is not checked, so the result is undefined if this is not the case.
//
// On successful return, ap, ah, ai and ax are empty, to indicate that the user application no longer
// owns them. They have instead been moved to the resulting matrix. If not successful, ap, ah, ai and ax
// are not modified.
//
// PackHyperCSC is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) PackHyperCSC(ap, ah, ai *SystemSlice[int], ax *SystemSlice[D], iso bool, nvec int, jumbled bool, desc *Descriptor) error {
	if err := checkType[D](matrix.Type()); err != nil {
		return err
	}
	uap, apCopied := makeGrBIndexSlice(ap)
	uah, ahCopied := makeGrBIndexSlice(ah)
	uai, aiCopied := makeGrBIndexSlice(ai)
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Matrix_pack_HyperCSC(
		matrix.grb,
		(**C.GrB_Index)(unsafe.Pointer(&uap.ptr)), (**C.GrB_Index)(unsafe.Pointer(&uah.ptr)),
		(**C.GrB_Index)(unsafe.Pointer(&uai.ptr)), &ax.ptr,
		C.GrB_Index(uap.size), C.GrB_Index(uah.size), C.GrB_Index(uai.size), C.GrB_Index(ax.size),
		C.bool(iso), C.GrB_Index(nvec), C.bool(jumbled), cdesc,
	))
	if info != success {
		if apCopied {
			uap.Free()
		}
		if ahCopied {
			uah.Free()
		}
		if aiCopied {
			uai.Free()
		}
		return makeError(info)
	}
	ap.Free()
	ah.Free()
	ai.Free()
	ax.size = 0
	return nil
}

// UnpackHyperCSC unpacks a matrix to user slices in hypersparse CSC format.
//
// No type casting is done, so the domain D must correctly reflect [Vector.Type]().
//
// Parameters:
//
//   - allowJumbled (IN): If false, the indices in ai appear in ascending order within each column.
//     If true, the indices may appear in any order within each column.
//
// Return Values:
//
//   - ap: A [SystemSlice] of integers that defines where the rows indices and
//     values appear in ai and ax, for each present column. The number of entries in column ah[i] is given
//     by the expression ap[i+1] - ap[i].
//
//   - ah: A [SystemSlice] of integers that defines which columns (column indices) are present. Columns
//     that are not present are considered completely empty.
//
//   - ai: A [SystemSlice] of row indices of the corresponding values in ax in present columns,
//     without duplicates per columns.
//
//   - ax: A [SystemSlice] of values.
//
//   - iso: If true, the input matrix was iso.
//
//   - nvec: The number of columns that appear in ah.
//
//   - jumbled: If false, the indices in aj appear in sorted order within each column.
//
// On successful return, the input matrix has no entries anymore, and the user application now
// owns the resulting slices. If not successful, the input matrix is not modified.
//
// UnpackHyperCSC is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) UnpackHyperCSC(allowJumbled bool, desc *Descriptor) (ap, ah, ai SystemSlice[int], ax SystemSlice[D], iso bool, nvec int, jumbled bool, err error) {
	if err = checkType[D](matrix.Type()); err != nil {
		return
	}
	var cap, cah, cai *C.GrB_Index
	var cax unsafe.Pointer
	var capSize, cahSize, caiSize, caxSize, cnvec C.GrB_Index
	var ciso C.bool
	cjumbled := C.bool(false)
	cdesc := processDescriptor(desc)
	var info Info
	if allowJumbled {
		info = Info(C.GxB_Matrix_unpack_HyperCSC(matrix.grb, &cap, &cah, &cai, &cax, &capSize, &cahSize, &caiSize, &caxSize, &ciso, &cnvec, &cjumbled, cdesc))
	} else {
		info = Info(C.GxB_Matrix_unpack_HyperCSC(matrix.grb, &cap, &cah, &cai, &cax, &capSize, &cahSize, &caiSize, &caxSize, &ciso, &cnvec, nil, cdesc))
	}
	if info != success {
		err = makeError(info)
		return
	}
	uap := AsSystemSlice[uint64](unsafe.Pointer(cap), int(capSize))
	ap, uapCopied := makeGoIndexSlice(&uap)
	if uapCopied {
		uap.Free()
	}
	uah := AsSystemSlice[uint64](unsafe.Pointer(cah), int(cahSize))
	ah, uahCopied := makeGoIndexSlice(&uah)
	if uahCopied {
		uah.Free()
	}
	uai := AsSystemSlice[uint64](unsafe.Pointer(cai), int(caiSize))
	ai, uaiCopied := makeGoIndexSlice(&uai)
	if uaiCopied {
		uai.Free()
	}
	ax = AsSystemSlice[D](cax, int(caxSize))
	iso = bool(ciso)
	nvec = int(cnvec)
	jumbled = bool(cjumbled)
	return
}

// UnpackHyperHash unpacks the hyper-hash from the hypersparse matrix. The hyper-hash
// of a hypersparse matrix provides quick access to the inverse of ah.
//
// If the matrix is not hypersparse or does not yet have a hyper-hash, then the resulting
// hash is not [Valid]. This is not an error condition. To ensure that a hyper-hash is
// constructed, call [Matrix.Wait] with [Materialize].
//
// On successful return, the hyper-hash is removed from the matrix, and it is the responsibility
// of the user application to free it.
//
// UnpackHyperHash is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) UnpackHyperHash(desc *Descriptor) (hash Matrix[int], err error) {
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_unpack_HyperHash(matrix.grb, &hash.grb, cdesc))
	if info != success {
		err = makeError(info)
	}
	return
}

// PackHyperHash assigns the input matrix as the hyper-hash of the matrix. The hyper-hash
// of a hypersparse matrix provides quick access to the inverse of ah.
//
// If the matrix is not hypersparse, or if it already has a hyper-hash, then nothing happens,
// and the input matrix hash is unchanged. This is not an error condition. The input matrix
// hash is still owned by the user application and freeing it is the responsibility of the user
// application.
//
// If the input matrix hash is moved into the matrix as its hyper-hash, then hash becomes
// not [Matrix.Valid] to indicate that it has been moved into the matrix. It is no longer
// owned by the caller.
//
// Results are undefined if the input matrix hash was not created by [Matrix.UnpackHyperHash],
// or if the matrix was modified after it was unpacked by [Matrix.UnpackHyperCSR] or
// [Matrix.UnpackHyperCSC].
//
// PackHyperHash is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) PackHyperHash(hash *Matrix[int], desc *Descriptor) error {
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_pack_HyperHash(matrix.grb, &hash.grb, cdesc))
	if info != success {
		return makeError(info)
	}
	return nil
}

// Matrix BitmapR

// PackBitmapR packs a matrix from two user slices in bitmap format.
//
// The matrix must exist on input with the right type and dimensions. No type casting is done,
// so the domain D must correctly reflect [Matrix.Type]().
//
// Parameters:
//
//   - ab (INOUT): A [SystemSlice] that indicates which indices are present: If ab is true
//     at i*ncols + j, then the entry at a(i, j) is present with value given by ax at
//     i*ncols + j. If ab is false at a given index, the entry at that index is not present,
//     and the value given by ax at that index is ignored.
//
//   - ax (INOUT): A [SystemSlice] of values.
//
//   - iso (IN): If true, the resulting matrix is iso.
//
//   - nvals (IN): The number of entries in the resulting vector.
//
// On successful return, ab and ax are empty, to indicate that the user application no longer
// owns them. They have instead been moved to the resulting matrix. If not successful, ab and ax
// are not modified.
//
// PackBitmapR is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) PackBitmapR(ab *SystemSlice[bool], ax *SystemSlice[D], iso bool, nvals int, desc *Descriptor) error {
	if err := checkType[D](matrix.Type()); err != nil {
		return err
	}
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Matrix_pack_BitmapR(
		matrix.grb,
		(**C.int8_t)(unsafe.Pointer(&ab.ptr)), &ax.ptr,
		C.GrB_Index(ab.size), C.GrB_Index(ax.size),
		C.bool(iso), C.GrB_Index(nvals), cdesc,
	))
	if info != success {
		return makeError(info)
	}
	ab.size = 0
	ax.size = 0
	return nil
}

// UnpackBitmapR unpacks a matrix to user slices in bitmap format.
//
// No type casting is done, so the domain D must correctly reflect [Matrix.Type]().
//
// Return Values:
//
//   - ab: A [SystemSlice] that indicates which indices are present: If ab is true
//     at i*ncols + j, then the entry at a(i, j) is present with value given by vx at
//     i*ncols + j. If ab is false at a given index, the entry at that index is not present,
//     and the value given by ax at that index is ignored.
//
//   - ax: A [SystemSlice] of values.
//
//   - iso: If true, the input vector was iso.
//
//   - nvals: The number of entries in the resulting slices.
//
// On successful return, the input matrix has no entries anymore, and the user application now
// owns the resulting slices. If not successful, the input matrix is not modified.
//
// UnpackBitmapR is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) UnpackBitmapR(desc *Descriptor) (ab SystemSlice[bool], ax SystemSlice[D], iso bool, nvals int, err error) {
	if err = checkType[D](matrix.Type()); err != nil {
		return
	}
	var cab *C.int8_t
	var cax unsafe.Pointer
	var cabSize, caxSize, cnvals C.GrB_Index
	var ciso C.bool
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Matrix_unpack_BitmapR(matrix.grb, &cab, &cax, &cabSize, &caxSize, &ciso, &cnvals, cdesc))
	if info != success {
		err = makeError(info)
		return
	}
	ab = AsSystemSlice[bool](unsafe.Pointer(cab), int(cabSize))
	ax = AsSystemSlice[D](cax, int(caxSize))
	iso = bool(ciso)
	nvals = int(cnvals)
	return
}

// Matrix BitmapC

// PackBitmapC packs a matrix from two user slices in bitmap format.
//
// The matrix must exist on input with the right type and dimensions. No type casting is done,
// so the domain D must correctly reflect [Matrix.Type]().
//
// Parameters:
//
//   - ab (INOUT): A [SystemSlice] that indicates which indices are present: If ab is true
//     at i + j*nrows, then the entry at a(i, j) is present with value given by ax at
//     i + j*nrows. If ab is false at a given index, the entry at that index is not present,
//     and the value given by ax at that index is ignored.
//
//   - ax (INOUT): A [SystemSlice] of values.
//
//   - iso (IN): If true, the resulting matrix is iso.
//
//   - nvals (IN): The number of entries in the resulting vector.
//
// On successful return, ab and ax are empty, to indicate that the user application no longer
// owns them. They have instead been moved to the resulting matrix. If not successful, ab and ax
// are not modified.
//
// PackBitmapC is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) PackBitmapC(ab *SystemSlice[bool], ax *SystemSlice[D], iso bool, nvals int, desc *Descriptor) error {
	if err := checkType[D](matrix.Type()); err != nil {
		return err
	}
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Matrix_pack_BitmapC(
		matrix.grb,
		(**C.int8_t)(unsafe.Pointer(&ab.ptr)), &ax.ptr,
		C.GrB_Index(ab.size), C.GrB_Index(ax.size),
		C.bool(iso), C.GrB_Index(nvals), cdesc,
	))
	if info != success {
		return makeError(info)
	}
	ab.size = 0
	ax.size = 0
	return nil
}

// UnpackBitmapC unpacks a matrix to user slices in bitmap format.
//
// No type casting is done, so the domain D must correctly reflect [Matrix.Type]().
//
// Return Values:
//
//   - ab: A [SystemSlice] that indicates which indices are present: If ab is true
//     at i + j*nrows, then the entry at a(i, j) is present with value given by vx at
//     i + j*nrows. If ab is false at a given index, the entry at that index is not present,
//     and the value given by ax at that index is ignored.
//
//   - ax: A [SystemSlice] of values.
//
//   - iso: If true, the input vector was iso.
//
//   - nvals: The number of entries in the resulting slices.
//
// On successful return, the input matrix has no entries anymore, and the user application now
// owns the resulting slices. If not successful, the input matrix is not modified.
//
// UnpackBitmapC is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) UnpackBitmapC(desc *Descriptor) (ab SystemSlice[bool], ax SystemSlice[D], iso bool, nvals int, err error) {
	if err = checkType[D](matrix.Type()); err != nil {
		return
	}
	var cab *C.int8_t
	var cax unsafe.Pointer
	var cabSize, caxSize, cnvals C.GrB_Index
	var ciso C.bool
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Matrix_unpack_BitmapC(matrix.grb, &cab, &cax, &cabSize, &caxSize, &ciso, &cnvals, cdesc))
	if info != success {
		err = makeError(info)
		return
	}
	ab = AsSystemSlice[bool](unsafe.Pointer(cab), int(cabSize))
	ax = AsSystemSlice[D](cax, int(caxSize))
	iso = bool(ciso)
	nvals = int(cnvals)
	return
}

// Matrix FullR

// PackFullR packs a matrix from a user slice in full format.
//
// The matrix must exist on input with the right type and size. No type casting is done,
// so the domain D must correctly reflect [Matrix.Type]().
//
// Parameters:
//
//   - ax (INOUT): A [SystemSlice] of values. All entries with index < nrows*ncols are present.
//     Values at index i*ncols + j correspond to matrix entries a(i, j).
//
//   - iso (IN): If true, the resulting matrix is iso.
//
// On successful return, ax is empty, to indicate that the user application no longer
// owns it. It has instead been moved to the resulting matrix. If not successful, ax
// is not modified.
//
// PackFullR is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) PackFullR(ax *SystemSlice[D], iso bool, desc *Descriptor) error {
	if err := checkType[D](matrix.Type()); err != nil {
		return err
	}
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Matrix_pack_FullR(matrix.grb, &ax.ptr, C.GrB_Index(ax.size), C.bool(iso), cdesc))
	if info != success {
		return makeError(info)
	}
	ax.size = 0
	return nil
}

// UnpackFullR unpacks a matrix to a user slice in full format.
//
// No type casting is done, so the domain D must correctly reflect [Matrix.Type]().
//
// Return Values:
//
//   - ax: A [SystemSlice] of values. All entries with index < nrows*ncols are present.
//     Values at index i*ncols + j correspond to matrix entries a(i, j).
//
//   - iso: If true, the input matrix was iso.
//
// On successful return, the input matrix has no entries anymore, and the user application now
// owns the resulting slice. If not successful, the input matrix is not modified.
//
// UnpackFullR is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) UnpackFullR(desc *Descriptor) (ax SystemSlice[D], iso bool, err error) {
	if err = checkType[D](matrix.Type()); err != nil {
		return
	}
	var cax unsafe.Pointer
	var caxSize C.GrB_Index
	var ciso C.bool
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Matrix_unpack_FullR(matrix.grb, &cax, &caxSize, &ciso, cdesc))
	if info != success {
		err = makeError(info)
		return
	}
	ax = AsSystemSlice[D](cax, int(caxSize))
	iso = bool(ciso)
	return
}

// Matrix FullC

// PackFullC packs a matrix from a user slice in full format.
//
// The matrix must exist on input with the right type and size. No type casting is done,
// so the domain D must correctly reflect [Matrix.Type]().
//
// Parameters:
//
//   - ax (INOUT): A [SystemSlice] of values. All entries with index < nrows*ncols are present.
//     Values at index i + j*nrows correspond to matrix entries a(i, j).
//
//   - iso (IN): If true, the resulting matrix is iso.
//
// On successful return, ax is empty, to indicate that the user application no longer
// owns it. It has instead been moved to the resulting matrix. If not successful, ax
// is not modified.
//
// PackFullC is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) PackFullC(ax *SystemSlice[D], iso bool, desc *Descriptor) error {
	if err := checkType[D](matrix.Type()); err != nil {
		return err
	}
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Matrix_pack_FullC(matrix.grb, &ax.ptr, C.GrB_Index(ax.size), C.bool(iso), cdesc))
	if info != success {
		return makeError(info)
	}
	ax.size = 0
	return nil
}

// UnpackFullC unpacks a matrix to a user slice in full format.
//
// No type casting is done, so the domain D must correctly reflect [Matrix.Type]().
//
// Return Values:
//
//   - ax: A [SystemSlice] of values. All entries with index < nrows*ncols are present.
//     Values at index i + j*nrows correspond to matrix entries a(i, j).
//
//   - iso: If true, the input matrix was iso.
//
// On successful return, the input matrix has no entries anymore, and the user application now
// owns the resulting slice. If not successful, the input matrix is not modified.
//
// UnpackFullC is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) UnpackFullC(desc *Descriptor) (ax SystemSlice[D], iso bool, err error) {
	if err = checkType[D](matrix.Type()); err != nil {
		return
	}
	var cax unsafe.Pointer
	var caxSize C.GrB_Index
	var ciso C.bool
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Matrix_unpack_FullC(matrix.grb, &cax, &caxSize, &ciso, cdesc))
	if info != success {
		err = makeError(info)
		return
	}
	ax = AsSystemSlice[D](cax, int(caxSize))
	iso = bool(ciso)
	return
}
