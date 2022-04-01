package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/pipeline"
	"sort"
)

// sparseVector has no duplicate cols
type sparseVector[T any] struct {
	nsize  int
	cols   []int
	values []T
}

func newSparseVector[T any](size int, cols []int, values []T) sparseVector[T] {
	return sparseVector[T]{
		nsize:  size,
		cols:   cols,
		values: values,
	}
}

func (vector sparseVector[T]) resize(ref *vectorReference[T], newSize int) *vectorReference[T] {
	switch {
	case newSize == vector.nsize:
		return ref
	case newSize > vector.nsize:
		return newVectorReference[T](newSparseVector[T](newSize, vector.cols, vector.values), int64(len(vector.values)))
	default:
		index := sort.SearchInts(vector.cols, newSize)
		if index == 0 {
			return newVectorReference[T](newSparseVector[T](newSize, nil, nil), 0)
		}
		return newVectorReference[T](newSparseVector[T](
			newSize,
			vector.cols[:index],
			vector.values[:index],
		),
			int64(index),
		)
	}
}

func (vector sparseVector[T]) size() int {
	return vector.nsize
}

func (vector sparseVector[T]) nvals() int {
	return len(vector.cols)
}

func (vector sparseVector[T]) setElement(ref *vectorReference[T], value T, index int) *vectorReference[T] {
	return setVectorElement[T](vector, ref, len(vector.values), value, index)
}

func (vector sparseVector[T]) removeElement(ref *vectorReference[T], index int) *vectorReference[T] {
	return removeVectorElement[T](vector, ref, len(vector.values), index, true)
}

func (vector sparseVector[T]) extractElement(index int) (result T, ok bool) {
	i := sort.SearchInts(vector.cols, index)
	if i < len(vector.cols) && vector.cols[i] == index {
		return vector.values[i], true
	}
	return
}

func (vector sparseVector[T]) getPipeline() *pipeline.Pipeline[any] {
	index := 0
	var p pipeline.Pipeline[any]
	p.Source(pipeline.NewFunc[any](len(vector.values), func(size int) (data any, fetched int, err error) {
		var result vectorSlice[T]
		if index >= len(vector.values) {
			return result, 0, nil
		}
		if index+size > len(vector.values) {
			size = len(vector.values) - index
		}
		result = vectorSlice[T]{
			cow:     cow0 | cowv,
			indices: vector.cols[index : index+size],
			values:  vector.values[index : index+size],
		}
		index += size
		return result, size, nil
	}))
	return &p
}

func (_ sparseVector[T]) optimized() bool {
	return true
}

func (vector sparseVector[T]) optimize() functionalVector[T] {
	return vector
}
