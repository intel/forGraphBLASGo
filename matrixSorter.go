package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/psort"
	"sort"
)

type matrixSorter[T any] struct {
	rowFirst   bool
	rows, cols []int
	values     []T
}

func (s matrixSorter[T]) Len() int {
	return len(s.rows)
}

func (s matrixSorter[T]) Less(i, j int) bool {
	if s.rowFirst {
		ri, rj := s.rows[i], s.rows[j]
		if ri < rj {
			return true
		}
		if ri > rj {
			return false
		}
		return s.cols[i] < s.cols[j]
	}
	ci, cj := s.cols[i], s.cols[j]
	if ci < cj {
		return true
	}
	if ci > cj {
		return false
	}
	return s.rows[i] < s.rows[j]
}

func (s matrixSorter[T]) SequentialSort(i, j int) {
	sort.Stable(matrixSorter[T]{
		rowFirst: s.rowFirst,
		rows:     s.rows[i:j],
		cols:     s.cols[i:j],
		values:   s.values[i:j],
	})
}

func (s matrixSorter[T]) NewTemp() psort.StableSorter {
	return matrixSorter[T]{
		rowFirst: s.rowFirst,
		rows:     make([]int, len(s.rows)),
		cols:     make([]int, len(s.cols)),
		values:   make([]T, len(s.values)),
	}
}

func (dst matrixSorter[T]) Assign(source psort.StableSorter) func(i, j, len int) {
	src := source.(matrixSorter[T])
	return func(i, j, len int) {
		copy(dst.rows[i:i+len], src.rows[j:j+len])
		copy(dst.cols[i:i+len], src.cols[j:j+len])
		copy(dst.values[i:i+len], src.values[j:j+len])
	}
}

func (s matrixSorter[T]) Swap(i, j int) {
	s.rows[i], s.rows[j] = s.rows[j], s.rows[i]
	s.cols[i], s.cols[j] = s.cols[j], s.cols[i]
	s.values[i], s.values[j] = s.values[j], s.values[i]
}

func matrixSort[T any](rows, cols []int, values []T) {
	psort.StableSort(matrixSorter[T]{
		rowFirst: true,
		rows:     rows,
		cols:     cols,
		values:   values,
	})
}
