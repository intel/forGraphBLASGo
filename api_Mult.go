package forGraphBLASGo

func MxM[DC, DA, DB any](C *Matrix[DC], mask *Matrix[bool], accum BinaryOp[DC, DC, DC], op Semiring[DC, DA, DB], A *Matrix[DA], B *Matrix[DB], desc Descriptor) error {
	nrows, ncols, err := C.Size()
	if err != nil {
		return err
	}
	Anrows, Ancols, err := A.Size()
	if err != nil {
		return err
	}
	Bnrows, Bncols, err := B.Size()
	if err != nil {
		return err
	}
	AIsTran, err := desc.Is(Inp0, Tran)
	if err != nil {
		panic(err)
	}
	if AIsTran {
		Anrows, Ancols = Ancols, Anrows
	}
	BIsTran, err := desc.Is(Inp1, Tran)
	if err != nil {
		panic(err)
	}
	if BIsTran {
		Bnrows, Bncols = Bncols, Bnrows
	}
	if nrows != Anrows || ncols != Bncols || Ancols != Bnrows {
		return DimensionMismatch
	}
	maskAsStructure, err := matrixMask(mask, nrows, ncols)
	if err != nil {
		return err
	}
	C.ref = newMatrixReference[DC](newComputedMatrix[DC](
		nrows, ncols, C.ref, maskAsStructure, accum,
		newMatrixMult[DC, DA, DB](op, maybeTran(A.ref, AIsTran), maybeTran(B.ref, BIsTran)),
		desc,
	), -1)
	return nil
}

func VxM[Dw, Du, DA any](w *Vector[Dw], mask *Vector[bool], accum BinaryOp[Dw, Dw, Dw], op Semiring[Dw, Du, DA], u *Vector[Du], A *Matrix[DA], desc Descriptor) error {
	wsize, err := w.Size()
	if err != nil {
		return err
	}
	usize, err := u.Size()
	if err != nil {
		return err
	}
	AIsTran, err := A.expectSizeTran(usize, wsize, desc, Inp1)
	if err != nil {
		return err
	}
	maskAsStructure, err := vectorMask(mask, wsize)
	if err != nil {
		return err
	}
	w.ref = newVectorReference[Dw](newComputedVector[Dw](
		wsize, w.ref, maskAsStructure, accum,
		newVxM[Dw](op, u.ref, maybeTran(A.ref, AIsTran)),
		desc,
	), -1)
	return nil
}

func MxV[Dw, DA, Du any](w *Vector[Dw], mask *Vector[bool], accum BinaryOp[Dw, Dw, Dw], op Semiring[Dw, DA, Du], A *Matrix[DA], u *Vector[Du], desc Descriptor) error {
	wsize, err := w.Size()
	if err != nil {
		return err
	}
	usize, err := u.Size()
	if err != nil {
		return err
	}
	AIsTran, err := A.expectSizeTran(wsize, usize, desc, Inp0)
	if err != nil {
		return err
	}
	maskAsStructure, err := vectorMask(mask, wsize)
	if err != nil {
		return err
	}
	w.ref = newVectorReference[Dw](newComputedVector[Dw](
		wsize, w.ref, maskAsStructure, accum,
		newMxV[Dw](op, maybeTran(A.ref, AIsTran), u.ref),
		desc,
	), -1)
	return nil
}

func KroneckerBinaryOp[DC, DA, DB any](C *Matrix[DC], mask *Matrix[bool], accum BinaryOp[DC, DC, DC], op BinaryOp[DC, DA, DB], A *Matrix[DA], B *Matrix[DB], desc Descriptor) error {
	Anrows, Ancols, err := A.Size()
	if err != nil {
		return err
	}
	Bnrows, Bncols, err := B.Size()
	if err != nil {
		return err
	}
	AIsTran, err := desc.Is(Inp0, Tran)
	if err != nil {
		panic(err)
	}
	if AIsTran {
		Anrows, Ancols = Ancols, Anrows
	}
	BIsTran, err := desc.Is(Inp1, Tran)
	if err != nil {
		panic(err)
	}
	if BIsTran {
		Bnrows, Bncols = Bncols, Bnrows
	}
	nrows, ncols := Anrows*Bnrows, Ancols*Bncols
	if err = C.expectSize(nrows, ncols); err != nil {
		return err
	}
	maskAsStructure, err := matrixMask(mask, nrows, ncols)
	if err != nil {
		return err
	}
	C.ref = newMatrixReference[DC](newComputedMatrix[DC](
		nrows, ncols, C.ref, maskAsStructure,
		accum, newKroneckerBinaryOp[DC, DA, DB](op, maybeTran(A.ref, AIsTran), maybeTran(B.ref, BIsTran)),
		desc,
	), -1)
	return nil
}

func KroneckerMonoid[D any](C *Matrix[D], mask *Matrix[bool], accum BinaryOp[D, D, D], op Monoid[D], A, B *Matrix[D], desc Descriptor) error {
	return KroneckerBinaryOp(C, mask, accum, op.operator(), A, B, desc)
}

func KroneckerSemiring[DC, DA, DB any](C *Matrix[DC], mask *Matrix[bool], accum BinaryOp[DC, DC, DC], op Semiring[DC, DA, DB], A *Matrix[DA], B *Matrix[DB], desc Descriptor) error {
	return KroneckerBinaryOp(C, mask, accum, op.multiplication(), A, B, desc)
}
