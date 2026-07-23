package fungsi

import (
	"fmt"
	"os"
	"time"

	"main/models"

	"gorm.io/gorm"
)

// IsStoreOpen mengecek apakah toko Ahmad Computer sedang buka berdasarkan waktu lokal saat ini.
// Mengembalikan status buka (bool), pesan informasi operasional, dan error (jika ada).
func IsStoreOpen(db *gorm.DB) (bool, string, error) {
	// 1. Ambil timezone dari environment (.env)
	tz := os.Getenv("TIMEZONE")
	if tz == "" {
		tz = "Asia/Jakarta" // Default timezone
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return false, "", fmt.Errorf("gagal memuat timezone %s: %v", tz, err)
	}

	// 2. Ambil waktu saat ini disesuaikan dengan timezone toko
	now := time.Now().In(loc)
	currentDay := now.Weekday() // 0 = Sunday, 1 = Monday, dst.
	currentTimeStr := now.Format("15:04")

	// 3. Ambil jadwal operasional hari ini dari database
	var schedule models.StoreSchedule
	errDb := db.Where("day_of_week = ?", int(currentDay)).First(&schedule).Error
	if errDb != nil {
		return false, "", fmt.Errorf("gagal mengambil jadwal toko dari database: %v", errDb)
	}

	// Jika secara manual ditandai tutup pada hari ini
	if !schedule.IsOpen {
		return false, "Toko kami sedang libur hari ini.", nil
	}

	// 4. Bandingkan waktu saat ini dengan jam buka dan tutup
	openTime, err := time.ParseInLocation("15:04", schedule.OpenTime, loc)
	if err != nil {
		return false, "", fmt.Errorf("gagal mem-parsing jam buka %s: %v", schedule.OpenTime, err)
	}

	closeTime, err := time.ParseInLocation("15:04", schedule.CloseTime, loc)
	if err != nil {
		return false, "", fmt.Errorf("gagal mem-parsing jam tutup %s: %v", schedule.CloseTime, err)
	}

	// Buat objek time.Time hari ini untuk perbandingan jam saja
	todayOpen := time.Date(now.Year(), now.Month(), now.Day(), openTime.Hour(), openTime.Minute(), 0, 0, loc)
	todayClose := time.Date(now.Year(), now.Month(), now.Day(), closeTime.Hour(), closeTime.Minute(), 0, 0, loc)

	if now.Before(todayOpen) || now.After(todayClose) {
		infoPesan := fmt.Sprintf("Jam operasional: %s - %s. Saat ini jam %s %s.", schedule.OpenTime, schedule.CloseTime, currentTimeStr, tz)
		return false, infoPesan, nil
	}

	return true, "", nil
}
