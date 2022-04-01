package forGraphBLASGo

func VectorEWiseMultBinaryOp[Dw, Du, Dv any](w *Vector[Dw], mask *Vector[bool], accum BinaryOp[Dw, Dw, Dw], op BinaryOp[Dw, Du, Dv], u *Vector[Du], v *Vector[Dv], desc Descriptor) error {
	size, err := w.Size()
	if err != nil {
		return err
	}
	if err = u.expectSize(size); err != nil {
		return err
	}
	if err = v.expectSize(size); err != nil {
		return err
	}
	maskAsStructure, err := vectorMask(mask, size)
	if err != nil {
		return err
	}
	w.ref = newVectorReference[Dw](newComputedVector[Dw](
		size, w.ref, maskAsStructure, accum,
		newVectorEWiseMultBinaryOp[Dw, Du, Dv](op, u.ref, v.ref),
		desc,
	), -1)
	return nil
}

func VectorEWiseMultMonoid[D any](w *Vector[D], mask *Vector[bool], accum BinaryOp[D, D, D], op Monoid[D], u, v *Vector[D], desc Descriptor) error {
	return VectorEWiseMultBinaryOp(w, mask, accum, op.operator(), u, v, desc)
}

func VectorEWiseMultSemiring[Dw, Du, Dv any](w *Vector[Dw], mask *Vector[bool], accum BinaryOp[Dw, Dw, Dw], op Semiring[Dw, Du, Dv], u *Vector[Du], v *Vector[Dv], desc Descriptor) error {
	return VectorEWiseMultBinaryOp(w, mask, accum, op.multiplication(), u, v, desc)
}

func VectorEWiseAddBinaryOp[D any](w *Vector[D], mask *Vector[bool], accum, op BinaryOp[D, D, D], u, v *Vector[D], desc Descriptor) error {
	size, err := w.Size()
	if err != nil {
		return err
	}
	if err = u.expectSize(size); err != nil {
		return err
	}
	if err = v.expectSize(size); err != nil {
		return err
	}
	maskAsStructure, err := vectorMask(mask, size)
	if err != nil {
		return err
	}
	w.ref = newVectorReference[D](newComputedVector[D](
		size, w.ref, maskAsStructure, accum,
		newVectorEWiseAddBinaryOp[D](op, u.ref, v.ref),
		desc,
	), -1)
	return nil
}

func VectorEWiseAddMonoid[D any](w *Vector[D], mask *Vector[bool], accum BinaryOp[D, D, D], op Monoid[D], u, v *Vector[D], desc Descriptor) error {
	return VectorEWiseAddBinaryOp(w, mask, accum, op.operator(), u, v, desc)
}

func VectorEWiseAddSemiring[D, Din1, Din2 any](w *Vector[D], mask *Vector[bool], accum BinaryOp[D, D, D], op Semiring[D, Din1, Din2], u, v *Vector[D], desc Descriptor) error {
	return VectorEWiseAddBinaryOp(w, mask, accum, op.addition().operator(), u, v, desc)
}

func MatrixEWiseMultBinaryOp[DC, DA, DB any](C *Matrix[DC], mask *Matrix[bool], accum BinaryOp[DC, DC, DC], op BinaryOp[DC, DA, DB], A *Matrix[DA], B *Matrix[DB], desc Descriptor) error {
	nrows, ncols, err := C.Size()
	if err != nil {
		return err
	}
	AIsTran, err := A.expectSizeTran(nrows, ncols, desc, Inp0)
	if err != nil {
		return err
	}
	BIsTran, err := B.expectSizeTran(nrows, ncols, desc, Inp1)
	if err != nil {
		return err
	}
	maskAsStructure, err := matrixMask(mask, nrows, ncols)
	if err != nil {
		return err
	}
	C.ref = newMatrixReference[DC](newComputedMatrix[DC](
		nrows, ncols, C.ref, maskAsStructure, accum,
		newMatrixEWiseMultBinaryOp[DC, DA, DB](op, maybeTran(A.ref, AIsTran), maybeTran(B.ref, BIsTran)),
		desc,
	), -1)
	return nil
}

func MatrixEWiseMultMonoid[D any](C *Matrix[D], mask *Matrix[bool], accum BinaryOp[D, D, D], op Monoid[D], A, B *Matrix[D], desc Descriptor) error {
	return MatrixEWiseMultBinaryOp(C, mask, accum, op.operator(), A, B, desc)
}

func MatrixEWiseMultSemiring[DC, DA, DB any](C *Matrix[DC], mask *Matrix[bool], accum BinaryOp[DC, DC, DC], op Semiring[DC, DA, DB], A *Matrix[DA], B *Matrix[DB], desc Descriptor) error {
	return MatrixEWiseMultBinaryOp(C, mask, accum, op.multiplication(), A, B, desc)
}

func MatrixEWiseAddBinaryOp[D any](C *Matrix[D], mask *Matrix[bool], accum, op BinaryOp[D, D, D], A, B *Matrix[D], desc Descriptor) error {
	nrows, ncols, err := C.Size()
	if err != nil {
		return err
	}
	AIsTran, err := A.expectSizeTran(nrows, ncols, desc, Inp0)
	if err != nil {
		return err
	}
	BIsTran, err := B.expectSizeTran(nrows, ncols, desc, Inp1)
	if err != nil {
		return err
	}
	maskAsStructure, err := matrixMask(mask, nrows, ncols)
	if err != nil {
		return err
	}
	C.ref = newMatrixReference[D](newComputedMatrix[D](
		nrows, ncols, C.ref, maskAsStructure, accum,
		newMatrixEWiseAddBinaryOp[D](op, maybeTran(A.ref, AIsTran), maybeTran(B.ref, BIsTran)),
		desc,
	), -1)
	return nil
}

func MatrixEWiseAddMonoid[D any](C *Matrix[D], mask *Matrix[bool], accum BinaryOp[D, D, D], op Monoid[D], A, B *Matrix[D], desc Descriptor) error {
	return MatrixEWiseAddBinaryOp(C, mask, accum, op.operator(), A, B, desc)
}

func MatrixEWiseAddSemiring[D, Din1, Din2 any](C *Matrix[D], mask *Matrix[bool], accum BinaryOp[D, D, D], op Semiring[D, Din1, Din2], A, B *Matrix[D], desc Descriptor) error {
	return MatrixEWiseAddBinaryOp(C, mask, accum, op.addition().operator(), A, B, desc)
}
