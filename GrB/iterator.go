package GrB

// #include "GraphBLAS.h"
import "C"
import "unsafe"

type iterator[D any] struct {
	grb    C.GxB_Iterator
	getter func(C.GxB_Iterator) D
}

func (it *iterator[D]) init() {
	var d D
	switch any(d).(type) {
	case bool:
		it.getter = func(grb C.GxB_Iterator) (result D) {
			*(*bool)(unsafe.Pointer(&result)) = bool(C.GxB_Iterator_get_BOOL(grb))
			return
		}
	case int:
		if unsafe.Sizeof(0) == 4 {
			it.getter = func(grb C.GxB_Iterator) (result D) {
				*(*int)(unsafe.Pointer(&result)) = int(C.GxB_Iterator_get_INT32(grb))
				return
			}
		} else {
			it.getter = func(grb C.GxB_Iterator) (result D) {
				*(*int)(unsafe.Pointer(&result)) = int(C.GxB_Iterator_get_INT64(grb))
				return
			}
		}
	case int8:
		it.getter = func(grb C.GxB_Iterator) (result D) {
			*(*int8)(unsafe.Pointer(&result)) = int8(C.GxB_Iterator_get_INT8(grb))
			return
		}
	case int16:
		it.getter = func(grb C.GxB_Iterator) (result D) {
			*(*int16)(unsafe.Pointer(&result)) = int16(C.GxB_Iterator_get_INT32(grb))
			return
		}
	case int32:
		it.getter = func(grb C.GxB_Iterator) (result D) {
			*(*int32)(unsafe.Pointer(&result)) = int32(C.GxB_Iterator_get_INT32(grb))
			return
		}
	case int64:
		it.getter = func(grb C.GxB_Iterator) (result D) {
			*(*int64)(unsafe.Pointer(&result)) = int64(C.GxB_Iterator_get_INT64(grb))
			return
		}
	case uint:
		if unsafe.Sizeof(0) == 4 {
			it.getter = func(grb C.GxB_Iterator) (result D) {
				*(*uint)(unsafe.Pointer(&result)) = uint(C.GxB_Iterator_get_UINT32(grb))
				return
			}
		} else {
			it.getter = func(grb C.GxB_Iterator) (result D) {
				*(*uint)(unsafe.Pointer(&result)) = uint(C.GxB_Iterator_get_UINT64(grb))
				return
			}
		}
	case uint8:
		it.getter = func(grb C.GxB_Iterator) (result D) {
			*(*uint8)(unsafe.Pointer(&result)) = uint8(C.GxB_Iterator_get_UINT8(grb))
			return
		}
	case uint16:
		it.getter = func(grb C.GxB_Iterator) (result D) {
			*(*uint16)(unsafe.Pointer(&result)) = uint16(C.GxB_Iterator_get_UINT32(grb))
			return
		}
	case uint32:
		it.getter = func(grb C.GxB_Iterator) (result D) {
			*(*uint32)(unsafe.Pointer(&result)) = uint32(C.GxB_Iterator_get_UINT32(grb))
			return
		}
	case uint64:
		it.getter = func(grb C.GxB_Iterator) (result D) {
			*(*uint64)(unsafe.Pointer(&result)) = uint64(C.GxB_Iterator_get_UINT64(grb))
			return
		}
	case float32:
		it.getter = func(grb C.GxB_Iterator) (result D) {
			*(*float32)(unsafe.Pointer(&result)) = float32(C.GxB_Iterator_get_FP32(grb))
			return
		}
	case float64:
		it.getter = func(grb C.GxB_Iterator) (result D) {
			*(*float64)(unsafe.Pointer(&result)) = float64(C.GxB_Iterator_get_FP64(grb))
			return
		}
	case complex64:
		it.getter = func(grb C.GxB_Iterator) (result D) {
			*(*complex64)(unsafe.Pointer(&result)) = complex64(C.GxB_Iterator_get_FC32(grb))
			return
		}
	case complex128:
		it.getter = func(grb C.GxB_Iterator) (result D) {
			*(*complex128)(unsafe.Pointer(&result)) = complex128(C.GxB_Iterator_get_FC64(grb))
			return
		}
	default:
		it.getter = func(grb C.GxB_Iterator) (result D) {
			C.GxB_Iterator_get_UDT(grb, unsafe.Pointer(&result))
			return
		}
	}
}

