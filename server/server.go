package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"net"
	"os"
)

var encryptionKey = []byte("1a2b3c4d5e6f7g8h9i10j11k12m13n1x") // 32-byte key for AES-256

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	encryptedFile := "received_encrypted.txt"
	decryptedFile := "decrypted_output.txt"

	err := receiveFile(encryptedFile, conn)
	if err != nil {
		fmt.Println("Failed to receive file:", err)
		return
	}

	fmt.Println("Encrypted file received successfully!")

	err = decryptFile(encryptedFile, decryptedFile)
	if err != nil {
		fmt.Println("Decryption failed:", err)
		return
	}

	fmt.Println("File decrypted successfully! Saved as:", decryptedFile)
}

func receiveFile(filePath string, conn net.Conn) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		_, err = file.Write(buffer[:n])
		if err != nil {
			return err
		}
	}
	return nil
}

func decryptFile(inputFile, outputFile string) error {
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

	// Read IV from the beginning of the file
	iv := make([]byte, aes.BlockSize)
	_, err = inFile.Read(iv)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return err
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	reader := &cipher.StreamReader{S: stream, R: inFile}

	// Copy decrypted data to output file
	_, err = io.Copy(outFile, reader)
	return err
}
