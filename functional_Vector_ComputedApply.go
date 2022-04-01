package forGraphBLASGo

import "github.com/intel/forGoParallel/pipeline"

type vectorApply[Dw, Du any] struct {
	op UnaryOp[Dw, Du]
	u  *vectorReference[Du]
}

func newVectorApply[Dw, Du any](op UnaryOp[Dw, Du], u *vectorReference[Du]) computeVectorT[Dw] {
	return vectorApply[Dw, Du]{op: op, u: u}
}

func (compute vectorApply[Dw, Du]) resize(newSize int) computeVectorT[Dw] {
	return vectorApply[Dw, Du]{
		op: compute.op,
		u:  compute.u.resize(newSize),
	}
}

func (compute vectorApply[Dw, Du]) computeElement(index int) (result Dw, ok bool) {
	if value, ok := compute.u.extractElement(index); ok {
		return compute.op(value), true
	}
	return
}

func (compute vectorApply[Dw, Du]) computePipeline() *pipeline.Pipeline[any] {
	p := compute.u.getPipeline()
	if p == nil {
		return nil
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[Du])
			result := vectorSlice[Dw]{
				cow:     slice.cow &^ cowv,
				indices: slice.indices,
				values:  make([]Dw, len(slice.values)),
			}
			for i, value := range slice.values {
				result.values[i] = compute.op(value)
			}
			return result
		})),
	)
	return p
}

type vectorApplyBinaryOp1st[Dw, Ds, Du any] struct {
	op    BinaryOp[Dw, Ds, Du]
	value Ds
	u     *vectorReference[Du]
}

func newVectorApplyBinaryOp1st[Dw, Ds, Du any](op BinaryOp[Dw, Ds, Du], value Ds, u *vectorReference[Du]) computeVectorT[Dw] {
	return vectorApplyBinaryOp1st[Dw, Ds, Du]{op: op, value: value, u: u}
}

func (compute vectorApplyBinaryOp1st[Dw, Ds, Du]) resize(newSize int) computeVectorT[Dw] {
	return vectorApplyBinaryOp1st[Dw, Ds, Du]{
		op:    compute.op,
		value: compute.value,
		u:     compute.u.resize(newSize),
	}
}

func (compute vectorApplyBinaryOp1st[Dw, Ds, Du]) computeElement(index int) (result Dw, ok bool) {
	if u, uok := compute.u.extractElement(index); uok {
		return compute.op(compute.value, u), true
	}
	return
}

func (compute vectorApplyBinaryOp1st[Dw, Ds, Du]) computePipeline() *pipeline.Pipeline[any] {
	p := compute.u.getPipeline()
	if p == nil {
		return nil
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[Du])
			result := vectorSlice[Dw]{
				cow:     slice.cow &^ cowv,
				indices: slice.indices,
				values:  make([]Dw, len(slice.values)),
			}
			for i, value := range slice.values {
				result.values[i] = compute.op(compute.value, value)
			}
			return result
		})),
	)
	return p
}

type vectorApplyBinaryOp1stScalar[Dw, Ds, Du any] struct {
	op    BinaryOp[Dw, Ds, Du]
	value *scalarReference[Ds]
	u     *vectorReference[Du]
}

func newVectorApplyBinaryOp1stScalar[Dw, Ds, Du any](op BinaryOp[Dw, Ds, Du], value *scalarReference[Ds], u *vectorReference[Du]) computeVectorT[Dw] {
	return vectorApplyBinaryOp1stScalar[Dw, Ds, Du]{op: op, value: value, u: u}
}

func (compute vectorApplyBinaryOp1stScalar[Dw, Ds, Du]) resize(newSize int) computeVectorT[Dw] {
	return vectorApplyBinaryOp1stScalar[Dw, Ds, Du]{
		op:    compute.op,
		value: compute.value,
		u:     compute.u.resize(newSize),
	}
}

