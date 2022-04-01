package forGraphBLASGo

import (
	"context"
	"github.com/intel/forGoParallel/parallel"
	"github.com/intel/forGoParallel/pipeline"
	"runtime"
	"sort"
	"sync"
)

type parallelArray[T any] struct {
	lock     sync.RWMutex
	contents []T
	bitset   vectorBitset
}

func newParallelArray[T any](size int) *parallelArray[T] {
	return &parallelArray[T]{
		contents: make([]T, size),
		bitset:   newVectorBitset(size),
	}
}

func (p *parallelArray[T]) set(index int, value T) {
	p.lock.RLock()
	if index < len(p.contents) {
		p.contents[index] = value
		p.bitset.atomicSet(index)
		p.lock.RUnlock()
		return
	}
	p.lock.RUnlock()
	p.lock.Lock()
	newSize := len(p.contents)
	for index >= newSize {
		newSize *= 2
	}
	if newSize > len(p.contents) {
		newContents := make([]T, newSize)
		copy(newContents, p.contents)
		p.contents = newContents
		newBitset := newVectorBitset(newSize)
		copy(newBitset, p.bitset)
		p.bitset = newBitset
	}
	p.contents[index] = value
	p.bitset.set(index)
	p.lock.Unlock()
}

type opt[T any] struct {
	value T
	ok    bool
}

func reduceParallelArray[T, R any](p *parallelArray[T], m func(T) R, f func(x, y R) R) (R, bool) {
	result := parallel.RangeReduce(0, len(p.bitset), func(low, high int) (result opt[R]) {
		for i := low; i < high; i++ {
			for b, j := p.bitset[i], 0; b > 0; b, j = b>>1, j+1 {
				if b&1 != 0 {
					if result.ok {
						result.value = f(result.value, m(p.contents[i*64+j]))
					} else {
						result.value = m(p.contents[i*64+j])
						result.ok = true
					}
				}
			}
		}
		return
	}, func(x, y opt[R]) (result opt[R]) {
		if x.ok {
			if y.ok {
				result.value = f(x.value, y.value)
				result.ok = true
				return
			}
			return x
		}
		return y
	})
	return result.value, result.ok
}

func (p *parallelArray[T]) forEachConsume(f func(T)) {
	var empty T
	for i, b := range p.bitset {
		for j := 0; b > 0; b, j = b>>1, j+1 {
			if b&1 != 0 {
				index := i*64 + j
				element := p.contents[index]
				p.contents[index] = empty
				f(element)
			}
		}
	}
}

const (
	cow0 = 1 << iota
	cow1
	cowv
)

type (
	vectorSlice[T any] struct {
		cow     uint
		indices []int
		values  []T
	}

	matrixSlice[T any] struct {
		cow        uint
		rows, cols []int
		values     []T
	}
)

func (dst *vectorSlice[T]) split(src *vectorSlice[T], index int) {
	dst.cow = src.cow
	dst.indices = src.indices[:index:index]
	dst.values = src.values[:index:index]
	src.indices = src.indices[index:]
	src.values = src.values[index:]
}

func (slice *vectorSlice[T]) filter(predicate func(index int, value T) (resultIndex int, resultValue T, ok bool)) {
	src := *slice
	if src.cow&cow0 != 0 {
		slice.cow &^= cow0
		slice.indices = nil
	} else {
		slice.indices = src.indices[:0]
	}
	if src.cow&cowv != 0 {
		slice.cow &^= cowv
		slice.values = nil
	} else {
		slice.values = src.values[:0]
	}
	for i := range src.values {
		if index, value, ok := predicate(src.indices[i], src.values[i]); ok {
			slice.indices = append(slice.indices, index)
			slice.values = append(slice.values, value)
		}
	}
}

func (slice *vectorSlice[T]) collect(p *pipeline.Pipeline[any]) {
	a := newParallelArray[vectorSlice[T]](16)
	p.Add(
		pipeline.Par(pipeline.Receive(func(seq int, data any) any {
			a.set(seq, data.(vectorSlice[T]))
			return nil
		})),
	)
	p.Run()
	if err := p.Err(); err != nil {
		panic(err)
	}
	size, _ := reduceParallelArray(a, func(slice vectorSlice[T]) int {
		return len(slice.values)
	}, func(x, y int) int {
		return x + y
	})
	slice.cow = 0
	slice.indices = make([]int, size)
	slice.values = make([]T, size)
	index := 0
	a.forEachConsume(func(s vectorSlice[T]) {
		copy(slice.indices[index:], s.indices)
		copy(slice.values[index:], s.values)
		index += len(s.values)
	})
}

func (dst *matrixSlice[T]) split(src *matrixSlice[T], index int) {
	dst.cow = src.cow
	dst.rows = src.rows[:index:index]
	dst.cols = src.cols[:index:index]
	dst.values = src.values[:index:index]
	src.rows = src.rows[index:]
	src.cols = src.cols[index:]
	src.values = src.values[index:]
}

func (slice *matrixSlice[T]) filter(predicate func(row, col int, value T) (resultRow, resultCol int, resultValue T, ok bool)) {
	src := *slice
	if src.cow&cow0 != 0 {
		slice.cow &^= cow0
		slice.rows = nil
	} else {
		slice.rows = src.rows[:0]
	}
	if src.cow&cow1 != 0 {
		slice.cow &^= cow1
		slice.cols = nil
	} else {
		slice.cols = src.cols[:0]
	}
	if src.cow&cowv != 0 {
		slice.cow &^= cowv
		slice.values = nil
	} else {
		slice.values = src.values[:0]
	}
	for i := range src.values {
		if row, col, value, ok := predicate(src.rows[i], src.cols[i], src.values[i]); ok {
			slice.rows = append(slice.rows, row)
			slice.cols = append(slice.cols, col)
			slice.values = append(slice.values, value)
		}
	}
}

