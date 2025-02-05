# 🚀 SecureSFTP-TLS  
A **secure file transfer protocol (SFTP) implementation in Go**, using **TLS encryption** for confidentiality and **authentication mechanisms** to ensure secure communication.

---

## 🔐 **Features**  
✅ **TLS Encryption** – Ensures secure communication between client and server  
✅ **Secure File Transfer** – Transfers files over an encrypted channel  
✅ **User Authentication** – Supports SSH keys & TLS certificates  
✅ **Integrity Verification** – Ensures files are not tampered with  
✅ **Cross-Platform Support** – Works on Linux, Windows, and macOS  

---

## 🛠 **Tech Stack**  
- **Go** – High-performance networking and security  
- **TLS (Transport Layer Security)** – Secure data transmission  
- **x509 Certificates** – Authentication mechanism  
- **AES Encryption** – File encryption before transfer  
- **OpenSSL** – Certificate generation  

---

## 📁 **Project Structure**  
```
SecureSFTP-TLS/
│── server/             # Server-side implementation  
│   ├── server.go       # TLS-enabled SFTP server  
│   ├── config.go       # TLS configuration setup  
│── client/             # Client-side implementation  
│   ├── client.go       # Secure SFTP client  
│── certs/              # TLS certificates  
│── README.md           # Project documentation  
│── .gitignore          # Git ignore file  
```

---

## ⚙️ **Installation & Setup**  

### **1️⃣ Generate TLS Certificates**  
Run the following command to generate self-signed certificates:  
```sh
openssl req -x509 -newkey rsa:4096 -keyout certs/server.key -out certs/server.crt -days 365 -nodes
```

### **2️⃣ Run the Secure SFTP Server**  
```sh
cd server
go run server.go
```

### **3️⃣ Run the Secure SFTP Client**  
```sh
cd client
go run client.go <server-ip> <file-to-send>
```

---

## 🔧 **How It Works**  
1️⃣ The **client** connects to the **server** using a **TLS handshake**  
2️⃣ Once authenticated, the file is **encrypted using AES** before sending  
3️⃣ The **server decrypts** and stores the file securely  

---

## 🛠 **Configuration**  
Modify the `config.go` file to:  
- Change **TLS certificate paths**  
- Update **server IP & port**  
- Set custom **encryption keys**  

---

## 🚀 **Future Enhancements**  
🔹 Support for **multi-user authentication**  
🔹 Implement **SFTP commands (list, delete, rename files, etc.)**  
🔹 Improve **logging & monitoring for security analysis**  

---

## 📜 **License**  
This project is licensed under the **MIT License**.  

---

## 💬 **Contributing**  
We welcome contributions! Feel free to submit pull requests or open issues.  

---
