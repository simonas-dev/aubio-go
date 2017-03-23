/*
 Copyright 2013 Jeremy Wall (jeremy@marzhillstudios.com)

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0
*/

package util

import (
	"flag"
	"fmt"
	"log"
	"os"

	"go.marzhillstudios.com/pkg/play/aubio"
)

var (
	srcPath    = flag.String("src", "", "Path to source file. Required")
	Samplerate = flag.Int("samplerate", 44100, "Sample rate to use for the audio file")
	Blocksize  = flag.Int("blocksize", 256, "Blocksize use for the audio file")
	Bufsize    = flag.Int("bufsize", 512, "Bufsize use for the audio file")
	Silence    = flag.Float64("silence", -90.0, "Threshold to use when detecting silence")
	Threshold  = flag.Float64("threshold", 0.0, "Detection threshold")
	Verbose    = flag.Bool("verbose", false, "Print verbose output")
	help       = flag.Bool("help", false, "Print this help")
)

func Init() *aubio.Source {
	flag.Parse()
	if *help {
		fmt.Println("usage:", os.Args[0], "--src=path/to/somefile.wav")
		fmt.Println("")
		fmt.Println("Flags:")
		os.Stdout.Sync()
		flag.PrintDefaults()
		os.Exit(0)
	}
	if *srcPath == "" {
		flag.PrintDefaults()
		log.Fatal("Must provide a src")
	}
	if *Verbose {
		fmt.Println("Input file: ", *srcPath)
	}
	src, err := aubio.OpenSource(*srcPath, uint(*Samplerate), uint(*Blocksize))
	if err != nil {
		log.Fatal(err)
	}
	return src
}
