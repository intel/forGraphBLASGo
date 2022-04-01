package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/pipeline"
)

type matrixAsMask[T Number] struct {
	v *matrixReference[T]
}

func newMatrixAsMask[T Number](v *matrixReference[T]) functionalMatrix[bool] {
	return matrixAsMask[T]{v: v}
}

func (matrix matrixAsMask[T]) resize(ref *matrixReference[bool], newNRows, newNCols int) *matrixReference[bool] {
	nrows, ncols := matrix.v.size()
	if newNRows == nrows && newNCols == ncols {
		return ref
	}
	return newMatrixReference[bool](newMatrixAsMask(matrix.v.resize(newNRows, newNCols)), -1)
}

func (matrix matrixAsMask[T]) size() (nrows, ncols int) {
	return matrix.v.size()
}

func (matrix matrixAsMask[T]) nvals() int {
	return matrix.v.nvals()
}

func (matrix matrixAsMask[T]) setElement(ref *matrixReference[bool], value bool, row, col int) *matrixReference[bool] {
	nrows, ncols := matrix.v.size()
	return newMatrixReference[bool](newListMatrix[bool](nrows, ncols, ref,
		&matrixValueList[bool]{
			row:   row,
			col:   col,
			value: value,
		},
	), -1)
}

func (matrix matrixAsMask[T]) removeElement(ref *matrixReference[bool], row, col int) *matrixReference[bool] {
	nrows, ncols := matrix.v.size()
	return newMatrixReference[bool](newListMatrix[bool](nrows, ncols, ref,
		&matrixValueList[bool]{
			row: -row,
			col: -col,
		},
	), -1)
}

func (matrix matrixAsMask[T]) extractElement(row, col int) (bool, bool) {
	if value, ok := matrix.v.extractElement(row, col); ok {
		return value != 0, true
	}
	return false, false
}

func (matrix matrixAsMask[T]) getPipeline() *pipeline.Pipeline[any] {
	p := matrix.v.getPipeline()
	if p == nil {
		return nil
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(matrixSlice[T])
			result := matrixSlice[bool]{
				cow:    slice.cow &^ cowv,
				rows:   slice.rows,
				cols:   slice.cols,
				values: make([]bool, len(slice.values)),
			}
			for i, value := range slice.values {
				result.values[i] = value != 0
			}
			return result
		})),
	)
	return p
}

func adaptAsMaskPipeline[T Number](p *pipeline.Pipeline[any]) {
	if p == nil {
		return
	}
	p.Add(
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
}

func (matrix matrixAsMask[T]) getRowPipeline(row int) *pipeline.Pipeline[any] {
	p := matrix.v.getRowPipeline(row)
	adaptAsMaskPipeline[T](p)
	return p
}

func (matrix matrixAsMask[T]) getColPipeline(col int) *pipeline.Pipeline[any] {
	p := matrix.v.getColPipeline(col)
	adaptAsMaskPipeline[T](p)
	return p
}

func (matrix matrixAsMask[T]) getRowPipelines() []matrix1Pipeline {
	ps := matrix.v.getRowPipelines()
	for _, p := range ps {
		adaptAsMaskPipeline[T](p.p)
	}
	return ps
}

func (matrix matrixAsMask[T]) getColPipelines() []matrix1Pipeline {
	ps := matrix.v.getColPipelines()
	for _, p := range ps {
		adaptAsMaskPipeline[T](p.p)
	}
	return ps
}

func (matrix matrixAsMask[T]) optimized() bool {
	return matrix.v.optimized()
}

func (matrix matrixAsMask[T]) optimize() functionalMatrix[bool] {
	matrix.v.optimize()
	return matrix
}

type matrixAsStructuralMask[T any] struct {
	v *matrixReference[T]
}

func newMatrixAsStructuralMask[T any](v *matrixReference[T]) functionalMatrix[bool] {
	return matrixAsStructuralMask[T]{v: v}
}

func (v matrixAsStructuralMask[T]) resize(ref *matrixReference[bool], newNRows, newNCols int) *matrixReference[bool] {
	nrows, ncols := v.v.size()
	if newNRows == nrows && newNCols == ncols {
		return ref
	}
	return newMatrixReference[bool](newMatrixAsStructuralMask(v.v.resize(newNRows, newNCols)), -1)
}

func (v matrixAsStructuralMask[T]) size() (nrows, ncols int) {
	return v.v.size()
}

func (v matrixAsStructuralMask[T]) nvals() int {
	return v.v.nvals()
}

func (v matrixAsStructuralMask[T]) setElement(ref *matrixReference[bool], value bool, row, col int) *matrixReference[bool] {
	nrows, ncols := v.v.size()
	return newMatrixReference[bool](newListMatrix[bool](nrows, ncols, ref,
		&matrixValueList[bool]{
			row:   row,
			col:   col,
			value: value,
		},
	), -1)
}

func (v matrixAsStructuralMask[T]) removeElement(ref *matrixReference[bool], row, col int) *matrixReference[bool] {
	nrows, ncols := v.v.size()
	return newMatrixReference[bool](newListMatrix[bool](nrows, ncols, ref,
		&matrixValueList[bool]{
			row: -row,
			col: -col,
		},
	), -1)
}

func (v matrixAsStructuralMask[T]) extractElement(row, col int) (bool, bool) {
	// todo: when accessing the first return value, this should normally
	// 	result in a DomainMismatch, so check the program flow for this
	_, ok := v.v.extractElement(row, col)
	return ok, ok
}

func (matrix matrixAsStructuralMask[T]) getPipeline() *pipeline.Pipeline[any] {
	p := matrix.v.getPipeline()
	if p == nil {
		return nil
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(matrixSlice[T])
			result := matrixSlice[bool]{
				cow:    slice.cow &^ cowv,
				rows:   slice.rows,
				cols:   slice.cols,
				values: make([]bool, len(slice.values)),
			}
			for i := range result.values {
				result.values[i] = true
			}
			return result
		})),
	)
	return p
}

func adaptAsStructuralMaskPipeline[T any](p *pipeline.Pipeline[any]) {
	if p == nil {
		return
	}
	p.Add(
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
}

func (matrix matrixAsStructuralMask[T]) getRowPipeline(row int) *pipeline.Pipeline[any] {
	p := matrix.v.getRowPipeline(row)
	adaptAsStructuralMaskPipeline[T](p)
	return p
}

func (matrix matrixAsStructuralMask[T]) getColPipeline(col int) *pipeline.Pipeline[any] {
	p := matrix.v.getRowPipeline(col)
	adaptAsStructuralMaskPipeline[T](p)
	return p
}

func (matrix matrixAsStructuralMask[T]) getRowPipelines() []matrix1Pipeline {
	ps := matrix.v.getRowPipelines()
	for _, p := range ps {
		adaptAsStructuralMaskPipeline[T](p.p)
	}
	return ps
}

func (matrix matrixAsStructuralMask[T]) getColPipelines() []matrix1Pipeline {
	ps := matrix.v.getColPipelines()
	for _, p := range ps {
		adaptAsStructuralMaskPipeline[T](p.p)
	}
	return ps
}

func (v matrixAsStructuralMask[T]) optimized() bool {
	return v.v.optimized()
}

func (v matrixAsStructuralMask[T]) optimize() functionalMatrix[bool] {
	v.v.optimize()
	return v
}
