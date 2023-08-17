package GrB

// #include "GraphBLAS.h"
import "C"
import (
	"math"
	"slices"
	"unsafe"
)

func processAD[D any](
	accum *BinaryOp[D, D, D],
	desc *Descriptor,
) (caccum C.GrB_BinaryOp, cdesc C.GrB_Descriptor) {
	if accum == nil {
		caccum = C.GrB_BinaryOp(C.GrB_NULL)
	} else {
		caccum = accum.grb
	}
	if desc == nil {
		cdesc = C.GrB_Descriptor(C.GrB_NULL)
	} else {
		cdesc = desc.grb
	}
	return
}

func processDescriptor(
	desc *Descriptor,
) (cdesc C.GrB_Descriptor) {
	if desc == nil {
		cdesc = C.GrB_Descriptor(C.GrB_NULL)
	} else {
		cdesc = desc.grb
	}
	return
}

func processMADM[D any](
	mask *Matrix[bool],
	accum *BinaryOp[D, D, D],
	desc *Descriptor,
) (cmask C.GrB_Matrix, caccum C.GrB_BinaryOp, cdesc C.GrB_Descriptor) {
	if mask == nil {
		cmask = C.GrB_Matrix(C.GrB_NULL)
	} else {
		cmask = mask.grb
	}
	if accum == nil {
		caccum = C.GrB_BinaryOp(C.GrB_NULL)
	} else {
		caccum = accum.grb
	}
	if desc == nil {
		cdesc = C.GrB_Descriptor(C.GrB_NULL)
	} else {
		cdesc = desc.grb
	}
	return
}

func processMADV[D any](
	mask *Vector[bool],
	accum *BinaryOp[D, D, D],
	desc *Descriptor,
) (cmask C.GrB_Vector, caccum C.GrB_BinaryOp, cdesc C.GrB_Descriptor) {
	if mask == nil {
		cmask = C.GrB_Vector(C.GrB_NULL)
	} else {
		cmask = mask.grb
	}
	if accum == nil {
		caccum = C.GrB_BinaryOp(C.GrB_NULL)
	} else {
		caccum = accum.grb
	}
	if desc == nil {
		cdesc = C.GrB_Descriptor(C.GrB_NULL)
	} else {
		cdesc = desc.grb
	}
	return
}

func cSlice[To, From any](s []From) *To {
	return (*To)(unsafe.Pointer(unsafe.SliceData(s)))
}

func growslice[T any](sptr *[]T, n int) (newSection []T) {
	if sptr == nil {
		return nil
	}
	s := *sptr
	t := slices.Grow(s, n)
	*sptr = t[:len(s)+n]
	return t[len(s) : len(s)+n]
}

func gotocbool(b bool) C.int32_t {
	if b {
		return 1
	}
	return 0
}

func ctogobool(b C.int32_t) bool {
	return b != 0
}

func goIndices(indices []C.GrB_Index) []int {
	if unsafe.Sizeof(0) == 4 {
		narrow := make([]int, len(indices))
		for i, index := range indices {
			if index > math.MaxInt32 {
				panic("overflow")
			}
			narrow[i] = int(index)
		}
		return narrow
	}
	return unsafe.Slice((*int)(unsafe.Pointer(unsafe.SliceData(indices))), len(indices))
}

func grbIndices(indices []int) *C.GrB_Index {
	if unsafe.Sizeof(0) == 4 {
		wide := make([]int64, len(indices))
		for i, index := range indices {
			wide[i] = int64(index)
		}
		return cSlice[C.GrB_Index, int64](wide)
	}
	return cSlice[C.GrB_Index, int](indices)
}

func growIndices(indicesPtr *[]int, n int) (newSection *C.GrB_Index, finalize func()) {
	if indicesPtr == nil {
		return nil, func() {}
	}
	indices := *indicesPtr
	extended := slices.Grow(indices, n)[:len(indices)+n]
	defer func() {
		*indicesPtr = extended
	}()
	if unsafe.Sizeof(0) == 4 {
		t := len(indices)
		ns := make([]uint64, n)
		newSection = cSlice[C.GrB_Index, uint64](ns)
		finalize = func() {
			for i, v := range ns {
				if v > math.MaxInt32 {
					panic("overflow")
				}
				extended[t+i] = int(v)
			}
		}
		return
	}
	newSection = cSlice[C.GrB_Index, int](extended[len(indices) : len(indices)+n])
	finalize = func() {}
	return
}
