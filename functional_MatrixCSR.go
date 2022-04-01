package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/parallel"
	"github.com/intel/forGoParallel/pipeline"
	"sort"
)

type csrMatrix[T any] struct {
	nrows, ncols         int
	rows, rowSpans, cols []int
	values               []T
}

func newCSRMatrix[T any](nrows, ncols int, rows, rowSpans, cols []int, values []T) csrMatrix[T] {
	return csrMatrix[T]{
		nrows:    nrows,
		ncols:    ncols,
		rows:     rows,
		rowSpans: rowSpans,
		cols:     cols,
		values:   values,
	}
}

func makeEmptyMatrix[T any](nrows, ncols int) functionalMatrix[T] {
	return newCSRMatrix[T](nrows, ncols, nil, []int{0}, nil, nil)
}

func makeSingletonMatrix[T any](nrows, ncols, row, col int, value T) csrMatrix[T] {
	return newCSRMatrix[T](nrows, ncols, []int{row}, []int{0, 1}, []int{col}, []T{value})
}

// todo: Resize seems expensive; do we need a resizeMatrix struct?
func (matrix csrMatrix[T]) resize(ref *matrixReference[T], newNRows, newNCols int) *matrixReference[T] {
	if newNRows == matrix.nrows && newNCols == matrix.ncols {
		return ref
	}
	if newNRows < matrix.nrows {
		newRowsSize := sort.SearchInts(matrix.rows, newNRows)
		matrix.rows = matrix.rows[:newRowsSize]
		matrix.rowSpans = matrix.rowSpans[:newRowsSize+1]
		nnz := matrix.rowSpans[len(matrix.rows)]
		matrix.cols = matrix.cols[:nnz]
		matrix.values = matrix.values[:nnz]
	}
	if newNCols < matrix.ncols {
		// todo: use speculative.RangeOr
		if parallel.RangeOr(0, len(matrix.cols), func(low, high int) bool {
			for i := low; i < high; i++ {
				if matrix.cols[i] >= newNCols {
					return true
				}
			}
			return false
		}) {
			newRows := make([]int, 0, len(matrix.rows))
			newRowSpans := make([]int, 0, len(matrix.rowSpans))
			newRowSpans = append(newRowSpans, 0)
			newCols := make([]int, 0, len(matrix.cols)-1)
			newValues := make([]T, 0, len(matrix.values)-1)
			for i, row := range matrix.rows {
				rowStart, rowEnd := matrix.rowSpans[i], matrix.rowSpans[i+1]
				rowNNZ := 0
				for j, col := range matrix.cols[rowStart:rowEnd] {
					if col < newNCols {
						rowNNZ++
						newCols = append(newCols, col)
						newValues = append(newValues, matrix.values[rowStart+j])
					}
				}
				if rowNNZ > 0 {
					newRows = append(newRows, row)
					newRowSpans = append(newRowSpans, newRowSpans[len(newRowSpans)-1]+rowNNZ)
				}
			}
			return newMatrixReference[T](newCSRMatrix[T](newNRows, newNCols, newRows, newRowSpans, newCols, newValues), int64(len(newValues)))
		}
	}
	return newMatrixReference[T](newCSRMatrix[T](newNRows, newNCols, matrix.rows, matrix.rowSpans, matrix.cols, matrix.values), int64(len(matrix.values)))
}

func (matrix csrMatrix[T]) size() (nrows, ncols int) {
	return matrix.nrows, matrix.ncols
}

func (matrix csrMatrix[T]) nvals() int {
	return len(matrix.values)
}

func (matrix csrMatrix[T]) setElement(ref *matrixReference[T], value T, row, col int) *matrixReference[T] {
	return setMatrixElement[T](matrix, ref, len(matrix.values), value, row, col)
}

func (matrix csrMatrix[T]) removeElement(ref *matrixReference[T], row, col int) *matrixReference[T] {
	return removeMatrixElement[T](matrix, ref, len(matrix.values), row, col, true)
}

func (matrix csrMatrix[T]) extractElement(row, col int) (result T, ok bool) {
	i := sort.SearchInts(matrix.rows, row)
	if i >= len(matrix.rows) || matrix.rows[i] != row {
		return
	}
	rowStart, rowEnd := matrix.rowSpans[i], matrix.rowSpans[i+1]
	j := sort.SearchInts(matrix.cols[rowStart:rowEnd], col)
	if j >= rowEnd-rowStart || matrix.cols[rowStart+j] != col {
		return
	}
	return matrix.values[rowStart+j], true
}

func (matrix csrMatrix[T]) getPipeline() *pipeline.Pipeline[any] {
	if len(matrix.values) == 0 {
		return nil
	}
	var p pipeline.Pipeline[any]
	index := 0
	rowIndex := 0
	rowStart, rowEnd := matrix.rowSpans[0], matrix.rowSpans[1]
	p.Source(pipeline.NewFunc[any](len(matrix.values), func(size int) (data any, fetched int, err error) {
		if index >= len(matrix.values) {
			return
		}
		if index+size > len(matrix.values) {
			size = len(matrix.values) - index
		}
		var result matrixSlice[T]
		result.cow = cow1 | cowv
		result.rows = make([]int, 0, size)
		for i := 0; i < size; i++ {
			result.rows = append(result.rows, matrix.rows[rowIndex])
			if rowStart++; rowStart == rowEnd {
				if rowIndex++; rowIndex < len(matrix.rows) {
					rowStart, rowEnd = matrix.rowSpans[rowIndex], matrix.rowSpans[rowIndex+1]
				}
			}
		}
		result.cols = matrix.cols[index : index+size]
		result.values = matrix.values[index : index+size]
		index += size
		return result, size, nil
	}))
	return &p
}

