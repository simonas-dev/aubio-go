/*
 Copyright 2013 Jeremy Wall (jeremy@marzhillstudios.com)

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0
*/

// sink is an aubio example application.
// Run it: sink --src=file.wav
package main

import (
	"flag"
	"log"

	"go.marzhillstudios.com/pkg/play/aubio"
)

var (
	srcPath    = flag.String("src", "", "Path to source file")
	sinkPath   = flag.String("sink", "", "Path to sink file")
	samplerate = flag.Int("samplerate", 0, "Sample rate to use for the audio file")
	blockSize  = flag.Int("blocksize", 256, "Blocksize use for the audio file")
	bufSize    = flag.Int("bufsize", 512, "Blocksize use for the audio file")
)

func main() {
	flag.Parse()

	if *srcPath == "" {
		log.Fatal("Must provide a src")
	}
	if *sinkPath == "" {
		log.Fatal("Must provide a sink")
	}
	p, err := aubio.PipelineFromUris(*srcPath, *sinkPath,
		uint(*samplerate), uint(*bufSize), uint(*blockSize))
	if err != nil {
		log.Fatalf("err: %s", err)
	}
	defer p.Close()
	log.Println("Wrote: ", p.DoAll())
}
