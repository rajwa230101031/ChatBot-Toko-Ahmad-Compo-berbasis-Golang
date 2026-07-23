package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main/models"
)

type StrukturPesan struct {
	Kode    string `json:"kode" binding:"required"`
	Balasan string `json:"balasan" binding:"required"`
}

func PesanTampil(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var modelPesan []models.Pesan

	hasil := db.Find(&modelPesan)
	kesalahan := hasil.Error

	if kesalahan == nil {
		c.JSON(http.StatusOK, gin.H{"status": true, "pesan": "Berhasil tampil data pesan", "data": modelPesan})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": false, "pesan": "Gagal tampil data", "kesalahan": kesalahan.Error()})
	}
}

func PesanTambah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataPesan StrukturPesan

	if err := c.ShouldBindJSON(&dataPesan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "pesan": "Gagal membaca Data", "kesalahan": err.Error()})
		return
	}

	modelPesan := models.Pesan{
		Kode:    dataPesan.Kode,
		Balasan: dataPesan.Balasan,
	}

	hasil := db.Create(&modelPesan)
	kesalahan := hasil.Error

	if kesalahan == nil {
		c.JSON(http.StatusOK, gin.H{"status": true, "pesan": "Berhasil tambah pesan", "data": modelPesan})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": false, "pesan": "Gagal tambah pesan", "kesalahan": kesalahan.Error()})
	}
}

func PesanUbah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataPesan StrukturPesan

	if err := c.ShouldBindJSON(&dataPesan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "pesan": "Gagal membaca Data", "kesalahan": err.Error()})
		return
	}

	var modelPesan models.Pesan
	if err := db.First(&modelPesan, "kode = ?", dataPesan.Kode).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "pesan": "Data tidak ditemukan"})
		return
	}

	modelPesan.Balasan = dataPesan.Balasan

	hasil := db.Save(&modelPesan)
	kesalahan := hasil.Error

	if kesalahan == nil {
		c.JSON(http.StatusOK, gin.H{"status": true, "pesan": "Berhasil ubah pesan", "data": modelPesan})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": false, "pesan": "Gagal ubah pesan", "kesalahan": kesalahan.Error()})
	}
}

func PesanHapus(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataPesan struct {
		Kode string `json:"kode" binding:"required"`
	}

	if err := c.ShouldBindJSON(&dataPesan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "pesan": "Gagal membaca Data", "kesalahan": err.Error()})
		return
	}

	var modelPesan models.Pesan

	hasil := db.Delete(&modelPesan, "kode = ?", dataPesan.Kode)
	kesalahan := hasil.Error

	if kesalahan == nil {
		c.JSON(http.StatusOK, gin.H{"status": true, "pesan": "Berhasil hapus pesan", "data": dataPesan})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": false, "pesan": "Gagal hapus pesan", "kesalahan": kesalahan.Error()})
	}
}
