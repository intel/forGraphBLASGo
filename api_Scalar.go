package forGraphBLASGo

type Scalar[T any] struct {
	ref *scalarReference[T]
}

func ScalarNew[T any]() (result *Scalar[T], err error) {
	return &Scalar[T]{newScalarReference[T](newEmptyScalar[T]())}, nil
}

func (s *Scalar[T]) Dup() (result *Scalar[T], err error) {
	if s == nil || s.ref == nil {
		err = UninitializedObject
		return
	}
	return &Scalar[T]{s.ref}, nil
}

func (s *Scalar[T]) Clear() error {
	if s == nil || s.ref == nil {
		return UninitializedObject
	}
	s.ref = newScalarReference[T](newEmptyScalar[T]())
	return nil
}

func (s *Scalar[T]) NVals() (int, error) {
	if s == nil || s.ref == nil {
		return 0, UninitializedObject
	}
	_, valid := s.ref.extractElement()
	if valid {
		return 1, nil
	}
	return 0, nil
}

func (s *Scalar[T]) SetElement(value T) error {
	if s == nil || s.ref == nil {
		return UninitializedObject
	}
	s.ref = newScalarReference[T](newFullScalar[T](value))
	return nil
}

func (s *Scalar[T]) ExtractElement() (value T, err error) {
	if s == nil || s.ref == nil {
		err = UninitializedObject
		return
	}
	if value, ok := s.ref.extractElement(); ok {
		return value, nil
	}
	err = NoValue
	return
}

func (s *Scalar[T]) Wait(mode WaitMode) error {
	if s == nil || s.ref == nil {
		return UninitializedObject
	}
	if mode == Complete {
		return nil
	}
	s.ref.optimize()
	return nil
}
