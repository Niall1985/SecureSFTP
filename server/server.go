package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"net"
	"os"
)

var encryptionKey = []byte("1a2b3c4d5e6f7g8h9i10j11k12m13n14") // 32-byte AES-256 key

func main() {
	// Ensure the directory exists
	os.MkdirAll("server_files", os.ModePerm)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		fmt.Println("Client connected")

		encryptedFile := "server_files/encrypted_received.txt"
		decryptedFile := "server_files/decrypted_output.txt"

		err = receiveFile(encryptedFile, conn)
		if err != nil {
			fmt.Println("Failed to receive file:", err)
			conn.Close()
			continue
		}
		fmt.Println("✅ Encrypted file saved as", encryptedFile)

		err = decryptFile(encryptedFile, decryptedFile)
		if err != nil {
			fmt.Println("❌ Decryption failed:", err)
		} else {
			fmt.Println("✅ Decryption successful! File saved as", decryptedFile)
		}

		conn.Close()
	}
}

// **Receives an encrypted file from the client and saves it**
func receiveFile(filePath string, conn net.Conn) error {
	fmt.Println("Receiving file...")

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	buffer := make([]byte, 4096) // Increased buffer size
	totalBytes := 0

	for {
		n, err := conn.Read(buffer)
		if n > 0 {
			_, writeErr := file.Write(buffer[:n])
			if writeErr != nil {
				fmt.Println("Error writing to file:", writeErr)
				return writeErr
			}
			totalBytes += n
			fmt.Println("📥 Received", totalBytes, "bytes")
		}

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error reading from connection:", err)
			return err
		}
	}

	fmt.Println("✅ File received successfully! Total bytes:", totalBytes)
	return nil
}

// **Decrypts the received file, prints its content, and saves it**
func decryptFile(inputFile, outputFile string) error {
	fmt.Println("🔍 DEBUG: Entering decryptFile function")

	inFile, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("❌ Error opening encrypted file:", err)
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("❌ Error creating decrypted file:", err)
		return err
	}
	defer outFile.Close()

	// Read IV
	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(inFile, iv)
	if err != nil {
		fmt.Println("❌ Error reading IV:", err)
		return err
	}
	fmt.Println("✅ IV read successfully")

	// Create AES cipher
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		fmt.Println("❌ Error creating cipher block:", err)
		return err
	}
	fmt.Println("✅ AES cipher block created")

	// Create decryption stream
	stream := cipher.NewCFBDecrypter(block, iv)
	reader := &cipher.StreamReader{S: stream, R: inFile}

	// Read decrypted data
	decryptedData, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("❌ Error decrypting file:", err)
		return err
	}

	// Print decrypted content to console
	fmt.Println("✅ Decryption complete. Decrypted content:")
	fmt.Println(string(decryptedData))

	// Write decrypted content to file
	_, err = outFile.Write(decryptedData)
	if err != nil {
		fmt.Println("❌ Error saving decrypted file:", err)
		return err
	}

	fmt.Println("✅ Decrypted file saved as", outputFile)
	return nil
}
