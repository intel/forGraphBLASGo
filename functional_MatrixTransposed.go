package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/pipeline"
	"runtime"
	"sync/atomic"
)

type transposedMatrix[T any] struct {
	nrows, ncols int
	base         *matrixReference[T]
}

func newTransposedMatrix[T any](ref *matrixReference[T]) *matrixReference[T] {
	if baseReferent, ok := ref.get().(transposedMatrix[T]); ok {
		return baseReferent.base
	}
	ncols, nrows := ref.size()
	n := atomic.LoadInt64(&ref.nvalues)
	return newMatrixReference[T](transposedMatrix[T]{
		nrows: nrows,
		ncols: ncols,
		base:  ref,
	}, n)
}

func newTransposedMatrixRaw[T any](nrows, ncols int, base *matrixReference[T]) transposedMatrix[T] {
	return transposedMatrix[T]{nrows: nrows, ncols: ncols, base: base}
}

func (m transposedMatrix[T]) resize(_ *matrixReference[T], newNRows, newNCols int) *matrixReference[T] {
	return newTransposedMatrix[T](m.base.resize(newNCols, newNRows))
}

func (m transposedMatrix[T]) size() (nrows, ncols int) {
	return m.nrows, m.ncols
}

func (m transposedMatrix[T]) nvals() int {
	return m.base.nvals()
}

func (m transposedMatrix[T]) setElement(_ *matrixReference[T], value T, row, col int) *matrixReference[T] {
	return newMatrixReference[T](newTransposedMatrixRaw[T](m.nrows, m.ncols, m.base.setElement(value, col, row)), -1)
}

func (m transposedMatrix[T]) removeElement(_ *matrixReference[T], row, col int) *matrixReference[T] {
	return newMatrixReference[T](newTransposedMatrixRaw[T](m.nrows, m.ncols, m.base.removeElement(col, row)), -1)
}

func (m transposedMatrix[T]) extractElement(row, col int) (T, bool) {
	return m.base.extractElement(col, row)
}

func transposeMatrixPipeline[T any](base *matrixReference[T]) *pipeline.Pipeline[any] {
	colPipelines := base.getColPipelines()
	ch := make(chan any, runtime.GOMAXPROCS(0))
	var np pipeline.Pipeline[any]
	np.Source(pipeline.NewChan(ch))
	np.Notify(func() {
		for pi, p := range colPipelines {
			p.p.Add(
				pipeline.Par(pipeline.Receive(func(_ int, data any) any {
					slice := data.(vectorSlice[T])
					ncow := slice.cow & cowv
					if slice.cow&cow0 != 0 {
						ncow |= cow1
					}
					result := matrixSlice[T]{
						cow:    ncow,
						rows:   make([]int, len(slice.values)),
						cols:   slice.indices,
						values: slice.values,
					}
					for i := range result.rows {
						result.rows[i] = p.index
					}
					return result
				})),
				pipeline.Ord(pipeline.Receive(func(_ int, data any) any {
					select {
					case <-p.p.Context().Done():
					case <-np.Context().Done():
					case ch <- data:
					}
					return nil
				})),
			)
			p.p.Run()
			if err := p.p.Err(); err != nil {
				panic(err)
			}
			colPipelines[pi].p = nil
		}
		close(ch)
	})
	return &np
}

func (m transposedMatrix[T]) getPipeline() *pipeline.Pipeline[any] {
	return transposeMatrixPipeline(m.base)
}

func (m transposedMatrix[T]) getRowPipeline(row int) *pipeline.Pipeline[any] {
	return m.base.getColPipeline(row)
}

func (m transposedMatrix[T]) getColPipeline(col int) *pipeline.Pipeline[any] {
	return m.base.getRowPipeline(col)
}

func (m transposedMatrix[T]) getRowPipelines() []matrix1Pipeline {
	return m.base.getColPipelines()
}

func (m transposedMatrix[T]) getColPipelines() []matrix1Pipeline {
	return m.base.getRowPipelines()
}

func (m transposedMatrix[T]) optimized() bool {
	return m.base.optimized()
}

func (m transposedMatrix[T]) optimize() functionalMatrix[T] {
	m.base.optimize()
	return m
}
