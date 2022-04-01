package forGraphBLASGo

func ScalarApply[Dt, Df any](t *Scalar[Dt], accum BinaryOp[Dt, Dt, Dt], op UnaryOp[Dt, Df], f *Scalar[Df], _ Descriptor) error {
	if t == nil || t.ref == nil || f == nil || f.ref == nil {
		return UninitializedObject
	}
	t.ref = newScalarReference[Dt](newComputedScalar[Dt](t.ref, accum, newScalarApply(op, f.ref)))
	return nil
}

func ScalarApplyBinaryOp1st[Dt, Ds, Df any](t *Scalar[Dt], accum BinaryOp[Dt, Dt, Dt], op BinaryOp[Dt, Ds, Df], value Ds, f *Scalar[Df], _ Descriptor) error {
	if t == nil || t.ref == nil || f == nil || f.ref == nil {
		return UninitializedObject
	}
	t.ref = newScalarReference[Dt](newComputedScalar[Dt](t.ref, accum, newScalarApplyBinaryOp1st(op, value, f.ref)))
	return nil
}

func ScalarApplyBinaryOp2nd[Dt, Df, Ds any](t *Scalar[Dt], accum BinaryOp[Dt, Dt, Dt], op BinaryOp[Dt, Df, Ds], f *Scalar[Df], value Ds, _ Descriptor) error {
	if t == nil || t.ref == nil || f == nil || f.ref == nil {
		return UninitializedObject
	}
	t.ref = newScalarReference[Dt](newComputedScalar[Dt](t.ref, accum, newScalarApplyBinaryOp2nd(op, f.ref, value)))
	return nil
}

func ScalarApplyBinary[Dt, Df1, Df2 any](t *Scalar[Dt], accum BinaryOp[Dt, Dt, Dt], op BinaryOp[Dt, Df1, Df2], f1 *Scalar[Df1], f2 *Scalar[Df2], _ Descriptor) error {
	if t == nil || t.ref == nil || f1 == nil || f1.ref == nil || f2 == nil || f2.ref == nil {
		return UninitializedObject
	}
	t.ref = newScalarReference[Dt](newComputedScalar[Dt](t.ref, accum, newScalarApplyBinary(op, f1.ref, f2.ref)))
	return nil
}

func VectorApply[Dw, Du any](w *Vector[Dw], mask *Vector[bool], accum BinaryOp[Dw, Dw, Dw], op UnaryOp[Dw, Du], u *Vector[Du], desc Descriptor) error {
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
	w.ref = newVectorReference[Dw](newComputedVector[Dw](
		size, w.ref, maskAsStructure, accum,
		newVectorApply(op, u.ref),
		desc,
	), -1)
	return nil
}

func MatrixApply[DC, DA any](C *Matrix[DC], mask *Matrix[bool], accum BinaryOp[DC, DC, DC], op UnaryOp[DC, DA], A *Matrix[DA], desc Descriptor) error {
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
	C.ref = newMatrixReference[DC](newComputedMatrix[DC](
		nrows, ncols, C.ref,
		maskAsStructure, accum,
		newMatrixApply(op, maybeTran(A.ref, isTran)),
		desc,
	), -1)
	return nil
}

func VectorApplyBinaryOp1st[Dw, Ds, Du any](w *Vector[Dw], mask *Vector[bool], accum BinaryOp[Dw, Dw, Dw], op BinaryOp[Dw, Ds, Du], value Ds, u *Vector[Du], desc Descriptor) error {
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
	w.ref = newVectorReference[Dw](newComputedVector[Dw](
		size, w.ref,
		maskAsStructure, accum,
		newVectorApplyBinaryOp1st(op, value, u.ref),
		desc,
	), -1)
	return nil
}

