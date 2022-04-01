package forGraphBLASGo

import "github.com/intel/forGoParallel/pipeline"

type kroneckerBinaryOp[DC, DA, DB any] struct {
	op BinaryOp[DC, DA, DB]
	A  *matrixReference[DA]
	B  *matrixReference[DB]
}

func newKroneckerBinaryOp[DC, DA, DB any](
	op BinaryOp[DC, DA, DB],
	A *matrixReference[DA],
	B *matrixReference[DB],
) computeMatrixT[DC] {
	return kroneckerBinaryOp[DC, DA, DB]{
		op: op,
		A:  A,
		B:  B,
	}
}

// todo
func (compute kroneckerBinaryOp[DC, DA, DB]) resize(newNRows, newNCols int) computeMatrixT[DC] {
	panic("todo")
}

func (compute kroneckerBinaryOp[DC, DA, DB]) computeElement(row, col int) (result DC, ok bool) {
	Bnrows, Bncols := compute.B.size()
	if a, aok := compute.A.extractElement(row/Bnrows, col/Bncols); aok {
		if b, bok := compute.B.extractElement(row%Bnrows, col%Bncols); bok {
			return compute.op(a, b), true
		}
	}
	return
}

// todo
func (compute kroneckerBinaryOp[DC, DA, DB]) computePipeline() *pipeline.Pipeline[any] {
	panic("todo")
}

// todo
func (compute kroneckerBinaryOp[DC, DA, DB]) computeRowPipeline(row int) *pipeline.Pipeline[any] {
	panic("todo")
}

// todo
func (compute kroneckerBinaryOp[DC, DA, DB]) computeColPipeline(col int) *pipeline.Pipeline[any] {
	panic("todo")
}

// todo
func (compute kroneckerBinaryOp[DC, DA, DB]) computeRowPipelines() []matrix1Pipeline {
	panic("todo")
}

// todo
func (compute kroneckerBinaryOp[DC, DA, DB]) computeColPipelines() []matrix1Pipeline {
	panic("todo")
}
