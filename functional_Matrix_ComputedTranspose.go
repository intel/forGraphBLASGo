package forGraphBLASGo

import "github.com/intel/forGoParallel/pipeline"

type transposeMatrix[D any] struct {
	A *matrixReference[D]
}

func newTransposeMatrix[D any](A *matrixReference[D]) computeMatrixT[D] {
	return transposeMatrix[D]{A: A}
}

func (compute transposeMatrix[D]) resize(newNRows, newNCols int) computeMatrixT[D] {
	return newTransposeMatrix(compute.A.resize(newNCols, newNRows))
}

func (m transposeMatrix[D]) computeElement(row, col int) (result D, ok bool) {
	return m.A.extractElement(col, row)
}

func (m transposeMatrix[D]) computePipeline() *pipeline.Pipeline[any] {
	return transposeMatrixPipeline(m.A)
}

func (m transposeMatrix[D]) computeRowPipeline(row int) *pipeline.Pipeline[any] {
	return m.A.getColPipeline(row)
}

func (m transposeMatrix[D]) computeColPipeline(col int) *pipeline.Pipeline[any] {
	return m.A.getRowPipeline(col)
}

func (m transposeMatrix[D]) computeRowPipelines() []matrix1Pipeline {
	return m.A.getColPipelines()
}

func (m transposeMatrix[D]) computeColPipelines() []matrix1Pipeline {
	return m.A.getRowPipelines()
}
