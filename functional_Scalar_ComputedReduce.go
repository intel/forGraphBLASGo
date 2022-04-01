package forGraphBLASGo

type vectorReduceBinaryOpScalar[D any] struct {
	op BinaryOp[D, D, D]
	u  *vectorReference[D]
}

func newVectorReduceBinaryOpScalar[D any](op BinaryOp[D, D, D], u *vectorReference[D]) computeScalarT[D] {
	return vectorReduceBinaryOpScalar[D]{op: op, u: u}
}

func (compute vectorReduceBinaryOpScalar[D]) computeElement() (result D, ok bool) {
	return vectorPipelineReduce(compute.u.getPipeline(), compute.op)
}

type matrixReduceBinaryOpScalar[D any] struct {
	op BinaryOp[D, D, D]
	A  *matrixReference[D]
}

func newMatrixReduceBinaryOpScalar[D any](op BinaryOp[D, D, D], A *matrixReference[D]) computeScalarT[D] {
	return matrixReduceBinaryOpScalar[D]{op: op, A: A}
}

func (compute matrixReduceBinaryOpScalar[D]) computeElement() (result D, ok bool) {
	return matrixPipelineReduce(compute.A.getPipeline(), compute.op)
}