// Valid returns true if the iterator has been created by a successful call to [Vector.IteratorNew],
// [Matrix.IteratorNew], [Matrix.RowIteratorNew], or [Matrix.ColIteratorNew].
//
// Valid is a forGraphBLASGo extension. It is used in place of comparing against GrB_INVALID_HANDLE.
func (it iterator[D]) Valid() bool {
	return it.grb != C.GxB_Iterator(C.GrB_INVALID_HANDLE)
}

// Free destroys a previously created [VectorIterator], [RowIterator], [ColIterator], or
// [EntryIterator] and releases any resources associated with it. Calling Free on an object that is
// not valid is legal.
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func (it *iterator[D]) Free() error {
	info := Info(C.GxB_Iterator_free(&it.grb))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Wait until function calls in a sequence put the iterator into a state of completion or
// materialization.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [IndexOutOfBounds], [OutOfMemory], [Panic]
func (it iterator[D]) Wait(WaitMode) error {
	return nil
}

// Err returns an error message about any errors encountered during the processing associated with
// the iterator.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func (it iterator[D]) Err() (string, error) {
	return "", nil
}

type (
	// A RowIterator iterates across the rows of a matrix, and then within each row to access
	// the entries in a given row. Accessing all the entries of a matrix using a row iterator
	// requires an outer loop (for the rows) and an inner loop (for the entries in each row).
	// A matrix can be accessed via a row iterator only if its layout (determined by
	// [Matrix.GetLayout] is [ByRow].
	//
	// RowIterator is a SuiteSparse:GraphBLAS extension.
	RowIterator[D any] struct {
		iterator[D]
	}

	// A ColIterator iterates across the columns of a matrix, and then within each column to access
	// the entries in a given row. Accessing all the entries of a matrix using a column iterator
	// requires an outer loop (for the columns) and an inner loop (for the entries in each column).
	// A matrix can be accessed via a column iterator only if its layout (determined by
	// [Matrix.GetLayout] is [ByCol].
	//
	// ColIterator is a SuiteSparse:GraphBLAS extension.
	ColIterator[D any] struct {
		iterator[D]
	}

	// An EntryIterator iterates across the entries of a matrix. Accessing all of the entries of
	// a matrix requires just a single loop. Any matrix can be accessed with an entry iterator.
	//
	// EntryIterator is a SuiteSparse:GraphBLAS extension.
	EntryIterator[D any] struct {
		iterator[D]
	}

	// A VectorIterator iterates across the entries of a vector. Accessing all of the entries of
	// a vector requires just a single loop. Any vector can be accessed with an entry iterator.
	//
	// VectorIterator is a SuiteSparse:GraphBLAS extension.
	VectorIterator[D any] struct {
		iterator[D]
	}
)

// SeekRow moves a row iterator to the specified row.
//
// For hypersparse matrices, if the requested row is implicit, the iterator is moved
// to the first explicit row following it. If no such row exists, the iterator is
// exhausted.
//
// Return Values:
//   - true, false, nil: The iterator has been moved to a row that contains at least one entry.
//   - false, false, nil: The iterator has been moved to a row with no entries.
//   - false, true, nil: row is out of bounds (>= nrows).
//
// SeekRow is a SuiteSparse:GraphBLAS extension.
func (it RowIterator[D]) SeekRow(row int) (ok, exhausted bool, err error) {
	info := Info(C.GxB_rowIterator_seekRow(it.grb, C.GrB_Index(row)))
	switch info {
	case success:
		return true, false, nil
	case noValue:
		return false, false, nil
	case isExhausted:
		return false, true, nil
	}
	err = makeError(info)
	return
}

