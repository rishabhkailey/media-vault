package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"os"

	"github.com/awnumar/memguard"
)

func main() {
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
