package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"os"
)

var encryptionKey = []byte("1a2b3c4d5e6f7g8h9i10j11k12m13n1x") // 32-byte key for AES-256

func main() {
	config := &tls.Config{InsecureSkipVerify: true} // For testing

	conn, err := tls.Dial("tcp", "localhost:8080", config)
	if err != nil {
		fmt.Println("Failed to connect:", err)
		return
	}
	defer conn.Close()

	inputFile := "sample.txt"
	encryptedFile := "encrypted_output.txt"

	err = encryptFile(inputFile, encryptedFile)
	if err != nil {
		fmt.Println("Encryption failed:", err)
		return
	}

	err = sendFile(encryptedFile, conn)
	if err != nil {
		fmt.Println("Failed to send file:", err)
		return
	}

	fmt.Println("Encrypted file sent successfully!")
}

func encryptFile(inputFile, outputFile string) error {
	inFile, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return err
	}
	_, err = outFile.Write(iv) // Write IV at the beginning
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return err
	}
	stream := cipher.NewCFBEncrypter(block, iv)

	writer := &cipher.StreamWriter{S: stream, W: outFile}
	_, err = io.Copy(writer, inFile)
	return err
}

func sendFile(filePath string, conn net.Conn) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		_, err = conn.Write(buffer[:n])
		if err != nil {
			return err
		}
	}
	return nil
}
