package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/parallel"
	"github.com/intel/forGoParallel/pipeline"
	"sync/atomic"
)

// See comments in the header of api_Matrix.go, which also apply here.

type Vector[T any] struct {
	ref *vectorReference[T]
}

func VectorNew[T any](size int) (result *Vector[T], err error) {
	if size <= 0 {
		err = InvalidValue
		return
	}
	return &Vector[T]{newVectorReference[T](newSparseVector[T](size, nil, nil), 0)}, nil
}

func (v *Vector[T]) Dup() (result *Vector[T], err error) {
	if v == nil || v.ref == nil {
		err = UninitializedObject
		return
	}
	return &Vector[T]{v.ref}, nil
}

func (v *Vector[T]) Resize(newSize int) error {
	if newSize <= 0 {
		return InvalidValue
	}
	if v == nil || v.ref == nil {
		return UninitializedObject
	}
	v.ref = v.ref.resize(newSize)
	return nil
}

func (v *Vector[T]) Clear() error {
	if v == nil || v.ref == nil {
		return UninitializedObject
	}
	v.ref = newVectorReference[T](newSparseVector[T](v.ref.size(), nil, nil), 0)
	return nil
}

func (v *Vector[T]) Size() (int, error) {
	if v == nil || v.ref == nil {
		return 0, UninitializedObject
	}
	return v.ref.size(), nil
}

func (v *Vector[T]) NVals() (int, error) {
	if v == nil || v.ref == nil {
		return 0, UninitializedObject
	}
	return v.ref.nvals(), nil
}

func (v *Vector[T]) Build(indices []int, values []T, dup BinaryOp[T, T, T]) error {
	if v == nil || v.ref == nil {
		return UninitializedObject
	}
	if len(indices) != len(values) {
		return IndexOutOfBounds
	}
	if v.ref.nvals() > 0 {
		return OutputNotEmpty
	}
	size := v.ref.size()
	// todo: use speculative.RangeOr
	if parallel.RangeOr(0, len(indices), func(low, high int) bool {
		for i := low; i < high; i++ {
			index := indices[i]
			if index < 0 || index >= size {
				return true
			}
		}
		return false
	}) {
		return IndexOutOfBounds
	}
	indexCopies, valueCopies := fpcopy2(indices, values)
	if dup == nil {
		vectorSort(indexCopies, valueCopies)
		// todo: use speculative.RangeOr
		if parallel.RangeOr(0, len(indices), func(low, high int) bool {
			for i := low; i < high-1; i++ {
				if indexCopies[i] == indexCopies[i+1] {
					return true
				}
			}
			return high < len(indices) && indexCopies[high-1] == indexCopies[high]
		}) {
			return InvalidValue
		}
		v.ref = newVectorReference[T](newSparseVector[T](size, indexCopies, valueCopies), int64(len(valueCopies)))
		return nil
	}
	v.ref = newDelayedVectorReference[T](func() (functionalVector[T], int64) {
		vectorSort(indexCopies, valueCopies)
		var dups [][2]int
		var p pipeline.Pipeline[any]
		p.Source(newIntervalSource(len(valueCopies)))
		p.Add(
			pipeline.Par(pipeline.Receive(func(_ int, data any) any {
				batch := data.(interval)
				low, high := batch.start, batch.end
				var result [][2]int
				if low > 0 {
					low--
				}
				if high < len(indexCopies) {
					high++
				}
				for i := low; i < high; {
					index := indexCopies[i]
					j := i + 1
					for j < high && indexCopies[j] == index {
						j++
					}
					if j-i > 1 {
						result = append(result, [2]int{i, j})
						i = j
					} else {
						i++
					}
				}
				return result
			})),
			pipeline.Ord(pipeline.Receive(func(_ int, data any) any {
				ndups := data.([][2]int)
				if len(ndups) == 0 {
					return nil
				}
				lx := len(dups)
				if lx == 0 {
					dups = ndups
					return nil
				}
				lx--
				if i, j := dups[lx][0], ndups[0][0]; indexCopies[i] == indexCopies[j] {
					ndups[0][0] = i
					if lx == 0 {
						dups = ndups
						return nil
					}
					dups = dups[:lx]
				}
				dups = append(dups, ndups...)
				return nil
			})),
		)
		p.Run()
		if err := p.Err(); err != nil {
			panic(err)
		}
		parallel.Range(0, len(dups), func(low, high int) {
			for i := low; i < high; i++ {
				dp := dups[i]
				start, end := dp[0], dp[1]
				for j := start + 1; j < end; j++ {
					valueCopies[start] = dup(valueCopies[start], valueCopies[j])
				}
			}
		})
		dups = append(dups, [2]int{len(indexCopies) - 1, len(indexCopies) - 1})
		delta := 0
		for i := 0; i < len(dups)-1; i++ {
			dstStart := dups[i][0] + 1 - delta
			srcStart := dups[i][1]
			srcEnd := dups[i+1][0] + 1
			copy(indexCopies[dstStart:], indexCopies[srcStart:srcEnd])
			copy(valueCopies[dstStart:], valueCopies[srcStart:srcEnd])
			delta += dups[i][1] - dups[i][0] - 1
		}
		indexCopies = indexCopies[:len(indexCopies)-delta]
		valueCopies = valueCopies[:len(valueCopies)-delta]
		return newSparseVector[T](size, indexCopies, valueCopies), int64(len(valueCopies))
	})
	return nil
}

