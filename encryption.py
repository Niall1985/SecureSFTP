from Crypto.Cipher import AES
from Crypto.Util.Padding import pad
import os
from dotenv import load_dotenv
import binascii

load_dotenv()

UPLOADS_DIR = "uploads"
ENCRYPTED_DIR = "encrypted_uploads"

os.makedirs(ENCRYPTED_DIR, exist_ok=True)

AES_KEY = binascii.unhexlify(os.getenv('AES_KEY'))

def encrypt_pdf(input_pdf, output_pdf, key):
    """Encrypt a PDF file using AES CBC mode."""
    iv = os.urandom(16)
    cipher = AES.new(key, AES.MODE_CBC, iv)

    with open(input_pdf, 'rb') as f:
        pdf_data = f.read()

    encrypted_data = cipher.encrypt(pad(pdf_data, AES.block_size))

    with open(output_pdf, 'wb') as f:
        f.write(iv + encrypted_data)

    print(f"🔒 Encrypted: {input_pdf} -> {output_pdf}")

def encrypt_all_pdfs():
    """Encrypt PDFs within email-named subdirectories and maintain structure."""
    for email_folder in os.listdir(UPLOADS_DIR):
        email_path = os.path.join(UPLOADS_DIR, email_folder)

        if os.path.isdir(email_path):
            for root, _, files in os.walk(email_path):
                for file_name in files:
                    if file_name.lower().endswith(".pdf"):
                        input_path = os.path.join(root, file_name)

                        relative_path = os.path.relpath(input_path, UPLOADS_DIR)
                        output_path = os.path.join(ENCRYPTED_DIR, relative_path)

                        os.makedirs(os.path.dirname(output_path), exist_ok=True)

                        encrypt_pdf(input_path, output_path, AES_KEY)

    print("✅ All PDFs encrypted successfully!")

if __name__ == "__main__":
    encrypt_all_pdfs()
