/*
 Copyright 2013 Jeremy Wall (jeremy@marzhillstudios.com)

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0
*/

package aubio

/*
#cgo LDFLAGS: -laubio
#include <aubio/aubio.h>
*/
import "C"

// Filter is a wrapper for the aubio_filter_t object.
type Filter struct {
	o   *C.aubio_filter_t
	buf *SimpleBuffer
}

// Constructs a Filter. Filters maintain their own working buffer
// which will get freed when the Filters Free method is called.
// The caller is responsible for calling Free on the constructed
// Filter or risk leaking memory.
func NewFilter(order, bufSize uint) (*Filter, error) {
	f, err := C.new_aubio_filter(C.uint_t(order))
	if f == nil {
		return nil, err
	}
	return &Filter{o: f, buf: NewSimpleBuffer(bufSize)}, nil
}

// Free frees up the memory allocatd by aubio for this Filter.
func (f *Filter) Free() {
	if f.o != nil {
		C.del_aubio_filter(f.o)
		f.o = nil
	}
	if f.buf != nil {
		f.buf.Free()
		f.buf = nil
	}
}

// Reset resets the memory for this Filter.
func (f *Filter) Reset() {
	if f.o != nil {
		C.aubio_filter_do_reset(f.o)
	}
}

// Buffer returns the output buffer for this Filter.
// The buffer is populated by calls to DoOutplace and is owned
// by the Filter object.
// Subsequent calls to DoOutplace may change the data contained in
// this buffer.
func (f *Filter) Buffer() *SimpleBuffer {
	return f.buf
}

// SetSamplerate sets the samplerate for this Filter.
func (f *Filter) SetSamplerate(rate uint) {
	if f.o != nil {
		C.aubio_filter_set_samplerate(f.o, C.uint_t(rate))
	}
}

// Do does an in-place filter on the input vector.
// The output buffer is not used.
func (f *Filter) Do(in *SimpleBuffer) {
	// Filter in-place
	if f.o != nil {
		C.aubio_filter_do(f.o, in.vec)
	}
}

// TODO(jwall): maybe the outplace filter should be a seperate type?
// DoOutPlace does filters the input vector into the Filter's output
// Buffer. Each call to this method will change the data contained
// in the output buffer. This buffer can be retrieved though the
// Buffer method.
func (f *Filter) DoOutplace(in *SimpleBuffer) {
	if f.o != nil {
		C.aubio_filter_do_outplace(f.o, in.vec, f.buf.vec)
	}
}

// DoFwdBack runs the aubio_filter_do_filtfilt function on this
// Filter.
func (f *Filter) DoFwdBack(in *SimpleBuffer, workBufSize uint) {
	if f.o != nil {
		tmp := NewSimpleBuffer(workBufSize)
		defer tmp.Free()
		C.aubio_filter_do_filtfilt(f.o, in.vec, tmp.vec)
	}
}

// Feedback returns the buffer containing the feedback coefficients.
func (f *Filter) Feedback() *LongSampleBuffer {
	if f.o != nil {
		return newLBufferFromVec(C.aubio_filter_get_feedback(f.o))
	}
	return nil
}

// Feedback returns the buffer containing the feedforward coefficients.
func (f *Filter) Feedforward() *LongSampleBuffer {
	if f.o != nil {
		return newLBufferFromVec(C.aubio_filter_get_feedforward(f.o))
	}
	return nil
}

// Order returns this Filters order.
func (f *Filter) Order() uint {
	if f.o != nil {
		return uint(C.aubio_filter_get_order(f.o))
	}
	return 0
}

// Samplerate returns this Filters samplerate.
func (f *Filter) Samplerate() uint {
	if f.o != nil {
		return uint(C.aubio_filter_get_samplerate(f.o))
	}
	return 0
}
