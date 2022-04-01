package forGraphBLASGo

import "github.com/intel/forGoParallel/pipeline"

// A functionalMatrix does not modify its contents, but always returns a new matrix as a result.

type (
	matrix1Pipeline struct {
		index int // row or col
		p     *pipeline.Pipeline[any]
	}

	functionalMatrix[T any] interface {
		resize(ref *matrixReference[T], newNRows, newNCols int) *matrixReference[T]
		size() (nrows, ncols int)
		nvals() int
		setElement(ref *matrixReference[T], value T, row, col int) *matrixReference[T]
		removeElement(ref *matrixReference[T], row, col int) *matrixReference[T]
		extractElement(row, col int) (T, bool)
		getPipeline() *pipeline.Pipeline[any]
		getRowPipeline(row int) *pipeline.Pipeline[any]
		getColPipeline(col int) *pipeline.Pipeline[any]
		getRowPipelines() []matrix1Pipeline
		getColPipelines() []matrix1Pipeline

		optimized() bool
		optimize() functionalMatrix[T]
	}

	homMatrix[T any] interface {
		functionalMatrix[T]
		homValue() (T, bool)
	}
)

func setMatrixElement[T any](
	m functionalMatrix[T],
	ref *matrixReference[T],
	nvals int,
	value T,
	row, col int,
) *matrixReference[T] {
	nrows, ncols := m.size()
	if nvals == 0 {
		return newMatrixReference[T](makeSingletonMatrix[T](nrows, ncols, row, col, value), 1)
	}
	return newMatrixReference[T](newListMatrix[T](
		nrows, ncols, ref,
		&matrixValueList[T]{
			row:   row,
			col:   col,
			value: value,
		},
	), -1)
}

func removeMatrixElement[T any](
	m functionalMatrix[T],
	ref *matrixReference[T],
	nvals int,
	row, col int,
	selfAsEmpty bool,
) *matrixReference[T] {
	nrows, ncols := m.size()
	if nvals == 0 {
		if selfAsEmpty {
			return ref
		}
		return newMatrixReference[T](makeEmptyMatrix[T](nrows, ncols), 0)
	}
	return newMatrixReference[T](newListMatrix[T](
		nrows, ncols, ref,
		&matrixValueList[T]{
			row: -row,
			col: -col,
		}), -1)
}
