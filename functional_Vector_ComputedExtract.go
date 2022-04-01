package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/pipeline"
	"sync"
)

type vectorExtract[D any] struct {
	u        *vectorReference[D]
	indices  []int
	index    func(int) (int, bool)
	searcher intSearcher
}

func newVectorExtract[D any](u *vectorReference[D], indices []int) vectorExtract[D] {
	var index func(i int) (int, bool)
	var searcher intSearcher
	if n, all := isAll(indices); all {
		index = func(i int) (int, bool) {
			return i, i < n
		}
	} else {
		index = func(i int) (int, bool) {
			if i < len(indices) {
				return indices[i], true
			}
			return -1, false
		}
		searcher = newIntSearcher(indices)
	}
	return vectorExtract[D]{
		u:        u,
		indices:  indices,
		index:    index,
		searcher: searcher,
	}
}

func (compute vectorExtract[D]) resize(newSize int) computeVectorT[D] {
	var newIndices []int
	n, all := isAll(compute.indices)
	if newSize < n {
		if all {
			newIndices = All(newSize)
		} else {
			newIndices = compute.indices[:n]
		}
	} else {
		newIndices = compute.indices
	}
	return newVectorExtract[D](compute.u, newIndices)
}

func (compute vectorExtract[D]) computeElement(index int) (result D, ok bool) {
	if i, k := compute.index(index); k {
		return compute.u.extractElement(i)
	}
	return
}

func (compute vectorExtract[D]) computePipeline() *pipeline.Pipeline[any] {
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
				newIndex, ok = compute.searcher.search(index)
				newValue = value
				return
			})
			return slice
		})),
	)
	if compute.searcher.valuesAreSorted() {
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

type colExtract[D any] struct {
	A          *matrixReference[D]
	rowIndices []int
	col        int
	index      func(int) (int, bool)
	searcher   intSearcher
}

func newColExtract[D any](A *matrixReference[D], rowIndices []int, col int) colExtract[D] {
	var index func(i int) (int, bool)
	var searcher intSearcher
	if n, all := isAll(rowIndices); all {
		index = func(i int) (int, bool) {
			return i, i < n
		}
	} else {
		index = func(i int) (int, bool) {
			if i < len(rowIndices) {
				return rowIndices[i], true
			}
			return -1, false
		}
		searcher = newIntSearcher(rowIndices)
	}
	return colExtract[D]{A: A, rowIndices: rowIndices, col: col, index: index, searcher: searcher}
}

func (compute colExtract[D]) resize(newSize int) computeVectorT[D] {
	var newIndices []int
	n, all := isAll(compute.rowIndices)
	if newSize < n {
		if all {
			newIndices = All(newSize)
		} else {
			newIndices = compute.rowIndices[:n]
		}
	} else {
		newIndices = compute.rowIndices
	}
	return newColExtract[D](compute.A, newIndices, compute.col)
}

func (compute colExtract[D]) computeElement(index int) (result D, ok bool) {
	if i, k := compute.index(index); k {
		return compute.A.extractElement(i, compute.col)
	}
	return
}

func (compute colExtract[D]) computePipeline() *pipeline.Pipeline[any] {
	p := compute.A.getColPipeline(compute.col)
	if p == nil {
		return nil
	}
	if n, all := isAll(compute.rowIndices); all {
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
				newIndex, ok = compute.searcher.search(index)
				newValue = value
				return
			})
			return slice
		})),
	)
	if compute.searcher.valuesAreSorted() {
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
