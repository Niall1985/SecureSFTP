package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"os"
)

var encryptionkey = []byte("1a2b3c4d5e6f7g8h9i10j11k12m13n1x")

func main() {
	cert, err := tls.LoadX509KeyPair("../certs/server.crt", "../certs/server.key") //load the tls certificate
	if err != nil {
		fmt.Println("Error loading tls certificates: ", err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}} //configure tls

	listener, err := tls.Listen("tcp", ":8080", config)

	if err != nil {
		fmt.Println("Failed to start server: ", err)
		return
	}
	defer listener.Close() //start the servr
	fmt.Println("Secure SFTP server running on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to connect: ", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) { //handle the content coming from the client
	defer conn.Close()
	fmt.Println("Connection estabilished: ", conn.RemoteAddr())

	file, err := os.Create("received_content.txt")
	if err != nil {
		fmt.Println("Error creating file: ", err)
		return
	}

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
	fmt.Println("Encrypted file received.")

	decryptFile("received_encrypted.txt", "decrypted_output.txt")
	fmt.Println("Decryption complete. File saved as decrypted_output.txt")
}

func decryptFile(inputFile, outputFile string) {
	inFile, _ := os.Open(inputFile)
	defer inFile.Close()
	outFile, _ := os.Create(outputFile)
	defer outFile.Close()

	iv := make([]byte, aes.BlockSize)
	inFile.Read(iv)

	block, _ := aes.NewCipher(encryptionkey) // Create AES cipher
	stream := cipher.NewCFBDecrypter(block, iv)

	io.Copy(outFile, &cipher.StreamReader{S: stream, R: inFile}) // Decrypt and write to file
}
