package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/pipeline"
)

// todo: this needs to become isoMatrixConstant
type homMatrixConstant[T any] struct {
	nrows, ncols int
	value        T
}

func newHomMatrixConstant[T any](nrows, ncols int, value T) homMatrixConstant[T] {
	return homMatrixConstant[T]{nrows: nrows, ncols: ncols, value: value}
}

func (m homMatrixConstant[T]) resize(ref *matrixReference[T], newNRows, newNCols int) *matrixReference[T] {
	if m.nrows == newNRows && m.ncols == newNCols {
		return ref
	}
	return newMatrixReference[T](newHomMatrixConstant[T](newNRows, newNCols, m.value), int64(newNRows*newNCols))
}

func (m homMatrixConstant[T]) size() (nrows, ncols int) {
	return m.nrows, m.ncols
}

func (m homMatrixConstant[T]) nvals() int {
	return m.nrows * m.ncols
}

func (m homMatrixConstant[T]) setElement(ref *matrixReference[T], value T, row, col int) *matrixReference[T] {
	if equal(m.value, value) {
		return ref
	}
	return newMatrixReference[T](newListMatrix[T](
		m.nrows, m.ncols, ref,
		&matrixValueList[T]{
			row:   row,
			col:   col,
			value: value,
		}),
		int64(m.nrows*m.ncols),
	)
}

func (m homMatrixConstant[T]) removeElement(ref *matrixReference[T], row, col int) *matrixReference[T] {
	return newMatrixReference[T](newListMatrix[T](
		m.nrows, m.ncols, ref,
		&matrixValueList[T]{
			row: -row,
			col: -col,
		}),
		int64(m.nrows*m.ncols-1),
	)
}

func (m homMatrixConstant[T]) extractElement(row, col int) (T, bool) {
	return m.value, true
}

func (m homMatrixConstant[T]) getPipeline() *pipeline.Pipeline[any] {
	var p pipeline.Pipeline[any]
	total := m.nrows * m.ncols
	index := 0
	var values []T
	p.Source(pipeline.NewFunc[any](total, func(size int) (data any, fetched int, err error) {
		if index >= total {
			return
		}
		if index+size > total {
			size = total - index
		}
		if size < len(values) {
			values = values[:size]
		} else {
			for len(values) < size {
				values = append(values, m.value)
			}
		}
		result := matrixSlice[T]{
			cow:    cowv,
			rows:   make([]int, 0, size),
			cols:   make([]int, 0, size),
			values: values,
		}
		for i := 0; i < size; i++ {
			row, col := indexToCoord(index+i, m.nrows, m.ncols)
			result.rows = append(result.rows, row)
			result.cols = append(result.cols, col)
		}
		index += size
		return result, size, nil
	}))
	return &p
}

func (m homMatrixConstant[T]) getSizedVectorPipeline(total int) *pipeline.Pipeline[any] {
	var p pipeline.Pipeline[any]
	index := 0
	var values []T
	p.Source(pipeline.NewFunc[any](total, func(size int) (data any, fetched int, err error) {
		if index >= total {
			return
		}
		if index+size > total {
			size = total - index
		}
		if size < len(values) {
			values = values[:size]
		} else {
			for len(values) < size {
				values = append(values, m.value)
			}
		}
		result := vectorSlice[T]{
			cow:     cowv,
			indices: make([]int, 0, size),
			values:  values,
		}
		for i := 0; i < size; i++ {
			result.indices = append(result.indices, index+i)
		}
		index += size
		return result, size, nil
	}))
	return &p
}

func (m homMatrixConstant[T]) getRowPipeline(_ int) *pipeline.Pipeline[any] {
	return m.getSizedVectorPipeline(m.ncols)
}

func (m homMatrixConstant[T]) getColPipeline(_ int) *pipeline.Pipeline[any] {
	return m.getSizedVectorPipeline(m.nrows)
}

func (m homMatrixConstant[T]) getRowPipelines() (result []matrix1Pipeline) {
	result = make([]matrix1Pipeline, m.nrows)
	for i := range result {
		result[i].p = m.getSizedVectorPipeline(m.ncols)
		result[i].index = i
	}
	return
}

func (m homMatrixConstant[T]) getColPipelines() (result []matrix1Pipeline) {
	result = make([]matrix1Pipeline, m.ncols)
	for i := range result {
		result[i].p = m.getSizedVectorPipeline(m.nrows)
		result[i].index = i
	}
	return
}

func (_ homMatrixConstant[T]) optimized() bool {
	return true
}

func (m homMatrixConstant[T]) optimize() functionalMatrix[T] {
	return m
}

// todo: this needs to become isoMatrixScalar
type homMatrixScalar[T any] struct {
	nrows, ncols int
	value        *scalarReference[T]
}

func newHomMatrixScalar[T any](nrows, ncols int, value *scalarReference[T]) functionalMatrix[T] {
	v := value.get()
	if v.optimized() {
		if val, ok := v.extractElement(value); ok {
			return newHomMatrixConstant[T](nrows, ncols, val)
		}
		return makeEmptyMatrix[T](nrows, ncols)
	}
	return homMatrixScalar[T]{nrows: nrows, ncols: ncols, value: value}
}

