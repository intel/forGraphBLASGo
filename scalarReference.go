package forGraphBLASGo

import "sync"

type scalarReference[T any] struct {
	mutex    sync.RWMutex
	referent functionalScalar[T]
}

func newScalarReference[T any](referent functionalScalar[T]) *scalarReference[T] {
	return &scalarReference[T]{referent: referent}
}

// todo: find possible uses
func newDelayedScalarReference[T any](make func() (referent functionalScalar[T])) *scalarReference[T] {
	ref := new(scalarReference[T])
	ref.mutex.Lock()
	go func() {
		defer ref.mutex.Unlock()
		ref.referent = make()
	}()
	return ref
}

func (s *scalarReference[T]) get() functionalScalar[T] {
	s.mutex.RLock()
	referent := s.referent
	s.mutex.RUnlock()
	return referent
}

func (s *scalarReference[T]) extractElement() (T, bool) {
	return s.get().extractElement(s)
}

func (s *scalarReference[T]) optimized() bool {
	return s.get().optimized()
}

func (s *scalarReference[T]) optimize() functionalScalar[T] {
	var referent functionalScalar[T]
	s.mutex.RLock()
	referent = s.referent
	s.mutex.RUnlock()
	if referent.optimized() {
		return referent
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	referent = s.referent
	if referent.optimized() {
		return referent
	}
	s.referent = s.referent.optimize()
	return s.referent
}
