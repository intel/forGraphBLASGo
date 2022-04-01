package forGraphBLASGo

import "github.com/intel/forGoParallel/pipeline"

type matrixReduceBinaryOp[D any] struct {
	op BinaryOp[D, D, D]
	A  *matrixReference[D]
}

func newMatrixReduceBinaryOp[D any](op BinaryOp[D, D, D], A *matrixReference[D]) computeVectorT[D] {
	return matrixReduceBinaryOp[D]{op: op, A: A}
}

func (compute matrixReduceBinaryOp[D]) resize(newSize int) computeVectorT[D] {
	_, ncols := compute.A.size()
	return newMatrixReduceBinaryOp[D](compute.op, compute.A.resize(newSize, ncols))
}

func (compute matrixReduceBinaryOp[D]) computeElement(index int) (result D, ok bool) {
	p := compute.A.getRowPipeline(index)
	if p == nil {
		return
	}
	return vectorPipelineReduce(p, compute.op)
}

func (compute matrixReduceBinaryOp[D]) computePipeline() *pipeline.Pipeline[any] {
	rowPipelines := compute.A.getRowPipelines()
	var p pipeline.Pipeline[any]
	p.Source(pipeline.NewFunc[any](-1, func(size int) (data any, fetched int, err error) {
		var result vectorSlice[D]
		for fetched < size && len(rowPipelines) > 0 {
			value, ok := vectorPipelineReduce(rowPipelines[0].p, compute.op)
			if ok {
				result.indices = append(result.indices, rowPipelines[0].index)
				result.values = append(result.values, value)
				fetched++
			}
			rowPipelines[0].p = nil
			rowPipelines = rowPipelines[1:]
		}
		return result, fetched, nil
	}))
	return &p
}
