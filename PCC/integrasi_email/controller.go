package integrasi_email

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main/models"
)

type StrukturInvoice struct {
	PesananID uint   `json:"pesanan_id" binding:"required"`
	Email     string `json:"email" binding:"required"`
}

func KirimInvoice(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var data StrukturInvoice

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "pesan": "Gagal membaca Data", "kesalahan": err.Error()})
		return
	}

	var pesanan models.Pesanan
	if err := db.First(&pesanan, data.PesananID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "pesan": "Pesanan tidak ditemukan"})
		return
	}

	// Panggil fungsi kirim email
	err := KirimEmailInvoice(data.Email, pesanan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "pesan": "Gagal mengirim email", "kesalahan": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "pesan": "Invoice berhasil dikirim ke email " + data.Email, "data": pesanan})
}