// Kount returns [Matrix.Nrows] for sparse, bitmap, and full matrices, and the
// number of explicit rows for hypersparse matrices.
//
// Kount is a SuiteSparse:GraphBLAS extension.
func (it RowIterator[D]) Kount() int {
	return int(C.GxB_rowIterator_kount(it.grb))
}

// KSeek moves a row iterator to the specified explicit row. For sparse, bitmap, and
// full matrices, this is the same as [RowIterator.SeekRow].
//
// Return Values:
//   - true, false, nil: The iterator has been moved to a row that contains at least one entry.
//   - false, false, nil: The iterator has been moved to a row with no entries.
//   - false, true, nil: k is out of bounds (>= kount).
//
// KSeek is a SuiteSparse:GraphBLAS extension.
func (it RowIterator[D]) KSeek(k int) (ok, exhausted bool, err error) {
	info := Info(C.GxB_rowIterator_kseek(it.grb, C.GrB_Index(k)))
	switch info {
	case success:
		return true, false, nil
	case noValue:
		return false, false, nil
	case isExhausted:
		return false, true, nil
	}
	err = makeError(info)
	return
}

// SeekCol moves a column iterator to the specified column.
//
// For hypersparse matrices, if the requested column is implicit, the iterator is moved
// to the first explicit column following it. If no such column exists, the iterator is
// exhausted.
//
// Return Values:
//   - true, false, nil: The iterator has been moved to a column that contains at least one entry.
//   - false, false, nil: The iterator has been moved to a column with no entries.
//   - false, true, nil: col is out of bounds (>= ncols).
//
// SeekCol is a SuiteSparse:GraphBLAS extension.
func (it ColIterator[D]) SeekCol(col int) (ok, exhausted bool, err error) {
	info := Info(C.GxB_colIterator_seekCol(it.grb, C.GrB_Index(col)))
	switch info {
	case success:
		return true, false, nil
	case noValue:
		return false, false, nil
	case isExhausted:
		return false, true, nil
	}
	err = makeError(info)
	return
}

// Kount returns [Matrix.Ncols] for sparse, bitmap, and full matrices, and the
// number of explicit columns for hypersparse matrices.
//
// Kount is a SuiteSparse:GraphBLAS extension.
func (it ColIterator[D]) Kount() int {
	return int(C.GxB_colIterator_kount(it.grb))
}

// KSeek moves a column iterator to the specified explicit column. For sparse, bitmap, and
// full matrices, this is the same as [ColIterator.SeekCol].
//
// Return Values:
//   - true, false, nil: The iterator has been moved to a column that contains at least one entry.
//   - false, false, nil: The iterator has been moved to a column with no entries.
//   - false, true, nil: k is out of bounds (>= kount).
//
// KSeek is a SuiteSparse:GraphBLAS extension.
func (it ColIterator[D]) KSeek(k int) (ok, exhausted bool, err error) {
	info := Info(C.GxB_colIterator_kseek(it.grb, C.GrB_Index(k)))
	switch info {
	case success:
		return true, false, nil
	case noValue:
		return false, false, nil
	case isExhausted:
		return false, true, nil
	}
	err = makeError(info)
	return
}

// Seek moves the entry iterator to the given position p, which is in the range
// 0 to pmax - 1, with pmax = [EntryIterator.Getpmax](). For sparse, hypersparse,
// and full matrices, pmax is the same as [Matrix.Nvals](). For bitmap matrices,
// pmax is equal to [Matrix.Nrows]() * [Matrix.Ncols]().
//
// If p >= pmax, the iterator is exhausted and false, nil is returned. Otherwise
// true, nil is returned.
//
// All entries in the matrix are given an ordinal position p. Seeking to position p
// will either move the iterator to that particular position, or to the next higher
// position containing an entry if there is no entry at position p. The latter case
// only occurs for bitmap matrices. Use [EntryIterator.Getp] to determine the current
// position of the iterator.
//
// Seek is a SuiteSparse:GraphBLAS extension.
func (it EntryIterator[D]) Seek(p int) (ok bool, err error) {
	info := Info(C.GxB_Matrix_Iterator_seek(it.grb, C.GrB_Index(p)))
	switch info {
	case success:
		return true, nil
	case isExhausted:
		return false, nil
	}
	err = makeError(info)
	return
}

