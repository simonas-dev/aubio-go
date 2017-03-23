/*
 Copyright 2013 Jeremy Wall (jeremy@marzhillstudios.com)

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0
*/

// tempo is an aubio example application.
// Run it: tempo --src=file.wav
package main

import (
	"fmt"

	"go.marzhillstudios.com/pkg/play/aubio"
	"go.marzhillstudios.com/pkg/play/aubio/examples/util"
)

func main() {
	src := util.Init()
	oa := aubio.OnsetOrDie(aubio.SpecDiff, uint(*util.Bufsize),
		uint(*util.Blocksize), uint(*util.Samplerate))
	oa.SetSilence(*util.Silence)
	oa.SetThreshold(*util.Threshold)
	ch := make(chan float64)
	p := aubio.NewSimplePipeline(src, nil, uint(*util.Bufsize))
	defer p.Close()
	go func() {
		n := p.DoAll(func(input *aubio.SimpleBuffer) {
			oa.Do(input)
			for _, f := range oa.Buffer().Slice() {
				ch <- f
			}
		})
		close(ch)
		fmt.Println("Processed:", n)
	}()
	for f := range ch {
		if f != 0 || *util.Verbose {
			fmt.Printf("Onset %.6f\n", f)
		}
	}
}
