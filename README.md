# Subdomain Breaker (SubBreaker)

**Subdomain Breaker** adalah tools untuk mencari subdomain dari domain tertentu dengan menggabungkan beberapa penyedia data publik. Alat ini juga memungkinkan pemindaian port untuk mengecek apakah subdomain tertentu memiliki port yang terbuka.

## Fitur Utama:
- **Pencarian Subdomain** dari berbagai penyedia data publik:
  - crt.sh
  - findsubdomains.com
  - (Untuk sementara hanya 2 provider)
- **Pemindaian Port** untuk mengecek port terbuka pada subdomain yang ditemukan.

## Cara Instalasi dan Penggunaan

### Sebelum menggunakan pastikan golang terinstall terlebih dahulu di pc / komputer / os kalian
### untuk step penginstalan dan setup golang cek sendiri  di google

### 1. Clone Repository
Untuk menggunakan, clone repository ini ke pc kalian:

git clone https://github.com/F1kro/Subbreaker.git
cd Subbreaker

### 2. Jalankan Subbreaker
-> Untuk Windows :
- ./subbreaker.exe
- ./subbreaker.exe -d domain.com (-d domain)
- ./subbreaker.exe -d domain.com -o output.txt (-o output .txt) 
- ./subbreaker.exe -h (-h help)
  
-> Untuk Linux   : 
- ./subbreaker (banner)
- ./subbreaker -d domain.com (-d domain)
- ./subbreaker -d domain.com -o output.txt (-o output .txt) 
- ./subbreaker -h (-h help) 

