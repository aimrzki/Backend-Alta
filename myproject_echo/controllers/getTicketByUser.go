package controllers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"myproject/middleware"
	"myproject/model"
	"net/http"
)

func GetTicketsByUser(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
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

		// Mendapatkan ID pengguna dari token
		var user model.User
		result := db.Where("username = ?", username).First(&user)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to fetch user data"})
		}

		// Mengambil tiket yang telah dibeli oleh pengguna berdasarkan UserID
		var tickets []model.Ticket
		result = db.Where("user_id = ?", user.ID).Find(&tickets)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to fetch user's tickets"})
		}

		// Membuat respons dengan data tiket yang telah dibeli
		var ticketDetails []map[string]interface{}
		for _, ticket := range tickets {
			// Mengambil detail event berdasarkan EventID yang ada pada tiket
			var event model.Event
			eventResult := db.First(&event, ticket.EventID)
			if eventResult.Error != nil {
				// Handle jika event tidak ditemukan
				continue
			}

			// Menambahkan informasi kode voucher yang digunakan
			var kodeVoucher string
			if ticket.KodeVoucher != "" {
				kodeVoucher = ticket.KodeVoucher
			}

			// Menambahkan nomor invoice ke detail tiket
			ticketDetail := map[string]interface{}{
				"user_id":        ticket.UserID,
				"event_id":       ticket.EventID,
				"event_title":    event.Title,
				"quantity":       ticket.Quantity,
				"total_cost":     ticket.TotalCost,
				"invoice_number": ticket.InvoiceNumber,
				"kode_voucher":   kodeVoucher, // Menambahkan informasi kode voucher ke respons
			}

			// Menambahkan objek tiket ke daftar ticketDetails
			ticketDetails = append(ticketDetails, ticketDetail)
		}

		// Mengembalikan respons dengan detail tiket yang telah dibeli
		return c.JSON(http.StatusOK, map[string]interface{}{
			"error":       false,
			"message":     "User's tickets retrieved successfully",
			"ticket_data": ticketDetails,
		})
	}
}
