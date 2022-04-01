package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/pipeline"
)

type vxM[Dw, Du, DA any] struct {
	op Semiring[Dw, Du, DA]
	u  *vectorReference[Du]
	A  *matrixReference[DA]
}

func newVxM[Dw, Du, DA any](
	op Semiring[Dw, Du, DA],
	u *vectorReference[Du],
	A *matrixReference[DA],
) computeVectorT[Dw] {
	return vxM[Dw, Du, DA]{
		op: op,
		u:  u,
		A:  A,
	}
}

func (compute vxM[Dw, Du, DA]) resize(newSize int) computeVectorT[Dw] {
	nrows, _ := compute.A.size()
	A := compute.A.resize(nrows, newSize)
	return newVxM[Dw, Du, DA](compute.op, compute.u, A)
}

func (compute vxM[Dw, Du, DA]) computeElement(index int) (result Dw, ok bool) {
	add := compute.op.addition().operator()
	mult := compute.op.multiplication()
	up := compute.u.getPipeline()
	if up == nil {
		return
	}
	ap := compute.A.getColPipeline(index)
	if ap == nil {
		return
	}
	return vectorPipelineReduce(makeVector2SourcePipeline(up, ap,
		func(index int, uValue Du, uok bool, aValue DA, aok bool) (result Dw, ok bool) {
			if uok && aok {
				return mult(uValue, aValue), true
			}
			return
		}), add)
}

func (compute vxM[Dw, Du, DA]) computePipeline() *pipeline.Pipeline[any] {
	add := compute.op.addition().operator()
	mult := compute.op.multiplication()
	go compute.u.optimize()
	if compute.u.nvals() == 0 {
		return nil
	}
	colPipelines := compute.A.getColPipelines()
	var p pipeline.Pipeline[any]
	p.Source(pipeline.NewFunc[any](-1, func(size int) (data any, fetched int, err error) {
		var result vectorSlice[Dw]
		for fetched < size && len(colPipelines) > 0 {
			value, ok := vectorPipelineReduce(makeVector2SourcePipeline(compute.u.getPipeline(), colPipelines[0].p,
				func(index int, uValue Du, uok bool, aValue DA, aok bool) (result Dw, ok bool) {
					if uok && aok {
						return mult(uValue, aValue), true
					}
					return
				},
			), add)
			if ok {
				result.indices = append(result.indices, colPipelines[0].index)
				result.values = append(result.values, value)
				fetched++
			}
			colPipelines[0].p = nil
			colPipelines = colPipelines[1:]
		}
		return result, fetched, nil
	}))
	return &p
}
