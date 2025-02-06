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

var encryptionkey = []byte("1a2b3c4d5e6f7g8h9i10j11k12m13n1x")

func main() {
	config := &tls.Config{InsecureSkipVerify: true} //for testing

	conn, err := tls.Dial("tcp", "localhost: 8080", config)
	if err != nil {
		fmt.Println("Failed to connect: ", err)
		return
	}
	defer conn.Close()
	encryptFile("sample.txt", "encrypted_output.txt")
	sendFile("encrypted_output.txt", conn)
	fmt.Println("Encrypted file sent successfully!")
}

func encryptFile(inputFile, outputFile string) {
	inFile, _ := os.Open(inputFile)
	defer inFile.Close()

	outFile, _ := os.Create(outputFile)
	defer outFile.Close()

	iv := make([]byte, aes.BlockSize)
	rand.Read(iv)
	outFile.Write(iv)

	//create AES cipher
	block, _ := aes.NewCipher(encryptionkey)
	stream := cipher.NewCFBEncrypter(block, iv)

	io.Copy(&cipher.StreamWriter{S: stream, W: outFile}, inFile)
}

func sendFile(filepath string, conn net.Conn) {
	file, _ := os.Open(filepath)
	defer file.Close()

	buffer := make([]byte, 1024) //created a buffer to store the encrypted data
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading from connection:", err)
			return
		}
		file.Write(buffer[:n])
	}
}
