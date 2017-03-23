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

// fft

// filterbank

// mfcc

// phasvoc

type PhaseVoc struct {
	o     *C.aubio_pvoc_t
	buf   *SimpleBuffer
	grain *ComplexBuffer
}

func NewPhaseVoc(bufSize, fftLen uint) (*PhaseVoc, error) {
	pvoc, err := C.new_aubio_pvoc(C.uint_t(bufSize), C.uint_t(fftLen))
	if err != nil {
		return nil, err
	}
	return &PhaseVoc{o: pvoc, grain: NewComplexBuffer(fftLen)}, nil
}

func (pv *PhaseVoc) Free() {
	if pv.o != nil {
		C.del_aubio_pvoc(pv.o)
		pv.o = nil
	}
	if pv.grain != nil {
		pv.grain.Free()
		pv.grain = nil
	}
}

func (pv *PhaseVoc) Grain() *ComplexBuffer {
	return pv.grain
}

func (pv *PhaseVoc) Do(in *SimpleBuffer) {
	if pv.o != nil {
		C.aubio_pvoc_do(pv.o, in.vec, pv.grain.data)
	} else {
		log.Println("Called Do on empty PhaseVoc. Maybe you called Free previously?")
	}
}

func (pv *PhaseVoc) ReverseDo(out *SimpleBuffer) {
	if pv.o != nil {
		C.aubio_pvoc_rdo(pv.o, pv.grain.data, out.vec)
	} else {
		log.Println("Called ReverseDo on empty PhaseVoc. Maybe you called Free previously?")
	}
}

// statistics

// tss
