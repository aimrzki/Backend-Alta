package routes

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"myproject/controllers"
)

var (
	secretKey = []byte("your-secret-key") // Gantilah dengan kunci rahasia yang kuat dalam produksi.
)

func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	e.Use(Logger())

	// Menggunakan routes yang telah dipisahkan
	e.POST("/signup", controllers.Signup(db, secretKey))
	e.POST("/signin", controllers.Signin(db, secretKey))
	e.GET("/user/:id", controllers.GetUserProfile(db, secretKey))
	e.PUT("/user/change-password/:id", controllers.ChangePassword(db, secretKey))
	e.POST("/event/create", controllers.CreateEvent(db, secretKey))
	e.GET("/event", controllers.GetEvents(db, secretKey))
	e.GET("/event/:id", controllers.GetEventByID(db, secretKey))
	e.PUT("/user/:id", controllers.EditUser(db, secretKey))
	e.PUT("/admin/user/:id", controllers.EditUserByAdmin(db, secretKey))
	e.DELETE("/admin/user/:id", controllers.DeleteUserByAdmin(db, secretKey))
	e.GET("/admin/user", controllers.GetAllUsersByAdmin(db, secretKey))
	e.POST("/user/buy", controllers.BuyTicket(db, secretKey))
	e.GET("/user/buy", controllers.GetOrderItemsByUserID(db, secretKey))
	e.GET("/user/ticket", controllers.GetTicketsByUser(db, secretKey))
	e.POST("/admin/promo", controllers.CreatePromo(db, secretKey))
	e.GET("/user/promo", controllers.GetPromos(db, secretKey))
	e.GET("/user/ticket/:invoiceNumber", controllers.GetTicketByInvoiceNumber(db, secretKey))
}
