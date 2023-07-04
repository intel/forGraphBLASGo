package GrB

// #include "GraphBLAS.h"
import "C"
import "unsafe"

// Vector CSC

// PackCSCBytes is like [Vector.PackCSC], except that vi and vx are passed as byte slices.
//
// PackCSCBytes is a SuiteSparse:GraphBLAS extension.
func (vector Vector[D]) PackCSCBytes(vi, vx *SystemSlice[byte], iso bool, nvals int, jumbled bool, desc *Descriptor) error {
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Vector_pack_CSC(
		vector.grb,
		(**C.GrB_Index)(unsafe.Pointer(&vi.ptr)), &vx.ptr,
		C.GrB_Index(vi.size), C.GrB_Index(vx.size),
		C.bool(iso), C.GrB_Index(nvals), C.bool(jumbled), cdesc,
	))
	if info != success {
		return makeError(info)
	}
	vi.size = 0
	vx.size = 0
	return nil
}

// UnpackCSCBytes is like [Vector.UnpackCSC], except that vi and vx are returned as byte slices.
//
// UnpackCSCBytes is a SuiteSparse:GraphBLAS extension.
func (vector Vector[D]) UnpackCSCBytes(allowJumbled bool, desc *Descriptor) (vi, vx SystemSlice[byte], iso bool, nvals int, jumbled bool, err error) {
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
	vi = AsSystemSlice[byte](unsafe.Pointer(cvi), int(cviSize))
	vx = AsSystemSlice[byte](cvx, int(cvxSize))
	iso = bool(ciso)
	nvals = int(cnvals)
	jumbled = bool(cjumbled)
	return
}

// Vector Bitmap

// PackBitmapBytes is like [Vector.PackBitmap], except that vb and vx are passed as byte slices.
//
// PackBitmapBytes is a SuiteSparse:GraphBLAS extension.
func (vector Vector[D]) PackBitmapBytes(vb, vx *SystemSlice[byte], iso bool, nvals int, desc *Descriptor) error {
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

// UnpackBitmapBytes is like [Vector.UnpackBitmap], except that vb and vx are returned as byte slices.
//
// UnpackBitmapBytes is a SuiteSparse:GraphBLAS extension.
func (vector Vector[D]) UnpackBitmapBytes(desc *Descriptor) (vb, vx SystemSlice[byte], iso bool, nvals int, err error) {
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
	vb = AsSystemSlice[byte](unsafe.Pointer(cvb), int(cvbSize))
	vx = AsSystemSlice[byte](cvx, int(cvxSize))
	iso = bool(ciso)
	nvals = int(cnvals)
	return
}

// Vector Full

// PackFullBytes is like [Vector.PackFull], except that vx is passed as a byte slice.
//
// PackFullBytes is a SuiteSparse:GraphBLAS extension.
func (vector Vector[D]) PackFullBytes(vx *SystemSlice[byte], iso bool, desc *Descriptor) error {
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Vector_pack_Full(
		vector.grb,
		&vx.ptr,
		C.GrB_Index(vx.size),
		C.bool(iso), cdesc,
	))
	if info != success {
		return makeError(info)
	}
	vx.size = 0
	return nil
}

// UnpackFullBytes is like [Vector.UnpackFull], except that vx is returned as a byte slice.
//
// UnpackFullBytes is a SuiteSparse:GraphBLAS extension.
func (vector Vector[D]) UnpackFullBytes(desc *Descriptor) (vx SystemSlice[byte], iso bool, err error) {
	var cvx unsafe.Pointer
	var cvxSize C.GrB_Index
	var ciso C.bool
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Vector_unpack_Full(vector.grb, &cvx, &cvxSize, &ciso, cdesc))
	if info != success {
		err = makeError(info)
		return
	}
	vx = AsSystemSlice[byte](cvx, int(cvxSize))
	iso = bool(ciso)
	return
}

// Matrix CSR

// PackCSRBytes is like [Matrix.PackCSRBytes], except that ap, aj and ax are passed as byte slices.
//
// PackCSRBytes is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) PackCSRBytes(ap, aj, ax *SystemSlice[byte], iso, jumbled bool, desc *Descriptor) error {
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Matrix_pack_CSR(
		matrix.grb,
		(**C.GrB_Index)(unsafe.Pointer(&ap.ptr)), (**C.GrB_Index)(unsafe.Pointer(&aj.ptr)), &ax.ptr,
		C.GrB_Index(ap.size), C.GrB_Index(aj.size), C.GrB_Index(ax.size),
		C.bool(iso), C.bool(jumbled), cdesc,
	))
	if info != success {
		return makeError(info)
	}
	ap.size = 0
	aj.size = 0
	ax.size = 0
	return nil
}

