package forGraphBLASGo

type scalarApply[Dt, Df any] struct {
	op UnaryOp[Dt, Df]
	s  *scalarReference[Df]
}

func newScalarApply[Dt, Df any](op UnaryOp[Dt, Df], s *scalarReference[Df]) computeScalarT[Dt] {
	return scalarApply[Dt, Df]{op: op, s: s}
}

func (compute scalarApply[Dt, Df]) computeElement() (result Dt, ok bool) {
	if v, vok := compute.s.extractElement(); vok {
		return compute.op(v), true
	}
	return
}

type scalarApplyBinaryOp1st[Dt, Ds, Df any] struct {
	op    BinaryOp[Dt, Ds, Df]
	value Ds
	s     *scalarReference[Df]
}

func newScalarApplyBinaryOp1st[Dt, Ds, Df any](op BinaryOp[Dt, Ds, Df], value Ds, s *scalarReference[Df]) computeScalarT[Dt] {
	return scalarApplyBinaryOp1st[Dt, Ds, Df]{op: op, value: value, s: s}
}

func (compute scalarApplyBinaryOp1st[Dt, Ds, Df]) computeElement() (result Dt, ok bool) {
	if v, vok := compute.s.extractElement(); vok {
		return compute.op(compute.value, v), true
	}
	return
}

type scalarApplyBinaryOp2nd[Dt, Df, Ds any] struct {
	op    BinaryOp[Dt, Df, Ds]
	s     *scalarReference[Df]
	value Ds
}

func newScalarApplyBinaryOp2nd[Dt, Df, Ds any](op BinaryOp[Dt, Df, Ds], s *scalarReference[Df], value Ds) computeScalarT[Dt] {
	return scalarApplyBinaryOp2nd[Dt, Df, Ds]{op: op, s: s, value: value}
}

func (compute scalarApplyBinaryOp2nd[Dt, Df, Ds]) computeElement() (result Dt, ok bool) {
	if v, vok := compute.s.extractElement(); vok {
		return compute.op(v, compute.value), true
	}
	return
}

type scalarApplyBinary[Dt, Df1, Df2 any] struct {
	op BinaryOp[Dt, Df1, Df2]
	s1 *scalarReference[Df1]
	s2 *scalarReference[Df2]
}

func newScalarApplyBinary[Dt, Df1, Df2 any](op BinaryOp[Dt, Df1, Df2], s1 *scalarReference[Df1], s2 *scalarReference[Df2]) computeScalarT[Dt] {
	return scalarApplyBinary[Dt, Df1, Df2]{op: op, s1: s1, s2: s2}
}

func (compute scalarApplyBinary[Dt, Df1, Df2]) computeElement() (result Dt, ok bool) {
	if v1, v1ok := compute.s1.extractElement(); v1ok {
		if v2, v2ok := compute.s2.extractElement(); v2ok {
			return compute.op(v1, v2), true
		}
	}
	return
}
