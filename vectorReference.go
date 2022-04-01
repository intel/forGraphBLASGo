package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/pipeline"
	"sync"
	"sync/atomic"
)

// vectorReference is the second level of indirection for the Vector type.
// See comments in the header of api_Vector.go (actually api_Matrix.go) for more details.

type vectorReference[T any] struct {
	mutex    sync.RWMutex
	referent functionalVector[T]
	nvalues  int64
}

func newVectorReference[T any](referent functionalVector[T], nvalues int64) *vectorReference[T] {
	return &vectorReference[T]{referent: referent, nvalues: nvalues}
}

// todo: find possible uses
func newDelayedVectorReference[T any](make func() (referent functionalVector[T], nvalues int64)) *vectorReference[T] {
	ref := new(vectorReference[T])
	ref.mutex.Lock()
	go func() {
		defer ref.mutex.Unlock()
		ref.referent, ref.nvalues = make()
	}()
	return ref
}

func (v *vectorReference[T]) get() (referent functionalVector[T]) {
	v.mutex.RLock()
	referent = v.referent
	v.mutex.RUnlock()
	return
}

func (v *vectorReference[T]) resize(newSize int) *vectorReference[T] {
	return v.get().resize(v, newSize)
}

func (v *vectorReference[T]) size() int {
	return v.get().size()
}

func (v *vectorReference[T]) nvals() int {
	if n := atomic.LoadInt64(&v.nvalues); n >= 0 {
		return int(n)
	}
	n := v.get().nvals()
	atomic.StoreInt64(&v.nvalues, int64(n))
	return n
}

func (v *vectorReference[T]) setElement(value T, index int) *vectorReference[T] {
	return v.get().setElement(v, value, index)
}

func (v *vectorReference[T]) removeElement(index int) *vectorReference[T] {
	return v.get().removeElement(v, index)
}

func (v *vectorReference[T]) extractElement(index int) (T, bool) {
	return v.get().extractElement(index)
}

func (v *vectorReference[T]) getPipeline() *pipeline.Pipeline[any] {
	return v.get().getPipeline()
}

func (v *vectorReference[T]) optimized() bool {
	return v.get().optimized()
}

func (v *vectorReference[T]) optimize() {
	var referent functionalVector[T]
	v.mutex.RLock()
	referent = v.referent
	v.mutex.RUnlock()
	if referent.optimized() {
		return
	}
	v.mutex.Lock()
	defer v.mutex.Unlock()
	referent = v.referent
	if referent.optimized() {
		return
	}
	v.referent = v.referent.optimize()
}
