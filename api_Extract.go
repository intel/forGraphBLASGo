package forGraphBLASGo

func VectorExtract[D any](w *Vector[D], mask *Vector[bool], accum BinaryOp[D, D, D], u *Vector[D], indices []int, desc Descriptor) error {
	usize, err := u.Size()
	if err != nil {
		return err
	}
	nindices, _, err := checkIndices(indices, usize, w.expectSize)
	if err != nil {
		return err
	}
	isComp, err := desc.Is(Mask, Comp)
	if err != nil {
		panic(err)
	}
	if usize == nindices && mask == nil && !isComp && accum == nil {
		w.ref = u.ref
		return nil
	}
	maskAsStructure, err := vectorMask(mask, nindices)
	if err != nil {
		return err
	}
	w.ref = newVectorReference[D](newComputedVector[D](
		nindices, w.ref,
		maskAsStructure, accum,
		newVectorExtract(u.ref, fpcopy(indices)),
		desc,
	), -1)
	return nil
}

func MatrixExtract[D any](C *Matrix[D], mask *Matrix[bool], accum BinaryOp[D, D, D], A *Matrix[D], rowIndices, colIndices []int, desc Descriptor) error {
	nrows, ncols, err := C.Size()
	if err != nil {
		return err
	}
	nRowIndices, rowIsAll := isAll(rowIndices)
	if nrows != nRowIndices {
		return DimensionMismatch
	}
	nColIndices, colIsAll := isAll(colIndices)
	if ncols != nColIndices {
		return DimensionMismatch
	}
	ANRows, ANCols, err := A.Size()
	if err != nil {
		return err
	}
	isTran, err := desc.Is(Inp0, Tran)
	if err != nil {
		panic(err)
	}
	if isTran {
		ANRows, ANCols = ANCols, ANRows
	}
	if rowIsAll {
		if ANRows < nrows {
			return IndexOutOfBounds
		}
	} else if isAnyIndexOutOfBounds(rowIndices, ANRows) {
		return IndexOutOfBounds
	}
	if colIsAll {
		if ANCols < ncols {
			return IndexOutOfBounds
		}
	} else if isAnyIndexOutOfBounds(colIndices, ANCols) {
		return IndexOutOfBounds
	}
	isComp, err := desc.Is(Mask, Comp)
	if err != nil {
		panic(err)
	}
	if rowIsAll && nrows == ANRows &&
		colIsAll && ncols == ANCols &&
		mask == nil && !isComp && accum == nil {
		C.ref = maybeTran(A.ref, isTran)
		return nil
	}
	maskAsStructure, err := matrixMask(mask, nrows, ncols)
	if err != nil {
		return err
	}
	rowIndicesCopy, colIndicesCopy := fpcopy2(rowIndices, colIndices)
	C.ref = newMatrixReference[D](newComputedMatrix[D](
		nrows, ncols,
		C.ref,
		maskAsStructure,
		accum,
		newMatrixExtract(maybeTran(A.ref, isTran), rowIndicesCopy, colIndicesCopy),
		desc,
	), -1)
	return nil
}

func ColExtract[D any](w *Vector[D], mask *Vector[bool], accum BinaryOp[D, D, D], A *Matrix[D], rowIndices []int, col int, desc Descriptor) error {
	nrows, ncols, err := A.Size()
	if err != nil {
		return nil
	}
	isTran, err := desc.Is(Inp0, Tran)
	if err != nil {
		panic(err)
	}
	if isTran {
		nrows, ncols = ncols, nrows
	}
	if col >= ncols {
		return InvalidIndex
	}
	nindices, _, err := checkIndices(rowIndices, nrows, w.expectSize)
	if err != nil {
		return err
	}
	maskAsStructure, err := vectorMask(mask, nindices)
	if err != nil {
		return err
	}
	w.ref = newVectorReference[D](newComputedVector[D](
		nindices, w.ref,
		maskAsStructure, accum,
		newColExtract(maybeTran(A.ref, isTran), fpcopy(rowIndices), col),
		desc,
	), -1)
	return nil
}
