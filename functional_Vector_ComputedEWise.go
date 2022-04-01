package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/parallel"
	"github.com/intel/forGoParallel/pipeline"
)

type vectorEWiseMultBinaryOp[Dw, Du, Dv any] struct {
	op BinaryOp[Dw, Du, Dv]
	u  *vectorReference[Du]
	v  *vectorReference[Dv]
}

func newVectorEWiseMultBinaryOp[Dw, Du, Dv any](
	op BinaryOp[Dw, Du, Dv],
	u *vectorReference[Du],
	v *vectorReference[Dv],
) computeVectorT[Dw] {
	return vectorEWiseMultBinaryOp[Dw, Du, Dv]{
		op: op,
		u:  u,
		v:  v,
	}
}

func (compute vectorEWiseMultBinaryOp[Dw, Du, Dv]) resize(newSize int) computeVectorT[Dw] {
	var u *vectorReference[Du]
	var v *vectorReference[Dv]
	parallel.Do(func() {
		u = compute.u.resize(newSize)
	}, func() {
		v = compute.v.resize(newSize)
	})
	return newVectorEWiseMultBinaryOp[Dw, Du, Dv](compute.op, u, v)
}

func (compute vectorEWiseMultBinaryOp[Dw, Du, Dv]) computeElement(index int) (result Dw, ok bool) {
	if uValue, uok := compute.u.extractElement(index); uok {
		if vValue, vok := compute.v.extractElement(index); vok {
			return compute.op(uValue, vValue), true
		}
	}
	return
}

func (compute vectorEWiseMultBinaryOp[Dw, Du, Dv]) computePipeline() *pipeline.Pipeline[any] {
	up := compute.u.getPipeline()
	if up == nil {
		return nil
	}
	vp := compute.v.getPipeline()
	if vp == nil {
		return nil
	}
	return makeVector2SourcePipeline(up, vp,
		func(index int, uValue Du, uok bool, vValue Dv, vok bool) (result Dw, ok bool) {
			if uok && vok {
				return compute.op(uValue, vValue), true
			}
			return
		})
}

type vectorEWiseAddBinaryOp[D any] struct {
	op   BinaryOp[D, D, D]
	u, v *vectorReference[D]
}

func newVectorEWiseAddBinaryOp[D any](
	op BinaryOp[D, D, D],
	u, v *vectorReference[D],
) computeVectorT[D] {
	return vectorEWiseAddBinaryOp[D]{
		op: op,
		u:  u,
		v:  v,
	}
}

func (compute vectorEWiseAddBinaryOp[D]) resize(newSize int) computeVectorT[D] {
	var u, v *vectorReference[D]
	parallel.Do(func() {
		u = compute.u.resize(newSize)
	}, func() {
		v = compute.v.resize(newSize)

	})
	return newVectorEWiseAddBinaryOp[D](compute.op, u, v)
}

func (compute vectorEWiseAddBinaryOp[D]) computeElement(index int) (result D, ok bool) {
	var uValue, vValue D
	var uok, vok bool
	parallel.Do(func() {
		uValue, uok = compute.u.extractElement(index)
	}, func() {
		vValue, vok = compute.v.extractElement(index)
	})
	if uok {
		if vok {
			return compute.op(uValue, vValue), true
		}
		return uValue, true
	}
	return vValue, vok
}

func (compute vectorEWiseAddBinaryOp[D]) computePipeline() *pipeline.Pipeline[any] {
	up := compute.u.getPipeline()
	vp := compute.v.getPipeline()
	if up == nil {
		return vp
	}
	if vp == nil {
		return up
	}
	return makeVector2SourcePipeline(up, vp,
		func(index int, uValue D, uok bool, vValue D, vok bool) (result D, ok bool) {
			if uok {
				if vok {
					return compute.op(uValue, vValue), true
				}
				return uValue, true
			}
			return vValue, vok
		})
}