func (slice *matrixSlice[T]) collect(p *pipeline.Pipeline[any]) {
	a := newParallelArray[matrixSlice[T]](16)
	p.Add(
		pipeline.Par(pipeline.Receive(func(seq int, data any) any {
			a.set(seq, data.(matrixSlice[T]))
			return nil
		})),
	)
	p.Run()
	if err := p.Err(); err != nil {
		panic(err)
	}
	size, _ := reduceParallelArray(a, func(m matrixSlice[T]) int {
		return len(m.values)
	}, func(x, y int) int {
		return x + y
	})
	slice.cow = 0
	slice.rows = make([]int, size)
	slice.cols = make([]int, size)
	slice.values = make([]T, size)
	index := 0
	a.forEachConsume(func(s matrixSlice[T]) {
		copy(slice.rows[index:], s.rows)
		copy(slice.cols[index:], s.cols)
		copy(slice.values[index:], s.values)
		index += len(s.values)
	})
}

type (
	interval struct {
		start, end int
	}

	intervalSource struct {
		size, start, end int
	}
)

func newIntervalSource(size int) *intervalSource {
	return &intervalSource{size: size}
}

func (_ *intervalSource) Err() error {
	return nil
}

func (src *intervalSource) Prepare(_ context.Context) (size int) {
	return src.size
}

func (src *intervalSource) Fetch(size int) (fetched int) {
	if src.end+size >= src.size {
		src.start = src.end
		src.end = src.size
		return src.end - src.start
	}
	src.start = src.end
	src.end += size
	return size
}

func (src *intervalSource) Data() any {
	return interval{start: src.start, end: src.end}
}

func pipelineReduce[T any](p *pipeline.Pipeline[any], op func(T, T) T, fetchValues func(any) []T) (result T, ok bool) {
	if p == nil {
		return
	}
	a := newParallelArray[T](16)
	p.Add(
		pipeline.Par(pipeline.Receive(func(seq int, data any) any {
			values := fetchValues(data)
			if len(values) == 0 {
				return nil
			}
			value := values[0]
			for _, v := range values[1:] {
				value = op(value, v)
			}
			a.set(seq, value)
			return nil
		})),
	)
	p.Run()
	if err := p.Err(); err != nil {
		panic(err)
	}
	return reduceParallelArray(a, func(x T) T {
		return x
	}, op)
}

func vectorPipelineReduce[T any](p *pipeline.Pipeline[any], op func(T, T) T) (result T, ok bool) {
	return pipelineReduce(p, op, func(data any) []T {
		return data.(vectorSlice[T]).values
	})
}

func matrixPipelineReduce[T any](p *pipeline.Pipeline[any], op func(T, T) T) (result T, ok bool) {
	return pipelineReduce(p, op, func(data any) []T {
		return data.(matrixSlice[T]).values
	})
}

func pipelineToChannel(p *pipeline.Pipeline[any]) <-chan any {
	ch := make(chan any, runtime.GOMAXPROCS(0))
	if p == nil {
		close(ch)
		return ch
	}
	p.Add(
		pipeline.Ord(pipeline.ReceiveAndFinalize(
			func(_ int, data any) any {
				select {
				case <-p.Context().Done():
				case ch <- data:
				}
				return nil
			}, func() {
				close(ch)
			})),
	)
	go func() {
		p.Run()
		if err := p.Err(); err != nil {
			panic(err)
		}
	}()
	return ch
}

func vectorSource[T any](indices []int, values []T) pipeline.Source[any] {
	index := 0
	return pipeline.NewFunc[any](len(values), func(size int) (data any, fetched int, err error) {
		var result vectorSlice[T]
		if index >= len(values) {
			return result, 0, nil
		}
		if index+size > len(values) {
			size = len(values) - index
		}
		result.indices = indices[index : index+size : index+size]
		result.values = values[index : index+size : index+size]
		index += size
		return result, size, nil
	})
}

func vectorSourceWithWaitGroup[T any](wg *sync.WaitGroup, indices *[]int, values *[]T) pipeline.Source[any] {
	index := 0
	return pipeline.NewFunc[any](len(*values), func(size int) (data any, fetched int, err error) {
		wg.Wait()
		var result vectorSlice[T]
		if index >= len(*values) {
			return result, 0, nil
		}
		if index+size > len(*values) {
			size = len(*values) - index
		}
		result.indices = (*indices)[index : index+size : index+size]
		result.values = (*values)[index : index+size : index+size]
		index += size
		return result, size, nil
	})
}

func matrixSource[T any](rows, cols []int, values []T) pipeline.Source[any] {
	index := 0
	return pipeline.NewFunc[any](len(values), func(size int) (data any, fetched int, err error) {
		var result matrixSlice[T]
		if index >= len(values) {
			return result, 0, nil
		}
		if index+size > len(values) {
			size = len(values) - index
		}
		result.rows = rows[index : index+size : index+size]
		result.cols = cols[index : index+size : index+size]
		result.values = values[index : index+size : index+size]
		index += size
		return result, size, nil
	})
}

