package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/parallel"
	"github.com/intel/forGoParallel/pipeline"
	"github.com/intel/forGoParallel/psort"
	"runtime"
	"sync"
)

// A listMatrix is used to represent single entry modifications (setElement / removeElement) over other matrix representations.

// todo: maybe unrolled linked list instead?

type (
	matrixValueList[T any] struct {
		row, col int
		value    T
		next     *matrixValueList[T]
	}

	listMatrix[T any] struct {
		nrows, ncols int
		base         *matrixReference[T]
		entries      *matrixValueList[T]
	}
)

func newListMatrix[T any](
	nrows, ncols int,
	base *matrixReference[T],
	entries *matrixValueList[T],
) listMatrix[T] {
	return listMatrix[T]{
		nrows:   nrows,
		ncols:   ncols,
		base:    base,
		entries: entries,
	}
}

func (list *matrixValueList[T]) resize(nrows, ncols int) *matrixValueList[T] {
	if list == nil {
		return nil
	}
	tail := list.next.resize(nrows, ncols)
	if absInt(list.row) < nrows && absInt(list.col) < ncols {
		return &matrixValueList[T]{
			row:   list.row,
			col:   list.col,
			value: list.value,
			next:  tail,
		}
	}
	return tail
}

func (matrix listMatrix[T]) resize(ref *matrixReference[T], newNRows, newNCols int) *matrixReference[T] {
	if newNRows == matrix.nrows && newNCols == matrix.ncols {
		return ref
	}
	var newBase *matrixReference[T]
	newEntries := matrix.entries
	parallel.Do(func() {
		newBase = matrix.base.resize(newNRows, newNCols)
	}, func() {
		if newNRows < matrix.nrows || newNCols < matrix.ncols {
			newEntries = newEntries.resize(newNRows, newNCols)
		}
	})
	if newEntries == nil {
		return newBase
	}
	return newMatrixReference[T](newListMatrix[T](newNRows, newNCols, newBase, newEntries), -1)
}

func (matrix listMatrix[T]) size() (nrows, ncols int) {
	return matrix.nrows, matrix.ncols
}

func (matrix listMatrix[T]) nvals() int {
	var baseNVals, listNVals int
	parallel.Do(func() {
		baseNVals = matrix.base.nvals()
	}, func() {
		set := newMatrixBitset(matrix.nrows, matrix.ncols)
		for entry := matrix.entries; entry != nil; entry = entry.next {
			row := absInt(entry.row)
			col := absInt(entry.col)
			if !set.test(row, col) {
				_, ok := matrix.base.extractElement(row, col)
				if entry.row < 0 {
					if ok {
						listNVals--
					}
				} else {
					if !ok {
						listNVals++
					}
				}
				set.set(row, col)
			}
		}
	})
	return baseNVals + listNVals
}

func (matrix listMatrix[T]) setElement(_ *matrixReference[T], value T, row, col int) *matrixReference[T] {
	return newMatrixReference[T](newListMatrix[T](
		matrix.nrows,
		matrix.ncols,
		matrix.base,
		&matrixValueList[T]{
			row:   row,
			col:   col,
			value: value,
			next:  matrix.entries,
		}), -1)
}

func (matrix listMatrix[T]) removeElement(_ *matrixReference[T], row, col int) *matrixReference[T] {
	return newMatrixReference[T](newListMatrix[T](
		matrix.nrows,
		matrix.ncols,
		matrix.base,
		&matrixValueList[T]{
			row:  -row,
			col:  -col,
			next: matrix.entries,
		}), -1)
}

func (matrix listMatrix[T]) extractElement(row, col int) (result T, ok bool) {
	for entry := matrix.entries; entry != nil; entry = entry.next {
		if absInt(entry.row) == row && absInt(entry.col) == col {
			if entry.row < 0 {
				return
			}
			return entry.value, true
		}
	}
	return matrix.base.extractElement(row, col)
}