func (compute vectorApplyBinaryOp1stScalar[Dw, Ds, Du]) computeElement(index int) (result Dw, ok bool) {
	if u, uok := compute.u.extractElement(index); uok {
		if s, sok := compute.value.extractElement(); sok {
			return compute.op(s, u), true
		} else {
			panic(EmptyObject)
		}
	}
	return
}

func (compute vectorApplyBinaryOp1stScalar[Dw, Ds, Du]) computePipeline() *pipeline.Pipeline[any] {
	p := compute.u.getPipeline()
	if p == nil {
		return nil
	}
	s, sok := compute.value.extractElement()
	if !sok {
		panic(EmptyObject)
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[Du])
			result := vectorSlice[Dw]{
				cow:     slice.cow &^ cowv,
				indices: slice.indices,
				values:  make([]Dw, len(slice.values)),
			}
			for i, value := range slice.values {
				result.values[i] = compute.op(s, value)
			}
			return result
		})),
	)
	return p
}

type vectorApplyBinaryOp2nd[Dw, Du, Ds any] struct {
	op    BinaryOp[Dw, Du, Ds]
	u     *vectorReference[Du]
	value Ds
}

func newVectorApplyBinaryOp2nd[Dw, Du, Ds any](op BinaryOp[Dw, Du, Ds], u *vectorReference[Du], value Ds) computeVectorT[Dw] {
	return vectorApplyBinaryOp2nd[Dw, Du, Ds]{op: op, u: u, value: value}
}

func (compute vectorApplyBinaryOp2nd[Dw, Du, Ds]) resize(newSize int) computeVectorT[Dw] {
	return vectorApplyBinaryOp2nd[Dw, Du, Ds]{
		op:    compute.op,
		u:     compute.u.resize(newSize),
		value: compute.value,
	}
}

func (compute vectorApplyBinaryOp2nd[Dw, Du, Ds]) computeElement(index int) (result Dw, ok bool) {
	if u, uok := compute.u.extractElement(index); uok {
		return compute.op(u, compute.value), true
	}
	return
}

func (compute vectorApplyBinaryOp2nd[Dw, Du, Ds]) computePipeline() *pipeline.Pipeline[any] {
	p := compute.u.getPipeline()
	if p == nil {
		return nil
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[Du])
			result := vectorSlice[Dw]{
				cow:     slice.cow &^ cowv,
				indices: slice.indices,
				values:  make([]Dw, len(slice.values)),
			}
			for i, value := range slice.values {
				result.values[i] = compute.op(value, compute.value)
			}
			return result
		})),
	)
	return p
}

type vectorApplyBinaryOp2ndScalar[Dw, Du, Ds any] struct {
	op    BinaryOp[Dw, Du, Ds]
	u     *vectorReference[Du]
	value *scalarReference[Ds]
}

func newVectorApplyBinaryOp2ndScalar[Dw, Du, Ds any](op BinaryOp[Dw, Du, Ds], u *vectorReference[Du], value *scalarReference[Ds]) computeVectorT[Dw] {
	return vectorApplyBinaryOp2ndScalar[Dw, Du, Ds]{op: op, u: u, value: value}
}

func (compute vectorApplyBinaryOp2ndScalar[Dw, Du, Ds]) resize(newSize int) computeVectorT[Dw] {
	return vectorApplyBinaryOp2ndScalar[Dw, Du, Ds]{
		op:    compute.op,
		u:     compute.u.resize(newSize),
		value: compute.value,
	}
}

func (compute vectorApplyBinaryOp2ndScalar[Dw, Du, Ds]) computeElement(index int) (result Dw, ok bool) {
	if u, uok := compute.u.extractElement(index); uok {
		if s, sok := compute.value.extractElement(); sok {
			return compute.op(u, s), true
		} else {
			panic(EmptyObject)
		}
	}
	return
}

