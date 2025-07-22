# Aplikasi Chat UDP (Golang)

Aplikasi ini adalah implementasi sederhana dari sistem chat berbasis **UDP socket** menggunakan bahasa pemrograman Go (Golang). Terdiri dari dua bagian utama: **client** dan **server**.

## Fitur

- Client dapat bergabung ke server dengan nama pengguna.
- Client dapat mengirim dan menerima pesan dari user lain.
- Server akan mengatur koneksi, mengirim notifikasi JOIN dan LEAVE, serta melakukan broadcast pesan ke semua client yang terhubung.

## Struktur File

- `main.go` (Client)
- `server.go` (Server)

## Cara Menjalankan

### 1. Jalankan Server

Buka terminal pertama:

```bash
go run server.go
```

### 2. Jalankan Client

Buka terminal kedua atau lebih (untuk banyak pengguna):

```bash
go run main.go
```

Masukkan nama Anda saat diminta, lalu mulai mengirim pesan ke pengguna lain yang juga terhubung ke server.

## Catatan

- Pastikan port `8080` tersedia dan tidak digunakan oleh aplikasi lain.
- Komunikasi hanya berjalan dalam satu jaringan (localhost atau jaringan LAN), karena menggunakan alamat `localhost`.

## Lisensi

MIT License

Copyright (c) 2025
