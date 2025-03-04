package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"os"
)

var encryptionKey = []byte("1a2b3c4d5e6f7g8h9i10j11k12m13n14") // 32-byte AES-256 key

func main() {
	serverAddress := "192.168.36.199:8080"
	filePath := "C:\\Users\\Niall Dcunha\\SecureSFTP\\hi.txt"
	encryptedFilePath := "encrypted_test.txt"

	// Ensure the input file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("❌ Error: File not found:", filePath)
		return
	}

	err := encryptFile(filePath, encryptedFilePath)
	if err != nil {
		fmt.Println("Encryption failed:", err)
		return
	}
	fmt.Println("✅ Encryption successful! File saved as", encryptedFilePath)

	// Connect to server
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("❌ Connection failed:", err)
		return
	}
	defer conn.Close()

	err = sendFile(encryptedFilePath, conn)
	if err != nil {
		fmt.Println("❌ File transfer failed:", err)
	} else {
		fmt.Println("✅ File sent successfully!")
	}
}

// **Encrypts a file and saves the encrypted version**
func encryptFile(inputFile, outputFile string) error {
	fmt.Println("🔒 Encrypting file...")

	inFile, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("error opening input file: %v", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outFile.Close()

	// Create AES cipher
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return fmt.Errorf("error creating cipher block: %v", err)
	}

	// Generate IV (16 bytes for AES)
	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(rand.Reader, iv) // FIX: Use `rand.Reader` instead of os.Open("/dev/urandom")
	if err != nil {
		return fmt.Errorf("error generating IV: %v", err)
	}
	_, err = outFile.Write(iv) // Save IV at the beginning of the file
	if err != nil {
		return fmt.Errorf("error writing IV: %v", err)
	}

	// Create encryption stream
	stream := cipher.NewCFBEncrypter(block, iv)
	writer := &cipher.StreamWriter{S: stream, W: outFile}

	// Encrypt and write to file
	_, err = io.Copy(writer, inFile)
	if err != nil {
		return fmt.Errorf("error encrypting file: %v", err)
	}

	fmt.Println("✅ File encrypted successfully!")
	return nil
}

// **Sends the encrypted file to the server**
func sendFile(filePath string, conn net.Conn) error {
	fmt.Println("📤 Sending file to server...")

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	buffer := make([]byte, 4096) // Increased buffer size
	totalBytes := 0

	for {
		n, err := file.Read(buffer)
		if n > 0 {
			_, writeErr := conn.Write(buffer[:n])
			if writeErr != nil {
				return fmt.Errorf("error writing to connection: %v", writeErr)
			}
			totalBytes += n
			fmt.Println("📤 Sent", totalBytes, "bytes")
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("error reading file: %v", err)
		}
	}

	fmt.Println("✅ File transfer complete! Total bytes:", totalBytes)
	return nil
}
