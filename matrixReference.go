package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/pipeline"
	"sync"
	"sync/atomic"
)

// matrixReference is the second level of indirection for the Matrix type.
// See comments in the header of api_Matrix.go for more details.

type matrixReference[T any] struct {
	mutex    sync.RWMutex
	referent functionalMatrix[T]
	nvalues  int64
}

func newMatrixReference[T any](referent functionalMatrix[T], nvalues int64) *matrixReference[T] {
	return &matrixReference[T]{referent: referent, nvalues: nvalues}
}

// todo: find possible uses
func newDelayedMatrixReference[T any](make func() (referent functionalMatrix[T], nvalues int64)) *matrixReference[T] {
	ref := new(matrixReference[T])
	ref.mutex.Lock()
	go func() {
		defer ref.mutex.Unlock()
		ref.referent, ref.nvalues = make()
	}()
	return ref
}

func (m *matrixReference[T]) get() (referent functionalMatrix[T]) {
	m.mutex.RLock()
	referent = m.referent
	m.mutex.RUnlock()
	return
}

func (m *matrixReference[T]) resize(newNRows, newNCols int) *matrixReference[T] {
	return m.get().resize(m, newNRows, newNCols)
}

func (m *matrixReference[T]) size() (nrows, ncols int) {
	return m.get().size()
}

func (m *matrixReference[T]) nvals() int {
	if n := atomic.LoadInt64(&m.nvalues); n >= 0 {
		return int(n)
	}
	n := m.get().nvals()
	atomic.StoreInt64(&m.nvalues, int64(n))
	return n
}

func (m *matrixReference[T]) setElement(value T, row, col int) *matrixReference[T] {
	return m.get().setElement(m, value, row, col)
}

func (m *matrixReference[T]) removeElement(row, col int) *matrixReference[T] {
	return m.get().removeElement(m, row, col)
}

func (m *matrixReference[T]) extractElement(row, col int) (T, bool) {
	return m.get().extractElement(row, col)
}

func (m *matrixReference[T]) getPipeline() *pipeline.Pipeline[any] {
	return m.get().getPipeline()
}

func (m *matrixReference[T]) getRowPipeline(row int) *pipeline.Pipeline[any] {
	return m.get().getRowPipeline(row)
}

func (m *matrixReference[T]) getColPipeline(col int) *pipeline.Pipeline[any] {
	return m.get().getColPipeline(col)
}

func (m *matrixReference[T]) getRowPipelines() []matrix1Pipeline {
	return m.get().getRowPipelines()
}

func (m *matrixReference[T]) getColPipelines() []matrix1Pipeline {
	return m.get().getColPipelines()
}

func (m *matrixReference[T]) optimized() bool {
	return m.get().optimized()
}

func (m *matrixReference[T]) optimize() {
	var referent functionalMatrix[T]
	m.mutex.RLock()
	referent = m.referent
	m.mutex.RUnlock()
	if referent.optimized() {
		return
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	referent = m.referent
	if referent.optimized() {
		return
	}
	m.referent = referent.optimize()
}
