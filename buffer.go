package aubio

/*
#cgo LDFLAGS: -laubio
#include <aubio/aubio.h>
*/
import "C"

// SimpleBuffer is a wrapper for the aubio fvec_t type. It is used
// as the buffer for processing audio data in an aubio pipeline.
// It is a short sample buffer (32 or 64 bits in size).
type SimpleBuffer struct {
	vec *C.fvec_t
}

// NewSimpleBuffer constructs a new SimpleBuffer.
//
// The caller is responsible for calling Free on the returned
// SimpleBuffer to release memory when done.
//
//     buf := NewSimpleBuffer(bufSize)
//     defer buf.Free()
func NewSimpleBuffer(size uint) *SimpleBuffer {
	return &SimpleBuffer{C.new_fvec(C.uint_t(size))}
}

// Returns the contents of this buffer as a slice.
// The data is copied so the slices are still valid even
// after the buffer has changed.
func (b *SimpleBuffer) Slice() []float64 {
	sl := make([]float64, b.Size())
	for i := uint(0); i < b.Size(); i++ {
		sl[int(i)] = float64(C.fvec_read_sample(b.vec, C.uint_t(i)))
	}
	return sl
}

// Size returns the size of this buffer.
func (b *SimpleBuffer) Size() uint {
	if b.vec == nil {
		return 0
	}
	return uint(b.vec.length)
}

// Free frees the memory aubio allocated for this buffer.
func (b *SimpleBuffer) Free() {
	if b.vec == nil {
		return
	}
	C.del_fvec(b.vec)
	b.vec = nil
}

// ComplexBuffer is a wrapper for the aubio cvec_t type.
// It contains complex sample data.
type ComplexBuffer struct {
	data *C.cvec_t
}

// NewComplexBuffer constructs a buffer.
//
// The caller is responsible for calling Free on the returned
// ComplexBuffer to release memory when done.
//
//     buf := NewComplexBuffer(bufSize)
//     defer buf.Free()
func NewComplexBuffer(size uint) *ComplexBuffer {
	return &ComplexBuffer{data: C.new_cvec(C.uint_t(size))}
}

// Free frees the memory aubio has allocated for this buffer.
func (cb *ComplexBuffer) Free() {
	if cb.data != nil {
		C.del_cvec(cb.data)
	}
}

// Size returns the size of this ComplexBuffer.
func (cb *ComplexBuffer) Size() uint {
	if cb.data == nil {
		return 0
	}
	return uint(cb.data.length)
}

// Norm returns the slice of norm data.
// The data is copies so the slice is still
// valid after the buffer has changed.
func (cb *ComplexBuffer) Norm() []float64 {
	sl := make([]float64, cb.Size())
	for i := uint(0); i < cb.Size(); i++ {
		sl[int(i)] = float64(C.cvec_read_norm(cb.data, C.uint_t(i)))
	}
	return sl
}

// Norm returns the slice of phase data.
// The data is copies so the slice is still
// valid after the buffer has changed.
func (cb *ComplexBuffer) Phase() []float64 {
	sl := make([]float64, cb.Size())
	for i := uint(0); i < cb.Size(); i++ {
		sl[int(i)] = float64(C.cvec_read_phas(cb.data, C.uint_t(i)))
	}
	return sl
}

// Buffer for Long sample data (64 bits)
type LongSampleBuffer struct {
	vec *C.lvec_t
}

// NewLBuffer constructs a *LongSampleBuffer.
//
// The caller is responsible for calling Free on the returned
// LongSampleBuffer to release memory when done.
//
//     buf := NewLBuffer(bufSize)
//     defer buf.Free()
func NewLBuffer(size uint) *LongSampleBuffer {
	return newLBufferFromVec(C.new_lvec(C.uint_t(size)))
}

func newLBufferFromVec(v *C.lvec_t) *LongSampleBuffer {
	return &LongSampleBuffer{vec: v}
}

// Free frees the memory allocated by aubio for this buffer.
func (lb *LongSampleBuffer) Free() {
	if lb.vec != nil {
		C.del_lvec(lb.vec)
		lb.vec = nil
	}
}

// Size returns this buffers size.
func (lb *LongSampleBuffer) Size() uint {
	return uint(lb.vec.length)
}

// Returns the contents of this buffer as a slice.
// The data is copied so the slices are still valid even
// after the buffer has changed.
func (lb *LongSampleBuffer) Slice() []float64 {
	sl := make([]float64, lb.Size())
	for i := uint(0); i < lb.Size(); i++ {
		sl[int(i)] = float64(C.lvec_read_sample(lb.vec, C.uint_t(i)))
	}
	return sl
}
