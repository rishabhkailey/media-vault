package main

import (
	"log"
	"os"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// ffmpeg -i input/test.mp4 -c:v libvpx -f webm -movflags faststart pipe:
// ffmpeg -i input/test.mp4 -c:a libopus -c:v libvpx -f webm pipe: -y
func main() {
	outputFile, err := os.Create("output/test.webm")
	if err != nil {
		log.Fatal("os.Create()", err)
	}
	// inputFile, err := os.Open("input/test.mp4")
	// if err != nil {
	// 	log.Fatal("os.Open()", err)
	// }
	// err := ffmpeg.Input("input/test.mp4").
	// 	Output("output/test.mp4", ffmpeg.KwArgs{"c:v": "libx265", "c:a": "aac"}).
	// 	OverWriteOutput().ErrorToStdOut().Run()
	err = ffmpeg.Input("input/test.mp4").
		Output("pipe:", ffmpeg.KwArgs{
			"c:v":      "libvpx",
			"c:a":      "libopus",
			"f":        "webm",
			"movflags": "faststart",
		}).
		WithOutput(outputFile).
		OverWriteOutput().
		ErrorToStdOut().
		Run()
	if err != nil {
		log.Fatal("Transcoding failed", err)
	}
}

// ffmpeg -i input/test.mp4 -c:v libx265 -c:a aac output/transcoded.mp4
// ffmpeg -i input/test.mp4 -c:a aac -c:v libx265 output/test.mp4
// ffmpeg -i input/test.mp4 -c:v libx265 -c:a aac

// :)
// easy to use
// good documentations and features

// :(
// uses ffmpeg pipe: so no faststart with custom io.reader/writers
