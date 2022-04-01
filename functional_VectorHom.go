package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/pipeline"
)

// todo: this needs to become isoVectorConstant, because of resize with size > nsize
type homVectorConstant[T any] struct {
	nsize int
	value T
}

func newHomVectorConstant[T any](size int, value T) homVectorConstant[T] {
	return homVectorConstant[T]{nsize: size, value: value}
}

func (v homVectorConstant[T]) homValue() (T, bool) {
	return v.value, true
}

func (v homVectorConstant[T]) resize(ref *vectorReference[T], newSize int) *vectorReference[T] {
	if v.nsize == newSize {
		return ref
	}
	return newVectorReference[T](newHomVectorConstant[T](newSize, v.value), int64(newSize))
}

func (v homVectorConstant[T]) size() int {
	return v.nsize
}

func (v homVectorConstant[T]) nvals() int {
	return v.nsize
}

func (v homVectorConstant[T]) setElement(ref *vectorReference[T], value T, index int) *vectorReference[T] {
	if equal(v.value, value) {
		return ref
	}
	return newVectorReference[T](newListVector[T](
		v.nsize, ref,
		&vectorValueList[T]{
			col:   index,
			value: value,
		}),
		int64(v.nsize),
	)
}

func (v homVectorConstant[T]) removeElement(ref *vectorReference[T], index int) *vectorReference[T] {
	return newVectorReference[T](newListVector[T](
		v.nsize, ref,
		&vectorValueList[T]{col: -index},
	),
		int64(v.nsize-1),
	)
}

func (v homVectorConstant[T]) extractElement(_ int) (T, bool) {
	return v.value, true
}

func (v homVectorConstant[T]) getPipeline() *pipeline.Pipeline[any] {
	index := 0
	var values []T
	var p pipeline.Pipeline[any]
	p.Source(pipeline.NewFunc[any](v.nsize, func(size int) (data any, fetched int, err error) {
		var result vectorSlice[T]
		if index >= v.nsize {
			return
		}
		if index+size > v.nsize {
			size = v.nsize - index
		}
		if size < len(values) {
			values = values[:size]
		} else {
			for len(values) < size {
				values = append(values, v.value)
			}
		}
		result.cow = cowv
		result.indices = make([]int, size)
		result.values = values
		for i := 0; i < size; i++ {
			result.indices[i] = index + i
		}
		index += size
		return result, size, nil
	}))
	return &p
}

func (_ homVectorConstant[T]) optimized() bool {
	return true
}

func (v homVectorConstant[T]) optimize() functionalVector[T] {
	return v
}

// todo: this needs to become isoVectorScalar
type homVectorScalar[T any] struct {
	nsize int
	value *scalarReference[T]
}

func newHomVectorScalar[T any](size int, value *scalarReference[T]) functionalVector[T] {
	v := value.get()
	if v.optimized() {
		if val, ok := v.extractElement(value); ok {
			return newHomVectorConstant[T](size, val)
		}
		return newSparseVector[T](size, nil, nil)
	}
	return homVectorScalar[T]{nsize: size, value: value}
}

func (v homVectorScalar[T]) homValue() (T, bool) {
	return v.value.extractElement()
}

func (v homVectorScalar[T]) resize(ref *vectorReference[T], newSize int) *vectorReference[T] {
	if v.nsize == newSize {
		return ref
	}
	return newVectorReference[T](newHomVectorScalar[T](newSize, v.value), -1)
}

func (v homVectorScalar[T]) size() int {
	return v.nsize
}

func (v homVectorScalar[T]) nvals() int {
	if _, ok := v.value.extractElement(); ok {
		return v.nsize
	}
	return 0
}

func (v homVectorScalar[T]) setElement(ref *vectorReference[T], value T, index int) *vectorReference[T] {
	s := v.value.get()
	if s.optimized() {
		if val, ok := s.extractElement(v.value); ok {
			if equal(val, value) {
				return ref
			}
		} else {
			return newVectorReference[T](newSparseVector[T](v.nsize, []int{index}, []T{value}), 1)
		}
	}
	return newVectorReference[T](newListVector[T](
		v.nsize, ref,
		&vectorValueList[T]{
			col:   index,
			value: value,
		}), -1)
}

func (v homVectorScalar[T]) removeElement(ref *vectorReference[T], index int) *vectorReference[T] {
	s := v.value.get()
	if s.optimized() && !s.valid() {
		return ref
	}
	return newVectorReference[T](newListVector[T](v.nsize, ref, &vectorValueList[T]{col: -index}), -1)
}

func (v homVectorScalar[T]) extractElement(_ int) (T, bool) {
	return v.value.extractElement()
}

func (v homVectorScalar[T]) getPipeline() *pipeline.Pipeline[any] {
	value, ok := v.value.extractElement()
	if !ok {
		return nil
	}
	index := 0
	var values []T
	var p pipeline.Pipeline[any]
	p.Source(pipeline.NewFunc[any](v.nsize, func(size int) (data any, fetched int, err error) {
		var result vectorSlice[T]
		if index >= v.nsize {
			return
		}
		if index+size > v.nsize {
			size = v.nsize - index
		}
		if size < len(values) {
			values = values[:size]
		} else {
			for len(values) < size {
				values = append(values, value)
			}
		}
		result.cow = cowv
		result.indices = make([]int, size)
		result.values = values
		for i := 0; i < size; i++ {
			result.indices[i] = index + i
		}
		index += size
		return result, size, nil
	}))
	return &p
}

func (_ homVectorScalar[T]) optimized() bool {
	return false
}

func (v homVectorScalar[T]) optimize() functionalVector[T] {
	value, ok := v.value.extractElement()
	if !ok {
		return newSparseVector[T](v.nsize, nil, nil)
	}
	return newHomVectorConstant[T](v.nsize, value)
}
