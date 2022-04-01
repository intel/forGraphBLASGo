package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/parallel"
	"github.com/intel/forGoParallel/psort"
	"sort"
	"sync/atomic"
)

type (
	intSearchPair struct{ key, value int }
	intSearcher   []intSearchPair
)

func (s intSearcher) Len() int {
	return len(s)
}

func (s intSearcher) Less(i, j int) bool {
	return s[i].key < s[j].key
}

func (s intSearcher) SequentialSort(i, j int) {
	sort.Stable(s[i:j])
}

func (s intSearcher) NewTemp() psort.StableSorter {
	return make(intSearcher, len(s))
}

func (dst intSearcher) Assign(source psort.StableSorter) func(i, j, len int) {
	src := source.(intSearcher)
	return func(i, j, len int) {
		copy(dst[i:i+len], src[j:j+len])
	}
}

func (s intSearcher) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func newIntSearcher(ints []int) intSearcher {
	s := make(intSearcher, 0, len(ints))
	for index, i := range ints {
		s = append(s, intSearchPair{i, index})
	}
	psort.StableSort(s)
	return s
}

func (s intSearcher) search(x int) (int, bool) {
	n := sort.Search(len(s), func(i int) bool {
		return s[i].key >= x
	})
	if n == len(s) {
		return -1, false
	}
	if s[n].key == x {
		return s[n].value, true
	}
	return -1, false
}

func (s intSearcher) valuesAreSorted() bool {
	// todo: use speculative
	return parallel.RangeAnd(0, len(s), func(low, high int) bool {
		if high < len(s) {
			high++
		}
		for i := 1; i < high; i++ {
			if !(s[i-1].value < s[i].value) {
				return false
			}
		}
		return true
	})
}

func (v *Vector[T]) expectSize(size int) error {
	if v == nil || v.ref == nil {
		return UninitializedObject
	}
	if v.ref.size() != size {
		return DimensionMismatch
	}
	return nil
}

func vectorMask(mask *Vector[bool], size int) (structure *vectorReference[bool], err error) {
	if mask == nil {
		return
	}
	if err = mask.expectSize(size); err != nil {
		return
	}
	return mask.ref, nil
}

func (m *Matrix[T]) expectSize(nrows, ncols int) error {
	if m == nil || m.ref == nil {
		return UninitializedObject
	}
	if mnrows, mncols := m.ref.size(); mnrows != nrows || mncols != ncols {
		return DimensionMismatch
	}
	return nil
}

func (m *Matrix[T]) expectSizeTran(nrows, ncols int, desc Descriptor, field DescField) (isTran bool, err error) {
	if isTran, err = desc.Is(field, Tran); err != nil {
		panic(err)
	}
	if m == nil || m.ref == nil {
		return isTran, UninitializedObject
	}
	mnrows, mncols := m.ref.size()
	if isTran {
		mnrows, mncols = mncols, mnrows
	}
	if mnrows != nrows || mncols != ncols {
		return isTran, DimensionMismatch
	}
	return
}

func matrixMask(mask *Matrix[bool], nrows, ncols int) (structure *matrixReference[bool], err error) {
	if mask == nil {
		return
	}
	if err = mask.expectSize(nrows, ncols); err != nil {
		return
	}
	return mask.ref, nil
}

func maybeTran[T any](m *matrixReference[T], tran bool) *matrixReference[T] {
	if tran {
		return newTransposedMatrix[T](m)
	}
	return m
}

func csrRows(cooRows []int) (rows, rowSpans []int) {
	rowSpans = []int{0}
	for i := 0; i < len(cooRows); {
		row := cooRows[i]
		rowNNZ := sort.SearchInts(cooRows[i:], row+1)
		rows = append(rows, row)
		rowSpans = append(rowSpans, rowSpans[len(rowSpans)-1]+rowNNZ)
		i += rowNNZ
	}
	return
}

type vectorBitset []uint64

func newVectorBitset(size int) vectorBitset {
	return make([]uint64, (size+63)/64)
}

func (b *vectorBitset) set(index int) {
	(*b)[index/64] |= uint64(1) << (index % 64)
}

func (b *vectorBitset) atomicSet(index int) {
	addr := &(*b)[index/64]
	bit := uint64(1) << (index % 64)
	for {
		v := atomic.LoadUint64(addr)
		if v&bit != 0 {
			return
		}
		if atomic.CompareAndSwapUint64(addr, v, v|bit) {
			break
		}
	}
}

func (b *vectorBitset) clr(index int) {
	(*b)[index/64] &^= uint64(1) << (index % 64)
}

func (b *vectorBitset) atomicClr(index int) {
	addr := &(*b)[index/64]
	bit := uint64(1) << (index % 64)
	for {
		v := atomic.LoadUint64(addr)
		if v&bit == 0 {
			return
		}
		if atomic.CompareAndSwapUint64(addr, v, v&^bit) {
			break
		}
	}
}

func (b vectorBitset) test(index int) bool {
	return (b[index/64] & (uint64(1) << (index % 64))) != 0
}

func (b *vectorBitset) or(a vectorBitset) {
	parallel.Range(0, len(a), func(low, high int) {
		for i := low; i < high; i++ {
			(*b)[i] |= a[i]
		}
	})
}

func (b vectorBitset) toSlice() (result []int) {
	for i, a := range b {
		major := i * 64
		for j, s := 0, uint64(1); j < 64; j, s = j+1, s<<1 {
			if a&s != 0 {
				result = append(result, major+j)
			}
		}
	}
	return
}

type matrixBitset struct {
	nrows, ncols int
	rows         map[int]vectorBitset
}

func newMatrixBitset(nrows, ncols int) matrixBitset {
	return matrixBitset{
		nrows: nrows,
		ncols: ncols,
		rows:  make(map[int]vectorBitset),
	}
}

func (b *matrixBitset) set(row, col int) {
	rowSet := b.rows[row]
	if rowSet == nil {
		rowSet = newVectorBitset(b.ncols)
		b.rows[row] = rowSet
	}
	rowSet.set(col)
}

func (b *matrixBitset) clr(row, col int) {
	rowSet := b.rows[row]
	if rowSet == nil {
		return
	}
	rowSet.clr(col)
}

func (b matrixBitset) test(row, col int) bool {
	rowSet := b.rows[row]
	if rowSet == nil {
		return false
	}
	return rowSet.test(col)
}

func (b *matrixBitset) or(a matrixBitset) {
	for arow, arowSet := range a.rows {
		if browSet := b.rows[arow]; browSet == nil {
			b.rows[arow] = fpcopy(arowSet)
		} else {
			browSet.or(arowSet)
		}
	}
}

func (b matrixBitset) toSlices() (rows, cols []int) {
	rowIndices := make([]int, 0, len(b.rows))
	for row := range b.rows {
		rowIndices = append(rowIndices, row)
	}
	psort.StableSort(psort.IntSlice(rowIndices))
	for _, row := range rowIndices {
		cols = append(cols, b.rows[row].toSlice()...)
		for len(rows) < len(cols) {
			rows = append(rows, row)
		}
	}
	return
}

func coordToIndex(row, col, _, ncols int) int {
	return row*ncols + col
}

func indexToCoord(index, _, ncols int) (row, col int) {
	return index / ncols, index % ncols
}