func (v *Vector[T]) SetElement(value T, index int) error {
	if index < 0 {
		return InvalidIndex
	}
	if v == nil || v.ref == nil {
		return UninitializedObject
	}
	if index >= v.ref.size() {
		return InvalidIndex
	}
	v.ref = v.ref.setElement(value, index)
	return nil
}

func (v *Vector[T]) RemoveElement(index int) error {
	if index < 0 {
		return InvalidIndex
	}
	if v == nil || v.ref == nil {
		return UninitializedObject
	}
	if index >= v.ref.size() {
		return InvalidIndex
	}
	v.ref = v.ref.removeElement(index)
	return nil
}

func (v *Vector[T]) ExtractElement(index int) (result T, err error) {
	if index < 0 {
		err = InvalidIndex
		return
	}
	if v == nil || v.ref == nil {
		err = UninitializedObject
		return
	}
	if index >= v.ref.size() {
		err = InvalidIndex
		return
	}
	if value, ok := v.ref.extractElement(index); ok {
		return value, nil
	}
	err = NoValue
	return
}

func (v *Vector[T]) ExtractTuples() (indices []int, values []T, err error) {
	if v == nil || v.ref == nil {
		err = UninitializedObject
		return
	}
	p := v.ref.getPipeline()
	if p == nil {
		atomic.StoreInt64(&v.ref.nvalues, 0)
		return
	}
	var result vectorSlice[T]
	result.collect(p)
	indices = result.indices
	values = result.values
	atomic.StoreInt64(&v.ref.nvalues, int64(len(values)))
	return
}

func (v *Vector[T]) Wait(mode WaitMode) error {
	if v == nil || v.ref == nil {
		return UninitializedObject
	}
	if mode == Complete {
		return nil
	}
	v.ref.optimize()
	return nil
}

func (v *Vector[T]) AsMask() *Vector[bool] {
	if v == nil || v.ref == nil {
		return nil
	}
	n := atomic.LoadInt64(&v.ref.nvalues)
	switch v := any(v).(type) {
	case *Vector[bool]:
		return v
	case *Vector[int8]:
		return &Vector[bool]{newVectorReference[bool](newVectorAsMask[int8](v.ref), n)}
	case *Vector[int16]:
		return &Vector[bool]{newVectorReference[bool](newVectorAsMask[int16](v.ref), n)}
	case *Vector[int32]:
		return &Vector[bool]{newVectorReference[bool](newVectorAsMask[int32](v.ref), n)}
	case *Vector[int64]:
		return &Vector[bool]{newVectorReference[bool](newVectorAsMask[int64](v.ref), n)}
	case *Vector[uint8]:
		return &Vector[bool]{newVectorReference[bool](newVectorAsMask[uint8](v.ref), n)}
	case *Vector[uint16]:
		return &Vector[bool]{newVectorReference[bool](newVectorAsMask[uint16](v.ref), n)}
	case *Vector[uint32]:
		return &Vector[bool]{newVectorReference[bool](newVectorAsMask[uint32](v.ref), n)}
	case *Vector[uint64]:
		return &Vector[bool]{newVectorReference[bool](newVectorAsMask[uint64](v.ref), n)}
	case *Vector[float32]:
		return &Vector[bool]{newVectorReference[bool](newVectorAsMask[float32](v.ref), n)}
	case *Vector[float64]:
		return &Vector[bool]{newVectorReference[bool](newVectorAsMask[float64](v.ref), n)}
	}
	return &Vector[bool]{newVectorReference[bool](newVectorAsStructuralMask[T](v.ref), n)}
}
