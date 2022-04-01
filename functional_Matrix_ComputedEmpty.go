package forGraphBLASGo

import "github.com/intel/forGoParallel/pipeline"

type emptyComputedMatrix[D any] struct{}

func newEmptyComputedMatrix[D any]() emptyComputedMatrix[D] {
	return emptyComputedMatrix[D]{}
}

func (_ emptyComputedMatrix[D]) resize(_, _ int) computeMatrixT[D] {
	return newEmptyComputedMatrix[D]()
}

func (_ emptyComputedMatrix[D]) computeElement(row, col int) (result D, ok bool) {
	return
}

func (_ emptyComputedMatrix[D]) computePipeline() *pipeline.Pipeline[any] {
	return nil
}

func (_ emptyComputedMatrix[D]) computeRowPipeline(_ int) *pipeline.Pipeline[any] {
	return nil
}

func (_ emptyComputedMatrix[D]) computeColPipeline(_ int) *pipeline.Pipeline[any] {
	return nil
}

func (_ emptyComputedMatrix[D]) computeRowPipelines() []matrix1Pipeline {
	return nil
}

func (_ emptyComputedMatrix[D]) computeColPipelines() []matrix1Pipeline {
	return nil
}
