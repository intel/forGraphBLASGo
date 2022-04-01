package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/parallel"
	"github.com/intel/forGoParallel/pipeline"
	"github.com/intel/forGoParallel/psort"
	"sort"
	"sync"
)

type vectorAssign[D any] struct {
	u       *vectorReference[D]
	indices []int
	index   func(j int) (i int, ok bool)
}

func vectorIndexLookup(indices []int) func(int) (int, bool) {
	if nindices, allIndices := isAll(indices); allIndices {
		return func(i int) (int, bool) {
			return i, i < nindices
		}
	}
	return newIntSearcher(indices).search
}

func newVectorAssign[D any](u *vectorReference[D], indices []int) computeVectorT[D] {
	return vectorAssign[D]{
		u:       u,
		indices: indices,
		index:   vectorIndexLookup(indices)}
}

func resizeAssignIndices(newSize int, indices []int) []int {
	n, all := isAll(indices)
	if newSize < n {
		if all {
			return All(newSize)
		}
		newIndices := make([]int, n)
		parallel.Range(0, n, func(low, high int) {
			for i := low; i < high; i++ {
				if index := indices[i]; index < newSize {
					newIndices[i] = index
				} else {
					newIndices[i] = -1
				}
			}
		})
		return newIndices
	}
	return indices
}

func (compute vectorAssign[D]) resize(newSize int) computeVectorT[D] {
	return newVectorAssign[D](compute.u, resizeAssignIndices(newSize, compute.indices))
}

func (compute vectorAssign[D]) assignIndex(index int) (int, bool) {
	return compute.index(index)
}

func (compute vectorAssign[D]) computeElement(index int) (result D, ok bool) {
	return compute.u.extractElement(index)
}

func (compute vectorAssign[D]) computePipeline() *pipeline.Pipeline[any] {
	p := compute.u.getPipeline()
	if p == nil {
		return nil
	}
	if n, all := isAll(compute.indices); all {
		p.Add(
			pipeline.Par(pipeline.Receive(func(_ int, data any) any {
				slice := data.(vectorSlice[D])
				slice.filter(func(index int, value D) (newIndex int, newValue D, ok bool) {
					return index, value, index < n
				})
				return slice
			})),
		)
		return p
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[D])
			slice.filter(func(index int, value D) (newIndex int, newValue D, ok bool) {
				if index < len(compute.indices) {
					newIndex = compute.indices[index]
					newValue = value
					ok = newIndex >= 0
				}
				return
			})
			return slice
		})),
	)
	if psort.IntsAreSorted(compute.indices) {
		return p
	}
	var result vectorSlice[D]
	var wg sync.WaitGroup
	wg.Add(1)
	var np pipeline.Pipeline[any]
	np.Source(vectorSourceWithWaitGroup(&wg, &result.indices, &result.values))
	np.Notify(func() {
		defer wg.Done()
		result.collect(p)
		vectorSort(result.indices, result.values)
	})
	return &np
}

func computeIndexValid(indices []int) func(int) bool {
	if nindices, ok := isAll(indices); ok {
		return func(index int) bool {
			return index < nindices
		}
	}
	psort.StableSort(psort.IntSlice(indices))
	return func(index int) bool {
		s := sort.SearchInts(indices, index)
		return s < len(indices) && indices[s] == index
	}
}

type vectorAssignConstant[D any] struct {
	value      D
	indices    []int
	indexValid func(int) bool
}

func newVectorAssignConstant[D any](value D, indices []int) computeVectorT[D] {
	return vectorAssignConstant[D]{
		value:      value,
		indices:    indices,
		indexValid: computeIndexValid(indices),
	}
}

func (compute vectorAssignConstant[D]) resize(newSize int) computeVectorT[D] {
	return newVectorAssignConstant[D](compute.value, resizeAssignIndices(newSize, compute.indices))
}

