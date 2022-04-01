package forGraphBLASGo

func MatrixReduceBinaryOp[D any](w *Vector[D], mask *Vector[bool], accum, op BinaryOp[D, D, D], A *Matrix[D], desc Descriptor) error {
	nrows, ncols, err := A.Size()
	if err != nil {
		return err
	}
	isTran, err := desc.Is(Inp0, Tran)
	if err != nil {
		panic(err)
	}
	if isTran {
		nrows, ncols = ncols, nrows
	}
	if err = w.expectSize(nrows); err != nil {
		return err
	}
	maskAsStructure, err := vectorMask(mask, nrows)
	if err != nil {
		return err
	}
	w.ref = newVectorReference[D](newComputedVector[D](
		nrows, w.ref, maskAsStructure, accum,
		newMatrixReduceBinaryOp[D](op, maybeTran(A.ref, isTran)),
		desc,
	), -1)
	return nil
}

func MatrixReduceMonoid[D any](w *Vector[D], mask *Vector[bool], accum BinaryOp[D, D, D], op Monoid[D], A *Matrix[D], desc Descriptor) error {
	return MatrixReduceBinaryOp(w, mask, accum, op.operator(), A, desc)
}

func VectorReduceBinaryOpScalar[D any](s *Scalar[D], accum, op BinaryOp[D, D, D], u *Vector[D], _ Descriptor) error {
	if s == nil || s.ref == nil || u == nil || u.ref == nil {
		return UninitializedObject
	}
	s.ref = newScalarReference[D](newComputedScalar[D](s.ref, accum, newVectorReduceBinaryOpScalar[D](op, u.ref)))
	return nil
}

func VectorReduceMonoidScalar[D any](s *Scalar[D], accum BinaryOp[D, D, D], op Monoid[D], u *Vector[D], desc Descriptor) error {
	return VectorReduceBinaryOpScalar(s, accum, op.operator(), u, desc)
}

func VectorReduce[D any](value *D, accum BinaryOp[D, D, D], op Monoid[D], u *Vector[D], _ Descriptor) error {
	if u == nil || u.ref == nil {
		return UninitializedObject
	}
	result, ok := vectorPipelineReduce(u.ref.getPipeline(), op.operator())
	if !ok {
		result = op.identity()
	}
	if accum == nil {
		*value = result
		return nil
	}
	*value = accum(*value, result)
	return nil
}

func MatrixReduceBinaryOpScalar[D any](s *Scalar[D], accum, op BinaryOp[D, D, D], A *Matrix[D], _ Descriptor) error {
	if s == nil || s.ref == nil || A == nil || A.ref == nil {
		return UninitializedObject
	}
	s.ref = newScalarReference[D](newComputedScalar[D](s.ref, accum, newMatrixReduceBinaryOpScalar[D](op, A.ref)))
	return nil
}

func MatrixReduceMonoidScalar[D any](s *Scalar[D], accum BinaryOp[D, D, D], op Monoid[D], A *Matrix[D], desc Descriptor) error {
	return MatrixReduceBinaryOpScalar(s, accum, op.operator(), A, desc)
}

func MatrixReduce[D any](value *D, accum BinaryOp[D, D, D], op Monoid[D], A *Matrix[D], _ Descriptor) error {
	if A == nil || A.ref == nil {
		return UninitializedObject
	}
	result, ok := matrixPipelineReduce(A.ref.getPipeline(), op.operator())
	if !ok {
		result = op.identity()
	}
	if accum == nil {
		*value = result
		return nil
	}
	*value = accum(*value, result)
	return nil
}
