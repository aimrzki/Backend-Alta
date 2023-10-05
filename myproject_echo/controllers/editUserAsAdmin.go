package controllers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"myproject/middleware"
	"myproject/model"
	"net/http"
)

func EditUserByAdmin(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Middleware Autentikasi
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Authorization token is missing"})
		}

		// Memverifikasi token dan mendapatkan informasi admin yang diautentikasi
		username, err := middleware.VerifyToken(tokenString, secretKey)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Invalid token"})
		}

		// Mendapatkan data admin dari token
		var adminUser model.User
		result := db.Where("username = ?", username).First(&adminUser)
		if result.Error != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": true, "message": "Anda bukan admin !"})
		}

		// Memeriksa apakah admin yang diautentikasi memiliki status IsAdmin yang true
		if !adminUser.IsAdmin {
			return c.JSON(http.StatusForbidden, map[string]interface{}{"error": true, "message": "Anda bukan admin, tidak bisa mengedit data user lain!"})
		}

		// Mendapatkan ID pengguna dari parameter URL
		userID := c.Param("id")

		// Mencari pengguna berdasarkan ID
		var user model.User
		result = db.First(&user, userID)
		if result.Error != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": true, "message": "User not found"})
		}

		// Mendapatkan data permintaan perubahan dari body
		var req model.User
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": true, "message": err.Error()})
		}

		// Validasi data permintaan (username, email, PhoneNumber)

		// Mengubah data pengguna
		if req.Username != "" {
			user.Username = req.Username
		}
		if req.Email != "" {
			user.Email = req.Email
		}
		if req.PhoneNumber != "" {
			user.PhoneNumber = req.PhoneNumber
		}

		// Simpan perubahan ke basis data
		db.Save(&user)

		return c.JSON(http.StatusOK, map[string]interface{}{"message": "User updated successfully"})
	}
}