func (compute vectorAssignConstant[D]) assignIndex(index int) (int, bool) {
	return index, compute.indexValid(index)
}

func (compute vectorAssignConstant[D]) computeElement(_ int) (result D, ok bool) {
	return compute.value, true
}

func newVectorAssignConstantPipeline[D any](value D, indices []int) *pipeline.Pipeline[any] {
	var values []D
	if n, all := isAll(indices); all {
		var p pipeline.Pipeline[any]
		index := 0
		p.Source(pipeline.NewFunc[any](n, func(size int) (data any, fetched int, err error) {
			var result vectorSlice[D]
			if index >= n {
				return result, 0, nil
			}
			if index+size > n {
				size = n - index
			}
			if size < len(values) {
				values = values[:size]
			} else {
				for len(values) < size {
					values = append(values, value)
				}
			}
			result.cow = cowv
			result.indices = make([]int, size)
			for i := range result.indices {
				result.indices[i] = index + i
			}
			result.values = values
			index += size
			return result, size, nil
		}))
		return &p
	}
	var p pipeline.Pipeline[any]
	index := 0
	p.Source(pipeline.NewFunc[any](len(indices), func(size int) (data any, fetched int, err error) {
		var result vectorSlice[D]
		if index >= len(indices) {
			return result, 0, nil
		}
		if index+size > len(indices) {
			size = len(indices) - index
		}
		if size < len(values) {
			values = values[:size]
		} else {
			for len(values) < size {
				values = append(values, value)
			}
		}
		result.cow = cow0 | cowv
		result.indices = indices[index : index+size]
		result.values = values
		index += size
		return result, size, nil
	}))
	return &p
}

func (compute vectorAssignConstant[D]) computePipeline() *pipeline.Pipeline[any] {
	return newVectorAssignConstantPipeline(compute.value, compute.indices)
}

type deleteVector[D any] struct {
	indices    []int
	indexValid func(int) bool
}

func newDeleteVector[D any](indices []int) computeVectorT[D] {
	return deleteVector[D]{
		indices:    indices,
		indexValid: computeIndexValid(indices),
	}
}

func (compute deleteVector[D]) resize(newSize int) computeVectorT[D] {
	return newDeleteVector[D](resizeAssignIndices(newSize, compute.indices))
}

func (compute deleteVector[D]) assignIndex(index int) (int, bool) {
	return index, compute.indexValid(index)
}

func (_ deleteVector[D]) computeElement(_ int) (result D, ok bool) {
	return
}

func (_ deleteVector[D]) computePipeline() *pipeline.Pipeline[any] {
	return nil
}

type vectorAssignConstantScalar[D any] struct {
	value      *scalarReference[D]
	indices    []int
	indexValid func(int) bool
}

func newVectorAssignConstantScalar[D any](value *scalarReference[D], indices []int) computeVectorT[D] {
	s := value.get()
	if s.optimized() {
		if val, ok := s.extractElement(value); ok {
			return newVectorAssignConstant[D](val, indices)
		}
		return newDeleteVector[D](indices)
	}
	return vectorAssignConstantScalar[D]{
		value:      value,
		indices:    indices,
		indexValid: computeIndexValid(indices),
	}
}

func (compute vectorAssignConstantScalar[D]) resize(newSize int) computeVectorT[D] {
	return newVectorAssignConstantScalar[D](compute.value, resizeAssignIndices(newSize, compute.indices))
}

func (compute vectorAssignConstantScalar[D]) assignIndex(index int) (int, bool) {
	return index, compute.indexValid(index)
}

func (compute vectorAssignConstantScalar[D]) computeElement(_ int) (result D, ok bool) {
	return compute.value.extractElement()
}

func (compute vectorAssignConstantScalar[D]) computePipeline() *pipeline.Pipeline[any] {
	if s, sok := compute.value.extractElement(); sok {
		return newVectorAssignConstantPipeline(s, compute.indices)
	}
	return nil
}
