package routes

import (
	"rentalin/internal/handler"
	"rentalin/internal/middleware"

	"github.com/labstack/echo/v4"
)

type RouteConfig struct {
	AuthHandler      *handler.AuthHandler
	UserHandler      *handler.UserHandler
	ProductHandler   *handler.ProductHandler
	RentalHandler    *handler.RentalHandler
	PaymentHandler   *handler.PaymentHandler
	DashboardHandler *handler.DashboardHandler
	JWTSecret        string
}

func SetupRoutes(e *echo.Echo, cfg RouteConfig) {
	jwtMiddleware := middleware.JWTMiddleware(cfg.JWTSecret)

	// ================= AUTH =================
	auth := e.Group("/auth")
	auth.POST("/register", cfg.AuthHandler.Register)
	auth.POST("/login", cfg.AuthHandler.Login)

	// ================= PUBLIC =================
	public := e.Group("")
	public.GET("/products", cfg.ProductHandler.GetAllProducts)
	public.GET("/products/:id", cfg.ProductHandler.GetProduct)
	public.POST("/payments/webhook/xendit", cfg.PaymentHandler.XenditWebhook)

	// ================= USER =================
	user := e.Group("")
	user.Use(jwtMiddleware)
	user.GET("/users/me", cfg.UserHandler.GetProfile)
	user.GET("/users/me/rentals", cfg.RentalHandler.GetRentalsByCustomerID)
	user.GET("/users/me/payments", cfg.PaymentHandler.GetPaymentsByCustomerID)
	user.PATCH("/users/me", cfg.UserHandler.UpdateProfile)

	// ================= RENTAL (USER & ADMIN) =================
	rental := e.Group("/rental")
	rental.Use(jwtMiddleware)

	rental.POST("", cfg.RentalHandler.CreateRental)
	rental.GET("/:id", cfg.RentalHandler.GetRentalByID)
	rental.POST("/:id/cancel", cfg.RentalHandler.CancelRental)

	// ================= PAYMENT =================
	payment := e.Group("/payments")
	payment.Use(jwtMiddleware)
	payment.POST("", cfg.PaymentHandler.CreatePayment)

	// // ================= ADMIN =================
	admin := e.Group("/admin")
	admin.Use(jwtMiddleware)
	admin.Use(middleware.AdminOnly)

	// // admin manage users
	admin.GET("/users", cfg.UserHandler.GetAllUsers)
	admin.DELETE("/users/:id", cfg.UserHandler.DeleteUser)

	// admin manage products
	admin.POST("/product", cfg.ProductHandler.CreateProduct)
	admin.PATCH("/product/:id", cfg.ProductHandler.UpdateProduct)
	admin.DELETE("/product/:id", cfg.ProductHandler.DeleteProduct)

	// admin manage rentals
	admin.GET("/rentals", cfg.RentalHandler.GetAllRentals)
	admin.POST("/rental/:id/complete", cfg.RentalHandler.CompleteRental)

	//admin dashboard
	admin.GET("/dashboard", cfg.DashboardHandler.GetDashboard)
}