func (matrix csrMatrix[T]) getRowPipeline(row int) *pipeline.Pipeline[any] {
	i := sort.SearchInts(matrix.rows, row)
	if i >= len(matrix.rows) || matrix.rows[i] != row {
		return nil
	}
	rowStart, rowEnd := matrix.rowSpans[i], matrix.rowSpans[i+1]
	var p pipeline.Pipeline[any]
	p.Source(pipeline.NewFunc[any](rowEnd-rowStart, func(size int) (data any, fetched int, err error) {
		if rowStart >= rowEnd {
			return
		}
		if rowStart+size > rowEnd {
			size = rowEnd - rowStart
		}
		result := vectorSlice[T]{
			cow:     cow0 | cowv,
			indices: matrix.cols[rowStart : rowStart+size],
			values:  matrix.values[rowStart : rowStart+size],
		}
		rowStart += size
		return result, size, nil
	}))
	return &p
}

func (matrix csrMatrix[T]) getColPipeline(col int) *pipeline.Pipeline[any] {
	if len(matrix.values) == 0 {
		return nil
	}
	rowIndex := 0
	var p pipeline.Pipeline[any]
	p.Source(pipeline.NewFunc[any](-1, func(size int) (data any, fetched int, err error) {
		if rowIndex >= len(matrix.rows) {
			return
		}
		var result vectorSlice[T]
		for fetched < size {
			rowStart, rowEnd := matrix.rowSpans[rowIndex], matrix.rowSpans[rowIndex+1]
			j := sort.SearchInts(matrix.cols[rowStart:rowEnd], col)
			if j < rowEnd-rowStart && matrix.cols[rowStart+j] == col {
				result.indices = append(result.indices, matrix.rows[rowIndex])
				result.values = append(result.values, matrix.values[rowStart+j])
				fetched++
			}
			if rowIndex++; rowIndex == len(matrix.rows) {
				return result, fetched, nil
			}
		}
		return result, fetched, nil
	}))
	return &p
}

func (matrix csrMatrix[T]) getRowPipelines() (result []matrix1Pipeline) {
	if len(matrix.rows) == 0 {
		return
	}
	result = make([]matrix1Pipeline, len(matrix.rows))
	for i := range result {
		row := matrix.rows[i]
		result[i].index = row
		rowStart, rowEnd := matrix.rowSpans[i], matrix.rowSpans[i+1]
		var p pipeline.Pipeline[any]
		p.Source(pipeline.NewFunc[any](rowEnd-rowStart, func(size int) (data any, fetched int, err error) {
			if rowStart >= rowEnd {
				return
			}
			if rowStart+size > rowEnd {
				size = rowEnd - rowStart
			}
			result := vectorSlice[T]{
				cow:     cow0 | cowv,
				indices: matrix.cols[rowStart : rowStart+size],
				values:  matrix.values[rowStart : rowStart+size],
			}
			rowStart += size
			return result, size, nil
		}))
		result[i].p = &p
	}
	return
}

func (matrix csrMatrix[T]) getColPipelines() (result []matrix1Pipeline) {
	if len(matrix.values) == 0 {
		return
	}
	columns := parallel.RangeReduce(0, len(matrix.rows), func(low, high int) (result vectorBitset) {
		result = newVectorBitset(matrix.ncols)
		for i := low; i < high; i++ {
			for _, col := range matrix.cols[matrix.rowSpans[i]:matrix.rowSpans[i+1]] {
				result.set(col)
			}
		}
		return result
	}, func(left, right vectorBitset) vectorBitset {
		left.or(right)
		return left
	}).toSlice()
	result = make([]matrix1Pipeline, len(columns))
	for i := range result {
		col := columns[i]
		result[i].index = col
		var p pipeline.Pipeline[any]
		rowIndex := 0
		p.Source(pipeline.NewFunc[any](-1, func(size int) (data any, fetched int, err error) {
			if rowIndex >= len(matrix.rows) {
				return
			}
			var result vectorSlice[T]
			for fetched < size {
				rowStart, rowEnd := matrix.rowSpans[rowIndex], matrix.rowSpans[rowIndex+1]
				j := sort.SearchInts(matrix.cols[rowStart:rowEnd], col)
				if j < rowEnd-rowStart && matrix.cols[rowStart+j] == col {
					result.indices = append(result.indices, matrix.rows[rowIndex])
					result.values = append(result.values, matrix.values[rowStart+j])
					fetched++
				}
				if rowIndex++; rowIndex == len(matrix.rows) {
					return result, fetched, nil
				}
			}
			return result, fetched, nil
		}))
		result[i].p = &p
	}
	return
}

func (_ csrMatrix[T]) optimized() bool {
	return true
}

func (matrix csrMatrix[T]) optimize() functionalMatrix[T] {
	return matrix
}
