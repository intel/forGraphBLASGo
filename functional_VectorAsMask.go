package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/pipeline"
)

type vectorAsMask[T Number] struct {
	v *vectorReference[T]
}

func newVectorAsMask[T Number](v *vectorReference[T]) functionalVector[bool] {
	return vectorAsMask[T]{v: v}
}

func (v vectorAsMask[T]) resize(ref *vectorReference[bool], newSize int) *vectorReference[bool] {
	if newSize == v.v.size() {
		return ref
	}
	return newVectorReference[bool](newVectorAsMask(v.v.resize(newSize)), -1)
}

func (v vectorAsMask[T]) size() int {
	return v.v.size()
}

func (v vectorAsMask[T]) nvals() int {
	return v.v.nvals()
}

func (v vectorAsMask[T]) setElement(ref *vectorReference[bool], value bool, index int) *vectorReference[bool] {
	return newVectorReference[bool](newListVector[bool](v.v.size(), ref,
		&vectorValueList[bool]{
			col:   index,
			value: value,
		},
	), -1)
}

func (v vectorAsMask[T]) removeElement(ref *vectorReference[bool], index int) *vectorReference[bool] {
	return newVectorReference[bool](newListVector[bool](v.v.size(), ref,
		&vectorValueList[bool]{
			col: -index,
		},
	), -1)
}

func (v vectorAsMask[T]) extractElement(index int) (bool, bool) {
	if value, ok := v.v.extractElement(index); ok {
		return value != 0, true
	}
	return false, false
}

func (v vectorAsMask[T]) getPipeline() *pipeline.Pipeline[any] {
	base := v.v.getPipeline()
	if base == nil {
		return nil
	}
	base.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[T])
			result := vectorSlice[bool]{
				cow:     slice.cow &^ cowv,
				indices: slice.indices,
				values:  make([]bool, len(slice.values)),
			}
			for i, value := range slice.values {
				result.values[i] = value != 0
			}
			return result
		})),
	)
	return base
}

func (v vectorAsMask[T]) optimized() bool {
	return v.v.optimized()
}

func (v vectorAsMask[T]) optimize() functionalVector[bool] {
	v.v.optimize()
	return v
}

type vectorAsStructuralMask[T any] struct {
	v *vectorReference[T]
}

func newVectorAsStructuralMask[T any](v *vectorReference[T]) functionalVector[bool] {
	return vectorAsStructuralMask[T]{v: v}
}

func (v vectorAsStructuralMask[T]) resize(ref *vectorReference[bool], newSize int) *vectorReference[bool] {
	if newSize == v.v.size() {
		return ref
	}
	return newVectorReference[bool](newVectorAsStructuralMask(v.v.resize(newSize)), -1)
}

func (v vectorAsStructuralMask[T]) size() int {
	return v.v.size()
}

func (v vectorAsStructuralMask[T]) nvals() int {
	return v.v.nvals()
}

func (v vectorAsStructuralMask[T]) setElement(ref *vectorReference[bool], value bool, index int) *vectorReference[bool] {
	return newVectorReference[bool](newListVector[bool](v.v.size(), ref,
		&vectorValueList[bool]{
			col:   index,
			value: value,
		},
	), -1)
}

func (v vectorAsStructuralMask[T]) removeElement(ref *vectorReference[bool], index int) *vectorReference[bool] {
	return newVectorReference[bool](newListVector[bool](v.v.size(), ref,
		&vectorValueList[bool]{
			col: -index,
		},
	), -1)
}

func (v vectorAsStructuralMask[T]) extractElement(index int) (bool, bool) {
	// todo: when accessing the first return value, this should normally
	// 	result in a DomainMismatch, so check the program flow for this
	_, ok := v.v.extractElement(index)
	return ok, ok
}

func (v vectorAsStructuralMask[T]) getPipeline() *pipeline.Pipeline[any] {
	base := v.v.getPipeline()
	if base == nil {
		return nil
	}
	base.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[T])
			result := vectorSlice[bool]{
				cow:     slice.cow &^ cowv,
				indices: slice.indices,
				values:  make([]bool, len(slice.values)),
			}
			for i := range result.values {
				result.values[i] = true
			}
			return result
		})),
	)
	return base
}

func (v vectorAsStructuralMask[T]) optimized() bool {
	return v.v.optimized()
}

func (v vectorAsStructuralMask[T]) optimize() functionalVector[bool] {
	v.v.optimize()
	return v
}
