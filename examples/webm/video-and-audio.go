package main

import (
	"fmt"
	"log"

	"github.com/3d0c/gmf"
)

// ffmpeg -i input/test.mp4 -c:v libx264 -c:a aac output/transcoded.mp4

func main() {
	inputFile := "input/test.mp4"
	outputFile := "output/try.mp4"
	inputCtx, err := gmf.NewInputCtx(inputFile)
	if err != nil {
		log.Fatalf("unable to create inputCtx: %v", err)
	}
	defer inputCtx.Free()
	outputCtx, err := gmf.NewOutputCtx(outputFile)
	if err != nil {
		log.Fatalf("unable to create outputCtx: %v", err)
	}
	defer outputCtx.Free()
	inputVideoStream, err := inputCtx.GetBestStream(gmf.AVMEDIA_TYPE_VIDEO)
	if err != nil {
		log.Fatalf("unalbe to get video stream from inputCtx: %v", err)
	}
	defer inputVideoStream.Free()
	inputAudioStream, err := inputCtx.GetBestStream(gmf.AVMEDIA_TYPE_AUDIO)
	if err != nil {
		log.Fatalf("unalbe to get video stream from inputCtx: %v", err)
	}
	defer inputAudioStream.Free()
	// create video and audio output streams
	outputVideoStream, err := createStream("libx264", outputCtx, inputVideoStream)
	if err != nil {
		log.Fatalf("unable to create stream in outputCtx: %v", err)
	}
	defer outputVideoStream.Free()
	defer outputVideoStream.CodecCtx().Free()

	outputAudioStream, err := createStream("aac", outputCtx, inputVideoStream)
	if err != nil {
		log.Fatalf("unable to create stream in outputCtx: %v", err)
	}
	defer outputAudioStream.Free()
	defer outputAudioStream.CodecCtx().Free()

	if err := outputCtx.WriteHeader(); err != nil {
		log.Fatalf("unable to write header: %v", err)
	}
	var inputStream *gmf.Stream
	var outputStream *gmf.Stream
	for packet := range inputCtx.GetNewPackets() {
		if packet == nil {
			// flush/drain?
		}
		// we are only doing video for now
		switch packet.StreamIndex() {
		case inputVideoStream.Index():
			inputStream = inputVideoStream
			outputStream = outputVideoStream
		case inputAudioStream.Index():
			inputStream = inputAudioStream
			outputStream = outputAudioStream
		default:
			continue
		}

		frames, err := inputStream.CodecCtx().Decode(packet)
		packet.Free()
		if err != nil {
			log.Fatalf("unable to decode packet: %v", err)
		}
		// todo try this by creating destination frame that doesn't require drain/flush argument
		encodedPackets, err := outputStream.CodecCtx().Encode(frames, -1)
		if err != nil {
			log.Fatalf("unable to encode decoded frames: %v", err)
		}
		for _, encodedPacket := range encodedPackets {
			// doesn't work properly without this
			gmf.RescaleTs(encodedPacket, outputStream.CodecCtx().TimeBase(), outputStream.TimeBase())
			// not doing anything
			encodedPacket.SetStreamIndex(outputStream.Index())

			if err := outputCtx.WritePacket(encodedPacket); err != nil {
				// exit on error?
				log.Fatalf("unable to write encodedPacket to ouputCtx: %v", err)
			}
			encodedPacket.Free()
		}
		for _, frame := range frames {
			if frame != nil {
				frame.Free()
			}
		}
	}
	// flush/drain packets from buffer
	{
		// video
		// todo try this by creating destination frame that doesn't require drain/flush argument
		encodedPackets, err := outputVideoStream.CodecCtx().Encode([]*gmf.Frame{}, 1)
		if err != nil {
			log.Printf("unable to encode decoded frames: %v", err)
		}
		log.Printf("flushing %v packets", len(encodedPackets))
		for _, encodedPacket := range encodedPackets {
			// doesn't work properly without this
			gmf.RescaleTs(encodedPacket, outputVideoStream.CodecCtx().TimeBase(), outputVideoStream.TimeBase())
			// not doing anything
			encodedPacket.SetStreamIndex(outputVideoStream.Index())

			if err := outputCtx.WritePacket(encodedPacket); err != nil {
				// exit on error?
				log.Fatalf("unable to write encodedPacket to ouputCtx: %v", err)
			}
			encodedPacket.Free()
		}
		// before flush?
		// outputVideoStream.CodecCtx().FlushBuffers()

		// audio
		// todo try this by creating destination frame that doesn't require drain/flush argument
		encodedPackets, err = outputAudioStream.CodecCtx().Encode([]*gmf.Frame{}, 1)
		if err != nil {
			log.Printf("unable to encode decoded frames: %v", err)
		}
		log.Printf("flushing %v packets", len(encodedPackets))
		for _, encodedPacket := range encodedPackets {
			// doesn't work properly without this
			gmf.RescaleTs(encodedPacket, outputAudioStream.CodecCtx().TimeBase(), outputAudioStream.TimeBase())
			// not doing anything
			encodedPacket.SetStreamIndex(outputAudioStream.Index())

			if err := outputCtx.WritePacket(encodedPacket); err != nil {
				// exit on error?
				log.Fatalf("unable to write encodedPacket to ouputCtx: %v", err)
			}
			encodedPacket.Free()
		}
		// before flush?
		// outputVideoStream.CodecCtx().FlushBuffers()

	}
	outputCtx.WriteTrailer()
}

