package integrasi_email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"main/models"
)

func KirimEmailInvoice(emailPenerima string, pesanan models.Pesanan) error {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	senderEmail := os.Getenv("SENDGRID_SENDER_EMAIL")

	if apiKey == "" || senderEmail == "" {
		return fmt.Errorf("SENDGRID_API_KEY atau SENDGRID_SENDER_EMAIL belum dikonfigurasi di file .env")
	}

	subject := fmt.Sprintf("Invoice Lunas Pembayaran Pesanan ORD-%d", pesanan.ID)

	// Susun HTML Invoice yang profesional & bersih (Tanpa QRIS)
	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #dddddd; border-radius: 5px; }
		.header { text-align: center; border-bottom: 2px solid #10b981; padding-bottom: 10px; margin-bottom: 20px; }
		.header h2 { color: #10b981; margin: 0; }
		.invoice-info { margin-bottom: 20px; }
		.invoice-table { width: 100%%; border-collapse: collapse; margin-bottom: 20px; }
		.invoice-table th, .invoice-table td { border: 1px solid #dddddd; padding: 10px; text-align: left; }
		.invoice-table th { background-color: #f3f4f6; }
		.total-section { font-size: 16px; font-weight: bold; text-align: right; margin-bottom: 20px; }
		.payment-status { background-color: #ecfdf5; border-left: 4px solid #10b981; padding: 15px; border-radius: 4px; margin-bottom: 20px; text-align: center; }
		.footer { text-align: center; font-size: 12px; color: #777777; border-top: 1px solid #dddddd; padding-top: 10px; margin-top: 20px; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h2>INVOICE PEMBAYARAN (LUNAS)</h2>
			<p>Toko Ahmad Computer</p>
		</div>
		
		<div class="invoice-info">
			<p><strong>ID Pesanan:</strong> ORD-%d</p>
			<p><strong>Tanggal:</strong> %s</p>
			<p><strong>Tipe Pelanggan:</strong> %s</p>
			<p><strong>Email Penerima:</strong> %s</p>
		</div>

		<table class="invoice-table">
			<thead>
				<tr>
					<th>Nama Jasa / Produk</th>
					<th>Harga Satuan</th>
					<th>Jumlah</th>
					<th>Diskon</th>
				</tr>
			</thead>
			<tbody>
				<tr>
					<td>%s</td>
					<td>Rp %.0f</td>
					<td>%d</td>
					<td>Rp %.0f (10%%)</td>
				</tr>
			</tbody>
		</table>

		<div class="total-section">
			Total Pembayaran: Rp %.0f
		</div>

		<div class="payment-status">
			<strong style="color: #10b981; font-size: 16px;">✓ STATUS PEMBAYARAN: LUNAS</strong><br>
			<p style="margin: 5px 0 0 0; font-size: 13px; color: #047857;">Pembayaran Anda telah berhasil diverifikasi. Pesanan Anda sedang diproses oleh tim Toko Ahmad Computer.</p>
		</div>

		<div class="footer">
			<p>Terima kasih atas kepercayaan Anda berbelanja di Toko Ahmad Computer!</p>
			<p>Ini adalah email otomatis yang dikirim oleh server Gin-Go Cloud Integration.</p>
		</div>
	</div>
</body>
</html>
`, pesanan.ID, pesanan.CreatedAt.Format("02 January 2006 15:04 MST"), pesanan.TipePelanggan, emailPenerima, pesanan.NamaProduk, pesanan.Harga, pesanan.Jumlah, pesanan.Diskon, pesanan.TotalHarga)

	// Definisikan struktur payload SendGrid v3 API
	type EmailUser struct {
		Email string `json:"email"`
		Name  string `json:"name,omitempty"`
	}

	type Personalization struct {
		To      []EmailUser `json:"to"`
		Subject string      `json:"subject"`
	}

	type Content struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}

	type Payload struct {
		Personalizations []Personalization `json:"personalizations"`
		From             EmailUser         `json:"from"`
		Content          []Content         `json:"content"`
	}

	dataPayload := Payload{
		Personalizations: []Personalization{
			{
				To: []EmailUser{
					{Email: emailPenerima},
				},
				Subject: subject,
			},
		},
		From: EmailUser{
			Email: senderEmail,
			Name:  "Toko Ahmad Computer",
		},
		Content: []Content{
			{
				Type:  "text/html",
				Value: htmlContent,
			},
		},
	}

	payloadBytes, err := json.Marshal(dataPayload)
	if err != nil {
		return fmt.Errorf("gagal marshalling payload: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("gagal membuat request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("koneksi API SendGrid gagal: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API SendGrid mengembalikan status error: %d", resp.StatusCode)
	}

	return nil
}
