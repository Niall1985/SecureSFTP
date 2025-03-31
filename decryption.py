from Crypto.Cipher import AES
from Crypto.Util.Padding import unpad
import os
from dotenv import load_dotenv
import binascii

load_dotenv()

ENCRYPTED_DIR = "encrypted_uploads"
DECRYPTED_DIR = "decrypted_uploads"

os.makedirs(DECRYPTED_DIR, exist_ok=True)

AES_KEY = binascii.unhexlify(os.getenv('AES_KEY'))

def decrypt_pdf(input_pdf, output_pdf, key):
    """Decrypt a PDF file using AES CBC mode."""
    with open(input_pdf, 'rb') as f:
        iv = f.read(16) 
        encrypted_data = f.read()  

    cipher = AES.new(key, AES.MODE_CBC, iv)
    decrypted_data = unpad(cipher.decrypt(encrypted_data), AES.block_size)

    with open(output_pdf, 'wb') as f:
        f.write(decrypted_data)

    print(f"🔓 Decrypted: {input_pdf} -> {output_pdf}")

def decrypt_all_pdfs():
    """Decrypt PDFs within encrypted directories and maintain structure."""
    for email_folder in os.listdir(ENCRYPTED_DIR):
        email_path = os.path.join(ENCRYPTED_DIR, email_folder)
        
        if os.path.isdir(email_path):
            for root, _, files in os.walk(email_path):
                for file_name in files:
                    if file_name.lower().endswith(".pdf"):
                        input_path = os.path.join(root, file_name)

                        relative_path = os.path.relpath(input_path, ENCRYPTED_DIR)
                        output_path = os.path.join(DECRYPTED_DIR, relative_path)

                        os.makedirs(os.path.dirname(output_path), exist_ok=True)

                        decrypt_pdf(input_path, output_path, AES_KEY)

    print("✅ All PDFs decrypted successfully!")

if __name__ == "__main__":
    decrypt_all_pdfs()
