package controllers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"myproject/middleware"
	"myproject/model"
	"net/http"
)

func Signin(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user model.User
		if err := c.Bind(&user); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		// Mengecek apakah username ada dalam database
		var existingUser model.User
		result := db.Where("username = ?", user.Username).First(&existingUser)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
			} else {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check username"})
			}
		}

		// Membandingkan password yang dimasukkan dengan password yang di-hash
		err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
		}

		// Generate JWT token
		tokenString, err := middleware.GenerateToken(existingUser.Username, secretKey)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
		}

		// Menyertakan ID pengguna dalam respons
		return c.JSON(http.StatusOK, map[string]interface{}{"message": "Login successful", "token": tokenString, "id": existingUser.ID})
	}
}
