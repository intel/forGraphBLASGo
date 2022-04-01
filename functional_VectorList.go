package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/parallel"
	"github.com/intel/forGoParallel/pipeline"
	"runtime"
)

// A listVector is used to represent single entry modifications (setElement / removeElement) over other vector representations.

type (
	vectorValueList[T any] struct {
		col   int
		value T
		next  *vectorValueList[T]
	}

	listVector[T any] struct {
		nsize   int
		base    *vectorReference[T]
		entries *vectorValueList[T]
	}
)

func newListVector[T any](
	size int,
	base *vectorReference[T],
	entries *vectorValueList[T],
) listVector[T] {
	return listVector[T]{
		nsize:   size,
		base:    base,
		entries: entries,
	}
}

func (list *vectorValueList[T]) resize(nsize int) *vectorValueList[T] {
	if list == nil {
		return nil
	}
	tail := list.next.resize(nsize)
	if absInt(list.col) < nsize {
		return &vectorValueList[T]{
			col:   list.col,
			value: list.value,
			next:  tail,
		}
	}
	return tail
}

func (vector listVector[T]) resize(ref *vectorReference[T], newSize int) *vectorReference[T] {
	if newSize == vector.nsize {
		return ref
	}
	var newBase *vectorReference[T]
	var newEntries *vectorValueList[T]
	parallel.Do(func() {
		newBase = vector.base.resize(newSize)
	}, func() {
		if newSize < vector.nsize {
			newEntries = vector.entries.resize(newSize)
		}
	})
	if newEntries == nil {
		return newBase
	}
	return newVectorReference[T](newListVector[T](newSize, newBase, newEntries), -1)
}

func (vector listVector[T]) size() int {
	return vector.nsize
}

func (vector listVector[T]) nvals() int {
	var baseNVals, listNVals int
	parallel.Do(func() {
		baseNVals = vector.base.nvals()
	}, func() {
		set := newVectorBitset(vector.nsize)
		for entry := vector.entries; entry != nil; entry = entry.next {
			col := absInt(entry.col)
			if !set.test(col) {
				_, ok := vector.base.extractElement(col)
				if entry.col < 0 {
					if ok {
						listNVals--
					}
				} else {
					if !ok {
						listNVals++
					}
				}
				set.set(col)
			}
		}
	})
	return baseNVals + listNVals
}

func (vector listVector[T]) setElement(_ *vectorReference[T], value T, index int) *vectorReference[T] {
	return newVectorReference[T](newListVector[T](
		vector.nsize,
		vector.base,
		&vectorValueList[T]{
			col:   index,
			value: value,
			next:  vector.entries,
		}), -1)
}

func (vector listVector[T]) removeElement(_ *vectorReference[T], index int) *vectorReference[T] {
	return newVectorReference[T](newListVector[T](
		vector.nsize,
		vector.base,
		&vectorValueList[T]{
			col:  -index,
			next: vector.entries,
		}), -1)
}

func (vector listVector[T]) extractElement(index int) (result T, ok bool) {
	for entry := vector.entries; entry != nil; entry = entry.next {
		if absInt(entry.col) == index {
			if entry.col < 0 {
				return
			}
			return entry.value, true
		}
	}
	return vector.base.extractElement(index)
}

