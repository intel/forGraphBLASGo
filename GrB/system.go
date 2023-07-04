package GrB

// #include <string.h>
import "C"
import "unsafe"

// A SystemSlice is similar to a regular Go slice, except that the memory it occupies is allocated
// by the allocator registered with [InitWithMalloc], and must be freed explicitly.
//
// The type T must not directly or indirectly contain any Go pointers.
// This is not checked.
//
// SystemSlice is a forGraphBLAS Go extension.
type SystemSlice[T any] struct {
	ptr  unsafe.Pointer
	size int
}

// MakeSystemSlice allocates a [SystemSlice].
// All entries are initialized to the type's zero value.
//
// MakeSystemSlice is a forGraphBLASGo extension.
func MakeSystemSlice[T any](size int) SystemSlice[T] {
	var x T
	sz := size * int(unsafe.Sizeof(x))
	return SystemSlice[T]{
		ptr:  calloc(sz),
		size: sz,
	}
}

// UnsafeSlice returns a go slice pointing to the same memory as the system slice.
// This allows for reading and writing the system slice, and performing other
// slice operations. Note, however, that the slice returned points to memory
// not handled by the Go memory manager, which implies that some slice operations are
// unsafe and might lead to undefined situations.
//
// UnsafeSlice is a forGraphBLASGo extension.
func (s SystemSlice[T]) UnsafeSlice() []T {
	var x T
	return unsafe.Slice((*T)(s.ptr), uintptr(s.size)/unsafe.Sizeof(x))
}

// Free deallocates the memory pointed to by this slice, and initializes
// the content of s to an empty slice.
//
// Note that this function is unsafe, because other instance of SystemSlice,
// or Go slices returned by [SystemSlice.UnsafeSlice], that share memory with this
// slice are now dangling.
//
// Free succeeds even if s is nil.
//
// Free is a forGraphBLASGo extension.
func (s *SystemSlice[T]) Free() {
	if s == nil {
		return
	}
	free(s.ptr)
	s.ptr = nil
	s.size = 0
}

// AsSystemSlice constructs a SystemSlice from a raw pointer and the given size in bytes.
// (size is not the number of T entries!)
//
// AsSystemSlice is a forGraphBLASGo extension.
func AsSystemSlice[T any](ptr unsafe.Pointer, size int) SystemSlice[T] {
	return SystemSlice[T]{
		ptr:  ptr,
		size: size,
	}
}

func makeGrBIndexSlice(s *SystemSlice[int]) (result SystemSlice[uint64], copied bool) {
	if unsafe.Sizeof(0) == 4 {
		src := s.UnsafeSlice()
		result = MakeSystemSlice[uint64](len(src))
		dst := result.UnsafeSlice()
		for i, v := range src {
			dst[i] = uint64(v)
		}
		copied = true
		return
	}
	result = SystemSlice[uint64]{
		ptr:  s.ptr,
		size: s.size,
	}
	s.ptr = nil
	s.size = 0
	return
}

func makeGoIndexSlice(s *SystemSlice[uint64]) (result SystemSlice[int], copied bool) {
	if unsafe.Sizeof(0) == 4 {
		src := s.UnsafeSlice()
		result = MakeSystemSlice[int](len(src))
		dst := result.UnsafeSlice()
		for i, v := range src {
			dst[i] = int(v)
		}
		copied = true
		return
	}
	result = SystemSlice[int]{
		ptr:  s.ptr,
		size: s.size,
	}
	s.ptr = nil
	s.size = 0
	return
}