func matrixSourceWithWaitGroup[T any](wg *sync.WaitGroup, rows, cols *[]int, values *[]T) pipeline.Source[any] {
	index := 0
	return pipeline.NewFunc[any](len(*values), func(size int) (data any, fetched int, err error) {
		wg.Wait()
		var result matrixSlice[T]
		if index >= len(*values) {
			return result, 0, nil
		}
		if index+size > len(*values) {
			size = len(*values) - index
		}
		result.rows = (*rows)[index : index+size : index+size]
		result.cols = (*cols)[index : index+size : index+size]
		result.values = (*values)[index : index+size : index+size]
		index += size
		return result, size, nil
	})
}

type (
	vectorSlicePair[DL, DR any] struct {
		vl vectorSlice[DL]
		vr vectorSlice[DR]
	}

	matrixSlicePair[DL, DR any] struct {
		ml matrixSlice[DL]
		mr matrixSlice[DR]
	}
)

func minimumIndex(indices1 []int, index1 int, indices2 []int, index2 int) int {
	if index1 < len(indices1) {
		if index2 < len(indices2) {
			return Min(indices1[index1], indices2[index2])
		}
		return indices1[index1]
	}
	if index2 < len(indices2) {
		return indices2[index2]
	}
	return -1
}

func minimumMatrixIndex(rows1, cols1 []int, index1 int, rows2, cols2 []int, index2 int) (row, col int) {
	if index1 < len(rows1) {
		if index2 < len(rows2) {
			if row1, col1, row2, col2 := rows1[index1], cols1[index1], rows2[index2], cols2[index2]; row1 < row2 || (row1 == row2 && col1 < col2) {
				return row1, col1
			} else {
				return row2, col2
			}
		}
		return rows1[index1], cols1[index1]
	}
	if index2 < len(rows2) {
		return rows2[index2], rows2[index2]
	}
	return -1, -1
}

type (
	vectorSliceTriple[D1, D2, D3 any] struct {
		v1 vectorSlice[D1]
		v2 vectorSlice[D2]
		v3 vectorSlice[D3]
	}

	matrixSliceTriple[D1, D2, D3 any] struct {
		m1 matrixSlice[D1]
		m2 matrixSlice[D2]
		m3 matrixSlice[D3]
	}
)

func minimumIndex3(indices1 []int, index1 int, indices2 []int, index2 int, indices3 []int, index3 int) int {
	if index1 < len(indices1) {
		if index2 < len(indices2) {
			if index3 < len(indices3) {
				return Min(indices1[index1], Min(indices2[index2], indices3[index3]))
			}
			return Min(indices1[index1], indices2[index2])
		}
		if index3 < len(indices3) {
			return Min(indices1[index1], indices3[index3])
		}
		return indices1[index1]
	}
	return minimumIndex(indices2, index2, indices3, index3)
}

func mininmumMatrixIndex3(rows1, cols1 []int, index1 int, rows2, cols2 []int, index2 int, rows3, cols3 []int, index3 int) (row, col int) {
	row, col = -1, -1
	if index1 < len(rows1) {
		row = rows1[index1]
		col = cols1[index1]
	}
	if index2 < len(rows2) {
		if row2, col2 := rows2[index2], cols2[index2]; row == -1 || row2 < row || (row2 == row && col2 < col) {
			row = row2
			col = col2
		}
	}
	if index3 < len(rows3) {
		if row3, col3 := rows3[index3], cols3[index3]; row == -1 || row3 < row || (row3 == row && col3 < col) {
			row = row3
			col = col3
		}
	}
	return
}

func makeVector2SourcePipeline[D, DL, DR any](
	src1, src2 *pipeline.Pipeline[any],
	processEntry func(index int, leftValue DL, leftOk bool, rightValue DR, rightOk bool) (D, bool),
) *pipeline.Pipeline[any] {
	ch1 := pipelineToChannel(src1)
	ch2 := pipelineToChannel(src2)
	ok1 := true
	ok2 := true
	var slice1 vectorSlice[DL]
	var slice2 vectorSlice[DR]
	var p pipeline.Pipeline[any]
	p.Source(pipeline.NewFunc[any](-1,
		func(_ int) (data any, fetched int, err error) {
			for ok1 && len(slice1.values) == 0 {
				var d1 any
				d1, ok1 = <-ch1
				if ok1 {
					slice1 = d1.(vectorSlice[DL])
				}
			}
			for ok2 && len(slice2.values) == 0 {
				var d2 any
				d2, ok2 = <-ch2
				if ok2 {
					slice2 = d2.(vectorSlice[DR])
				}
			}
			if len(slice1.values)+len(slice2.values) == 0 {
				return
			}
			if len(slice1.values) > 0 {
				if len(slice2.values) > 0 {
					var sliceLeft vectorSlice[DL]
					var sliceRight vectorSlice[DR]
					stop := Min(slice1.indices[len(slice1.indices)-1], slice2.indices[len(slice2.indices)-1]) + 1
					sliceLeft.split(&slice1, sort.SearchInts(slice1.indices, stop))
					sliceRight.split(&slice2, sort.SearchInts(slice2.indices, stop))
					return vectorSlicePair[DL, DR]{sliceLeft, sliceRight}, len(sliceLeft.values) + len(sliceRight.values), nil
				}
				result := vectorSlicePair[DL, DR]{vl: slice1}
				slice1.indices = nil
				slice1.values = nil
				return result, len(result.vl.values), nil
			}
			result := vectorSlicePair[DL, DR]{vr: slice2}
			slice2.indices = nil
			slice2.values = nil
			return result, len(result.vr.values), nil
		},
	))
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			in := data.(vectorSlicePair[DL, DR])
			var result vectorSlice[D]
			addEntry := func(index int, lvalue DL, lok bool, rvalue DR, rok bool) {
				if value, ok := processEntry(index, lvalue, lok, rvalue, rok); ok {
					result.indices = append(result.indices, index)
					result.values = append(result.values, value)
				}
			}
			var emptyLeft DL
			var emptyRight DR
			var il, ir int
			for {
				minIndex := minimumIndex(in.vl.indices, il, in.vr.indices, ir)
				if minIndex == -1 {
					return result
				}
				if il < len(in.vl.indices) && in.vl.indices[il] == minIndex {
					if ir < len(in.vr.indices) && in.vr.indices[ir] == minIndex {
						addEntry(minIndex, in.vl.values[il], true, in.vr.values[ir], true)
						ir++
					} else {
						addEntry(minIndex, in.vl.values[il], true, emptyRight, false)
					}
					il++
				} else {
					addEntry(minIndex, emptyLeft, false, in.vr.values[ir], true)
					ir++
				}
			}
		})),
	)
	return &p
}

