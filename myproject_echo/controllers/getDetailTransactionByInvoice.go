package controllers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"myproject/middleware"
	"myproject/model"
	"net/http"
)

func GetTicketByInvoiceNumber(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
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

		// Mendapatkan informasi pengguna dari basis data
		var user model.User
		result := db.Where("username = ?", username).First(&user)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to fetch user data"})
		}

		// Memeriksa apakah pengguna memiliki status admin
		if !user.IsAdmin {
			return c.JSON(http.StatusForbidden, map[string]interface{}{"error": true, "message": "Access forbidden for non-admin users"})
		}

		// Mendapatkan nomor invoice dari parameter URL
		invoiceNumber := c.Param("invoiceNumber")

		// Mencari tiket berdasarkan nomor invoice
		var ticket model.Ticket
		result = db.Where("invoice_number = ?", invoiceNumber).First(&ticket)
		if result.Error != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": true, "message": "Ticket not found"})
		}

		// Mengambil detail event berdasarkan EventID yang ada pada tiket
		var event model.Event
		eventResult := db.First(&event, ticket.EventID)
		if eventResult.Error != nil {
			// Handle jika event tidak ditemukan
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to fetch event data"})
		}

		// Membuat respons dengan detail pembelian tiket
		ticketDetail := map[string]interface{}{
			"ticketID":       ticket.ID,
			"user_id":        ticket.UserID,
			"event_id":       ticket.EventID,
			"event_title":    event.Title,
			"quantity":       ticket.Quantity,
			"total_cost":     ticket.TotalCost,
			"invoice_number": ticket.InvoiceNumber,
			"kode_voucher":   ticket.KodeVoucher, // Menambahkan kode voucher ke respons jika ada
		}

		// Mengembalikan respons dengan detail pembelian tiket
		return c.JSON(http.StatusOK, map[string]interface{}{
			"error":       false,
			"message":     "Ticket details retrieved successfully",
			"ticket_data": ticketDetail,
		})
	}
}
