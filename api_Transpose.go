package forGraphBLASGo

func Transpose[D any](C *Matrix[D], mask *Matrix[bool], accum BinaryOp[D, D, D], A *Matrix[D], desc Descriptor) error {
	nrows, ncols, err := C.Size()
	if err != nil {
		return err
	}
	isTran, err := A.expectSizeTran(ncols, nrows, desc, Inp0)
	if err != nil {
		return err
	}
	if isReplace, err := desc.Is(Outp, Replace); err != nil {
		panic(err)
	} else if isComp, err := desc.Is(Mask, Comp); err != nil {
		panic(err)
	} else if accum == nil && (isReplace || (mask == nil && !isComp)) {
		if isTran {
			C.ref = A.ref
		} else {
			C.ref = newTransposedMatrix(A.ref)
		}
		return nil
	}
	maskAsStructure, err := matrixMask(mask, nrows, ncols)
	if err != nil {
		return err
	}
	C.ref = newMatrixReference[D](newComputedMatrix[D](
		nrows, ncols, C.ref, maskAsStructure, accum,
		newTransposeMatrix[D](maybeTran(A.ref, isTran)),
		desc,
	), -1)
	return nil
}
