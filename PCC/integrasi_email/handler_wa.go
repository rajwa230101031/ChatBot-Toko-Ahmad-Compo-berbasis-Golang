package integrasi_email

import (
	"fmt"
	"strconv"
	"strings"

	"main/models"

	"gorm.io/gorm"
)

// HitungHargaCommission menghitung total harga sesuai rincian pesanan Toko Ahmad Computer
func HitungHargaCommission(jenis, bagian string, jumlah int) float64 {
	jenisLower := strings.ToLower(jenis)
	bagianLower := strings.ToLower(bagian)
	var hargaSatuan float64 = 0

	switch {
	case strings.Contains(jenisLower, "mono") || strings.Contains(jenisLower, "monocrome"):
		switch {
		case strings.Contains(bagianLower, "head") || strings.Contains(bagianLower, "bust") || strings.Contains(bagianLower, "up"):
			hargaSatuan = 15000
		case strings.Contains(bagianLower, "half") || strings.Contains(bagianLower, "setengah"):
			hargaSatuan = 20000
		case strings.Contains(bagianLower, "full") || strings.Contains(bagianLower, "seluruh"):
			hargaSatuan = 25000
		default:
			hargaSatuan = 20000 // default
		}

	case strings.Contains(jenisLower, "simple"):
		switch {
		case strings.Contains(bagianLower, "head") || strings.Contains(bagianLower, "bust") || strings.Contains(bagianLower, "up"):
			hargaSatuan = 15000
		case strings.Contains(bagianLower, "half") || strings.Contains(bagianLower, "setengah"):
			hargaSatuan = 20000
		case strings.Contains(bagianLower, "full") || strings.Contains(bagianLower, "seluruh"):
			hargaSatuan = 25000
		default:
			hargaSatuan = 20000
		}

	case strings.Contains(jenisLower, "animal"):
		switch {
		case strings.Contains(bagianLower, "full") || strings.Contains(bagianLower, "seluruh"):
			hargaSatuan = 25000
		default:
			hargaSatuan = 20000 // headshot/bustup/halfbody
		}

	case strings.Contains(jenisLower, "chibi"):
		switch {
		case strings.Contains(bagianLower, "head") || strings.Contains(bagianLower, "bust") || strings.Contains(bagianLower, "up"):
			hargaSatuan = 25000
		case strings.Contains(bagianLower, "half") || strings.Contains(bagianLower, "setengah"):
			hargaSatuan = 35000
		case strings.Contains(bagianLower, "full") || strings.Contains(bagianLower, "seluruh"):
			hargaSatuan = 45000
		default:
			hargaSatuan = 35000
		}

	case strings.Contains(jenisLower, "stiker") || strings.Contains(jenisLower, "sticker"):
		if jumlah >= 4 {
			// Paket minimal 4 stiker seharga 30.000, selebihnya 8.000 per stiker
			return 30000 + float64(jumlah-4)*8000
		}
		return float64(jumlah) * 8000

	default:
		hargaSatuan = 20000 // default fallback
	}

	return hargaSatuan * float64(jumlah)
}

// ProsesFormPemesananWA memproses formulir pendaftaran pesanan dari chat WA
func ProsesFormPemesananWA(db *gorm.DB, pesanRaw string) (models.Pesanan, error) {
	var nama, email, jenis, bagian, catatan string
	var jumlah int = 1

	lines := strings.Split(pesanRaw, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])

		switch {
		case strings.HasPrefix(key, "nama"):
			nama = value
		case strings.HasPrefix(key, "email"):
			email = value
		case strings.HasPrefix(key, "jenis commission") || strings.HasPrefix(key, "jenis"):
			jenis = value
		case strings.HasPrefix(key, "bagian tubuh") || strings.HasPrefix(key, "bagian"):
			bagian = value
		case strings.HasPrefix(key, "jumlah"):
			if v, err := strconv.Atoi(value); err == nil {
				jumlah = v
			}
		case strings.HasPrefix(key, "catatan"):
			catatan = value
		}
	}

	if nama == "" || email == "" || jenis == "" {
		return models.Pesanan{}, fmt.Errorf("Formulir tidak lengkap! Pastikan Nama, Email, dan Jenis Commission terisi dengan benar.")
	}

	if jumlah <= 0 {
		jumlah = 1
	}

	totalHarga := HitungHargaCommission(jenis, bagian, jumlah)
	hargaSatuan := HitungHargaCommission(jenis, bagian, 1)
	if strings.Contains(strings.ToLower(jenis), "stiker") {
		hargaSatuan = 8000
	}

	// Simpan ke database pesanan
	pesananBaru := models.Pesanan{
		NamaProduk:    fmt.Sprintf("%s - %s (Catatan: %s)", jenis, bagian, catatan),
		Harga:         hargaSatuan,
		Jumlah:        jumlah,
		Diskon:        0,
		TotalHarga:    totalHarga,
		TipePelanggan: email, // Menyimpan email di field TipePelanggan
	}

	if err := db.Create(&pesananBaru).Error; err != nil {
		return models.Pesanan{}, fmt.Errorf("Gagal menyimpan data pesanan: %v", err)
	}

	return pesananBaru, nil
}

// ProsesKonfirmasiWA memproses chat konfirmasi pembayaran: [konfirmasi] <ID>
func ProsesKonfirmasiWA(db *gorm.DB, pesanRaw string) (string, error) {
	bagian := strings.Fields(pesanRaw)
	if len(bagian) < 2 {
		return "", fmt.Errorf("Format konfirmasi salah. Gunakan format: *[konfirmasi] <ID_Pesanan>*. Contoh: *[konfirmasi] 12*")
	}

	pesananID, err := strconv.Atoi(bagian[1])
	if err != nil {
		return "", fmt.Errorf("ID Pesanan harus berupa angka")
	}

	// Cari pesanan di database
	var pesanan models.Pesanan
	if err := db.First(&pesanan, pesananID).Error; err != nil {
		return "", fmt.Errorf("Pesanan dengan ID %d tidak ditemukan di database", pesananID)
	}

	emailPenerima := pesanan.TipePelanggan
	if emailPenerima == "" || !strings.Contains(emailPenerima, "@") {
		return "", fmt.Errorf("Email penerima untuk pesanan ini tidak ditemukan atau tidak valid")
	}

	// Kirim email invoice lunas via SendGrid
	if err := KirimEmailInvoice(emailPenerima, pesanan); err != nil {
		return "", fmt.Errorf("Gagal mengirim email invoice: %v", err)
	}

	responseText := fmt.Sprintf("📧 *INVOICE LUNAS TERKIRIM!*\n\n"+
		"Invoice Lunas untuk pesanan *ORD-%d* telah dikirim ke email *%s*.\n\n"+
		"Pesanan Anda akan langsung kami proses ya kak! Terima kasih! ✨",
		pesanan.ID, emailPenerima)

	return responseText, nil
}
