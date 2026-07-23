# Panduan Bot WhatsApp Toko Ahmad Computer

Dokumen ini merangkum fitur chatbot yang saat ini sudah diimplementasikan pada bot WhatsApp, berdasarkan alur di [wa/bot_flow.go](wa/bot_flow.go) dan [wa/wa.go](wa/wa.go).

## 1. Fitur utama chatbot

1. Menu interaktif
   - Bot bisa menampilkan menu utama saat customer mengirim kata seperti "Menu", "Halo", "Hi", atau "P".
   - Menu berisi:
     - 1. Cek Status
     - 2. Konsultasi & Estimasi Biaya
     - 3. Booking Servis
     - 4. Lokasi & Jam
     - 5. Hubungi CS

2. Cek status servis
   - Customer bisa mengirim nomor tiket/nota servis, misalnya TCK-20260723-123 atau ORD-123.
   - Bot akan mencoba menampilkan status perangkat dan estimasi biaya.

3. Konsultasi dan estimasi biaya
   - Customer bisa memilih menu 2 atau mengirim kata seperti "Konsultasi", "Servis", "Repair", atau "Perbaiki".
   - Bot akan menanyakan jenis perangkat (Laptop atau PC).
   - Setelah itu bot akan menanyakan gejala masalah.
   - Bot lalu memberikan perkiraan biaya awal.

4. Booking servis
   - Customer bisa memilih menu 3 atau mengirim kata "Booking".
   - Bot akan meminta tanggal, jam, dan apakah perlu pick-up.
   - Bot lalu membuat tiket booking otomatis.

5. Informasi lokasi dan jam operasional
   - Customer bisa memilih menu 4 atau mengirim kata "Lokasi".
   - Bot membalas dengan alamat toko dan jam operasional.

6. Hubungi customer service
   - Customer bisa memilih menu 5 atau mengirim kata "CS".
   - Bot akan memberikan nomor kontak CS.

7. Sambutan otomatis untuk customer baru
   - Jika customer mengirim sapaan seperti "Halo", "Hallo", "Hi", "Permisi", atau menyebut nama toko, bot akan memberi sambutan dan menawarkan bantuan.

8. Form pemesanan
   - Customer bisa mengetik "Pesan" untuk memulai form pemesanan.
   - Bot mengirim template form yang harus diisi.

9. Konfirmasi pembayaran
   - Setelah order dibuat, customer bisa mengirim pesan dengan awalan "[konfirmasi]" atau kata "Konfirmasi".
   - Bot akan memproses konfirmasi pembayaran.

10. AI assistant
   - Customer bisa mengirim pertanyaan dengan prefix "[ai]".
   - Bot akan mencoba menjawab dengan bantuan AI.

11. Penanganan jam tutup
   - Jika customer datang di luar jam operasional, bot akan memberi pesan tutup dan menyimpan chat untuk ditindaklanjuti saat toko buka.

12. Notifikasi ke grup
   - Saat ada customer yang mengirim sapaan atau konsultasi, bot bisa mengirim notifikasi ke grup tertentu jika grup tersebut tersedia.

## 2. Chat yang bisa dibalas chatbot

Berikut contoh input customer dan respons bot yang sesuai:

- "Halo" / "Hi" / "Permisi" / "P"
  - Bot akan memberi sambutan dan menu bantuan.

- "Menu"
  - Bot akan menampilkan daftar menu utama.

- "1" atau "Cek Status"
  - Bot meminta nomor tiket/nota servis.

- "2" atau "Konsultasi" atau "Servis"
  - Bot memulai alur konsultasi dan menanyakan jenis perangkat.

- "3" atau "Booking"
  - Bot memulai alur booking servis.

- "4" atau "Lokasi"
  - Bot membalas alamat dan jam operasional.

- "5" atau "CS"
  - Bot membalas nomor customer service.

- "TCK-20260723-123" atau "ORD-123"
  - Bot memeriksa status tiket.

- "Laptop" atau "PC"
  - Bot melanjutkan alur konsultasi dan menanyakan gejala masalah.

- "Laptop lemot", "layar mati", "PC tidak bisa nyala"
  - Bot akan menganggap ini sebagai keluhan dan menanyakan perangkat yang dimaksud.

- "Pesan"
  - Bot mengirim template formulir pemesanan.

- "[ai] apa saja jasa yang ada?"
  - Bot menjawab pertanyaan dengan mode AI.

- "[konfirmasi] 123" atau "Konfirmasi"
  - Bot memproses bukti transfer atau konfirmasi pembayaran.

## 3. Contoh kondisi jika ada customer

### A. Customer baru
- Input: "Halo"
- Reaksi bot: menyapa dan menampilkan menu bantuan.

### B. Customer ingin cek status servis
- Input: "1"
- Reaksi bot: "Silakan kirim nomor tiket/nota servis Anda."
- Jika customer kirim: "TCK-20260723-123"
- Reaksi bot: menampilkan status tiket dan estimasi biaya.

### C. Customer ingin konsultasi kerusakan
- Input: "Konsultasi"
- Reaksi bot: "Pilih perangkat: Laptop atau PC?"
- Jika customer balas: "Laptop"
- Reaksi bot: "Baik, perangkat Anda adalah Laptop. Silakan jelaskan gejala/masalah perangkat Anda."
- Jika customer lalu kirim: "Layar mati"
- Reaksi bot: memberikan perkiraan biaya awal dan opsi booking.

### D. Customer ingin booking servis
- Input: "Booking"
- Reaksi bot: "Silakan kirim tanggal & jam yang diinginkan, serta apakah pick-up diperlukan."
- Setelah customer mengirim detail
- Reaksi bot: membuat tiket booking dan memberi nomor tiket.

### E. Customer ingin tahu lokasi toko
- Input: "Lokasi"
- Reaksi bot: memberi alamat toko dan jam operasional.

### F. Customer ingin pesan barang/jasa
- Input: "Pesan"
- Reaksi bot: mengirim template formulir.
- Setelah customer mengirim data formulir
- Reaksi bot: mencatat pesanan dan mengirim QRIS untuk pembayaran.

### G. Customer di luar jam operasional
- Input: masuk saat toko tutup
- Reaksi bot: mengirim pesan tutup dan menyimpan chat untuk diproses nanti.

## 4. Kesimpulan

Bot ini bisa menangani:
- sambutan customer,
- menu layanan,
- cek status servis,
- konsultasi dan estimasi biaya,
- booking servis,
- informasi lokasi dan jam,
- form pemesanan,
- konfirmasi pembayaran,
- serta jawaban singkat via AI.
