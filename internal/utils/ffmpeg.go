package utils

import (
	"fmt"
	"io"
	"log"

	"github.com/3d0c/gmf"
)

func customReader(r io.Reader) ([]byte, int) {
	b := make([]byte, gmf.IO_BUFFER_SIZE)
	n, err := r.Read(b)
	if err != nil {
		// how to handle error/tell parent about error?
		fmt.Println("section.Read():", err)
	}
	return b, n
}

func GenerateThumbnail(r io.Reader) (b []byte, err error) {
	ctx := gmf.NewCtx()
	defer ctx.Free()

	avioCtx, err := gmf.NewAVIOContext(ctx, &gmf.AVIOHandlers{ReadPacket: func() ([]byte, int) { return customReader(r) }})
	if err != nil {
		return
	}
	ctx.SetPb(avioCtx)
	// 50 mb
	// ctx.SetProbeSize(50000000)
	// ctx.SetOptions([]*gmf.Option{
	// 	{
	// 		Key: "analyzeduration",
	// 		Val: 1,
	// 	},
	// })
	// required for intialization as we are handling the file read
	if err = ctx.OpenInput(""); err != nil {
		return
	}
	stream, err := ctx.GetBestStream(gmf.AVMEDIA_TYPE_VIDEO)
	if err != nil {
		return
	}

	fmt.Println("ictx.Duration:", stream.Duration())
	fmt.Printf("bitrate: %d/sec\n", stream.CodecCtx().BitRate())
	jpegCodec, err := gmf.FindEncoder(gmf.AV_CODEC_ID_PNG)
	// jpegCodec, err := gmf.FindEncoder(gmf.AV_CODEC_ID_MJPEG)
	if err != nil {
		return
	}
	jpegCodecCtx := gmf.NewCodecCtx(jpegCodec)
	defer jpegCodecCtx.Free()
	// ffmpeg -h encoder=mjpeg -v quiet # get supported pixel formats for a codec
	// ffmpeg -h encoder=png -v quiet # get supported pixel formats for a codec
	jpegCodecCtx.SetPixFmt(gmf.AV_PIX_FMT_RGBA).
		SetWidth(stream.CodecCtx().Width()).
		SetHeight(stream.CodecCtx().Height()).
		SetTimeBase(gmf.AVR{Num: 1, Den: 1})
	if jpegCodecCtx.Codec().IsExperimental() {
		jpegCodecCtx.SetStrictCompliance(gmf.FF_COMPLIANCE_EXPERIMENTAL)
	}
	if err = jpegCodecCtx.Open(nil); err != nil {
		return
	}
	defer jpegCodecCtx.Free()

	inputStream, err := ctx.GetStream(stream.Index())
	if err != nil {
		return
	}
	inputStreamCodec := inputStream.CodecCtx()
	// get source pix format using ffprobe
	swsCtx, err := gmf.NewSwsCtx(inputStreamCodec.Width(), inputStreamCodec.Height(), inputStreamCodec.PixFmt(), jpegCodecCtx.Width(), jpegCodecCtx.Height(), jpegCodecCtx.PixFmt(), gmf.SWS_BICUBIC)
	if err != nil {
		return
	}
	defer swsCtx.Free()
	var (
		pkt    *gmf.Packet
		frames []*gmf.Frame
		drain  int = -1
		// frameCount int = 0
	)
	for {
		if drain >= 0 {
			break
		}

		pkt, err = ctx.GetNextPacket()
		if err != nil && err != io.EOF {
			if pkt != nil {
				pkt.Free()
			}
			log.Printf("error getting next packet - %s", err)
			break
		} else if err != nil && pkt == nil {
			drain = 0
		}
		defer pkt.Free()

		if pkt != nil && pkt.StreamIndex() != stream.Index() {
			continue
		}

		frames, err = inputStreamCodec.Decode(pkt)
		if err != nil {
			log.Printf("Fatal error during decoding - %s\n", err)
			break
		}

		// Decode() method doesn't treat EAGAIN and EOF as errors
		// it returns empty frames slice instead. Countinue until
		// input EOF or frames received.
		if len(frames) == 0 && drain < 0 {
			continue
		}

		if frames, err = gmf.DefaultRescaler(swsCtx, frames); err != nil {
			panic(err)
		}

		b = encode(jpegCodecCtx, frames, drain)
		return b, nil
		// for i, _ := range frames {
		// 	frames[i].Free()
		// 	frameCount++
		// }

		// if pkt != nil {
		// 	pkt.Free()
		// 	pkt = nil
		// }
		// break
	}

	return
}

func encode(cc *gmf.CodecCtx, frames []*gmf.Frame, drain int) (b []byte) {
	packets, err := cc.Encode(frames, drain)
	if err != nil {
		log.Fatalf("Error encoding - %s\n", err)
	}
	if len(packets) == 0 {
		return
	}

	for _, p := range packets {
		defer p.Free()
		return p.Data()
	}

	return
}
