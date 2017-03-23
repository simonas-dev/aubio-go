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
	"fmt"
)

// Tempo is a wrapper for the aubio_tempo_t tempo detection object.
type Tempo struct {
	o   *C.aubio_tempo_t
	buf *SimpleBuffer
}

// TempoOrDie constructs a new Tempo object.
// It panics on any errors.
func TempoOrDie(mode onsetMode, bufSize, blocksize, samplerate uint) *Tempo {
	if t, err := NewTempo(mode, bufSize, blocksize, samplerate); err == nil {
		return t
	} else {
		panic(err)
	}
	panic("Unreachable")
}

// NewTempo constructs a new Tempo object.
// It is the Callers responsibility to call Free on the returned
// Tempo object or leak memory.
//     t, err := NewTempo(mode, bufSize, blockSize, samplerate)
//     if err != nil {
//         // handle error
//     }
//     defer t.Free()
func NewTempo(
	mode onsetMode, bufSize, blockSize, samplerate uint) (*Tempo, error) {
	t, err := C.new_aubio_tempo(toCharTPtr(string(mode)),
		C.uint_t(bufSize), C.uint_t(blockSize), C.uint_t(samplerate))
	if t == nil {
		return nil, fmt.Errorf("Failure creating Tempo object %q", err)
	}
	return &Tempo{o: t, buf: NewSimpleBuffer(blockSize)}, nil
}

func (t *Tempo) Buffer() *SimpleBuffer {
	return t.buf
}

// Do executes the tempo detection on an input Buffer.
// It returns the estimated beat locations in a new buffer.
func (t *Tempo) Do(input *SimpleBuffer) {
	if t.o == nil {
		return
	}
	C.aubio_tempo_do(t.o, input.vec, t.buf.vec)
}

// SetSilence sets the tempo detection silence threshold.
func (t *Tempo) SetSilence(silence float64) {
	if t.o == nil {
		return
	}
	C.aubio_tempo_set_silence(t.o, C.smpl_t(silence))
}

// SetThreshold sets the tempo detection peak picking threshold.
func (t *Tempo) SetThreshold(threshold float64) {
	if t.o == nil {
		return
	}
	C.aubio_tempo_set_threshold(t.o, C.smpl_t(threshold))
}

// GetBpm returns the bpm after running Do on an input Buffer
//     t, err := NewTempo(mode, bufSize, blockSize, samplerate)
//      if err != nil {
//      }
//      defer t.Close()
//      t.Do(buf)
//      fmt.Println("BPM: ", t.GetBpm())
func (t *Tempo) GetBpm() float64 {
	if t.o == nil {
		return 0
	}
	return float64(C.aubio_tempo_get_bpm(t.o))
}

// GetConfidence returns the confidence after running Do on an input Buffer
//     t, err := NewTempo(mode, bufSize, blockSize, samplerate)
//      if err != nil {
//      }
//      defer t.Close()
//      t.Do(buf)
//      fmt.Println("Confidence: ", t.GetConfidence())
func (t *Tempo) GetConfidence() float64 {
	if t.o == nil {
		return 0
	}
	return float64(C.aubio_tempo_get_confidence(t.o))
}

// Free frees the aubio_temp_t object's memory.
func (t *Tempo) Free() {
	if t.o == nil {
		return
	}
	C.del_aubio_tempo(t.o)
	t.o = nil
}

/* Only available in AUBIO_UNSTABLE

// BeatTracker is a wrapper for the aubio_beattracking_t beattracking
// detection object. See https://github.com/piem/aubio/blob/develop/src/tempo/beattracking.h
// for more details.
type BeatTracker struct {
	o *C.aubio_beattracking_t
}

// NewBeatTracker constructs a new BeatTracker object.
// It is the Callers responsibility to call Free on the returned
// BeatTracker object or leak memory.
//     t, err := NewBeatTracker(mode, bufSize, blockSize, samplerate)
//     if err != nil {
//         // handle error
//     }
//     defer t.Free()
func NewBeatTracker(blockSize uint) (*BeatTracker, error) {
	t, err := C.new_aubio_beattracking(C.uint_t(blockSize))
    if t == nil {
		return nil, fmt.Errorf("Failure creating BeatTracker object %q", err)
	}
	return &BeatTracker{t}, nil
}

// Do executes the beattracking detection on an input Buffer.
// It returns the estimated beat locations in a new buffer.
func (t *BeatTracker) Do(input *Buffer, out *Buffer) {
	if t.o == nil {
		return
	}
	C.aubio_beattracking_do(t.o, input.vec, out.vec)
}
// GetBpm returns the bpm after running Do on an input Buffer
//     t, err := NewBeatTracker(mode, bufSize, blockSize, samplerate)
//      if err != nil {
//      }
//      defer t.Close()
//      t.Do(buf)
//      fmt.Println("BPM: ", t.GetBpm())
func (t *BeatTracker) GetBpm() float64 {
	if t.o == nil {
		return 0
	}
	return float64(C.aubio_beattracking_get_bpm(t.o))
}

// GetConfidence returns the confidence after running Do on an input Buffer
//     t, err := NewBeatTracker(mode, bufSize, blockSize, samplerate)
//      if err != nil {
//      }
//      defer t.Close()
//      t.Do(buf)
//      fmt.Println("Confidence: ", t.GetConfidence())
func (t *BeatTracker) GetConfidence() float64 {
	if t.o == nil {
		return 0
	}
	return float64(C.aubio_beattracking_get_confidence(t.o))
}

// Free frees the aubio_temp_t object's memory.
func (t *BeatTracker) Free() {
	if t.o == nil {
		return
	}
	C.del_aubio_beattracking(t.o)
	t.o = nil
}
*/ // AUBIO_UNSTABLE