func (vector listVector[T]) getPipeline() *pipeline.Pipeline[any] {
	// todo: can we actually use this representation by default?
	var tmp vectorSlice[T]
	set := newVectorBitset(vector.nsize)
	for entry := vector.entries; entry != nil; entry = entry.next {
		set.set(absInt(entry.col))
	}
	tmp.indices = set.toSlice()
	tmp.values = make([]T, len(tmp.indices))
	indexToPos := make(map[int]int, len(tmp.indices))
	for i, index := range tmp.indices {
		indexToPos[index] = i
	}
	for entry := vector.entries; entry != nil; entry = entry.next {
		col := absInt(entry.col)
		if set.test(col) {
			i := indexToPos[col]
			tmp.indices[i] = entry.col
			tmp.values[i] = entry.value
			set.clr(col)
		}
	}
	ch := make(chan any, runtime.GOMAXPROCS(0))
	var np pipeline.Pipeline[any]
	p := vector.base.getPipeline()
	if p == nil {
		if len(tmp.indices) > 0 {
			np.Source(pipeline.NewChan(ch))
			np.Notify(func() {
				for len(tmp.indices) > 0 {
					var batch vectorSlice[T]
					batch.split(&tmp, Min(len(tmp.indices), 512))
					batch.filter(func(index int, value T) (newIndex int, newValue T, ok bool) {
						return index, value, index >= 0
					})
					select {
					case <-np.Context().Done():
						close(ch)
						return
					case ch <- batch:
					}
				}
				close(ch)
			})
			return &np
		}
		return nil
	}
	p.Add(
		pipeline.Ord(pipeline.ReceiveAndFinalize(func(_ int, data any) any {
			slice := data.(vectorSlice[T])
			if len(slice.indices) == 0 {
				return data
			}
			maxIndex := slice.indices[len(slice.indices)-1] + 1
			if len(tmp.indices) == 0 || absInt(tmp.indices[0]) >= maxIndex {
				select {
				case <-p.Context().Done():
				case <-np.Context().Done():
				case ch <- slice:
				}
				return nil
			}
			minIndex := Min(absInt(tmp.indices[0]), slice.indices[0])
			set := newVectorBitset(maxIndex - minIndex)
			for _, index := range slice.indices {
				set.set(index - minIndex)
			}
			seen := newVectorBitset(maxIndex - minIndex)
			for _, index := range tmp.indices {
				absIndex := absInt(index)
				if absIndex >= maxIndex {
					break
				}
				if !seen.test(absIndex - minIndex) {
					if index < 0 {
						set.clr(absIndex - minIndex)
					} else {
						set.set(index - minIndex)
					}
					seen.set(absIndex - minIndex)
				}
			}
			var result vectorSlice[T]
			result.indices = set.toSlice()
			result.values = make([]T, len(result.indices))
			indexToPos := make(map[int]int)
			for i, index := range result.indices {
				index += minIndex
				result.indices[i] = index
				indexToPos[index] = i
			}
			for i, index := range tmp.indices {
				if absInt(index) >= maxIndex {
					tmp.indices = tmp.indices[i:]
					tmp.values = tmp.values[i:]
					break
				}
				if index >= 0 && set.test(index-minIndex) {
					result.values[indexToPos[index]] = tmp.values[i]
					set.clr(index - minIndex)
				}
			}
			for i, index := range slice.indices {
				if set.test(index - minIndex) {
					result.values[indexToPos[index]] = slice.values[i]
				}
			}
			select {
			case <-p.Context().Done():
			case <-np.Context().Done():
			case ch <- result:
			}
			return nil
		}, func() {
			for len(tmp.indices) > 0 {
				var batch vectorSlice[T]
				batch.split(&tmp, Min(len(tmp.indices), 512))
				batch.filter(func(index int, value T) (newIndex int, newValue T, ok bool) {
					return index, value, index >= 0
				})
				select {
				case <-p.Context().Done():
					close(ch)
					return
				case <-np.Context().Done():
					close(ch)
					return
				case ch <- batch:
				}
			}
			close(ch)
		})),
	)
	np.Source(pipeline.NewChan(ch))
	np.Notify(func() {
		p.Run()
		if err := p.Err(); err != nil {
			panic(err)
		}
	})
	return &np
}

func (vector listVector[T]) extractTuples(size int) (indices []int, values []T) {
	var result vectorSlice[T]
	result.collect(vector.getPipeline())
	indices = result.indices
	values = result.values
	return
}

func (_ listVector[T]) optimized() bool {
	return false
}

func (vector listVector[T]) optimize() (result functionalVector[T]) {
	indices, values := vector.extractTuples(vector.nsize)
	return newSparseVector[T](vector.nsize, indices, values)
}
