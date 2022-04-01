package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/pipeline"
)

// A functionalVector does not modify its contents, but always returns a new vector as a result.

type (
	functionalVector[T any] interface {
		resize(ref *vectorReference[T], newSize int) *vectorReference[T]
		size() int
		nvals() int
		setElement(ref *vectorReference[T], value T, index int) *vectorReference[T]
		removeElement(ref *vectorReference[T], index int) *vectorReference[T]
		extractElement(index int) (T, bool)
		getPipeline() *pipeline.Pipeline[any]

		optimized() bool
		optimize() functionalVector[T]
	}

	homVector[T any] interface {
		functionalVector[T]
		homValue() (T, bool)
	}
)

func setVectorElement[T any](
	v functionalVector[T], ref *vectorReference[T], nvals int,
	value T, index int,
) *vectorReference[T] {
	size := v.size()
	if nvals == 0 {
		return newVectorReference[T](newSparseVector[T](size, []int{index}, []T{value}), 1)
	}
	return newVectorReference[T](newListVector[T](
		size, ref,
		&vectorValueList[T]{col: index, value: value},
	), -1)
}

func removeVectorElement[T any](
	v functionalVector[T], ref *vectorReference[T], nvals int,
	index int, selfAsEmpty bool,
) *vectorReference[T] {
	size := v.size()
	if nvals == 0 {
		if selfAsEmpty {
			return ref
		}
		return newVectorReference[T](newSparseVector[T](size, nil, nil), 0)
	}
	return newVectorReference[T](newListVector[T](
		size, ref,
		&vectorValueList[T]{col: -index},
	), -1)
}
