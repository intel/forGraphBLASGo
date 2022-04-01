package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/parallel"
)

type (
	computeScalarT[D any] interface {
		computeElement() (D, bool)
	}

	computedScalar[D any] struct {
		computeElement func() (D, bool)
	}
)

func makeScalarComputeElement[D any](
	s *scalarReference[D],
	accum BinaryOp[D, D, D],
	compute computeScalarT[D],
) func() (D, bool) {
	if accum == nil {
		return compute.computeElement
	}
	return func() (D, bool) {
		var sValue, tValue D
		var sok, tok bool
		parallel.Do(func() {
			sValue, sok = s.extractElement()
		}, func() {
			tValue, tok = compute.computeElement()
		})
		if sok {
			if tok {
				return accum(sValue, tValue), true
			}
			return sValue, true
		}
		return tValue, tok
	}
}

func newComputedScalar[D any](
	s *scalarReference[D],
	accum BinaryOp[D, D, D],
	compute computeScalarT[D],
) *computedScalar[D] {
	return &computedScalar[D]{
		computeElement: makeScalarComputeElement[D](s, accum, compute),
	}
}

func (s *computedScalar[D]) extractElement(ref *scalarReference[D]) (result D, ok bool) {
	return ref.optimize().extractElement(ref)
}

func (_ *computedScalar[D]) optimized() bool {
	return false
}

func (_ *computedScalar[D]) valid() bool {
	panic("valid should only be called on optimized scalars")
}

func (s *computedScalar[D]) optimize() functionalScalar[D] {
	if value, ok := s.computeElement(); ok {
		return newFullScalar[D](value)
	}
	return newEmptyScalar[D]()
}
