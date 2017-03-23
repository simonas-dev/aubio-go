/*
 Copyright 2013 Jeremy Wall (jeremy@marzhillstudios.com)

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0
*/

// Package aubio is a Go binding to the aubio audio analysis library
// http://aubio.org/.
package aubio

/*
#cgo LDFLAGS: -laubio
#include <aubio/aubio.h>
*/
import "C"

import (
	"fmt"
	"log"
	"runtime"
	"syscall"
)

func newSink(uri string, sr uint) (*C.aubio_sink_t, error) {
	sink, err := C.new_aubio_sink(
		toCharTPtr(uri), C.uint_t(sr))
	return sink, err
}

func newSource(uri string, sr, hopSize uint) (*C.aubio_source_t, error) {
	src, err := C.new_aubio_source(
		toCharTPtr(uri), C.uint_t(sr), C.uint_t(hopSize))
	return src, err
}

// Source is a wrapper for an aubio_source_t object.
type Source struct {
	blockSize uint
	s         *C.aubio_source_t
}

// OpenSource opens an aubio_source_t from the uri.
// It uses the given samplerate and hopSize for processing
// the audio source stream.
//
// The caller is responsible for calling close on
// the returned Source to release memory.
//
//     s := OpenSource(uri, 44100, 1024)
//     defer s.Close()
func OpenSource(uri string, samplerate, hopSize uint) (*Source, error) {
	src, err := newSource(uri, samplerate, hopSize)
	if src == nil {
		return nil, fmt.Errorf("Failed to open source uri %q %s errno: %d", uri, err,
			int(err.(syscall.Errno)))
	}
	return &Source{
		blockSize: hopSize,
		s:         src,
	}, nil
}

// BlockSize returns the blockSize used by this Source.
func (s *Source) BlockSize() (n uint) {
	return s.blockSize
}

// Samplerate returns the sample rate of a Source.
func (s *Source) Samplerate() (n uint) {
	s.ifOpen(func() {
		n = uint(C.aubio_source_get_samplerate(s.s))
	})
	return
}

func (s *Source) ifOpen(f func()) {
	if s.s != nil {
		f()
	} else {
		if pc, _, _, ok := runtime.Caller(1); ok {
			log.Printf("Called %s on Closed Sink", runtime.FuncForPC(pc).Name())
		}
	}
}

// Do reads from a source into a buffer.
// It returns the amount of data read.
func (s *Source) Do(buf *SimpleBuffer) uint {
	var n C.uint_t = 0
	s.ifOpen(func() {
		C.aubio_source_do(s.s, buf.vec, &n)
	})
	return uint(n)
}

// Close closes the aubio_source_t and frees the memory.
func (s *Source) Close() {
	s.ifOpen(func() { C.del_aubio_source(s.s) })
	s.s = nil
}

// Sink is a wrapper for an aubio_sink_t object.
type Sink struct {
	samplerate uint
	s          *C.aubio_sink_t
}

// OpenSink opens an aubio_sink_t from the uri.
// It uses the samplerate to write data to the sink.
//
// The caller is responsible for calling close on
// the returned Sink to release memory.
//
//     s := OpenSink(uri, 44100)
//     defer s.Close()
func OpenSink(uri string, samplerate uint) (*Sink, error) {
	sink, err := newSink(uri, samplerate)
	if sink == nil {
		return nil, fmt.Errorf("Failed to open source uri %q %s errno: %d", uri, err,
			int(err.(syscall.Errno)))
	}
	return &Sink{
		samplerate: samplerate,
		s:          sink,
	}, nil
}

func (s *Sink) ifOpen(f func()) {
	if s.s != nil {
		f()
	} else {
		if pc, _, _, ok := runtime.Caller(1); ok {
			log.Printf("Called %s on Closed Sink", runtime.FuncForPC(pc).Name())
		}
	}
}

// Samplerate returns the samplerate for this Sink.
func (s *Sink) Samplerate() uint {
	return s.samplerate
}

// Close closes the aubio_sink_t and frees the memory.
func (s *Sink) Close() {
	s.ifOpen(func() { C.del_aubio_sink(s.s) })
	s.s = nil
}