// UnpackCSRBytes is like [Matrix.UnpackCSR], except that ap, aj and ax are returned as byte slices.
//
// UnpackCSRBytes is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) UnpackCSRBytes(allowJumbled bool, desc *Descriptor) (ap, aj, ax SystemSlice[byte], iso, jumbled bool, err error) {
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
	ap = AsSystemSlice[byte](unsafe.Pointer(cap), int(capSize))
	aj = AsSystemSlice[byte](unsafe.Pointer(caj), int(cajSize))
	ax = AsSystemSlice[byte](cax, int(caxSize))
	iso = bool(ciso)
	jumbled = bool(cjumbled)
	return
}

// Matrix CSC

// PackCSCBytes is like [Matrix.PackCSC], except that ap, ai and ax are passed as byte slices.
//
// PackCSCBytes is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) PackCSCBytes(ap, ai, ax *SystemSlice[byte], iso, jumbled bool, desc *Descriptor) error {
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Matrix_pack_CSC(
		matrix.grb,
		(**C.GrB_Index)(unsafe.Pointer(&ap.ptr)), (**C.GrB_Index)(unsafe.Pointer(&ai.ptr)), &ax.ptr,
		C.GrB_Index(ap.size), C.GrB_Index(ai.size), C.GrB_Index(ax.size),
		C.bool(iso), C.bool(jumbled), cdesc,
	))
	if info != success {
		return makeError(info)
	}
	ap.size = 0
	ai.size = 0
	ax.size = 0
	return nil
}

// UnpackCSCBytes is like [Matrix.UnpackCSC], except that ap, ai and ax are returned as byte slices.
//
// UnpackCSCBytes is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) UnpackCSCBytes(allowJumbled bool, desc *Descriptor) (ap, ai, ax SystemSlice[byte], iso, jumbled bool, err error) {
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
	ap = AsSystemSlice[byte](unsafe.Pointer(cap), int(capSize))
	ai = AsSystemSlice[byte](unsafe.Pointer(cai), int(caiSize))
	ax = AsSystemSlice[byte](cax, int(caxSize))
	iso = bool(ciso)
	jumbled = bool(cjumbled)
	return
}

// Matrix HyperCSR

// PackHyperCSRBytes is like [Matrix.PackHyperCSR], except that ap, ah, aj and ax are passed as byte slices.
//
// PackHyperCSRBytes is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) PackHyperCSRBytes(ap, ah, aj, ax *SystemSlice[byte], iso bool, nvec int, jumbled bool, desc *Descriptor) error {
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Matrix_pack_HyperCSR(
		matrix.grb,
		(**C.GrB_Index)(unsafe.Pointer(&ap.ptr)), (**C.GrB_Index)(unsafe.Pointer(&ah.ptr)),
		(**C.GrB_Index)(unsafe.Pointer(&aj.ptr)), &ax.ptr,
		C.GrB_Index(ap.size), C.GrB_Index(ah.size), C.GrB_Index(aj.size), C.GrB_Index(ax.size),
		C.bool(iso), C.GrB_Index(nvec), C.bool(jumbled), cdesc,
	))
	if info != success {
		return makeError(info)
	}
	ap.size = 0
	ah.size = 0
	aj.size = 0
	ax.size = 0
	return nil
}

// UnpackHyperCSRBytes is like [Matrix.UnpackHyperCSR], except that ap, ah, aj and ax are returned as byte slices.
//
// UnpackHyperCSRBytes is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) UnpackHyperCSRBytes(allowJumbled bool, desc *Descriptor) (ap, ah, aj, ax SystemSlice[byte], iso bool, nvec int, jumbled bool, err error) {
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
	ap = AsSystemSlice[byte](unsafe.Pointer(cap), int(capSize))
	ah = AsSystemSlice[byte](unsafe.Pointer(cah), int(cahSize))
	aj = AsSystemSlice[byte](unsafe.Pointer(caj), int(cajSize))
	ax = AsSystemSlice[byte](cax, int(caxSize))
	iso = bool(ciso)
	nvec = int(cnvec)
	jumbled = bool(cjumbled)
	return
}

// Matrix HyperCSC

// PackHyperCSCBytes is like [Matrix.PackHyperCSC], except that ap, ah, ai and ax are passed as byte slices.
//
// PackHyperCSCBytes is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) PackHyperCSCBytes(ap, ah, ai, ax *SystemSlice[byte], iso bool, nvec int, jumbled bool, desc *Descriptor) error {
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Matrix_pack_HyperCSC(
		matrix.grb,
		(**C.GrB_Index)(unsafe.Pointer(&ap.ptr)), (**C.GrB_Index)(unsafe.Pointer(&ah.ptr)),
		(**C.GrB_Index)(unsafe.Pointer(&ai.ptr)), &ax.ptr,
		C.GrB_Index(ap.size), C.GrB_Index(ah.size), C.GrB_Index(ai.size), C.GrB_Index(ax.size),
		C.bool(iso), C.GrB_Index(nvec), C.bool(jumbled), cdesc,
	))
	if info != success {
		return makeError(info)
	}
	ap.size = 0
	ah.size = 0
	ai.size = 0
	ax.size = 0
	return nil
}