func createStream(codecName string, outputCtx *gmf.FmtCtx, inputStream *gmf.Stream) (newStream *gmf.Stream, err error) {
	codec, err := gmf.FindEncoder(codecName)
	if err != nil {
		return nil, fmt.Errorf("unable to find codec for %v: %w", codecName, err)

	}
	// create and set codecCtx options
	codecCtx := gmf.NewCodecCtx(codec)
	if codecCtx == nil {
		return nil, fmt.Errorf("unable create codecCtx: %w", err)
	}
	if outputCtx.IsGlobalHeader() {
		outputCtx.SetFlag(gmf.CODEC_FLAG_GLOBAL_HEADER)
	}
	if codec.IsExperimental() {
		codecCtx.SetStrictCompliance(gmf.FF_COMPLIANCE_EXPERIMENTAL)
	}
	var options []gmf.Option
	if codecCtx.Type() == gmf.AVMEDIA_TYPE_AUDIO {
		sampleRate := inputStream.CodecCtx().SampleRate()
		if sampleRate <= 0 {
			sampleRate = codecCtx.SelectSampleRate()
		}
		channelLayout := inputStream.CodecCtx().ChannelLayout()
		if channelLayout <= 0 {
			channelLayout = codecCtx.SelectChannelLayout()
		}
		sampleFormat := inputStream.CodecCtx().SampleFmt()
		if sampleFormat <= 0 {
			sampleFormat = codecCtx.SelectSampleFmt()
		}
		// https://ffmpeg.org/doxygen/3.1/muxing_8c-example.html
		options = []gmf.Option{
			{Key: "time_base", Val: gmf.AVR{Num: 1, Den: sampleRate}},
			{Key: "ar", Val: sampleRate},
			{Key: "ac", Val: inputStream.CodecCtx().Channels()},
			{Key: "channel_layout", Val: channelLayout},
		}
		log.Print(options)
		codecCtx.SetSampleFmt(sampleFormat)
		// codecCtx.SetSampleFmt(gmf.AV_SAMPLE_FMT_FLTP)
		// codecCtx.SelectSampleRate()
	}
	if codecCtx.Type() == gmf.AVMEDIA_TYPE_VIDEO {
		options = []gmf.Option{
			// {Key: "time_base", Val: gmf.AVR{Num: 1, Den: 25}},
			{Key: "time_base", Val: inputStream.CodecCtx().TimeBase().AVR()},
			{Key: "pixel_format", Val: gmf.AV_PIX_FMT_YUV420P},
			// Save original
			{Key: "video_size", Val: inputStream.CodecCtx().GetVideoSize()},
			{Key: "b", Val: 500000},
		}
		codecCtx.SetProfile(inputStream.CodecCtx().GetProfile())
	}
	codecCtx.SetOptions(options)
	if err := codecCtx.Open(nil); err != nil {
		return nil, fmt.Errorf("unable to open codecCtx: %w", err)
	}
	// create new stream in ouputCtx and set parameters from codecCtx
	var outputStream *gmf.Stream
	if outputStream = outputCtx.NewStream(codec); outputStream == nil {
		return nil, fmt.Errorf("unable to create new stream in output context")
	}
	params := gmf.NewCodecParameters()
	if err = params.FromContext(codecCtx); err != nil {
		return nil, fmt.Errorf("error creating codec parameters from context: %w", err)
	}
	defer params.Free()
	outputStream.CopyCodecPar(params)
	outputStream.SetCodecCtx(codecCtx)
	// can we not set frame rate from inputStream
	if codecCtx.Type() == gmf.AVMEDIA_TYPE_VIDEO {
		outputStream.SetTimeBase(gmf.AVR{Num: 1, Den: 25})
		outputStream.SetRFrameRate(gmf.AVR{Num: 25, Den: 1})
	}
	return outputStream, nil
}

// 2 frames missing even after drain/flush
// ffmpeg -i output/try.mp4 -vcodec copy -acodec copy -f null /dev/null 2>&1 | grep 'frame='
// ffmpeg -i input/test.mp4 -vcodec copy -acodec copy -f null /dev/null 2>&1 | grep 'frame='
