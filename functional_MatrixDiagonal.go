package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/pipeline"
)

type diagonalMatrix[T any] struct {
	nrows, ncols int
	v            *vectorReference[T]
	k            int
}

func newDiagonalMatrix[T any](nrows, ncols int, v *vectorReference[T], k int) functionalMatrix[T] {
	return diagonalMatrix[T]{
		nrows: nrows,
		ncols: ncols,
		v:     v,
		k:     k,
	}
}

func diagonalOffset(k int) (int, int) {
	if k >= 0 {
		return 0, k
	}
	return k, 0
}

func (m diagonalMatrix[T]) resize(ref *matrixReference[T], newNRows, newNCols int) *matrixReference[T] {
	if newNRows == m.nrows && newNCols == m.ncols {
		return ref
	}
	i, j := diagonalOffset(m.k)
	// todo: compute nvalues
	return newMatrixReference(newDiagonalMatrix[T](
		newNRows, newNCols,
		m.v.resize(Min[int](newNRows+i, newNCols-j)),
		m.k,
	), -1)
}

func (m diagonalMatrix[T]) size() (nrows, ncols int) {
	return m.nrows, m.ncols
}

func (m diagonalMatrix[T]) nvals() int {
	return m.v.nvals()
}

func (m diagonalMatrix[T]) setElement(ref *matrixReference[T], value T, row, col int) *matrixReference[T] {
	return newMatrixReference[T](newListMatrix[T](
		m.nrows, m.ncols, ref,
		&matrixValueList[T]{
			row:   row,
			col:   col,
			value: value,
		},
	), -1)
}

func (m diagonalMatrix[T]) removeElement(ref *matrixReference[T], row, col int) *matrixReference[T] {
	return newMatrixReference[T](newListMatrix[T](
		m.nrows, m.ncols, ref,
		&matrixValueList[T]{
			row: -row,
			col: -col,
		},
	), -1)
}

func (m diagonalMatrix[T]) extractElement(row, col int) (result T, ok bool) {
	i, j := diagonalOffset(m.k)
	index := row + i
	if index != col-j {
		return
	}
	return m.v.extractElement(index)
}

func (m diagonalMatrix[T]) getPipeline() *pipeline.Pipeline[any] {
	i, j := diagonalOffset(m.k)
	p := m.v.getPipeline()
	if p == nil {
		return nil
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[T])
			var newCols []int
			if slice.cow&cow0 != 0 {
				newCols = make([]int, len(slice.indices))
			} else {
				newCols = slice.indices
			}
			result := matrixSlice[T]{
				cow:    cowv,
				rows:   make([]int, len(slice.indices)),
				cols:   newCols,
				values: slice.values,
			}
			for k, index := range slice.indices {
				result.rows[k] = index - i
				result.cols[k] = index + j
			}
			return result
		})),
	)
	return p
}

func makeSingleElementPipeline[T any](v *vectorReference[T], accessIndex, reportedIndex int) *pipeline.Pipeline[any] {
	done := false
	var p pipeline.Pipeline[any]
	p.Source(pipeline.NewFunc[any](1, func(size int) (data any, fetched int, err error) {
		if done {
			return
		}
		done = true
		value, ok := v.extractElement(accessIndex)
		if ok {
			return vectorSlice[T]{
				indices: []int{reportedIndex},
				values:  []T{value},
			}, 1, nil
		}
		return
	}))
	return &p
}

func (m diagonalMatrix[T]) getRowPipeline(row int) *pipeline.Pipeline[any] {
	i, j := diagonalOffset(m.k)
	accessIndex := row + i
	reportedIndex := accessIndex + j
	return makeSingleElementPipeline(m.v, accessIndex, reportedIndex)
}

func (m diagonalMatrix[T]) getColPipeline(col int) *pipeline.Pipeline[any] {
	i, j := diagonalOffset(m.k)
	accessIndex := col - j
	reportedIndex := accessIndex - i
	return makeSingleElementPipeline(m.v, accessIndex, reportedIndex)
}

func (m diagonalMatrix[T]) getRowPipelines() []matrix1Pipeline {
	i, j := diagonalOffset(m.k)
	size := m.v.size()
	result := make([]matrix1Pipeline, size)
	for index := range result {
		result[index].index = index - i
		result[index].p = makeSingleElementPipeline(m.v, index, index+j)
	}
	return result
}

func (m diagonalMatrix[T]) getColPipelines() []matrix1Pipeline {
	i, j := diagonalOffset(m.k)
	size := m.v.size()
	result := make([]matrix1Pipeline, size)
	for index := range result {
		result[index].index = index + j
		result[index].p = makeSingleElementPipeline(m.v, index, index-i)
	}
	return result
}

func (m diagonalMatrix[T]) optimized() bool {
	return m.v.optimized()
}

func (m diagonalMatrix[T]) optimize() functionalMatrix[T] {
	m.v.optimize()
	return m
}