// Getpmax returns [Matrix.Nvals]() for sparse, hypersparse, and full matrices; or
// [Matrix.Nrows]() * [Matrix.Ncols]() for bitmap matrices.
//
// Getpmax is a SuiteSparse:GraphBLAS extension.
func (it EntryIterator[D]) Getpmax() int {
	return int(C.GxB_Matrix_Iterator_getpmax(it.grb))
}

// Getp returns the current position of the iterator.
//
// Getp is a SuiteSparse:GraphBLAS extension.
func (it EntryIterator[D]) Getp() int {
	return int(C.GxB_Matrix_Iterator_getp(it.grb))
}

// Seek moves the vector iterator to the given position p, which is in the range
// 0 to pmax - 1, with pmax = [VectorIterator.Getpmax]().
//
// If p >= pmax, the iterator is exhausted and false, nil is returned. Otherwise
// true, nil is returned.
//
// All entries in the matrix are given an ordinal position p. Seeking to position p
// will either move the iterator to that particular position, or to the next higher
// position containing an entry if there is no entry at position p. The latter case
// only occurs for bitmap vectors. Use [EntryIterator.Getp] to determine the current
// position of the iterator.
//
// Seek is a SuiteSparse:GraphBLAS extension.
func (it VectorIterator[D]) Seek(p int) (ok bool, err error) {
	info := Info(C.GxB_Vector_Iterator_seek(it.grb, C.GrB_Index(p)))
	switch info {
	case success:
		return true, nil
	case isExhausted:
		return false, nil
	}
	err = makeError(info)
	return
}

// Getpmax returns [Vector.Nvals](); or [Vector.Size]() for bitmap vectors.
//
// Getpmax is a SuiteSparse:GraphBLAS extension.
func (it VectorIterator[D]) Getpmax() int {
	return int(C.GxB_Vector_Iterator_getpmax(it.grb))
}

// Getp returns the current position of the iterator.
//
// Getp is a SuiteSparse:GraphBLAS extension.
func (it VectorIterator[D]) Getp() int {
	return int(C.GxB_Vector_Iterator_getp(it.grb))
}

// NextRow moves the iterator to the next row.
//
// If the matrix is hypersparse, the next row is always an explicit row.
// Implicit rows are skipped.
//
// Return Values:
//   - true, false, nil: The iterator has been moved to a row that contains at least one entry.
//   - false, false, nil: The iterator has been moved to a row with no entries.
//   - false, true, nil: The iterator is exhausted.
//
// NextRow is a SuiteSparse:GraphBLAS extension.
func (it RowIterator[D]) NextRow() (ok, exhausted bool, err error) {
	info := Info(C.GxB_rowIterator_nextRow(it.grb))
	switch info {
	case success:
		return true, false, nil
	case noValue:
		return false, false, nil
	case isExhausted:
		return false, true, nil
	}
	err = makeError(info)
	return
}

// NextCol moves the iterator to the next entry in the current row.
//
// Return Values:
//   - true, nil: The iterator has been moved to the next entry.
//   - false, nil: The end of the row is reached. The iterator does not move to the next row.
//
// NextCol is a SuiteSparse:GraphBLAS extension.
func (it RowIterator[D]) NextCol() (ok bool, err error) {
	info := Info(C.GxB_rowIterator_nextCol(it.grb))
	switch info {
	case success:
		return true, nil
	case noValue:
		return false, nil
	}
	err = makeError(info)
	return
}

