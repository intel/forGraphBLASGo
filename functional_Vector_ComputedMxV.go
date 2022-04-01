package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/pipeline"
)

type mxV[Dw, DA, Du any] struct {
	op Semiring[Dw, DA, Du]
	A  *matrixReference[DA]
	u  *vectorReference[Du]
}

func newMxV[Dw, DA, Du any](
	op Semiring[Dw, DA, Du],
	A *matrixReference[DA],
	u *vectorReference[Du],
) computeVectorT[Dw] {
	return mxV[Dw, DA, Du]{
		op: op,
		A:  A,
		u:  u,
	}
}

func (compute mxV[Dw, DA, Du]) resize(newSize int) computeVectorT[Dw] {
	_, ncols := compute.A.size()
	A := compute.A.resize(newSize, ncols)
	return newMxV[Dw, DA, Du](compute.op, A, compute.u)
}

func (compute mxV[Dw, DA, Du]) computeElement(index int) (result Dw, ok bool) {
	add := compute.op.addition().operator()
	mult := compute.op.multiplication()
	ap := compute.A.getRowPipeline(index)
	if ap == nil {
		return
	}
	up := compute.u.getPipeline()
	if up == nil {
		return
	}
	return vectorPipelineReduce(makeVector2SourcePipeline(ap, up,
		func(index int, aValue DA, aok bool, uValue Du, uok bool) (result Dw, ok bool) {
			if aok && uok {
				return mult(aValue, uValue), true
			}
			return
		}), add)
}

func (compute mxV[Dw, DA, Du]) computePipeline() *pipeline.Pipeline[any] {
	add := compute.op.addition().operator()
	mult := compute.op.multiplication()
	go compute.u.optimize()
	if compute.u.nvals() == 0 {
		return nil
	}
	rowPipelines := compute.A.getRowPipelines()
	var p pipeline.Pipeline[any]
	p.Source(pipeline.NewFunc[any](-1, func(size int) (data any, fetched int, err error) {
		var result vectorSlice[Dw]
		for fetched < size && len(rowPipelines) > 0 {
			value, ok := vectorPipelineReduce(makeVector2SourcePipeline(rowPipelines[0].p, compute.u.getPipeline(),
				func(index int, aValue DA, aok bool, uValue Du, uok bool) (result Dw, ok bool) {
					if aok && uok {
						return mult(aValue, uValue), true
					}
					return
				},
			), add)
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
