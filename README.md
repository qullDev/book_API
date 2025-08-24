# Book API

Layanan RESTful untuk manajemen kategori dan buku dengan autentikasi JSON Web Token (JWT). Access Token digunakan untuk mengakses endpoint, sedangkan Refresh Token disimpan di Redis untuk rotasi/invalidasi sesi.

Teknologi:
- Go + Gin
- PostgreSQL (ORM: GORM)
- Redis (penyimpanan Refresh Token)
- JWT HS256

Daftar fitur:
- Autentikasi:
  - Login (Access Token + Refresh Token)
  - Refresh token (rotasi RT, RT lama dicabut)
  - Logout (revoke satu RT tertentu atau seluruh RT milik user)
- Kategori:
  - List, Create, Detail, Update, Delete
  - List buku per kategori
- Buku:
  - List, Create, Detail, Update, Delete
  - Validasi release_year (1980–2024)
  - Konversi otomatis kolom thickness berdasarkan total_page


## Prasyarat

- Go 1.24+
- PostgreSQL 13+ (atau kompatibel)
- Redis 6+


## Konfigurasi Environment

Buat file `.env` di root proyek (atau set environment variable di sistem):
