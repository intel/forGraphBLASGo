package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/psort"
	"sort"
)

type vectorSorter[T any] struct {
	cols   []int
	values []T
}

func (s vectorSorter[T]) Len() int {
	return len(s.cols)
}

func (s vectorSorter[T]) Less(i, j int) bool {
	return s.cols[i] < s.cols[j]
}

func (s vectorSorter[T]) SequentialSort(i, j int) {
	sort.Stable(vectorSorter[T]{
		cols:   s.cols[i:j],
		values: s.values[i:j],
	})
}

func (s vectorSorter[T]) NewTemp() psort.StableSorter {
	return vectorSorter[T]{
		cols:   make([]int, len(s.cols)),
		values: make([]T, len(s.values)),
	}
}

func (dst vectorSorter[T]) Assign(source psort.StableSorter) func(i, j, len int) {
	src := source.(vectorSorter[T])
	return func(i, j, len int) {
		copy(dst.cols[i:i+len], src.cols[j:j+len])
		copy(dst.values[i:i+len], src.values[j:j+len])
	}
}

func (s vectorSorter[T]) Swap(i, j int) {
	s.cols[i], s.cols[j] = s.cols[j], s.cols[i]
	s.values[i], s.values[j] = s.values[j], s.values[i]
}

func vectorSort[T any](cols []int, values []T) {
	psort.StableSort(vectorSorter[T]{
		cols:   cols,
		values: values,
	})
}