func (m homMatrixScalar[T]) resize(ref *matrixReference[T], newNRows, newNCols int) *matrixReference[T] {
	if m.nrows == newNRows && m.ncols == newNCols {
		return ref
	}
	return newMatrixReference[T](newHomMatrixScalar[T](newNRows, newNCols, m.value), -1)
}

func (m homMatrixScalar[T]) size() (nrows, ncols int) {
	return m.nrows, m.ncols
}

func (m homMatrixScalar[T]) nvals() int {
	if _, ok := m.value.extractElement(); ok {
		return m.nrows * m.ncols
	}
	return 0
}

func (m homMatrixScalar[T]) setElement(ref *matrixReference[T], value T, row, col int) *matrixReference[T] {
	v := m.value.get()
	if v.optimized() {
		if val, ok := v.extractElement(m.value); ok {
			if equal(val, value) {
				return ref
			}
			return newMatrixReference[T](newListMatrix[T](
				m.nrows, m.ncols, ref,
				&matrixValueList[T]{
					row:   row,
					col:   col,
					value: value,
				}),
				int64(m.nrows*m.ncols),
			)
		} else {
			return newMatrixReference[T](makeSingletonMatrix[T](m.nrows, m.ncols, row, col, value), 1)
		}
	}
	return newMatrixReference[T](newListMatrix[T](
		m.nrows, m.ncols, ref,
		&matrixValueList[T]{
			row:   row,
			col:   col,
			value: value,
		}), -1)
}

func (m homMatrixScalar[T]) removeElement(ref *matrixReference[T], row, col int) *matrixReference[T] {
	v := m.value.get()
	if v.optimized() && !v.valid() {
		return ref
	}
	return newMatrixReference[T](newListMatrix[T](
		m.nrows, m.ncols, ref,
		&matrixValueList[T]{
			row: -row,
			col: -col,
		}), -1)
}

func (m homMatrixScalar[T]) extractElement(_, _ int) (T, bool) {
	return m.value.extractElement()
}

func (m homMatrixScalar[T]) getPipeline() *pipeline.Pipeline[any] {
	var p pipeline.Pipeline[any]
	total := m.nrows * m.ncols
	index := 0
	var value T
	var ok, extracted bool
	var values []T
	p.Source(pipeline.NewFunc[any](total, func(size int) (data any, fetched int, err error) {
		if !extracted {
			value, ok = m.value.extractElement()
			extracted = true
		}
		if !ok || index >= total {
			return
		}
		if index+size > total {
			size = total - index
		}
		if size < len(values) {
			values = values[:size]
		} else {
			for len(values) < size {
				values = append(values, value)
			}
		}
		result := matrixSlice[T]{
			cow:    cowv,
			rows:   make([]int, 0, size),
			cols:   make([]int, 0, size),
			values: values,
		}
		for i := 0; i < size; i++ {
			row, col := indexToCoord(index+i, m.nrows, m.ncols)
			result.rows = append(result.rows, row)
			result.cols = append(result.cols, col)
		}
		index += size
		return result, size, nil
	}))
	return &p
}

func (m homMatrixScalar[T]) getSizedVectorPipeline(total int) *pipeline.Pipeline[any] {
	var p pipeline.Pipeline[any]
	index := 0
	var value T
	var ok, extracted bool
	var values []T
	p.Source(pipeline.NewFunc[any](total, func(size int) (data any, fetched int, err error) {
		if !extracted {
			value, ok = m.value.extractElement()
			extracted = true
		}
		if !ok || index >= total {
			return
		}
		if index+size > total {
			size = total - index
		}
		if size < len(values) {
			values = values[:size]
		} else {
			for len(values) < size {
				values = append(values, value)
			}
		}
		result := vectorSlice[T]{
			cow:     cowv,
			indices: make([]int, 0, size),
			values:  values,
		}
		for i := 0; i < size; i++ {
			result.indices = append(result.indices, index+i)
		}
		index += size
		return result, size, nil
	}))
	return &p
}

func (m homMatrixScalar[T]) getRowPipeline(_ int) *pipeline.Pipeline[any] {
	return m.getSizedVectorPipeline(m.ncols)
}

func (m homMatrixScalar[T]) getColPipeline(_ int) *pipeline.Pipeline[any] {
	return m.getSizedVectorPipeline(m.nrows)
}

func (m homMatrixScalar[T]) getRowPipelines() (result []matrix1Pipeline) {
	result = make([]matrix1Pipeline, m.nrows)
	for i := range result {
		result[i].p = m.getSizedVectorPipeline(m.ncols)
		result[i].index = i
	}
	return
}

func (m homMatrixScalar[T]) getColPipelines() (result []matrix1Pipeline) {
	result = make([]matrix1Pipeline, m.ncols)
	for i := range result {
		result[i].p = m.getSizedVectorPipeline(m.nrows)
		result[i].index = i
	}
	return
}

func (m homMatrixScalar[T]) optimized() bool {
	return m.value.optimized()
}

func (m homMatrixScalar[T]) optimize() functionalMatrix[T] {
	m.value.optimize()
	return m
}