func VectorApplyBinaryOp1stScalar[Dw, Ds, Du any](w *Vector[Dw], mask *Vector[bool], accum BinaryOp[Dw, Dw, Dw], op BinaryOp[Dw, Ds, Du], value *Scalar[Ds], u *Vector[Du], desc Descriptor) error {
	size, err := w.Size()
	if err != nil {
		return err
	}
	if value == nil || value.ref == nil {
		return UninitializedObject
	}
	if err = u.expectSize(size); err != nil {
		return err
	}
	maskAsStructure, err := vectorMask(mask, size)
	if err != nil {
		return err
	}
	w.ref = newVectorReference[Dw](newComputedVector[Dw](
		size, w.ref,
		maskAsStructure, accum,
		newVectorApplyBinaryOp1stScalar(op, value.ref, u.ref),
		desc,
	), -1)
	return nil
}

func VectorApplyBinaryOp2nd[Dw, Du, Ds any](w *Vector[Dw], mask *Vector[bool], accum BinaryOp[Dw, Dw, Dw], op BinaryOp[Dw, Du, Ds], u *Vector[Du], value Ds, desc Descriptor) error {
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
	w.ref = newVectorReference[Dw](newComputedVector[Dw](
		size, w.ref,
		maskAsStructure, accum,
		newVectorApplyBinaryOp2nd(op, u.ref, value),
		desc,
	), -1)
	return nil
}

func VectorApplyBinaryOp2ndScalar[Dw, Du, Ds any](w *Vector[Dw], mask *Vector[bool], accum BinaryOp[Dw, Dw, Dw], op BinaryOp[Dw, Du, Ds], u *Vector[Du], value *Scalar[Ds], desc Descriptor) error {
	size, err := w.Size()
	if err != nil {
		return err
	}
	if err = u.expectSize(size); err != nil {
		return err
	}
	if value == nil || value.ref == nil {
		return UninitializedObject
	}
	maskAsStructure, err := vectorMask(mask, size)
	if err != nil {
		return err
	}
	w.ref = newVectorReference[Dw](newComputedVector[Dw](
		size, w.ref,
		maskAsStructure, accum,
		newVectorApplyBinaryOp2ndScalar(op, u.ref, value.ref),
		desc,
	), -1)
	return nil
}

func MatrixApplyBinaryOp1st[DC, Ds, DA any](C *Matrix[DC], mask *Matrix[bool], accum BinaryOp[DC, DC, DC], op BinaryOp[DC, Ds, DA], value Ds, A *Matrix[DA], desc Descriptor) error {
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
	C.ref = newMatrixReference[DC](newComputedMatrix[DC](
		nrows, ncols,
		C.ref, maskAsStructure, accum,
		newMatrixApplyBinaryOp1st(op, value, maybeTran(A.ref, isTran)),
		desc,
	), -1)
	return nil
}

func MatrixApplyBinaryOp1stScalar[DC, Ds, DA any](C *Matrix[DC], mask *Matrix[bool], accum BinaryOp[DC, DC, DC], op BinaryOp[DC, Ds, DA], value *Scalar[Ds], A *Matrix[DA], desc Descriptor) error {
	nrows, ncols, err := C.Size()
	if err != nil {
		return err
	}
	if value == nil || value.ref == nil {
		return UninitializedObject
	}
	isTran, err := A.expectSizeTran(nrows, ncols, desc, Inp0)
	if err != nil {
		return err
	}
	maskAsStructure, err := matrixMask(mask, nrows, ncols)
	if err != nil {
		return err
	}
	C.ref = newMatrixReference[DC](newComputedMatrix[DC](
		nrows, ncols,
		C.ref, maskAsStructure, accum,
		newMatrixApplyBinaryOp1stScalar(op, value.ref, maybeTran(A.ref, isTran)),
		desc,
	), -1)
	return nil
}

func MatrixApplyBinaryOp2nd[DC, DA, Ds any](C *Matrix[DC], mask *Matrix[bool], accum BinaryOp[DC, DC, DC], op BinaryOp[DC, DA, Ds], A *Matrix[DA], value Ds, desc Descriptor) error {
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
	C.ref = newMatrixReference[DC](newComputedMatrix[DC](
		nrows, ncols,
		C.ref, maskAsStructure, accum,
		newMatrixApplyBinaryOp2nd(op, maybeTran(A.ref, isTran), value),
		desc,
	), -1)
	return nil
}

