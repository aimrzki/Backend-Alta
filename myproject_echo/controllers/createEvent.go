package controllers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"myproject/middleware"
	"myproject/model"
	"net/http"
)

func CreateEvent(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
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

		// Memeriksa apakah pengguna memiliki izin untuk membuat event (isAdmin=true)
		if !user.IsAdmin {
			return c.JSON(http.StatusForbidden, map[string]interface{}{"error": true, "message": "Hanya Admin yang dapat menambahkan"})
		}

		// Menguraikan data event dari JSON yang diterima
		var event struct {
			Title            string `json:"title"`
			Location         string `json:"location"`
			Description      string `json:"description"`
			Price            int    `json:"price"`
			AvailableTickets int    `json:"available_ticket"`
		}

		if err := c.Bind(&event); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": true, "message": err.Error()})
		}

		// Membuat event baru dan mengaitkannya dengan pengguna
		newEvent := model.Event{
			Title:            event.Title,
			Location:         event.Location,
			Description:      event.Description,
			Price:            event.Price,
			AvailableTickets: event.AvailableTickets,
			UserID:           user.ID, // Mengaitkan event dengan pengguna yang membuatnya
		}

		if err := db.Create(&newEvent).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to create event"})
		}

		// Mengembalikan respons sukses jika berhasil
		return c.JSON(http.StatusOK, map[string]interface{}{
			"error":     false,
			"message":   "Event created successfully",
			"eventData": newEvent, // Mengirim data event yang baru saja dibuat
		})
	}
}
