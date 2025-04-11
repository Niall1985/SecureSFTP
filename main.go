// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/exec"
// 	"path/filepath"

// 	"github.com/rs/cors"
// )

// type UploadResponse struct {
// 	Email string   `json:"email"`
// 	Files []string `json:"files"`
// }

// type ReceiveRequest struct {
// 	Email string `json:"email"`
// }

// type ReceiveResponse struct {
// 	Email string   `json:"email"`
// 	Files []string `json:"files"`
// 	URLs  []string `json:"urls"`
// }

// const (
// 	encryptionScript = "encryption.py"
// 	decryptionScript = "decryption.py"
// )

// // Handle uploading files and triggering encryption
// func uploadHandler(w http.ResponseWriter, r *http.Request) {
// 	err := r.ParseMultipartForm(100 << 20)
// 	if err != nil {
// 		http.Error(w, "Error parsing form", http.StatusBadRequest)
// 		return
// 	}

// 	email := r.FormValue("email")
// 	if email == "" {
// 		http.Error(w, "Email is required", http.StatusBadRequest)
// 		return
// 	}

// 	files := r.MultipartForm.File["files"]
// 	if len(files) == 0 {
// 		http.Error(w, "No files uploaded", http.StatusBadRequest)
// 		return
// 	}

// 	uploadDir := filepath.Join("uploads", email)
// 	os.MkdirAll(uploadDir, os.ModePerm)

// 	var uploadedFiles []string

// 	for _, fileHeader := range files {
// 		file, err := fileHeader.Open()
// 		if err != nil {
// 			http.Error(w, "Error opening file", http.StatusInternalServerError)
// 			return
// 		}
// 		defer file.Close()

// 		filePath := filepath.Join(uploadDir, fileHeader.Filename)
// 		outFile, err := os.Create(filePath)
// 		if err != nil {
// 			http.Error(w, "Error creating file", http.StatusInternalServerError)
// 			return
// 		}
// 		defer outFile.Close()

// 		_, err = io.Copy(outFile, file)
// 		if err != nil {
// 			http.Error(w, "Error saving file", http.StatusInternalServerError)
// 			return
// 		}
// 		uploadedFiles = append(uploadedFiles, fileHeader.Filename)
// 	}

// 	metadata := UploadResponse{Email: email, Files: uploadedFiles}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(metadata)

// 	// Encrypt uploaded files
// 	go runPythonScript(encryptionScript, uploadDir)
// }

// // Run encryption/decryption Python scripts
// func runPythonScript(script string, dir string) {
// 	cmd := exec.Command("python", script, dir)
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	err := cmd.Run()
// 	if err != nil {
// 		fmt.Println("❌ Error running Python script:", err)
// 	} else {
// 		fmt.Println("✅ Python script executed successfully.")
// 	}
// }

// // Handle receiving files (trigger decryption + return URLs)
// func receiveHandler(w http.ResponseWriter, r *http.Request) {
// 	var req ReceiveRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
// 		http.Error(w, "Invalid email provided", http.StatusBadRequest)
// 		return
// 	}

