# from Crypto.Cipher import AES
# from Crypto.Util.Padding import unpad
# import os
# from dotenv import load_dotenv
# import binascii
# import sys

# load_dotenv()

# ENCRYPTED_DIR = "encrypted_uploads"
# DECRYPTED_DIR = "decrypted_uploads"
# AES_KEY = binascii.unhexlify(os.getenv('AES_KEY'))

# def decrypt_pdf(input_pdf, output_pdf, key):
#     with open(input_pdf, 'rb') as f:
#         iv = f.read(16)
#         encrypted_data = f.read()

#     cipher = AES.new(key, AES.MODE_CBC, iv)
#     decrypted_data = unpad(cipher.decrypt(encrypted_data), AES.block_size)

#     with open(output_pdf, 'wb') as f:
#         f.write(decrypted_data)

#     print(f"🔓 Decrypted: {input_pdf} -> {output_pdf}")

# def decrypt_user_pdfs(email):
#     email_folder = os.path.join(ENCRYPTED_DIR, email)
#     if not os.path.exists(email_folder):
#         print(f"❌ No encrypted files found for {email}")
#         return

#     for root, _, files in os.walk(email_folder):
#         for file_name in files:
#             if file_name.lower().endswith(".pdf"):
#                 input_path = os.path.join(root, file_name)
#                 relative_path = os.path.relpath(input_path, ENCRYPTED_DIR)
#                 output_path = os.path.join(DECRYPTED_DIR, relative_path)

#                 os.makedirs(os.path.dirname(output_path), exist_ok=True)
#                 decrypt_pdf(input_path, output_path, AES_KEY)

#     print(f"✅ Decrypted all PDFs for user: {email}")

# if __name__ == "__main__":
#     if len(sys.argv) < 2:
#         print("Usage: python decryption.py <email>")
#     else:
#         decrypt_user_pdfs(sys.argv[1])

from Crypto.Cipher import AES
from Crypto.Util.Padding import unpad
import os
from dotenv import load_dotenv
import binascii
import sys

load_dotenv()

ENCRYPTED_DIR = "encrypted_uploads"
DECRYPTED_DIR = "decrypted_uploads"
AES_KEY = binascii.unhexlify(os.getenv('AES_KEY'))

def decrypt_file(input_file, output_file, key):
    """Decrypt a file (any type) using AES CBC mode."""
    with open(input_file, 'rb') as f:
        iv = f.read(16)
        encrypted_data = f.read()

    cipher = AES.new(key, AES.MODE_CBC, iv)
    decrypted_data = unpad(cipher.decrypt(encrypted_data), AES.block_size)

    with open(output_file, 'wb') as f:
        f.write(decrypted_data)

    print(f"🔓 Decrypted: {input_file} -> {output_file}")

def decrypt_user_files(email):
    """Decrypt all encrypted files under a user's folder."""
    email_folder = os.path.join(ENCRYPTED_DIR, email)
    if not os.path.exists(email_folder):
        print(f"❌ No encrypted files found for {email}")
        return

    for root, _, files in os.walk(email_folder):
        for file_name in files:
            input_path = os.path.join(root, file_name)
            relative_path = os.path.relpath(input_path, ENCRYPTED_DIR)
            output_path = os.path.join(DECRYPTED_DIR, relative_path)

            os.makedirs(os.path.dirname(output_path), exist_ok=True)
            decrypt_file(input_path, output_path, AES_KEY)

    print(f"✅ Decrypted all files for user: {email}")

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python decryption.py <email>")
    else:
        decrypt_user_files(sys.argv[1])
