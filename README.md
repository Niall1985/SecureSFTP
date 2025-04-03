# 🚀 SecureSFTP
A **secure file transfer system** built with **Go**, utilizing **backend on Render** and **frontend on Vercel**. The system ensures **encrypted file storage and secure retrieval**.

---

## 🔐 **Features**
✅ **Secure File Transfer** – Accepts PDF files from frontend and stores them securely
✅ **End-to-End Encryption** – Uses AES encryption for confidentiality
✅ **Decryption on Demand** – Users can request files, and they are decrypted only when needed
✅ **User-Based Storage** – Files are associated with user email for organized access
✅ **Dockerized Deployment** – Easy containerized deployment with **Docker**
✅ **Cloud-Hosted** – Backend on **Render**, Frontend on **Vercel**

---

## 🛠 **Tech Stack**
- **Go** – Backend for handling file uploads and downloads
- **Python** – Encryption & decryption scripts
- **Docker** – Containerized deployment
- **Render** – Backend hosting
- **Vercel** – Frontend hosting

---

## 📁 **Project Structure**
```
SecureSFTP/
│── backend/             # Backend API (Go + encryption scripts)
│   ├── main.go         # Handles file uploads, encryption, and retrieval
│   ├── encryption.py   # AES encryption script
│   ├── decryption.py   # AES decryption script
│── frontend/            # Frontend (Vercel-hosted UI)
│── Dockerfile           # Container setup
│── README.md            # Project documentation
│── .gitignore           # Git ignore file
```

---

## ⚙️ **Installation & Setup**  
### **1️⃣ Clone the Repository**  
```sh
git clone https://github.com/your-repo/SecureSFTP.git
cd SecureSFTP
```

### **2️⃣ Build and Run with Docker**  
```sh
docker build -t securesftp .
docker run -p 8080:8080 securesftp
```

---

## 🔧 **How It Works**  
1️⃣ **Frontend (Vercel)**: Accepts **PDF uploads** along with the **target user's email**
2️⃣ **Backend (Render)**:
   - Receives files and **encrypts them using AES**
   - Stores them under `uploads/{email}/`
3️⃣ **File Retrieval**:
   - Users can request their files using their email
   - The backend **decrypts the files on-demand**
   - Generates a **temporary download link**
4️⃣ **File Cleanup**:
   - Once a file is downloaded, it is **deleted from decrypted storage**
   - If no decrypted files remain for a user, their directory is also deleted

---

## 🚀 **Future Enhancements**  
🔹 Implement user authentication for file access  
🔹 Add support for file previews  
🔹 Enable multi-file batch processing  

---

## 📜 **License**  
This project is licensed under the **MIT License**.  

---

## 💬 **Contributing**  
We welcome contributions! Feel free to submit pull requests or open issues.

