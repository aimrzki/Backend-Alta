package controllers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"myproject/middleware"
	"myproject/model"
	"net/http"
)

func GetEvents(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
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

		// Mengambil daftar event dari database
		var events []model.Event
		if err := db.Find(&events).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to fetch events"})
		}

		// Mengembalikan daftar event dalam format yang diinginkan
		return c.JSON(http.StatusOK, map[string]interface{}{"error": false, "events": events})
	}
}
