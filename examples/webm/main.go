package main

// import (
// 	"io"
// 	"log"

// 	"github.com/3d0c/gmf"
// )

// // ffmpeg vs gmf
// // inputCtx (used for reading data) -> inputStream (used for metadata and codec info)
// // videoStream := ctx.GetBestStream(gmf.AVMEDIA_TYPE_VIDEO)
// // audioStream := ctx.GetBestStream(gmf.AVMEDIA_TYPE_AUDIO)
// // ctx.getStream(videoStream.Index()) == videoStream
// // outputCtx (for writing data) -> outputStream (setting metada and codec info)
// // outputVideoStream := outputCtx.GetBestStream(gmf.AVMEDIA_TYPE_VIDEO)
// // // packet in input/source format
// // pkt, err = inputCtx.GetNextPacket()
// // frames, err := videoStream.CodecCtx().Decode(pkt)
// // // packets in output/required format
// // packets, err := outputVideoStream.CodecCtx().Encode(frames, flush)
// // for _, packet := range packets {
// // 	// packet, source timebase, dest timebase
// // 	gmf.RescaleTs(packet, ost.CodecCtx().TimeBase(), ost.TimeBase())
// // 	// tell if the packet is for video or audio
// // 	packet.SetStreamIndex(outputVideoStream.Index())
// // 	if err = outputCtx.WritePacket(packet); err != nil {
// // 		break
// // 	}
// // 	packet.Free()
// // }

// // ffmpeg -i tears_of_steel_1080outputPacket.webm -c:v libvpx-vp9 -c:a libopus output.webm
// func main() {
// 	// inputCtx, err := gmf.NewInputCtxWithFormatName("testsrc=decimals=3", "lavfi")
// 	inputCtx, err := gmf.NewInputCtx("input/test.mp4")
// 	if err != nil {
// 		log.Fatal("NewInputCtxWithFormatName", err)
// 		return
// 	}
// 	inputStream, err := inputCtx.GetBestStream(gmf.AVMEDIA_TYPE_VIDEO)
// 	if err != nil {
// 		log.Fatal("GetBestStream", err)
// 		return
// 	}
// 	videoEncoder, err := gmf.FindEncoder("libx264")
// 	// videoEncoder, err := gmf.FindEncoder(gmf.AV_CODEC_ID_H264)
// 	if err != nil {
// 		log.Fatal("videoEncoder", err)
// 		return
// 	}
// 	videoEncCtx := gmf.NewCodecCtx(videoEncoder)
// 	if videoEncCtx == nil {
// 		log.Fatal("NewCodecCtx", videoEncCtx)
// 		return
// 	}
// 	var options []gmf.Option
// 	if videoEncCtx.Type() == gmf.AVMEDIA_TYPE_VIDEO {
// 		options = append(
// 			[]gmf.Option{
// 				{Key: "time_base", Val: gmf.AVR{Num: 1, Den: 25}},
// 				{Key: "pixel_format", Val: gmf.AV_PIX_FMT_YUV420P},
// 				// Save original
// 				{Key: "video_size", Val: inputStream.CodecCtx().GetVideoSize()},
// 				{Key: "b", Val: 500000},
// 			},
// 		)

// 		videoEncCtx.SetProfile(inputStream.CodecCtx().GetProfile())
// 	}
// 	videoEncCtx.SetOptions(options)
// 	videoEncCtx.
// 		SetWidth(inputStream.CodecCtx().Width()).
// 		SetHeight(inputStream.CodecCtx().Height()).
// 		SetTimeBase(inputStream.TimeBase().AVR()).
// 		SetPixFmt(gmf.AV_PIX_FMT_YUV420P)
// 	outputCtx, err := gmf.NewOutputCtx("output/transcoded.mp4")
// 	if err != nil {
// 		log.Fatal("NewOutputCtx", err)
// 		return
// 	}
// 	if outputCtx.IsGlobalHeader() {
// 		outputCtx.SetFlag(gmf.CODEC_FLAG_GLOBAL_HEADER)
// 	}
// 	outputStream := outputCtx.NewStream(videoEncoder)
// 	if outputStream == nil {
// 		log.Fatal("outputCtx.NewStream", nil)
// 		return
// 	}
// 	if err := videoEncCtx.Open(nil); err != nil {
// 		log.Fatal("videoEncCtx.Open ", err)
// 		return
// 	}

// 	par := gmf.NewCodecParameters()
// 	if err = par.FromContext(videoEncCtx); err != nil {
// 		log.Fatal("videoEncCtx.Open ", err)
// 	}
// 	outputStream.CopyCodecPar(par)
// 	outputStream.SetCodecCtx(videoEncCtx)
// 	outputStream.SetTimeBase(gmf.AVR{Num: 1, Den: 25})
// 	outputStream.SetRFrameRate(gmf.AVR{Num: 25, Den: 1})

