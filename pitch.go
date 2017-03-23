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

import (
	"log"
)

type pitchMode string

const (
	// Pitch detection function
	// see: https://github.com/piem/aubio/blob/develop/src/pitch/pitch.c
	PitchDefault pitchMode = "default"
	PitchYin     pitchMode = "yin"
	PitchMcomb   pitchMode = "mcomb"
	PitchSchmitt pitchMode = "schmitt"
	PitchFcomb   pitchMode = "fcomb"
	PitchYinfft  pitchMode = "default"
)

type pitchOutMode string

const (
	// Pitch detection output modes
	// see: https://github.com/piem/aubio/blob/develop/src/pitch/pitch.c
	PitchOutFreq    = "freq"
	PitchOutMidi    = "midi"
	PitchOutCent    = "cent"
	PitchOutBin     = "bin"
	PitchOutDefault = "default"
)

// Pitch is a wrapper for the aubio_pitch_t pitch detection object.
type Pitch struct {
	o   *C.aubio_pitch_t
	buf *SimpleBuffer
}

// TODO(jwall): Shared buffers?

// NewPitch constructs a new Pitch object.
// It is the Callers responsibility to call Free on the returned
// Pitch object or leak memory.
//     p := NewPitch(mode, bufSize, blockSize, samplerate)
//     defer p.Free()
func NewPitch(mode pitchMode, bufSize, blockSize, sampleRate uint) *Pitch {
	return &Pitch{
		o: C.new_aubio_pitch(
			toCharTPtr(string(mode)),
			C.uint_t(bufSize),
			C.uint_t(blockSize),
			C.uint_t(sampleRate)),
		buf: NewSimpleBuffer(blockSize),
	}
}

func (p *Pitch) Buffer() *SimpleBuffer {
	return p.buf
}

// SetTolerance sets the yin or yinfft tolerance threshold.
func (p *Pitch) SetTolerance(tol float64) {
	C.aubio_pitch_set_tolerance(p.o, C.smpl_t(tol))
}

// SetUnit sets the output unit.
func (p *Pitch) SetUnit(outMode pitchOutMode) {
	C.aubio_pitch_set_unit(p.o, toCharTPtr(string(outMode)))
}

/* TODO(jwall) Update to latest version of aubio and add this in
// GetConfidence returns the current confidence of the pitch detection
// algorithm.
func (p *Pitch) GetConfidence() float64 {
	return float64(C.aubio_pitch_get_confidence(p.o))
}
*/

// Do runs one step of the pitch detection as determined by the bufSize.
func (p *Pitch) Do(in *SimpleBuffer) {
	if p.o != nil {
		C.aubio_pitch_do(p.o, in.vec, p.buf.vec)
	} else {
		log.Println("Called Do on empty Pitch. Maybe you called Free previously?")
	}
}

// Free frees the memory allocated by the aubio library for this object.
func (p *Pitch) Free() {
	if p.o != nil {
		C.del_aubio_pitch(p.o)
		p.o = nil
	}
	if p.buf != nil {
		p.buf.Free()
		p.buf = nil
	}
}