// UnpackHyperCSCBytes is like [Matrix.UnpackHyperCSC], except that ap, ah, ai and ax are returned as byte slices.
//
// UnpackHyperCSCBytes is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) UnpackHyperCSCBytes(allowJumbled bool, desc *Descriptor) (ap, ah, ai, ax SystemSlice[byte], iso bool, nvec int, jumbled bool, err error) {
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
	ap = AsSystemSlice[byte](unsafe.Pointer(cap), int(capSize))
	ah = AsSystemSlice[byte](unsafe.Pointer(cah), int(cahSize))
	ai = AsSystemSlice[byte](unsafe.Pointer(cai), int(caiSize))
	ax = AsSystemSlice[byte](cax, int(caxSize))
	iso = bool(ciso)
	nvec = int(cnvec)
	jumbled = bool(cjumbled)
	return
}

// Matrix BitmapR

// PackBitmapRBytes is like [Matrix.PackBitmapR], except that ab and ax are passed as byte slices.
//
// PackBitmapRBytes is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) PackBitmapRBytes(ab, ax *SystemSlice[byte], iso bool, nvals int, desc *Descriptor) error {
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

// UnpackBitmapRBytes is like [Matrix.UnpackBitmapR], except that ab and ax are returned as byte slices.
//
// UnpackBitmapRBytes is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) UnpackBitmapRBytes(desc *Descriptor) (ab, ax SystemSlice[byte], iso bool, nvals int, err error) {
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
	ab = AsSystemSlice[byte](unsafe.Pointer(cab), int(cabSize))
	ax = AsSystemSlice[byte](cax, int(caxSize))
	iso = bool(ciso)
	nvals = int(cnvals)
	return
}

// Matrix BitmapC

// PackBitmapCBytes is like [Matrix.PackBitmapC], except that ab and ax are passed as byte slices.
//
// PackBitmapCBytes is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) PackBitmapCBytes(ab, ax *SystemSlice[byte], iso bool, nvals int, desc *Descriptor) error {
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

// UnpackBitmapCBytes is like [Matrix.UnpackBitmapC], except that ab and ax are returned as bytes slices.
//
// UnpackBitmapCBytes is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) UnpackBitmapCBytes(desc *Descriptor) (ab, ax SystemSlice[byte], iso bool, nvals int, err error) {
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
	ab = AsSystemSlice[byte](unsafe.Pointer(cab), int(cabSize))
	ax = AsSystemSlice[byte](cax, int(caxSize))
	iso = bool(ciso)
	nvals = int(cnvals)
	return
}

// Matrix FullR

// PackFullRBytes is like [Matrix.PackFullR], except that ax is passed as a byte slice.
//
// PackFullRBytes is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) PackFullRBytes(ax *SystemSlice[byte], iso bool, desc *Descriptor) error {
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Matrix_pack_FullR(matrix.grb, &ax.ptr, C.GrB_Index(ax.size), C.bool(iso), cdesc))
	if info != success {
		return makeError(info)
	}
	ax.size = 0
	return nil
}

// UnpackFullRBytes is like [Matrix.UnpackFullRBytes], except that ax is return as a byte slice.
//
// UnpackFullRBytes is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) UnpackFullRBytes(desc *Descriptor) (ax SystemSlice[byte], iso bool, err error) {
	var cax unsafe.Pointer
	var caxSize C.GrB_Index
	var ciso C.bool
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Matrix_unpack_FullR(matrix.grb, &cax, &caxSize, &ciso, cdesc))
	if info != success {
		err = makeError(info)
		return
	}
	ax = AsSystemSlice[byte](cax, int(caxSize))
	iso = bool(ciso)
	return
}

// Matrix FullC

// PackFullCBytes is like [Matrix.PackFullC], except that ax is passed as a byte slice.
//
// PackFullCBytes is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) PackFullCBytes(ax *SystemSlice[byte], iso bool, desc *Descriptor) error {
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Matrix_pack_FullC(matrix.grb, &ax.ptr, C.GrB_Index(ax.size), C.bool(iso), cdesc))
	if info != success {
		return makeError(info)
	}
	ax.size = 0
	return nil
}

// UnpackFullCBytes is like [Matrix.UnpackFullC], except that ax is returned as a byte slice.
//
// UnpackFullCBytes is a SuiteSparse:GraphBLAS extension.
func (matrix Matrix[D]) UnpackFullCBytes(desc *Descriptor) (ax SystemSlice[byte], iso bool, err error) {
	var cax unsafe.Pointer
	var caxSize C.GrB_Index
	var ciso C.bool
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Matrix_unpack_FullC(matrix.grb, &cax, &caxSize, &ciso, cdesc))
	if info != success {
		err = makeError(info)
		return
	}
	ax = AsSystemSlice[byte](cax, int(caxSize))
	iso = bool(ciso)
	return
}