func MatrixApplyBinaryOp2ndScalar[DC, DA, Ds any](C *Matrix[DC], mask *Matrix[bool], accum BinaryOp[DC, DC, DC], op BinaryOp[DC, DA, Ds], A *Matrix[DA], value *Scalar[Ds], desc Descriptor) error {
	nrows, ncols, err := C.Size()
	if err != nil {
		return err
	}
	isTran, err := A.expectSizeTran(nrows, ncols, desc, Inp0)
	if err != nil {
		return err
	}
	if value == nil || value.ref == nil {
		return UninitializedObject
	}
	maskAsStructure, err := matrixMask(mask, nrows, ncols)
	if err != nil {
		return err
	}
	C.ref = newMatrixReference[DC](newComputedMatrix[DC](
		nrows, ncols,
		C.ref, maskAsStructure, accum,
		newMatrixApplyBinaryOp2ndScalar(op, maybeTran(A.ref, isTran), value.ref),
		desc,
	), -1)
	return nil
}

func VectorApplyIndexOp[Dw, Du, Ds any](w *Vector[Dw], mask *Vector[bool], accum BinaryOp[Dw, Dw, Dw], op IndexUnaryOp[Dw, Du, Ds], u *Vector[Du], value Ds, desc Descriptor) error {
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
	w.ref = newVectorReference[Dw](newComputedVector[Dw](
		size, w.ref,
		maskAsStructure, accum,
		newVectorApplyIndexOp(op, u.ref, value),
		desc,
	), -1)
	return nil
}

func VectorApplyIndexOpScalar[Dw, Du, Ds any](w *Vector[Dw], mask *Vector[bool], accum BinaryOp[Dw, Dw, Dw], op IndexUnaryOp[Dw, Du, Ds], u *Vector[Du], value *Scalar[Ds], desc Descriptor) error {
	size, err := w.Size()
	if err != nil {
		return err
	}
	if err = u.expectSize(size); err != nil {
		return err
	}
	if value == nil || value.ref == nil {
		return UninitializedObject
	}
	maskAsStructure, err := vectorMask(mask, size)
	if err != nil {
		return err
	}
	w.ref = newVectorReference[Dw](newComputedVector[Dw](
		size, w.ref,
		maskAsStructure, accum,
		newVectorApplyIndexOpScalar(op, u.ref, value.ref),
		desc,
	), -1)
	return nil
}

func MatrixApplyIndexOp[DC, DA, Ds any](C *Matrix[DC], mask *Matrix[bool], accum BinaryOp[DC, DC, DC], op IndexUnaryOp[DC, DA, Ds], A *Matrix[DA], value Ds, desc Descriptor) error {
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
	C.ref = newMatrixReference[DC](newComputedMatrix[DC](
		nrows, ncols,
		C.ref, maskAsStructure, accum,
		newMatrixApplyIndexOp(op, maybeTran(A.ref, isTran), value),
		desc,
	), -1)
	return nil
}

func MatrixApplyIndexOpScalar[DC, DA, Ds any](C *Matrix[DC], mask *Matrix[bool], accum BinaryOp[DC, DC, DC], op IndexUnaryOp[DC, DA, Ds], A *Matrix[DA], value *Scalar[Ds], desc Descriptor) error {
	nrows, ncols, err := C.Size()
	if err != nil {
		return err
	}
	isTran, err := A.expectSizeTran(nrows, ncols, desc, Inp0)
	if err != nil {
		return err
	}
	if value == nil || value.ref == nil {
		return UninitializedObject
	}
	maskAsStructure, err := matrixMask(mask, nrows, ncols)
	if err != nil {
		return err
	}
	C.ref = newMatrixReference[DC](newComputedMatrix[DC](
		nrows, ncols,
		C.ref, maskAsStructure, accum,
		newMatrixApplyIndexOpScalar(op, maybeTran(A.ref, isTran), value.ref),
		desc,
	), -1)
	return nil
}
