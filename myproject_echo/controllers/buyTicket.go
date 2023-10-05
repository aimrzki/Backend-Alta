package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"math/rand"
	"myproject/middleware"
	"myproject/model"
	"net/http"
	"time"
)

func BuyTicket(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Mendapatkan token dari header Authorization
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Authorization token is missing"})
		}

		// Memverifikasi token
		username, err := middleware.VerifyToken(tokenString, secretKey)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Invalid token"})
		}

		// Mengambil informasi user yang terkait dengan token
		var user model.User
		result := db.Where("username = ?", username).First(&user)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to fetch user data"})
		}

		// Menguraikan data pembelian tiket dari JSON yang diterima
		var ticketPurchase struct {
			EventID     uint   `json:"event_id"`
			Quantity    int    `json:"quantity"`
			KodeVoucher string `json:"kode_voucher"`
		}

		if err := c.Bind(&ticketPurchase); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": true, "message": err.Error()})
		}

		// Memeriksa apakah event yang akan dibeli tiketnya ada
		var event model.Event
		eventResult := db.First(&event, ticketPurchase.EventID)
		if eventResult.Error != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": true, "message": "Event not found"})
		}

		// Mencari promo berdasarkan kode voucher yang dimasukkan (jika ada)
		var promo model.Promo
		if ticketPurchase.KodeVoucher != "" {
			promoResult := db.Where("kode_voucher = ?", ticketPurchase.KodeVoucher).First(&promo)
			if promoResult.Error != nil {
				// Jika promo tidak ditemukan, kirimkan pesan kesalahan
				return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": true, "message": "Invalid voucher code"})
			}
		}

		// Menghitung total biaya pembelian tiket
		totalCost := event.Price * ticketPurchase.Quantity

		// Jika promo ditemukan, menghitung potongan biaya tiket
		if promo.ID != 0 {
			potongan := (float64(promo.JumlahPotonganPersen) / 100) * float64(totalCost)
			totalCost -= int(potongan)
		}

		// Membuat entri baru dalam tabel Ticket
		ticket := model.Ticket{
			EventID:       event.ID,
			UserID:        user.ID,
			Quantity:      ticketPurchase.Quantity,
			TotalCost:     totalCost,                  // Total biaya tiket setelah potongan
			InvoiceNumber: generateInvoiceNumber(),    // Simpan nomor invoice dalam tiket
			KodeVoucher:   ticketPurchase.KodeVoucher, // Menyimpan kode voucher dalam entri tiket
		}

		if err := db.Create(&ticket).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to create ticket"})
		}

		// Mengurangi jumlah tiket yang tersedia
		event.DecrementTickets(ticketPurchase.Quantity)
		db.Save(&event)

		// Menyiapkan respons JSON yang mencakup kode voucher yang digunakan
		responseData := map[string]interface{}{
			"error":         false,
			"message":       "Ticket purchased successfully",
			"ticketID":      ticket.ID,                  // Mengirimkan ID tiket yang telah dibeli
			"invoiceNumber": ticket.InvoiceNumber,       // Mengirimkan nomor invoice
			"totalCost":     totalCost,                  // Mengirimkan total biaya tiket
			"kode_voucher":  ticketPurchase.KodeVoucher, // Mengirimkan kode voucher yang digunakan
		}

		return c.JSON(http.StatusOK, responseData)
	}
}

func generateInvoiceNumber() string {
	// Menggunakan waktu sekarang dan nomor acak untuk membuat nomor invoice unik
	timestamp := time.Now().Unix()
	randomNum := rand.Intn(1000) // Ganti dengan rentang nomor yang sesuai
	invoiceNumber := fmt.Sprintf("%d-%d", timestamp, randomNum)
	return invoiceNumber
}
