/*
 Copyright 2013 Jeremy Wall (jeremy@marzhillstudios.com)

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0
*/

// pitch is an aubio example application.
// Run it: pitch --src=file.wav
package main

import (
	"fmt"

	"go.marzhillstudios.com/pkg/play/aubio"
	"go.marzhillstudios.com/pkg/play/aubio/examples/util"
)

func main() {
	src := util.Init()
	pitch := aubio.NewPitch(aubio.PitchDefault, uint(*util.Bufsize), uint(*util.Blocksize), uint(*util.Samplerate))
	//pitch.SetUnit(aubio.PitchOutDefault)
	pitch.SetTolerance(0.7)
	p := aubio.NewSimplePipeline(src, nil, uint(*util.Bufsize))
	defer p.Close()
	ch := make(chan float64)
	go func() {
		n := p.DoAll(func(in *aubio.SimpleBuffer) {
			pitch.Do(in)
			for _, f := range pitch.Buffer().Slice() {
				ch <- f
			}
		})
		close(ch)
		if *util.Verbose {
			fmt.Println("Processed", n, "frames")
		}
	}()
	for f := range ch {
		if f != 0 || *util.Verbose {
			fmt.Printf("pitch %.6f\n", f)
		}
	}
}
