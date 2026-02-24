package et

import (
	"github.com/viant/igo/internal"
	"unsafe"
)

type mergeCompareIntSliceForStmt struct {
	init      func(ptr unsafe.Pointer) unsafe.Pointer
	iOffset   uintptr
	jOffset   uintptr
	kOffset   uintptr
	midOffset uintptr
	rightOff  uintptr
	dataOff   uintptr
	tmpOff    uintptr
}

func (s *mergeCompareIntSliceForStmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	if s.init != nil {
		s.init(ptr)
	}
	i := *(*int)(unsafe.Pointer(uintptr(ptr) + s.iOffset))
	j := *(*int)(unsafe.Pointer(uintptr(ptr) + s.jOffset))
	k := *(*int)(unsafe.Pointer(uintptr(ptr) + s.kOffset))
	mid := *(*int)(unsafe.Pointer(uintptr(ptr) + s.midOffset))
	right := *(*int)(unsafe.Pointer(uintptr(ptr) + s.rightOff))
	data := *(*[]int)(unsafe.Pointer(uintptr(ptr) + s.dataOff))
	tmp := *(*[]int)(unsafe.Pointer(uintptr(ptr) + s.tmpOff))
	for i < mid && j < right {
		di := data[i]
		dj := data[j]
		if di <= dj {
			tmp[k] = di
			i++
		} else {
			tmp[k] = dj
			j++
		}
		k++
	}
	*(*int)(unsafe.Pointer(uintptr(ptr) + s.iOffset)) = i
	*(*int)(unsafe.Pointer(uintptr(ptr) + s.jOffset)) = j
	*(*int)(unsafe.Pointer(uintptr(ptr) + s.kOffset)) = k
	return nil
}

// NewForMergeCompareIntSliceStmt creates a fused loop for
// for i < mid && j < right { if data[i] <= data[j] { tmp[k] = data[i]; i++ } else { tmp[k] = data[j]; j++ }; k++ }.
func NewForMergeCompareIntSliceStmt(init New, iOffset, jOffset, kOffset, midOffset, rightOffset, dataOffset, tmpOffset uintptr) (New, error) {
	return func(control *Control) (internal.Compute, error) {
		stmt := &mergeCompareIntSliceForStmt{
			iOffset:   iOffset,
			jOffset:   jOffset,
			kOffset:   kOffset,
			midOffset: midOffset,
			rightOff:  rightOffset,
			dataOff:   dataOffset,
			tmpOff:    tmpOffset,
		}
		if init != nil {
			initCompute, err := init(control)
			if err != nil {
				return nil, err
			}
			stmt.init = initCompute
		}
		return stmt.compute, nil
	}, nil
}

type mergeCopyIntSliceForStmt struct {
	init        func(ptr unsafe.Pointer) unsafe.Pointer
	idxOffset   uintptr
	limitOffset uintptr
	dstOff      uintptr
	srcOff      uintptr
}

func (s *mergeCopyIntSliceForStmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	if s.init != nil {
		s.init(ptr)
	}
	idx := *(*int)(unsafe.Pointer(uintptr(ptr) + s.idxOffset))
	limit := *(*int)(unsafe.Pointer(uintptr(ptr) + s.limitOffset))
	dst := *(*[]int)(unsafe.Pointer(uintptr(ptr) + s.dstOff))
	src := *(*[]int)(unsafe.Pointer(uintptr(ptr) + s.srcOff))
	for idx < limit {
		dst[idx] = src[idx]
		idx++
	}
	*(*int)(unsafe.Pointer(uintptr(ptr) + s.idxOffset)) = idx
	return nil
}

// NewForMergeCopyIntSliceStmt creates a fused loop for
// for i < n { dst[i] = src[i]; i++ }.
func NewForMergeCopyIntSliceStmt(init New, idxOffset, limitOffset, dstOffset, srcOffset uintptr) (New, error) {
	return func(control *Control) (internal.Compute, error) {
		stmt := &mergeCopyIntSliceForStmt{
			idxOffset:   idxOffset,
			limitOffset: limitOffset,
			dstOff:      dstOffset,
			srcOff:      srcOffset,
		}
		if init != nil {
			initCompute, err := init(control)
			if err != nil {
				return nil, err
			}
			stmt.init = initCompute
		}
		return stmt.compute, nil
	}, nil
}

type mergeCopyWithKIntSliceForStmt struct {
	init        func(ptr unsafe.Pointer) unsafe.Pointer
	idxOffset   uintptr
	limitOffset uintptr
	kOffset     uintptr
	dstOff      uintptr
	srcOff      uintptr
}

func (s *mergeCopyWithKIntSliceForStmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	if s.init != nil {
		s.init(ptr)
	}
	idx := *(*int)(unsafe.Pointer(uintptr(ptr) + s.idxOffset))
	limit := *(*int)(unsafe.Pointer(uintptr(ptr) + s.limitOffset))
	k := *(*int)(unsafe.Pointer(uintptr(ptr) + s.kOffset))
	dst := *(*[]int)(unsafe.Pointer(uintptr(ptr) + s.dstOff))
	src := *(*[]int)(unsafe.Pointer(uintptr(ptr) + s.srcOff))
	for idx < limit {
		dst[k] = src[idx]
		idx++
		k++
	}
	*(*int)(unsafe.Pointer(uintptr(ptr) + s.idxOffset)) = idx
	*(*int)(unsafe.Pointer(uintptr(ptr) + s.kOffset)) = k
	return nil
}

// NewForMergeCopyWithKIntSliceStmt creates a fused loop for
// for i < mid { tmp[k] = data[i]; i++; k++ }.
func NewForMergeCopyWithKIntSliceStmt(init New, idxOffset, limitOffset, kOffset, dstOffset, srcOffset uintptr) (New, error) {
	return func(control *Control) (internal.Compute, error) {
		stmt := &mergeCopyWithKIntSliceForStmt{
			idxOffset:   idxOffset,
			limitOffset: limitOffset,
			kOffset:     kOffset,
			dstOff:      dstOffset,
			srcOff:      srcOffset,
		}
		if init != nil {
			initCompute, err := init(control)
			if err != nil {
				return nil, err
			}
			stmt.init = initCompute
		}
		return stmt.compute, nil
	}, nil
}