func makeMaskVector1SourcePipeline[D any](
	mask, src1 *pipeline.Pipeline[any],
	processEntry func(index int, maskValue, maskOk bool, value D, valueOk bool) (D, bool),
) *pipeline.Pipeline[any] {
	chm := pipelineToChannel(mask)
	ch1 := pipelineToChannel(src1)
	okm := true
	ok1 := true
	var slicem vectorSlice[bool]
	var slice1 vectorSlice[D]
	var p pipeline.Pipeline[any]
	p.Source(pipeline.NewFunc[any](-1,
		func(_ int) (data any, fetched int, err error) {
			for okm && len(slicem.values) == 0 {
				var dm any
				dm, okm = <-chm
				if okm {
					slicem = dm.(vectorSlice[bool])
				}
			}
			for ok1 && len(slice1.values) == 0 {
				var d1 any
				d1, ok1 = <-ch1
				if ok1 {
					slice1 = d1.(vectorSlice[D])
				}
			}
			if len(slice1.values) == 0 {
				mask.Cancel()
				return
			}
			if len(slicem.values) > 0 {
				var sliceLeft vectorSlice[bool]
				var sliceRight vectorSlice[D]
				stop := Min(slicem.indices[len(slicem.indices)-1], slice1.indices[len(slice1.indices)-1]) + 1
				sliceLeft.split(&slicem, sort.SearchInts(slicem.indices, stop))
				sliceRight.split(&slice1, sort.SearchInts(slice1.indices, stop))
				return vectorSlicePair[bool, D]{sliceLeft, sliceRight}, len(sliceLeft.values) + len(sliceRight.values), nil
			}
			result := vectorSlicePair[bool, D]{vr: slice1}
			slice1.indices = nil
			slice1.values = nil
			return result, len(result.vr.values), nil
		},
	))
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			in := data.(vectorSlicePair[bool, D])
			var result vectorSlice[D]
			addEntry := func(index int, lvalue, lok bool, rvalue D, rok bool) {
				if value, ok := processEntry(index, lvalue, lok, rvalue, rok); ok {
					result.indices = append(result.indices, index)
					result.values = append(result.values, value)
				}
			}
			var il, ir int
			for {
				minIndex := minimumIndex(in.vl.indices, il, in.vr.indices, ir)
				if minIndex == -1 {
					return result
				}
				if il < len(in.vl.indices) && in.vl.indices[il] == minIndex {
					if ir < len(in.vr.indices) && in.vr.indices[ir] == minIndex {
						addEntry(minIndex, in.vl.values[il], true, in.vr.values[ir], true)
						ir++
					}
					il++
				} else {
					addEntry(minIndex, false, false, in.vr.values[ir], true)
					ir++
				}
			}
		})),
	)
	return &p
}