// NextCol moves the iterator to the next column.
//
// If the matrix is hypersparse, the next column is always an explicit column.
// Implicit column are skipped.
//
// Return Values:
//   - true, false, nil: The iterator has been moved to a column that contains at least one entry.
//   - false, false, nil: The iterator has been moved to a column with no entries.
//   - false, true, nil: The iterator is exhausted.
//
// NextCol is a SuiteSparse:GraphBLAS extension.
func (it ColIterator[D]) NextCol() (ok, exhausted bool, err error) {
	info := Info(C.GxB_colIterator_nextCol(it.grb))
	switch info {
	case success:
		return true, false, nil
	case noValue:
		return false, false, nil
	case isExhausted:
		return false, true, nil
	}
	err = makeError(info)
	return
}

// NextRow moves the iterator to the next entry in the current column.
//
// Return Values:
//   - true, nil: The iterator has been moved to the next entry.
//   - false, nil: The end of the column is reached. The iterator does not move to the next column.
//
// NextRow is a SuiteSparse:GraphBLAS extension.
func (it ColIterator[D]) NextRow() (ok bool, err error) {
	info := Info(C.GxB_colIterator_nextRow(it.grb))
	switch info {
	case success:
		return true, nil
	case noValue:
		return false, nil
	}
	err = makeError(info)
	return
}

// Next moves the iterator to the next entry. It returns true, nil if the iterator
// is at an entry that exists in the matrix, or false, nil otherwise.
//
// Next is a SuiteSparse:GraphBLAS extension.
func (it EntryIterator[D]) Next() (ok bool, err error) {
	info := Info(C.GxB_Matrix_Iterator_next(it.grb))
	switch info {
	case success:
		return true, nil
	case isExhausted:
		return false, nil
	}
	err = makeError(info)
	return
}

// Next moves the iterator to the next entry. It returns true, nil if the iterator
// is at an entry that exists in the vector, or false, nil otherwise.
//
// Next is a SuiteSparse:GraphBLAS extension.
func (it VectorIterator[D]) Next() (ok bool, err error) {
	info := Info(C.GxB_Vector_Iterator_next(it.grb))
	switch info {
	case success:
		return true, nil
	case isExhausted:
		return false, nil
	}
	err = makeError(info)
	return
}

// GetRowIndex returns nrows(a) if the iterator is exhausted, or the current
// row index otherwise.
//
// GetRowIndex is a SuiteSparse:GraphBLAS extension.
func (it RowIterator[D]) GetRowIndex() int {
	return int(C.GxB_rowIterator_getRowIndex(it.grb))
}

// GetColIndex returns the current column index.
//
// GetColIndex is a SuiteSparse:GraphBLAS extension.
func (it RowIterator[D]) GetColIndex() int {
	return int(C.GxB_rowIterator_getColIndex(it.grb))
}

// GetColIndex returns ncols(a) if the iterator is exhausted, or the current
// column index otherwise.
//
// GetColIndex is a SuiteSparse:GraphBLAS extension.
func (it ColIterator[D]) GetColIndex() int {
	return int(C.GxB_colIterator_getColIndex(it.grb))
}

// GetRowIndex returns the current row index.
//
// GetRowIndex is a SuiteSparse:GraphBLAS extension.
func (it ColIterator[D]) GetRowIndex() int {
	return int(C.GxB_colIterator_getRowIndex(it.grb))
}

// GetIndex returns the current row and column index.
//
// GetIndex is a SuiteSparse:GraphBLAS extension.
func (it EntryIterator[D]) GetIndex() (int, int) {
	var i, j C.GrB_Index
	C.GxB_Matrix_Iterator_getIndex(it.grb, &i, &j)
	return int(i), int(j)
}

// GetIndex returns the current index.
//
// GetIndex is a SuiteSparse:GraphBLAS extension.
func (it VectorIterator[D]) GetIndex() int {
	return int(C.GxB_Vector_Iterator_getIndex(it.grb))
}

// Get returns the value of the current iterator entry.
//
// Get is a SuiteSparse:GraphBLAS extension.
func (it iterator[D]) Get() D {
	return it.getter(it.grb)
}