func (compute vectorApplyBinaryOp2ndScalar[Dw, Du, Ds]) computePipeline() *pipeline.Pipeline[any] {
	p := compute.u.getPipeline()
	if p == nil {
		return nil
	}
	s, sok := compute.value.extractElement()
	if !sok {
		panic(EmptyObject)
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[Du])
			result := vectorSlice[Dw]{
				cow:     slice.cow &^ cowv,
				indices: slice.indices,
				values:  make([]Dw, len(slice.values)),
			}
			for i, value := range slice.values {
				result.values[i] = compute.op(value, s)
			}
			return result
		})),
	)
	return p
}

type vectorApplyIndexOp[Dw, Du, Ds any] struct {
	op IndexUnaryOp[Dw, Du, Ds]
	u  *vectorReference[Du]
	s  Ds
}

func newVectorApplyIndexOp[Dw, Du, Ds any](op IndexUnaryOp[Dw, Du, Ds], u *vectorReference[Du], s Ds) computeVectorT[Dw] {
	return vectorApplyIndexOp[Dw, Du, Ds]{op: op, u: u, s: s}
}

func (compute vectorApplyIndexOp[Dw, Du, Ds]) resize(newSize int) computeVectorT[Dw] {
	return vectorApplyIndexOp[Dw, Du, Ds]{
		op: compute.op,
		u:  compute.u.resize(newSize),
		s:  compute.s,
	}
}

func (compute vectorApplyIndexOp[Dw, Du, Ds]) computeElement(index int) (result Dw, ok bool) {
	if value, ok := compute.u.extractElement(index); ok {
		return compute.op(value, index, 0, compute.s), true
	}
	return
}

func (compute vectorApplyIndexOp[Dw, Du, Ds]) computePipeline() *pipeline.Pipeline[any] {
	p := compute.u.getPipeline()
	if p == nil {
		return nil
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[Du])
			result := vectorSlice[Dw]{
				cow:     slice.cow &^ cowv,
				indices: slice.indices,
				values:  make([]Dw, len(slice.values)),
			}
			for i, value := range slice.values {
				result.values[i] = compute.op(value, slice.indices[i], 0, compute.s)
			}
			return result
		})),
	)
	return p
}

type vectorApplyIndexOpScalar[Dw, Du, Ds any] struct {
	op IndexUnaryOp[Dw, Du, Ds]
	u  *vectorReference[Du]
	s  *scalarReference[Ds]
}

func newVectorApplyIndexOpScalar[Dw, Du, Ds any](op IndexUnaryOp[Dw, Du, Ds], u *vectorReference[Du], s *scalarReference[Ds]) computeVectorT[Dw] {
	return vectorApplyIndexOpScalar[Dw, Du, Ds]{op: op, u: u, s: s}
}

func (compute vectorApplyIndexOpScalar[Dw, Du, Ds]) resize(newSize int) computeVectorT[Dw] {
	return vectorApplyIndexOpScalar[Dw, Du, Ds]{
		op: compute.op,
		u:  compute.u.resize(newSize),
		s:  compute.s,
	}
}

func (compute vectorApplyIndexOpScalar[Dw, Du, Ds]) computeElement(index int) (result Dw, ok bool) {
	if u, uok := compute.u.extractElement(index); uok {
		if s, sok := compute.s.extractElement(); sok {
			return compute.op(u, index, 0, s), true
		} else {
			panic(EmptyObject)
		}
	}
	return
}

func (compute vectorApplyIndexOpScalar[Dw, Du, Ds]) computePipeline() *pipeline.Pipeline[any] {
	p := compute.u.getPipeline()
	if p == nil {
		return nil
	}
	s, sok := compute.s.extractElement()
	if !sok {
		panic(EmptyObject)
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[Du])
			result := vectorSlice[Dw]{
				cow:     slice.cow &^ cowv,
				indices: slice.indices,
				values:  make([]Dw, len(slice.values)),
			}
			for i, value := range slice.values {
				result.values[i] = compute.op(value, slice.indices[i], 0, s)
			}
			return result
		})),
	)
	return p
}