func makeMaskVector2SourcePipeline[D any](
	mask, src1, src2 *pipeline.Pipeline[any],
	processEntry func(index int, maskValue, maskOk bool, leftValue D, leftOk bool, rightValue D, rightOk bool) (D, bool),
) *pipeline.Pipeline[any] {
	chm := pipelineToChannel(mask)
	ch1 := pipelineToChannel(src1)
	ch2 := pipelineToChannel(src2)
	okm := true
	ok1 := true
	ok2 := true
	var slicem vectorSlice[bool]
	var slice1, slice2 vectorSlice[D]
	var p pipeline.Pipeline[any]
	p.Source(pipeline.NewFunc[any](-1,
		func(_ int) (data any, fetched int, err error) {
			for okm && len(slicem.values) == 0 {
				var dm any
				dm, okm = <-chm
				if okm {
					slicem = dm.(vectorSlice[bool])
				}
			}
			for ok1 && len(slice1.values) == 0 {
				var d1 any
				d1, ok1 = <-ch1
				if ok1 {
					slice1 = d1.(vectorSlice[D])
				}
			}
			for ok2 && len(slice2.values) == 0 {
				var d2 any
				d2, ok2 = <-ch2
				if ok2 {
					slice2 = d2.(vectorSlice[D])
				}
			}
			if len(slice1.values)+len(slice2.values) == 0 {
				mask.Cancel()
				return
			}
			if len(slicem.values) > 0 {
				var sliceMask vectorSlice[bool]
				if len(slice1.values) > 0 {
					var sliceLeft vectorSlice[D]
					if len(slice2.values) > 0 {
						var sliceRight vectorSlice[D]
						stop := Min(slicem.indices[len(slicem.indices)-1],
							Min(slice1.indices[len(slice1.indices)-1],
								slice2.indices[len(slice2.indices)-1])) + 1
						sliceMask.split(&slicem, sort.SearchInts(slicem.indices, stop))
						sliceLeft.split(&slice1, sort.SearchInts(slice1.indices, stop))
						sliceRight.split(&slice2, sort.SearchInts(slice2.indices, stop))
						return vectorSliceTriple[bool, D, D]{sliceMask, sliceLeft, sliceRight}, len(sliceMask.values) + len(sliceLeft.values) + len(sliceRight.values), nil
					}
					stop := Min(slicem.indices[len(slicem.indices)-1], slice1.indices[len(slice1.indices)-1]) + 1
					sliceMask.split(&slicem, sort.SearchInts(slicem.indices, stop))
					sliceLeft.split(&slice1, sort.SearchInts(slice1.indices, stop))
					return vectorSliceTriple[bool, D, D]{v1: sliceMask, v2: sliceLeft}, len(sliceMask.values) + len(sliceLeft.values), nil
				}
				var sliceRight vectorSlice[D]
				stop := Min(slicem.indices[len(slicem.indices)-1], slice2.indices[len(slice2.indices)-1]) + 1
				sliceMask.split(&slicem, sort.SearchInts(slicem.indices, stop))
				sliceRight.split(&slice2, sort.SearchInts(slice2.indices, stop))
				return vectorSliceTriple[bool, D, D]{v1: sliceMask, v3: sliceRight}, len(sliceMask.values) + len(sliceRight.values), nil
			}
			if len(slice1.values) > 0 {
				if len(slice2.values) > 0 {
					var sliceLeft, sliceRight vectorSlice[D]
					stop := Min(slice1.indices[len(slice1.indices)-1], slice2.indices[len(slice2.indices)-1]) + 1
					sliceLeft.split(&slice1, sort.SearchInts(slice1.indices, stop))
					sliceRight.split(&slice2, sort.SearchInts(slice2.indices, stop))
					return vectorSliceTriple[bool, D, D]{v2: sliceLeft, v3: sliceRight}, len(sliceLeft.values) + len(sliceRight.values), nil
				}
				result := vectorSliceTriple[bool, D, D]{v2: slice1}
				slice1.indices = nil
				slice1.values = nil
				return result, len(result.v2.values), nil
			}
			result := vectorSliceTriple[bool, D, D]{v3: slice2}
			slice2.indices = nil
			slice2.values = nil
			return result, len(result.v3.values), err
		},
	))
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			in := data.(vectorSliceTriple[bool, D, D])
			var result vectorSlice[D]
			addEntry := func(index int, value1, ok1 bool, value2 D, ok2 bool, value3 D, ok3 bool) {
				if value, ok := processEntry(index, value1, ok1, value2, ok2, value3, ok3); ok {
					result.indices = append(result.indices, index)
					result.values = append(result.values, value)
				}
			}
			var i1, i2, i3 int
			var empty D
			for {
				minIndex := minimumIndex3(in.v1.indices, i1, in.v2.indices, i2, in.v3.indices, i3)
				if minIndex == -1 {
					return result
				}
				if i1 < len(in.v1.indices) && in.v1.indices[i1] == minIndex {
					if i2 < len(in.v2.indices) && in.v2.indices[i2] == minIndex {
						if i3 < len(in.v3.indices) && in.v3.indices[i3] == minIndex {
							addEntry(minIndex, in.v1.values[i1], true, in.v2.values[i2], true, in.v3.values[i3], true)
							i3++
						} else {
							addEntry(minIndex, in.v1.values[i1], true, in.v2.values[i2], true, empty, false)
						}
						i2++
					} else {
						if i3 < len(in.v3.indices) && in.v3.indices[i3] == minIndex {
							addEntry(minIndex, in.v1.values[i1], true, empty, false, in.v3.values[i3], true)
							i3++
						}
					}
					i1++
				} else {
					if i2 < len(in.v2.indices) && in.v2.indices[i2] == minIndex {
						if i3 < len(in.v3.indices) && in.v3.indices[i3] == minIndex {
							addEntry(minIndex, false, false, in.v2.values[i2], true, in.v3.values[i3], true)
							i3++
						} else {
							addEntry(minIndex, false, false, in.v2.values[i2], true, empty, false)
						}
						i2++
					} else {
						addEntry(minIndex, false, false, empty, false, in.v3.values[i3], true)
						i3++
					}
				}
			}
		})),
	)
	return &p
}

