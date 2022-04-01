package forGraphBLASGo

func VectorSelect[D, Ds any](w *Vector[D], mask *Vector[bool], accum BinaryOp[D, D, D], op IndexUnaryOp[bool, D, Ds], u *Vector[D], value Ds, desc Descriptor) error {
	size, err := w.Size()
	if err != nil {
		return err
	}
	if err = u.expectSize(size); err != nil {
		return err
	}
	maskAsStructure, err := vectorMask(mask, size)
	if err != nil {
		return err
	}
	w.ref = newVectorReference[D](newComputedVector[D](
		size, w.ref, maskAsStructure, accum,
		newVectorSelect[D, Ds](op, u.ref, value),
		desc,
	), -1)
	return nil
}

func VectorSelectScalar[D, Ds any](w *Vector[D], mask *Vector[bool], accum BinaryOp[D, D, D], op IndexUnaryOp[bool, D, Ds], u *Vector[D], value *Scalar[Ds], desc Descriptor) error {
	size, err := w.Size()
	if err != nil {
		return err
	}
	if err = u.expectSize(size); err != nil {
		return err
	}
	maskAsStructure, err := vectorMask(mask, size)
	if err != nil {
		return err
	}
	w.ref = newVectorReference[D](newComputedVector[D](
		size, w.ref, maskAsStructure, accum,
		newVectorSelectScalar[D, Ds](op, u.ref, value.ref),
		desc,
	), -1)
	return nil
}

func MatrixSelect[D, Ds any](C *Matrix[D], mask *Matrix[bool], accum BinaryOp[D, D, D], op IndexUnaryOp[bool, D, Ds], A *Matrix[D], value Ds, desc Descriptor) error {
	nrows, ncols, err := C.Size()
	if err != nil {
		return err
	}
	isTran, err := A.expectSizeTran(nrows, ncols, desc, Inp0)
	if err != nil {
		return err
	}
	maskAsStructure, err := matrixMask(mask, nrows, ncols)
	if err != nil {
		return err
	}
	C.ref = newMatrixReference[D](newComputedMatrix[D](
		nrows, ncols, C.ref, maskAsStructure, accum,
		newMatrixSelect[D, Ds](op, maybeTran(A.ref, isTran), value),
		desc,
	), -1)
	return nil
}

func MatrixSelectScalar[D, Ds any](C *Matrix[D], mask *Matrix[bool], accum BinaryOp[D, D, D], op IndexUnaryOp[bool, D, Ds], A *Matrix[D], value *Scalar[Ds], desc Descriptor) error {
	nrows, ncols, err := C.Size()
	if err != nil {
		return err
	}
	isTran, err := A.expectSizeTran(nrows, ncols, desc, Inp0)
	if err != nil {
		return err
	}
	maskAsStructure, err := matrixMask(mask, nrows, ncols)
	if err != nil {
		return err
	}
	C.ref = newMatrixReference[D](newComputedMatrix[D](
		nrows, ncols, C.ref, maskAsStructure, accum,
		newMatrixSelectScalar[D, Ds](op, maybeTran(A.ref, isTran), value.ref),
		desc,
	), -1)
	return nil
}
