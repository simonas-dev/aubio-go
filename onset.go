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
#include <aubio/onset/onset.h>
*/
import "C"

import (
	"fmt"
)

type onsetMode string

const (
	// Onset Detection functions see: https://github.com/piem/aubio/blob/develop/src/spectral/specdesc.h
	// Energy based onset detection function
	Energy onsetMode = "energy"
	// High Frequency Content onset detection function
	HFC onsetMode = "hfc"
	// Complex Domain Method onset detection function
	Complex onsetMode = "complex"
	// Phase based Method onset detection function
	Phase onsetMode = "phase"
	// Spectral difference method onset detection function
	SpecDiff onsetMode = "specdiff"
	// Kullback-Liebler onset detection function
	K1 onsetMode = "k1"
	// Modified Kullback-Liebler onset detection function
	MK1 onsetMode = "mk1"
	// Spectral Flux
	SpecFlux onsetMode = "specflux"
)

// Tempo is a wrapper for the aubio_tempo_t tempo detection object.
type Onset struct {
	o   *C.aubio_onset_t
	buf *SimpleBuffer
}

// OnsetOrDie constructs a new Onset object.
// It panics on any errors.
func OnsetOrDie(mode onsetMode, bufSize, blocksize, samplerate uint) *Onset {
	if t, err := NewOnset(mode, bufSize, blocksize, samplerate); err == nil {
		return t
	} else {
		panic(err)
	}
	panic("Unreachable")
}

// NewOnset constructs a new Onset object.
// It is the Callers responsibility to call Free on the returned
// Onset object or leak memory.
//     t, err := NewOnset(mode, bufSize, blockSize, samplerate)
//     if err != nil {
//         // handle error
//     }
//     defer t.Free()
func NewOnset(
	onset_mode onsetMode, bufSize, blockSize, samplerate uint) (*Onset, error) {
	t, err := C.new_aubio_onset(toCharTPtr(string(onset_mode)),
		C.uint_t(bufSize), C.uint_t(blockSize), C.uint_t(samplerate))
	if t == nil {
		return nil, fmt.Errorf("Failure creating Onset object %q", err)
	}
	return &Onset{o: t, buf: NewSimpleBuffer(blockSize)}, nil
}

func (t *Onset) Buffer() *SimpleBuffer {
	return t.buf
}

// Do executes the onset detection on an input Buffer.
// It returns the estimated beat locations in a new buffer.
func (t *Onset) Do(input *SimpleBuffer) {
	if t.o == nil {
		return
	}
	C.aubio_onset_do(t.o, input.vec, t.buf.vec)
}

// SetSilence sets the onset detection silence threshold.
func (t *Onset) SetSilence(silence float64) {
	if t.o == nil {
		return
	}
	C.aubio_onset_set_silence(t.o, C.smpl_t(silence))
}

// SetThreshold sets the onset detection peak picking threshold.
func (t *Onset) SetThreshold(threshold float64) {
	if t.o == nil {
		return
	}
	C.aubio_onset_set_threshold(t.o, C.smpl_t(threshold))
}

/* TODO(jwall): Update to the latest version of the aubio library
// GetLastOnset returns the bpm after running Do on an input Buffer
//     t, err := NewOnset(mode, bufSize, blockSize, samplerate)
//      if err != nil {
//      }
//      defer t.Close()
//      t.Do(buf)
//      fmt.Println("BPM: ", t.GetBpm())
func (t *Onset) GetLastOnset() float64 {
	if t.o == nil {
		return 0
	}
	return float64(C.aubio_onset_get_last_onset(t.o))
}
*/

// Free frees the aubio_temp_t object's memory.
func (t *Onset) Free() {
	if t.o == nil {
		return
	}
	C.del_aubio_onset(t.o)
	t.o = nil
}