func makeMatrix2SourcePipeline[D, DL, DR any](
	src1, src2 *pipeline.Pipeline[any],
	processEntry func(row, col int, leftValue DL, leftOk bool, rightValue DR, rightOk bool) (D, bool),
) *pipeline.Pipeline[any] {
	ch1 := pipelineToChannel(src1)
	ch2 := pipelineToChannel(src2)
	ok1 := true
	ok2 := true
	var slice1 matrixSlice[DL]
	var slice2 matrixSlice[DR]
	var p pipeline.Pipeline[any]
	p.Source(pipeline.NewFunc[any](-1,
		func(_ int) (data any, fetched int, err error) {
			for ok1 && len(slice1.values) == 0 {
				var d1 any
				d1, ok1 = <-ch1
				if ok1 {
					slice1 = d1.(matrixSlice[DL])
				}
			}
			for ok2 && len(slice2.values) == 0 {
				var d2 any
				d2, ok2 = <-ch2
				if ok2 {
					slice2 = d2.(matrixSlice[DR])
				}
			}
			if len(slice1.values)+len(slice2.values) == 0 {
				return
			}
			if len(slice1.values) > 0 {
				if len(slice2.values) > 0 {
					var sliceLeft matrixSlice[DL]
					var sliceRight matrixSlice[DR]
					stop := [2]int{
						slice1.rows[len(slice1.rows)-1],
						slice1.cols[len(slice1.cols)-1],
					}
					if row2, col2 := slice2.rows[len(slice2.rows)-1], slice2.cols[len(slice2.cols)-1]; row2 < stop[0] || (row2 == stop[0] && col2 < stop[1]) {
						stop = [2]int{row2, col2}
					}
					parallel.Do(func() {
						sliceLeft.split(&slice1, sort.Search(len(slice1.values), func(i int) bool {
							return slice1.rows[i] > stop[0] || (slice1.rows[i] == stop[0] && slice1.cols[i] > stop[1])
						}))
					}, func() {
						sliceRight.split(&slice2, sort.Search(len(slice2.values), func(i int) bool {
							return slice2.rows[i] > stop[0] || (slice2.rows[i] == stop[0] && slice2.cols[i] > stop[1])
						}))
					})
					return matrixSlicePair[DL, DR]{sliceLeft, sliceRight}, len(sliceLeft.values) + len(sliceRight.values), nil
				}
				result := matrixSlicePair[DL, DR]{ml: slice1}
				slice1.rows = nil
				slice1.cols = nil
				slice1.values = nil
				return result, len(result.ml.values), nil
			}
			result := matrixSlicePair[DL, DR]{mr: slice2}
			slice2.rows = nil
			slice2.cols = nil
			slice2.values = nil
			return result, len(result.mr.values), nil
		},
	))
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			in := data.(matrixSlicePair[DL, DR])
			var result matrixSlice[D]
			addEntry := func(row, col int, lvalue DL, lok bool, rvalue DR, rok bool) {
				if value, ok := processEntry(row, col, lvalue, lok, rvalue, rok); ok {
					result.rows = append(result.rows, row)
					result.cols = append(result.cols, col)
					result.values = append(result.values, value)
				}
			}
			var emptyLeft DL
			var emptyRight DR
			var il, ir int
			for {
				minRow, minCol := minimumMatrixIndex(in.ml.rows, in.ml.cols, il, in.mr.rows, in.mr.cols, ir)
				if minRow == -1 {
					return result
				}
				if il < len(in.ml.rows) && in.ml.rows[il] == minRow && in.mr.cols[il] == minCol {
					if ir < len(in.mr.rows) && in.mr.rows[ir] == minRow && in.mr.cols[ir] == minCol {
						addEntry(minRow, minCol, in.ml.values[il], true, in.mr.values[ir], true)
						ir++
					} else {
						addEntry(minRow, minCol, in.ml.values[il], true, emptyRight, false)
					}
					il++
				} else {
					addEntry(minRow, minCol, emptyLeft, false, in.mr.values[ir], true)
					ir++
				}
			}
		})),
	)
	return &p
}

func makeMaskMatrix1SourcePipeline[D any](
	mask, src1 *pipeline.Pipeline[any],
	processEntry func(row, col int, maskValue, maskOk bool, value D, valueOk bool) (D, bool),
) *pipeline.Pipeline[any] {
	chm := pipelineToChannel(mask)
	ch1 := pipelineToChannel(src1)
	okm := true
	ok1 := true
	var slicem matrixSlice[bool]
	var slice1 matrixSlice[D]
	var p pipeline.Pipeline[any]
	p.Source(pipeline.NewFunc(-1,
		func(_ int) (data any, fetched int, err error) {
			for okm && len(slicem.values) == 0 {
				var dm any
				dm, okm = <-chm
				if okm {
					slicem = dm.(matrixSlice[bool])
				}
			}
			for ok1 && len(slice1.values) == 0 {
				var d1 any
				d1, ok1 = <-ch1
				if ok1 {
					slice1 = d1.(matrixSlice[D])
				}
			}
			if len(slice1.values) == 0 {
				mask.Cancel()
				return
			}
			if len(slicem.values) > 0 {
				var sliceLeft matrixSlice[bool]
				var sliceRight matrixSlice[D]
				stop := [2]int{
					slicem.rows[len(slicem.rows)-1],
					slicem.cols[len(slicem.cols)-1],
				}
				if row1, col1 := slice1.rows[len(slice1.rows)-1], slice1.cols[len(slice1.cols)-1]; row1 < stop[0] || (row1 == stop[0] && col1 < stop[1]) {
					stop = [2]int{row1, col1}
				}
				parallel.Do(func() {
					sliceLeft.split(&slicem, sort.Search(len(slicem.values), func(i int) bool {
						return slicem.rows[i] > stop[0] || (slicem.rows[i] == stop[0] && slicem.cols[i] > stop[1])
					}))
				}, func() {
					sliceRight.split(&slice1, sort.Search(len(slice1.values), func(i int) bool {
						return slice1.rows[i] > stop[0] || (slice1.rows[i] == stop[0] && slice1.cols[i] > stop[1])
					}))
				})
				return matrixSlicePair[bool, D]{sliceLeft, sliceRight}, len(sliceLeft.values) + len(sliceRight.values), nil
			}
			result := matrixSlicePair[bool, D]{mr: slice1}
			slice1.rows = nil
			slice1.cols = nil
			slice1.values = nil
			return result, len(result.mr.values), nil
		},
	))
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			in := data.(matrixSlicePair[bool, D])
			var result matrixSlice[D]
			addEntry := func(row, col int, lvalue, lok bool, rvalue D, rok bool) {
				if value, ok := processEntry(row, col, lvalue, lok, rvalue, rok); ok {
					result.rows = append(result.rows, row)
					result.cols = append(result.cols, col)
					result.values = append(result.values, value)
				}
			}
			var il, ir int
			for {
				minRow, minCol := minimumMatrixIndex(in.ml.rows, in.ml.cols, il, in.mr.rows, in.mr.cols, ir)
				if minRow == -1 {
					return result
				}
				if il < len(in.ml.rows) && in.ml.rows[il] == minRow && in.ml.cols[il] == minCol {
					if ir < len(in.mr.rows) && in.mr.rows[ir] == minRow && in.mr.cols[ir] == minCol {
						addEntry(minRow, minCol, in.ml.values[il], true, in.mr.values[ir], true)
						ir++
					}
					il++
				} else {
					addEntry(minRow, minCol, false, false, in.mr.values[ir], true)
					ir++
				}
			}
		})),
	)
	return &p
}

