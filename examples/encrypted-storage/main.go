package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/3d0c/gmf"
	"github.com/awnumar/memguard"
)

func main() {
	// encryptDecrypt()
	// readTranscode()
	fmt.Fprint(os.Stderr, encryptDecryptTranscode())
}

func encryptDecrypt() {
	// dd if=/dev/zero of=input/test.dat  bs=6G  count=1
	inputFilePath := "input/test.dat"
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file %v: %v", inputFilePath, err)
		return
	}
	// fileReader := bufio.NewReader(inputFile)
	if err := processFile(inputFile); err != nil {
		fmt.Fprintf(os.Stderr, "error processing file %v: %v", inputFilePath, err)
		return
	}
}

func processFile(reader io.Reader) error {
	memguard.CatchInterrupt()
	defer memguard.Purge()
	secureKey := memguard.NewEnclaveRandom(16)
	keyBuf, err := secureKey.Open()
	if err != nil {
		return fmt.Errorf("securekey.Open failed %v", err)
	}
	defer keyBuf.Destroy()
	outputEncryptedFilePath := "output/encrypted.mp4"
	encryptedFile, err := os.Create(outputEncryptedFilePath)
	if err != nil {
		return fmt.Errorf("%v file creation failed %v", outputEncryptedFilePath, err)
	}
	// encryptedFile.Close()
	defer encryptedFile.Close()
	if err = encryptAndSave(reader, encryptedFile, keyBuf.Bytes()); err != nil {
		return fmt.Errorf("file encryption failed %v", err)
	}
	outputDecryptedFilePath := "output/decrypted.mp4"
	decryptedFile, err := os.Create(outputDecryptedFilePath)
	if err != nil {
		return fmt.Errorf("%v file creation failed %v", outputDecryptedFilePath, err)
	}
	defer decryptedFile.Close()
	// encryptedFile, err = os.Open(outputEncryptedFilePath)
	// if err != nil {
	// 	return fmt.Errorf("failed to open %v file: %v", outputEncryptedFilePath, err)
	// }
	// defer encryptedFile.Close()
	if _, err := encryptedFile.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("reset seek failed: %v", err)
	}
	if err := decryptAndSave(encryptedFile, decryptedFile, keyBuf.Bytes()); err != nil {
		return fmt.Errorf("file decryption failed %v", err)
	}
	return nil
}

func encryptAndSave(r io.Reader, w io.Writer, key []byte) error {
	fmt.Println(string(key))
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// If the key is unique for each ciphertext, then it's ok to use a zero
	// IV.
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])

	writer := &cipher.StreamWriter{S: stream, W: w}
	// Copy the input to the output stream, decrypting as we go.
	if n, err := io.Copy(writer, r); err != nil {
		return err
	} else {
		fmt.Println(n)
	}

	return nil
}

func decryptAndSave(r io.Reader, w io.Writer, key []byte) error {
	fmt.Println(string(key))
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// If the key is unique for each ciphertext, then it's ok to use a zero
	// IV.
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])

	reader := &cipher.StreamReader{S: stream, R: r}
	// Copy the input to the output stream, decrypting as we go.
	if n, err := io.Copy(w, reader); err != nil {
		return err
	} else {
		fmt.Println(n)
	}
	return nil
}

func readTranscode() {
	// add a mp4 file in input/test.mp4
	inputFilePath := "input/test2.mp4"
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file %v: %v", inputFilePath, err)
		return
	}
	outputFile, err := os.Create("output/test.jpeg")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file %v: %v", inputFilePath, err)
		return
	}

	// fileInfo, err := inputFile.Stat()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	if err = readAndTranscode(inputFile, outputFile); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}
}