func (matrix listMatrix[T]) getPipeline() *pipeline.Pipeline[any] {
	var tmp matrixSlice[T]
	set := newMatrixBitset(matrix.nrows, matrix.ncols)
	for entry := matrix.entries; entry != nil; entry = entry.next {
		row := absInt(entry.row)
		col := absInt(entry.col)
		set.set(row, col)
	}
	tmp.rows, tmp.cols = set.toSlices()
	tmp.values = make([]T, len(tmp.rows))
	indexToPos := make(map[[2]int]int, len(tmp.rows))
	for i, row := range tmp.rows {
		col := tmp.cols[i]
		indexToPos[[2]int{row, col}] = i
	}
	for entry := matrix.entries; entry != nil; entry = entry.next {
		row := absInt(entry.row)
		col := absInt(entry.col)
		if set.test(row, col) {
			i := indexToPos[[2]int{row, col}]
			tmp.rows[i] = entry.row
			tmp.cols[i] = entry.col
			tmp.values[i] = entry.value
			set.clr(row, col)
		}
	}
	ch := make(chan any, runtime.GOMAXPROCS(0))
	var np pipeline.Pipeline[any]
	p := matrix.base.getPipeline()
	if p == nil {
		if len(tmp.values) > 0 {
			np.Source(pipeline.NewChan(ch))
			np.Notify(func() {
				for len(tmp.values) > 0 {
					var batch matrixSlice[T]
					batch.split(&tmp, Min(len(tmp.values), 512))
					batch.filter(func(row, col int, value T) (nrow, ncol int, nvalue T, ok bool) {
						return row, col, value, row >= 0
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
			slice := data.(matrixSlice[T])
			if len(slice.values) == 0 {
				return data
			}
			if len(tmp.values) == 0 {
				select {
				case <-p.Context().Done():
				case <-np.Context().Done():
				case ch <- slice:
				}
				return nil
			}
			maxIndex := [2]int{
				slice.rows[len(slice.rows)-1],
				slice.cols[len(slice.cols)-1],
			}
			if tmpRow, tmpCol := absInt(tmp.rows[0]), absInt(tmp.cols[0]); tmpRow > maxIndex[0] || (tmpRow == maxIndex[0] && tmpCol > maxIndex[1]) {
				select {
				case <-p.Context().Done():
				case <-np.Context().Done():
				case ch <- slice:
				}
				return nil
			}
			minIndex := [2]int{
				slice.rows[0],
				slice.cols[0],
			}
			if tmpRow, tmpCol := absInt(tmp.rows[0]), absInt(tmp.cols[0]); tmpRow < minIndex[0] || (tmpRow == minIndex[0] && tmpCol < minIndex[1]) {
				minIndex[0] = tmpRow
				minIndex[1] = tmpCol
			}
			set := newMatrixBitset(matrix.nrows, matrix.ncols)
			for i, row := range slice.rows {
				set.set(row, slice.cols[i])
			}
			seen := newMatrixBitset(matrix.nrows, matrix.ncols)
			for i, row := range tmp.rows {
				absRow := absInt(row)
				absCol := absInt(tmp.cols[i])
				if absRow > maxIndex[0] || (absRow == maxIndex[0] && absCol > maxIndex[1]) {
					break
				}
				if !seen.test(absRow, absCol) {
					if row < 0 {
						set.clr(absRow, absCol)
					} else {
						set.set(absRow, absCol)
					}
					seen.set(absRow, absCol)
				}
			}
			var result matrixSlice[T]
			result.rows, result.cols = set.toSlices()
			result.values = make([]T, len(result.rows))
			indexToPos := make(map[[2]int]int)
			for i, row := range result.rows {
				indexToPos[[2]int{row, result.cols[i]}] = i
			}
			for i, row := range tmp.rows {
				absRow := absInt(row)
				absCol := absInt(tmp.cols[i])
				if absRow > maxIndex[0] || (absRow == maxIndex[0] && absCol > maxIndex[1]) {
					tmp.rows = tmp.rows[i:]
					tmp.cols = tmp.cols[i:]
					tmp.values = tmp.values[i:]
					break
				}
				if row >= 0 && set.test(absRow, absCol) {
					result.values[indexToPos[[2]int{absRow, absCol}]] = tmp.values[i]
					set.clr(absRow, absCol)
				}
			}
			for i, row := range slice.rows {
				col := slice.cols[i]
				if set.test(row, col) {
					result.values[indexToPos[[2]int{row, col}]] = slice.values[i]
				}
			}
			select {
			case <-p.Context().Done():
			case <-np.Context().Done():
			case ch <- result:
			}
			return nil
		}, func() {
			for len(tmp.values) > 0 {
				var batch matrixSlice[T]
				batch.split(&tmp, Min(len(tmp.values), 512))
				batch.filter(func(row, col int, value T) (nrow, ncol int, nvalue T, ok bool) {
					return row, col, value, row >= 0
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

// todo: abstract this away against getColPipeline and vector.getPipeline
func (matrix listMatrix[T]) getRowPipeline(row int) *pipeline.Pipeline[any] {
	var tmp vectorSlice[T]
	set := newVectorBitset(matrix.ncols)
	for entry := matrix.entries; entry != nil; entry = entry.next {
		if absInt(entry.row) == row {
			set.set(absInt(entry.col))
		}
	}
	tmp.indices = set.toSlice()
	tmp.values = make([]T, len(tmp.indices))
	indexToPos := make(map[int]int, len(tmp.indices))
	for i, index := range tmp.indices {
		indexToPos[index] = i
	}
	for entry := matrix.entries; entry != nil; entry = entry.next {
		if absInt(entry.row) == row {
			col := absInt(entry.col)
			if set.test(col) {
				i := indexToPos[col]
				tmp.indices[i] = entry.col
				tmp.values[i] = entry.value
				set.clr(col)
			}
		}
	}
	ch := make(chan any, runtime.GOMAXPROCS(0))
	var np pipeline.Pipeline[any]
	p := matrix.base.getRowPipeline(row)
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

func (matrix listMatrix[T]) getColPipeline(col int) *pipeline.Pipeline[any] {
	var tmp vectorSlice[T]
	set := newVectorBitset(matrix.nrows)
	for entry := matrix.entries; entry != nil; entry = entry.next {
		if absInt(entry.col) == col {
			set.set(absInt(entry.row))
		}
	}
	tmp.indices = set.toSlice()
	tmp.values = make([]T, len(tmp.indices))
	indexToPos := make(map[int]int, len(tmp.indices))
	for i, index := range tmp.indices {
		indexToPos[index] = i
	}
	for entry := matrix.entries; entry != nil; entry = entry.next {
		if absInt(entry.col) == col {
			row := absInt(entry.row)
			if set.test(row) {
				i := indexToPos[row]
				tmp.indices[i] = entry.row
				tmp.values[i] = entry.value
				set.clr(col)
			}
		}
	}
	ch := make(chan any, runtime.GOMAXPROCS(0))
	var np pipeline.Pipeline[any]
	p := matrix.base.getColPipeline(col)
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

func makeTmpPipeline[T any](index int, tmpIndices []int, tmpValues []T) matrix1Pipeline {
	target := 0
	for i, index := range tmpIndices {
		if index >= 0 {
			tmpIndices[target] = index
			tmpValues[target] = tmpValues[i]
			target++
		}
	}
	tmpIndices = tmpIndices[:target]
	tmpValues = tmpValues[:target]
	var p pipeline.Pipeline[any]
	p.Source(pipeline.NewFunc[any](len(tmpValues), func(size int) (data any, fetched int, err error) {
		if len(tmpValues) == 0 {
			return
		}
		if size > len(tmpValues) {
			size = len(tmpValues)
		}
		result := vectorSlice[T]{
			indices: tmpIndices[:size:size],
			values:  tmpValues[:size:size],
		}
		tmpIndices = tmpIndices[size:]
		tmpValues = tmpValues[size:]
		return result, size, nil
	}))
	return matrix1Pipeline{
		index: index,
		p:     &p,
	}
}

func makeMergedPipeline[T any](base matrix1Pipeline, tmp vectorSlice[T]) matrix1Pipeline {
	ch := make(chan any, runtime.GOMAXPROCS(0))
	p := base.p
	var np pipeline.Pipeline[any]
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
						set.set(absIndex - minIndex)
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
	np.Notify(func() {
		p.Run()
		if err := p.Err(); err != nil {
			panic(err)
		}
	})
	np.Source(pipeline.NewChan(ch))
	return matrix1Pipeline{
		index: base.index,
		p:     &np,
	}
}

func (matrix listMatrix[T]) getRowPipelines() []matrix1Pipeline {
	set := newMatrixBitset(matrix.nrows, matrix.ncols)
	for entry := matrix.entries; entry != nil; entry = entry.next {
		row := absInt(entry.row)
		col := absInt(entry.col)
		set.set(row, col)
	}
	tmpIndices := make(map[int][]int)
	tmpValues := make(map[int][]T)
	var wg sync.WaitGroup
	for row, rowSet := range set.rows {
		indices := rowSet.toSlice()
		tmpIndices[row] = indices
		values := make([]T, len(indices))
		tmpValues[row] = values
		wg.Add(1)
		go func(row int, rowSet vectorBitset, indices []int, values []T) {
			defer wg.Done()
			indexToPos := make(map[int]int, len(indices))
			for i, index := range indices {
				indexToPos[index] = i
			}
			for entry := matrix.entries; entry != nil; entry = entry.next {
				if absInt(entry.row) == row {
					col := absInt(entry.col)
					if rowSet.test(col) {
						i := indexToPos[col]
						indices[i] = entry.col
						values[i] = entry.value
						rowSet.clr(col)
					}
				}
			}
		}(row, rowSet, indices, values)
	}
	var tmpRows []int
	for row := range tmpIndices {
		tmpRows = append(tmpRows, row)
	}
	psort.StableSort(psort.IntSlice(tmpRows))
	basePipelines := matrix.base.getRowPipelines()
	bi, ti := 0, 0
	var result []matrix1Pipeline
	wg.Wait()
	for {
		if bi >= len(basePipelines) {
			for ; ti < len(tmpRows); ti++ {
				row := tmpRows[ti]
				result = append(result, makeTmpPipeline(row, tmpIndices[row], tmpValues[row]))
			}
			return result
		}
		if ti >= len(tmpRows) {
			return append(result, basePipelines[bi:]...)
		}
		if basePipelines[bi].index == tmpRows[ti] {
			row := tmpRows[ti]
			result = append(result, makeMergedPipeline(basePipelines[bi], vectorSlice[T]{
				indices: tmpIndices[row],
				values:  tmpValues[row],
			}))
			bi++
			ti++
		} else if basePipelines[bi].index < tmpRows[ti] {
			result = append(result, basePipelines[bi])
			bi++
		} else {
			row := tmpRows[ti]
			result = append(result, makeTmpPipeline(row, tmpIndices[row], tmpValues[row]))
			ti++
		}
	}
}

func (matrix listMatrix[T]) getColPipelines() []matrix1Pipeline {
	set := newMatrixBitset(matrix.ncols, matrix.nrows)
	for entry := matrix.entries; entry != nil; entry = entry.next {
		row := absInt(entry.row)
		col := absInt(entry.col)
		set.set(col, row)
	}
	tmpIndices := make(map[int][]int)
	tmpValues := make(map[int][]T)
	var wg sync.WaitGroup
	for col, colSet := range set.rows {
		indices := colSet.toSlice()
		tmpIndices[col] = indices
		values := make([]T, len(indices))
		tmpValues[col] = values
		wg.Add(1)
		go func(col int, colSet vectorBitset, indices []int, values []T) {
			defer wg.Done()
			indexToPos := make(map[int]int, len(indices))
			for i, index := range indices {
				indexToPos[index] = i
			}
			for entry := matrix.entries; entry != nil; entry = entry.next {
				if absInt(entry.col) == col {
					row := absInt(entry.row)
					if colSet.test(row) {
						i := indexToPos[row]
						indices[i] = entry.row
						values[i] = entry.value
						colSet.clr(row)
					}
				}
			}
		}(col, colSet, indices, values)
	}
	var tmpCols []int
	for col := range tmpIndices {
		tmpCols = append(tmpCols, col)
	}
	psort.StableSort(psort.IntSlice(tmpCols))
	basePipelines := matrix.base.getColPipelines()
	bi, ti := 0, 0
	var result []matrix1Pipeline
	wg.Wait()
	for {
		if bi >= len(basePipelines) {
			for ; ti < len(tmpCols); ti++ {
				col := tmpCols[ti]
				result = append(result, makeTmpPipeline(col, tmpIndices[col], tmpValues[col]))
			}
			return result
		}
		if ti >= len(tmpCols) {
			return append(result, basePipelines[bi:]...)
		}
		if basePipelines[bi].index == tmpCols[ti] {
			col := tmpCols[ti]
			result = append(result, makeMergedPipeline(basePipelines[bi], vectorSlice[T]{
				indices: tmpIndices[col],
				values:  tmpValues[col],
			}))
			bi++
			ti++
		} else if basePipelines[bi].index < tmpCols[ti] {
			result = append(result, basePipelines[bi])
			bi++
		} else {
			col := tmpCols[ti]
			result = append(result, makeTmpPipeline(col, tmpIndices[col], tmpValues[col]))
			ti++
		}
	}
}

func (matrix listMatrix[T]) extractTuples() (rows, cols []int, values []T) {
	var result matrixSlice[T]
	result.collect(matrix.getPipeline())
	rows = result.rows
	cols = result.cols
	values = result.values
	return
}

func (_ listMatrix[T]) optimized() bool {
	return false
}

func (matrix listMatrix[T]) optimize() (result functionalMatrix[T]) {
	// todo: is there an efficient way to merge the functionality of the next two lines
	rows, cols, values := matrix.extractTuples()
	newRows, rowSpans := csrRows(rows)
	return newCSRMatrix(matrix.nrows, matrix.ncols, newRows, rowSpans, cols, values)
}