// Do writes to the sink from the buffer.
// It returns the amount of data written.
func (s *Sink) Do(buf *SimpleBuffer, n uint) uint {
	s.ifOpen(func() { C.aubio_sink_do(s.s, buf.vec, C.uint_t(n)) })
	return n
}

// Pipeline pipes data from a Source to a Sink.
type SimplePipeline struct {
	source *Source
	sink   *Sink
	buf    *SimpleBuffer
}

// NewPipeline constructs a Pipeline between a Source and an optional Sink
// using a Buffer of bufSize.  It will run any ProcessFuncs passed to it as it
// pulls from the source.
//
// The Pipeline assumes ownership of the source and sink so calling Close on
// the Pipeline will close the source as well as the sink.
//
//     pitch := NewPitch(...)
//     fn := func(in, out *Buffer) {
//        pitch.Do(in, out)
//        // do something with that data.
//     }
//     p := NewPipeline(OpenSource(sourceUri, samplerate, hopSize),
//                      OpenSink(sinkUri, samplerate),
//                      bufSize, fn)
//     defer p.Close()
//     p.DoAll() // pipe all the data in source to the sink
func NewSimplePipeline(in *Source, out *Sink, bufSize uint) *SimplePipeline {
	return &SimplePipeline{
		buf:    NewSimpleBuffer(bufSize),
		source: in,
		sink:   out,
	}
}

// PipelineFromUris constructs a pipe from two URIs.
// If samplerate is 0 it detects the samplerate from the source and
// uses that for the sink.
func PipelineFromUris(inUri, outUri string, samplerate, blockSize, bufSize uint) (*SimplePipeline, error) {
	src, err := OpenSource(inUri, samplerate, blockSize)
	if src == nil {
		return nil, err
	}
	if samplerate == 0 {
		samplerate = src.Samplerate()
	}
	sink, err := OpenSink(outUri, samplerate)
	if sink == nil {
		defer src.Close()
		return nil, err
	}
	return NewSimplePipeline(src, sink, bufSize), nil
}

// Close closes the the Source, Sink, and frees the Buffer.
func (p *SimplePipeline) Close() {
	p.source.Close()
	p.source = nil
	if p.sink != nil {
		p.sink.Close()
		p.sink = nil
	}
	p.buf.Free()
	p.buf = nil
}

// BlockSize returns the BlockSize used by this Pipeline.
func (p *SimplePipeline) BlockSize() uint {
	return p.source.BlockSize()
}

// BufSize returns the current buffer size the Pipeline is using.
func (p *SimplePipeline) BufSize() uint {
	return uint(p.buf.vec.length)
}

func (p *SimplePipeline) do(fs []ProcessFunc) (amt uint) {
	n := p.source.Do(p.buf)
	for _, f := range fs {
		f(p.buf)
	}
	if p.sink != nil {
		return p.sink.Do(p.buf, n)
	}
	return n
}

// Do pipes up to BufSize data from the source to a sink if there is
// one.
// It returns the number of frames processed.
func (p *SimplePipeline) Do(fs ...ProcessFunc) uint {
	out := NewSimpleBuffer(p.BufSize())
	defer out.Free()
	return p.do(fs)
}

// DoN runs Do up to n times.
// It returns the number of frames processed.
func (p *SimplePipeline) DoN(n int, fs ...ProcessFunc) (total uint) {
	total = uint(0)
	read := p.BlockSize()
	for i := 0; i < n && read != p.BlockSize(); i++ {
		// TODO(jwall): Is it safe to share the buffer?
		out := NewSimpleBuffer(p.BufSize())
		defer out.Free()
		read = p.do(fs)
		total += read
	}
	return
}

// DoAll runs Do until the source has been exhausted.
// It returns the number of frames processed.
func (p *SimplePipeline) DoAll(fs ...ProcessFunc) (total uint) {
	read := p.BlockSize()
	for read == p.BlockSize() {
		// TODO(jwall): Is it safe to share the buffer?
		out := NewSimpleBuffer(p.BufSize())
		defer out.Free()
		read = p.do(fs)
		total += read
	}
	return
}
