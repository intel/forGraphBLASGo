package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/parallel"
	"unsafe"
)

func signInt(i int) int {
	return i >> ((unsafe.Sizeof(i) << 3) - 1)
}

func absInt(i int) int {
	sign := signInt(i)
	return (i ^ sign) + (sign & 1)
}

func equal[T any](x, y T) (ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()
	return any(x) == any(y)
}

func pcopy[T any](dst, src []T) {
	n := Min(len(dst), len(src))
	parallel.Range(0, n, func(low, high int) {
		copy(dst[low:high], src[low:high])
	})
}

func fpcopy[T any](src []T) []T {
	result := make([]T, len(src))
	pcopy(result, src)
	return result
}

func fpcopy2[T1, T2 any](src1 []T1, src2 []T2) (res1 []T1, res2 []T2) {
	res1 = fpcopy(src1)
	res2 = fpcopy(src2)
	return
}

func fpcopy3[T1, T2, T3 any](src1 []T1, src2 []T2, src3 []T3) (res1 []T1, res2 []T2, res3 []T3) {
	res1 = fpcopy(src1)
	res2 = fpcopy(src2)
	res3 = fpcopy(src3)
	return
}
