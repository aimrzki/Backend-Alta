package controllers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"myproject/middleware"
	"myproject/model"
	"net/http"
	"strconv"
)

func GetEventByID(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Mendapatkan token dari header Authorization
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Authorization token is missing"})
		}

		// Memverifikasi token
		_, err := middleware.VerifyToken(tokenString, secretKey)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Invalid token"})
		}

		// Mengambil event ID dari path parameter
		eventIDParam := c.Param("id")

		// Mengonversi eventIDParam ke tipe data uint
		eventID, err := strconv.ParseUint(eventIDParam, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": true, "message": "Invalid event ID"})
		}

		// Mencari event berdasarkan event ID
		var event model.Event
		if err := db.First(&event, eventID).Error; err != nil {
			// Event tidak ditemukan
			if err == gorm.ErrRecordNotFound {
				return c.JSON(http.StatusNotFound, map[string]interface{}{"error": true, "message": "Event not found"})
			}
			// Terjadi kesalahan lain saat mengambil data dari database
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to fetch event"})
		}

		// Mengembalikan event yang ditemukan dalam format yang diinginkan
		return c.JSON(http.StatusOK, map[string]interface{}{"error": false, "event": event})
	}
}