// 	fmt.Println("🔄 Decrypting files for:", req.Email)
// 	cmd := exec.Command("python", decryptionScript, req.Email)
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	if err := cmd.Run(); err != nil {
// 		http.Error(w, "Error decrypting files", http.StatusInternalServerError)
// 		return
// 	}

// 	decryptedDir := filepath.Join("decrypted_uploads", req.Email)
// 	files, err := os.ReadDir(decryptedDir)
// 	if err != nil {
// 		http.Error(w, "Error reading decrypted files", http.StatusInternalServerError)
// 		return
// 	}

// 	var fileNames, urls []string
// 	fmt.Println("📁 Decrypted files found:")
// 	for _, file := range files {
// 		if !file.IsDir() {
// 			fileNames = append(fileNames, file.Name())
// 			urls = append(urls, fmt.Sprintf("https://securesftp.onrender.com/download/%s/%s", req.Email, file.Name()))
// 			fmt.Println("  -", file.Name())
// 		}
// 	}

// 	response := ReceiveResponse{Email: req.Email, Files: fileNames, URLs: urls}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

// // Serve decrypted files for download
// func downloadHandler(w http.ResponseWriter, r *http.Request) {
// 	email := r.PathValue("email")
// 	fileName := r.PathValue("file")

// 	filePath := filepath.Join("decrypted_uploads", email, fileName)

// 	if _, err := os.Stat(filePath); os.IsNotExist(err) {
// 		http.Error(w, "File not found", http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
// 	w.Header().Set("Content-Type", "application/octet-stream")
// 	w.Header().Set("Content-Transfer-Encoding", "binary")
// 	http.ServeFile(w, r, filePath)
// }

// // Main function
// func main() {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/upload", uploadHandler)
// 	mux.HandleFunc("/receive", receiveHandler)
// 	mux.HandleFunc("/download/{email}/{file}", downloadHandler)

// 	handler := cors.AllowAll().Handler(mux)
// 	port := "8080"
// 	fmt.Println("🚀 Server running on port", port)
// 	log.Fatal(http.ListenAndServe(":"+port, handler))
// }

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rs/cors"
)

type UploadResponse struct {
	Email string   `json:"email"`
	Files []string `json:"files"`
}

type ReceiveRequest struct {
	Email string `json:"email"`
}

type ReceiveResponse struct {
	Email string   `json:"email"`
	Files []string `json:"files"`
	URLs  []string `json:"urls"`
}

const (
	encryptionScript = "encryption.py"
	decryptionScript = "decryption.py"
)

// Handle uploading files and triggering encryption
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	uploadDir := filepath.Join("uploads", email)
	os.MkdirAll(uploadDir, os.ModePerm)

	var uploadedFiles []string

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Error opening file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		filePath := filepath.Join(uploadDir, fileHeader.Filename)
		outFile, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Error creating file", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, file)
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}
		uploadedFiles = append(uploadedFiles, fileHeader.Filename)
	}

	// Encrypt uploaded files
	go runPythonScript(encryptionScript, uploadDir)

	metadata := UploadResponse{Email: email, Files: uploadedFiles}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metadata)
}

// Run encryption/decryption Python scripts
func runPythonScript(script string, arg string) {
	cmd := exec.Command("python", script, arg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("❌ Error running Python script:", err)
	} else {
		fmt.Println("✅ Python script executed successfully.")
	}
}

// Return encrypted files & URLs to frontend (no decryption here!)
func receiveHandler(w http.ResponseWriter, r *http.Request) {
	var req ReceiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		http.Error(w, "Invalid email provided", http.StatusBadRequest)
		return
	}

	encryptedDir := filepath.Join("encrypted_uploads", req.Email)
	files, err := os.ReadDir(encryptedDir)
	if err != nil {
		http.Error(w, "Error reading encrypted files", http.StatusInternalServerError)
		return
	}

	var fileNames, urls []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
			urls = append(urls, fmt.Sprintf("https://securesftp.onrender.com/download/%s/%s", req.Email, file.Name()))
		}
	}

	response := ReceiveResponse{
		Email: req.Email,
		Files: fileNames,
		URLs:  urls,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Decrypt files only when user downloads
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	email := r.PathValue("email")
	fileName := r.PathValue("file")

	// Trigger decryption only for the requested file
	fmt.Println("🔐 Triggering decryption for:", fileName)
	cmd := exec.Command("python", decryptionScript, email)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		http.Error(w, "Error decrypting file", http.StatusInternalServerError)
		return
	}

	// Serve decrypted file
	filePath := filepath.Join("decrypted_uploads", email, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/pdf")
	http.ServeFile(w, r, filePath)
}

// Main function
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", uploadHandler)
	mux.HandleFunc("/receive", receiveHandler)
	mux.HandleFunc("/download/{email}/{file}", downloadHandler)

	handler := cors.AllowAll().Handler(mux)
	port := "8080"
	fmt.Println("🚀 Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