// 	dstFrame := gmf.NewFrame().
// 		SetWidth(inputStream.CodecCtx().Width()).
// 		SetHeight(inputStream.CodecCtx().Height()).
// 		SetFormat(gmf.AV_PIX_FMT_YUV420P)

// 	if err := dstFrame.ImgAlloc(); err != nil {
// 		log.Fatal("dstFrame.ImgAlloc", err)
// 	}

// 	var (
// 		pkt *gmf.Packet
// 		// streamIdx int
// 		// flush int   = -1
// 		drain int   = -1
// 		pts   int64 = 0
// 	)

// 	for {
// 		if drain >= 0 {
// 			break
// 		}

// 		pkt, err = inputCtx.GetNextPacket()
// 		if err != nil && err != io.EOF {
// 			if pkt != nil {
// 				pkt.Free()
// 			}
// 			log.Printf("error getting next packet - %s", err)
// 			break
// 		} else if err != nil && pkt == nil {
// 			drain = 0
// 		}
// 		defer pkt.Free()

// 		if pkt != nil && pkt.StreamIndex() != inputStream.Index() {
// 			continue
// 		}

// 		frames, err := inputStream.CodecCtx().Decode(pkt)
// 		if err != nil {
// 			log.Fatalf("error decoding - %s\n", err)
// 		}
// 		for _, frame := range frames {
// 			frame.SetPts(pts)
// 			pts++
// 		}
// 		packets, err := outputStream.CodecCtx().Encode(frames, drain)
// 		if err != nil {
// 			log.Print("outputStream.CodecCtx().Encode(frames, flush) ", err)
// 		}
// 		if len(packets) > 0 {
// 			log.Printf("encoding frame at ", pts)
// 		} else {
// 			log.Printf("skipping frame at ", pts)
// 		}
// 		for _, packet := range packets {
// 			gmf.RescaleTs(packet, outputStream.CodecCtx().TimeBase(), outputStream.TimeBase())
// 			packet.SetStreamIndex(outputStream.Index())

// 			if err = outputCtx.WritePacket(packet); err != nil {
// 				break
// 			}
// 			packet.Free()
// 		}

// 		for _, frame := range frames {
// 			if frame != nil {
// 				frame.Free()
// 			}
// 		}

// 		if pkt != nil {
// 			pkt.Free()
// 		}
// 	}

// 	// timestamp := int64(0)
// 	// for inputPacket := range inputCtx.GetNewPackets() {
// 	// 	if inputPacket.StreamIndex() != inputStream.Index() {
// 	// 		// for non video packets
// 	// 		continue
// 	// 	}
// 	// 	// todo use of decode vs decode2
// 	// 	frame, code := inputStream.CodecCtx().Decode2(inputPacket)
// 	// 	// if code == int(gmf.AvError())
// 	// 	if code == -11 {
// 	// 		log.Print("resource temporary unavailable", code, gmf.AvError(code))
// 	// 		inputPacket.Free()
// 	// 		continue
// 	// 	}
// 	// 	if code != 0 {
// 	// 		inputPacket.Free()
// 	// 		log.Fatal("inputStream.CodecCtx().Decode2(inputPacket)", code, gmf.AvError(code))
// 	// 		break
// 	// 	}
// 	// 	swsCtx.Scale(frame, dstFrame)
// 	// 	dstFrame.SetPts(timestamp)

// 	// 	if outputPacket, err := dstFrame.Encode(outputStream.CodecCtx()); outputPacket != nil {
// 	// 		outputPacket.SetStreamIndex(outputStream.Index())
// 	// 		if outputPacket.Pts() != gmf.AV_NOPTS_VALUE {
// 	// 			outputPacket.SetPts(gmf.RescaleQ(outputPacket.Pts(), videoEncCtx.TimeBase(), outputStream.TimeBase()))
// 	// 		}

// 	// 		if outputPacket.Dts() != gmf.AV_NOPTS_VALUE {
// 	// 			outputPacket.SetDts(gmf.RescaleQ(outputPacket.Dts(), videoEncCtx.TimeBase(), outputStream.TimeBase()))
// 	// 		}
// 	// 		if err := outputCtx.WritePacket(outputPacket); err != nil {
// 	// 			inputPacket.Free()
// 	// 			outputPacket.Free()
// 	// 			log.Fatal("outputCtx.WritePacket", err)
// 	// 			return
// 	// 		}
// 	// 		outputPacket.Free()
// 	// 		inputPacket.Free()
// 	// 		outputPacket = nil
// 	// 	} else if err != nil {
// 	// 		inputPacket.Free()
// 	// 		outputPacket.Free()
// 	// 		log.Fatal("dstFrame.Encode", err)
// 	// 		return
// 	// 	} else {
// 	// 		log.Printf("encode frame=%d frame=%d is not ready", timestamp, frame.Pts())
// 	// 	}
// 	// 	timestamp++
// 	// }

// }
