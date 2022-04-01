package forGraphBLASGo

import "github.com/intel/forGoParallel/pipeline"

type vectorSelect[D, Ds any] struct {
	op    IndexUnaryOp[bool, D, Ds]
	u     *vectorReference[D]
	value Ds
}

func newVectorSelect[D, Ds any](op IndexUnaryOp[bool, D, Ds], u *vectorReference[D], value Ds) computeVectorT[D] {
	return vectorSelect[D, Ds]{op: op, u: u, value: value}
}

func (compute vectorSelect[D, Ds]) resize(newSize int) computeVectorT[D] {
	return newVectorSelect[D, Ds](compute.op, compute.u.resize(newSize), compute.value)
}

func (compute vectorSelect[D, Ds]) computeElement(index int) (result D, ok bool) {
	if u, uok := compute.u.extractElement(index); uok {
		if compute.op(u, index, 0, compute.value) {
			return u, true
		}
	}
	return
}

func (compute vectorSelect[D, Ds]) computePipeline() *pipeline.Pipeline[any] {
	p := compute.u.getPipeline()
	if p == nil {
		return nil
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[D])
			slice.filter(func(index int, value D) (newIndex int, newValue D, ok bool) {
				return index, value, compute.op(value, index, 0, compute.value)
			})
			return slice
		})),
	)
	return p
}

type vectorSelectScalar[D, Ds any] struct {
	op    IndexUnaryOp[bool, D, Ds]
	u     *vectorReference[D]
	value *scalarReference[Ds]
}

func newVectorSelectScalar[D, Ds any](op IndexUnaryOp[bool, D, Ds], u *vectorReference[D], value *scalarReference[Ds]) computeVectorT[D] {
	return vectorSelectScalar[D, Ds]{op: op, u: u, value: value}
}

func (compute vectorSelectScalar[D, Ds]) resize(newSize int) computeVectorT[D] {
	return newVectorSelectScalar[D, Ds](compute.op, compute.u.resize(newSize), compute.value)
}

func (compute vectorSelectScalar[D, Ds]) computeElement(index int) (result D, ok bool) {
	if u, uok := compute.u.extractElement(index); uok {
		if v, vok := compute.value.extractElement(); vok {
			if compute.op(u, index, 0, v) {
				return u, true
			}
		} else {
			panic(EmptyObject)
		}
	}
	return
}

func (compute vectorSelectScalar[D, Ds]) computePipeline() *pipeline.Pipeline[any] {
	p := compute.u.getPipeline()
	if p == nil {
		return nil
	}
	v, vok := compute.value.extractElement()
	if !vok {
		panic(EmptyObject)
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[D])
			slice.filter(func(index int, value D) (newIndex int, newValue D, ok bool) {
				return index, value, compute.op(value, index, 0, v)
			})
			return slice
		})),
	)
	return p
}