func makeMaskMatrix2SourcePipeline[D any](
	mask, src1, src2 *pipeline.Pipeline[any],
	processEntry func(row, col int, maskValue, maskOk bool, leftValue D, leftOk bool, rightValue D, rightOk bool) (D, bool),
) *pipeline.Pipeline[any] {
	chm := pipelineToChannel(mask)
	ch1 := pipelineToChannel(src1)
	ch2 := pipelineToChannel(src2)
	okm := true
	ok1 := true
	ok2 := true
	var slicem matrixSlice[bool]
	var slice1, slice2 matrixSlice[D]
	var p pipeline.Pipeline[any]
	p.Source(pipeline.NewFunc(-1,
		func(_ int) (data any, fetched int, err error) {
			for okm && len(slicem.values) == 0 {
				var dm any
				dm, okm = <-chm
				if okm {
					slicem = dm.(matrixSlice[bool])
				}
			}
			for ok1 && len(slice1.values) == 0 {
				var d1 any
				d1, ok1 = <-ch1
				if ok1 {
					slice1 = d1.(matrixSlice[D])
				}
			}
			for ok2 && len(slice2.values) == 0 {
				var d2 any
				d2, ok2 = <-ch2
				if ok2 {
					slice2 = d2.(matrixSlice[D])
				}
			}
			if len(slice1.values)+len(slice2.values) == 0 {
				mask.Cancel()
				return
			}
			if len(slicem.values) > 0 {
				var sliceMask matrixSlice[bool]
				if len(slice1.values) > 0 {
					var sliceLeft matrixSlice[D]
					if len(slice2.values) > 0 {
						var sliceRight matrixSlice[D]
						stop := [2]int{
							slicem.rows[len(slicem.rows)-1],
							slicem.cols[len(slicem.cols)-1],
						}
						if row1, col1 := slice1.rows[len(slice1.rows)-1], slice1.cols[len(slice1.cols)-1]; row1 < stop[0] || (row1 == stop[0] && col1 < stop[1]) {
							stop = [2]int{row1, col1}
						}
						if row2, col2 := slice2.rows[len(slice2.rows)-1], slice2.cols[len(slice2.cols)-1]; row2 < stop[0] || (row2 == stop[0] && col2 < stop[1]) {
							stop = [2]int{row2, col2}
						}
						parallel.Do(func() {
							sliceMask.split(&slicem, sort.Search(len(slicem.values), func(i int) bool {
								return slicem.rows[i] > stop[0] || (slicem.rows[i] == stop[0] && slicem.cols[i] > stop[1])
							}))
						}, func() {
							sliceLeft.split(&slice1, sort.Search(len(slice1.values), func(i int) bool {
								return slice1.rows[i] > stop[0] || (slice1.rows[i] == stop[0] && slice1.cols[i] > stop[1])
							}))
						}, func() {
							sliceRight.split(&slice2, sort.Search(len(slice2.values), func(i int) bool {
								return slice2.rows[i] > stop[0] || (slice2.rows[i] == stop[0] && slice2.cols[i] > stop[1])
							}))
						})
						return matrixSliceTriple[bool, D, D]{sliceMask, sliceLeft, sliceRight}, len(sliceMask.values) + len(sliceLeft.values) + len(sliceRight.values), nil
					}
					stop := [2]int{
						slicem.rows[len(slicem.rows)-1],
						slicem.cols[len(slicem.cols)-1],
					}
					if row1, col1 := slice1.rows[len(slice1.rows)-1], slice1.cols[len(slice1.cols)-1]; row1 < stop[0] || (row1 == stop[0] && col1 < stop[1]) {
						stop = [2]int{row1, col1}
					}
					parallel.Do(func() {
						sliceMask.split(&slicem, sort.Search(len(slicem.values), func(i int) bool {
							return slicem.rows[i] > stop[0] || (slicem.rows[i] == stop[0] && slicem.cols[i] > stop[1])
						}))
					}, func() {
						sliceLeft.split(&slice1, sort.Search(len(slice1.values), func(i int) bool {
							return slice1.rows[i] > stop[0] || (slice1.rows[i] == stop[0] && slice1.cols[i] > stop[1])
						}))
					})
					return matrixSliceTriple[bool, D, D]{m1: sliceMask, m2: sliceLeft}, len(sliceMask.values) + len(sliceLeft.values), nil
				}
				var sliceRight matrixSlice[D]
				stop := [2]int{
					slicem.rows[len(slicem.rows)-1],
					slicem.cols[len(slicem.cols)-1],
				}
				if row2, col2 := slice2.rows[len(slice2.rows)-1], slice2.cols[len(slice2.cols)-1]; row2 < stop[0] || (row2 == stop[0] && col2 < stop[1]) {
					stop = [2]int{row2, col2}
				}
				parallel.Do(func() {
					sliceMask.split(&slicem, sort.Search(len(slicem.values), func(i int) bool {
						return slicem.rows[i] > stop[0] || (slicem.rows[i] == stop[0] && slicem.cols[i] > stop[1])
					}))
				}, func() {
					sliceRight.split(&slice2, sort.Search(len(slice2.values), func(i int) bool {
						return slice2.rows[i] > stop[0] || (slice2.rows[i] == stop[0] && slice2.cols[i] > stop[1])
					}))
				})
				return matrixSliceTriple[bool, D, D]{m1: sliceMask, m3: sliceRight}, len(sliceMask.values) + len(sliceRight.values), nil
			}
			if len(slice1.values) > 0 {
				if len(slice2.values) > 0 {
					var sliceLeft, sliceRight matrixSlice[D]
					stop := [2]int{
						slice1.rows[len(slice1.rows)-1],
						slice1.cols[len(slice1.cols)-1],
					}
					if row2, col2 := slice2.rows[len(slice2.rows)-1], slice2.cols[len(slice2.cols)-1]; row2 < stop[0] || (row2 == stop[0] && col2 < stop[1]) {
						stop = [2]int{row2, col2}
					}
					parallel.Do(func() {
						sliceLeft.split(&slice1, sort.Search(len(slice1.values), func(i int) bool {
							return slice1.rows[i] > stop[0] || (slice1.rows[i] == stop[0] && slice1.cols[i] > stop[1])
						}))
					}, func() {
						sliceRight.split(&slice2, sort.Search(len(slice2.values), func(i int) bool {
							return slice2.rows[i] > stop[0] || (slice2.rows[i] == stop[0] && slice2.cols[i] > stop[1])
						}))
					})
					return matrixSliceTriple[bool, D, D]{m2: sliceLeft, m3: sliceRight}, len(sliceLeft.values) + len(sliceRight.values), nil
				}
				result := matrixSliceTriple[bool, D, D]{m2: slice1}
				slice1.rows = nil
				slice1.cols = nil
				slice1.values = nil
				return result, len(result.m2.values), nil
			}
			result := matrixSliceTriple[bool, D, D]{m3: slice2}
			slice2.rows = nil
			slice2.cols = nil
			slice2.values = nil
			return result, len(result.m3.values), err
		},
	))
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			in := data.(matrixSliceTriple[bool, D, D])
			var result matrixSlice[D]
			addEntry := func(row, col int, value1, ok1 bool, value2 D, ok2 bool, value3 D, ok3 bool) {
				if value, ok := processEntry(row, col, value1, ok1, value2, ok2, value3, ok3); ok {
					result.rows = append(result.rows, row)
					result.cols = append(result.cols, col)
					result.values = append(result.values, value)
				}
			}
			var i1, i2, i3 int
			var empty D
			for {
				minRow, minCol := mininmumMatrixIndex3(in.m1.rows, in.m1.cols, i1, in.m2.rows, in.m2.cols, i2, in.m3.rows, in.m3.cols, i3)
				if minRow == -1 {
					return result
				}
				if i1 < len(in.m1.rows) && in.m1.rows[i1] == minRow && in.m1.cols[i1] == minCol {
					if i2 < len(in.m2.rows) && in.m2.rows[i2] == minRow && in.m2.cols[i2] == minCol {
						if i3 < len(in.m3.rows) && in.m3.rows[i3] == minRow && in.m3.cols[i3] == minCol {
							addEntry(minRow, minCol, in.m1.values[i1], true, in.m2.values[i2], true, in.m3.values[i3], true)
							i3++
						} else {
							addEntry(minRow, minCol, in.m1.values[i1], true, in.m2.values[i2], true, empty, false)
						}
						i2++
					} else {
						if i3 < len(in.m3.rows) && in.m3.rows[i3] == minRow && in.m3.cols[i3] == minCol {
							addEntry(minRow, minCol, in.m1.values[i1], true, empty, false, in.m3.values[i3], true)
							i3++
						}
					}
					i1++
				} else {
					if i2 < len(in.m2.rows) && in.m2.rows[i2] == minRow && in.m2.cols[i2] == minCol {
						if i3 < len(in.m3.rows) && in.m3.rows[i3] == minRow && in.m3.cols[i3] == minCol {
							addEntry(minRow, minCol, false, false, in.m2.values[i2], true, in.m3.values[i3], true)
							i3++
						} else {
							addEntry(minRow, minCol, false, false, in.m2.values[i2], true, empty, false)
						}
						i2++
					} else {
						addEntry(minRow, minCol, false, false, empty, false, in.m3.values[i3], true)
						i3++
					}
				}
			}
		})),
	)
	return &p
}