// read file -> encrypt and save -> decrypt +  transcode and save
func encryptDecryptTranscode() error {
	inputFilePath := "input/test2.mp4"
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file %v: %v", inputFilePath, err)
		return err
	}

	memguard.CatchInterrupt()
	defer memguard.Purge()
	secureKey := memguard.NewEnclaveRandom(16)
	keyBuf, err := secureKey.Open()
	if err != nil {
		return fmt.Errorf("securekey.Open failed %v", err)
	}
	defer keyBuf.Destroy()
	outputEncryptedFilePath := "output/encrypted.mp4"
	encryptedFile, err := os.Create(outputEncryptedFilePath)
	if err != nil {
		return fmt.Errorf("%v file creation failed %v", outputEncryptedFilePath, err)
	}
	// encryptedFile.Close()
	defer encryptedFile.Close()
	if err = encryptAndSave(inputFile, encryptedFile, keyBuf.Bytes()); err != nil {
		return fmt.Errorf("file encryption failed %v", err)
	}
	decryptedTranscodedFilePath := "output/test.png"
	decryptedTranscodedFile, err := os.Create(decryptedTranscodedFilePath)
	if err != nil {
		return fmt.Errorf("%v file creation failed %v", decryptedTranscodedFilePath, err)
	}
	defer decryptedTranscodedFile.Close()
	// encryptedFile, err = os.Open(outputEncryptedFilePath)
	// if err != nil {
	// 	return fmt.Errorf("failed to open %v file: %v", outputEncryptedFilePath, err)
	// }
	// defer encryptedFile.Close()
	if _, err := encryptedFile.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("reset seek failed: %v", err)
	}
	block, err := aes.NewCipher(keyBuf.Bytes())
	if err != nil {
		return err
	}

	// If the key is unique for each ciphertext, then it's ok to use a zero
	// IV.
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])

	encryptedFileReader := &cipher.StreamReader{S: stream, R: encryptedFile}
	if err = readAndTranscode(encryptedFileReader, decryptedTranscodedFile); err != nil {
		return err
	}
	return nil
}

func readAndTranscode(r io.Reader, w io.Writer) error {
	ctx := gmf.NewCtx()
	defer ctx.Free()

	avioCtx, err := gmf.NewAVIOContext(ctx, &gmf.AVIOHandlers{ReadPacket: func() ([]byte, int) {
		b := make([]byte, gmf.IO_BUFFER_SIZE)
		n, err := r.Read(b)
		if err != nil {
			// how to handle error/tell parent about error?
			fmt.Println("section.Read():", err)
		}
		return b, n
	}})
	if err != nil {
		return err
	}
	ctx.SetPb(avioCtx)
	// 50 mb
	ctx.SetProbeSize(50000000)
	ctx.SetOptions([]*gmf.Option{
		{
			Key: "analyzeduration",
			Val: 2147483647,
		},
	})
	// required for intialization as we are handling the file read
	if err := ctx.OpenInput(""); err != nil {
		return err
	}
	stream, err := ctx.GetBestStream(gmf.AVMEDIA_TYPE_VIDEO)
	if err != nil {
		return err
	}

	fmt.Println("ictx.Duration:", stream.Duration())
	fmt.Printf("bitrate: %d/sec\n", stream.CodecCtx().BitRate())
	jpegCodec, err := gmf.FindEncoder(gmf.AV_CODEC_ID_PNG)
	// jpegCodec, err := gmf.FindEncoder(gmf.AV_CODEC_ID_MJPEG)
	if err != nil {
		return err
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
	if err := jpegCodecCtx.Open(nil); err != nil {
		return err
	}
	defer jpegCodecCtx.Free()

	inputStream, err := ctx.GetStream(stream.Index())
	if err != nil {
		return err
	}
	inputStreamCodec := inputStream.CodecCtx()
	// get source pix format using ffprobe
	swsCtx, err := gmf.NewSwsCtx(inputStreamCodec.Width(), inputStreamCodec.Height(), inputStreamCodec.PixFmt(), jpegCodecCtx.Width(), jpegCodecCtx.Height(), jpegCodecCtx.PixFmt(), gmf.SWS_BICUBIC)
	if err != nil {
		return err
	}
	defer swsCtx.Free()
	var (
		pkt        *gmf.Packet
		frames     []*gmf.Frame
		drain      int = -1
		frameCount int = 0
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

		encode(w, jpegCodecCtx, frames, drain)

		for i, _ := range frames {
			frames[i].Free()
			frameCount++
		}

		if pkt != nil {
			pkt.Free()
			pkt = nil
		}
		break
	}

	return nil
}

func writeFile(w io.Writer, b []byte) {
	fmt.Println(io.Copy(w, bytes.NewReader(b)))
}

func encode(w io.Writer, cc *gmf.CodecCtx, frames []*gmf.Frame, drain int) {
	packets, err := cc.Encode(frames, drain)
	if err != nil {
		log.Fatalf("Error encoding - %s\n", err)
	}
	if len(packets) == 0 {
		return
	}

	for _, p := range packets {
		writeFile(w, p.Data())
		p.Free()
	}

	return
}

// todo update exmples
