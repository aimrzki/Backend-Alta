package controllers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"myproject/middleware"
	"myproject/model"
	"net/http"
)

func GetPromos(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Mendapatkan token dari header Authorization
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Authorization token is missing"})
		}

		// Memverifikasi token
		username, err := middleware.VerifyToken(tokenString, secretKey)
		if err != nil {
			// Token tidak valid tetapkan username menjadi string kosong
			username = ""
		}

		// Mengambil daftar promo dari database
		var promos []model.Promo
		if err := db.Find(&promos).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": "Failed to fetch promos"})
		}

		// Mengembalikan daftar promo dalam format yang diinginkan
		return c.JSON(http.StatusOK, map[string]interface{}{
			"error":    false,
			"username": username, // Mengirimkan username pengguna (kosong jika token tidak valid)
			"promos":   promos,
		})
	}
}
