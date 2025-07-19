package api

import (
	"github.com/geoo115/Ecommerce/api/handlers"
	"github.com/geoo115/Ecommerce/api/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Health check endpoints (no rate limiting)
	r.GET("/health", handlers.HealthCheck)
	r.GET("/health/detailed", handlers.DetailedHealthCheck)
	r.GET("/ready", handlers.ReadinessCheck)
	r.GET("/live", handlers.LivenessCheck)
	r.GET("/metrics", handlers.Metrics)

	// Authentication routes with stricter rate limiting
	authGroup := r.Group("/")
	authGroup.Use(middlewares.AuthRateLimit())
	{
		authGroup.POST("/signup", handlers.Signup)
		authGroup.POST("/login", handlers.Login)
	}

	r.POST("/logout", middlewares.AuthMiddleware(), handlers.Logout)

	// Admin routes with admin rate limiting
	adminGroup := r.Group("/admin")
	adminGroup.Use(middlewares.AdminMiddleware())
	adminGroup.Use(middlewares.AdminRateLimit())
	{
		adminGroup.GET("/reports/sales", handlers.SalesReport)
		adminGroup.GET("/reports/inventory", handlers.InventoryReport)
	}

	// Categories routes
	r.GET("/categories", handlers.ListCategories)
	r.POST("/categories", middlewares.AdminMiddleware(), handlers.AddCategory)
	r.DELETE("/categories/:id", middlewares.AdminMiddleware(), handlers.DeleteCategory)

	// Product routes
	r.GET("/products", handlers.ListProducts)
	r.GET("/product/:id", handlers.GetProduct)
	r.GET("/products/search", handlers.SearchProducts)

	// Admin product routes
	productAdminGroup := r.Group("/product")
	productAdminGroup.Use(middlewares.AdminMiddleware())
	productAdminGroup.Use(middlewares.AdminRateLimit())
	{
		productAdminGroup.POST("", middlewares.ValidateProduct(), handlers.AddProduct)
		productAdminGroup.PUT("/:id", handlers.EditProduct)
		productAdminGroup.DELETE("/:id", handlers.DeleteProduct)
	}

	// Order routes
	orderGroup := r.Group("/orders")
	orderGroup.Use(middlewares.AuthMiddleware())
	{
		orderGroup.POST("", handlers.PlaceOrder)
		orderGroup.GET("", handlers.ListOrders)
		orderGroup.GET("/:id", handlers.GetOrder)
		orderGroup.PUT("/:id/cancel", handlers.CancelOrder)
	}

	// Cart routes
	cartGroup := r.Group("/cart")
	cartGroup.Use(middlewares.AuthMiddleware())
	{
		cartGroup.POST("", handlers.AddToCart)
		cartGroup.GET("", handlers.ListCart)
		cartGroup.DELETE("/:id", handlers.RemoveFromCart)
	}

	// Address routes
	addressGroup := r.Group("/address")
	addressGroup.Use(middlewares.AuthMiddleware())
	{
		addressGroup.POST("", handlers.AddAddress)
		addressGroup.PUT("/:id", handlers.EditAddress)
		addressGroup.DELETE("/:id", handlers.DeleteAddress)
	}

	// Review routes
	reviewGroup := r.Group("/reviews")
	{
		reviewGroup.GET("/:product_id", handlers.ListReviews)
		reviewGroup.POST("", middlewares.AuthMiddleware(), handlers.AddReview)
	}

	// Wishlist routes
	wishlistGroup := r.Group("/wishlist")
	wishlistGroup.Use(middlewares.AuthMiddleware())
	{
		wishlistGroup.GET("", handlers.ListWishlist)
		wishlistGroup.POST("", handlers.AddToWishlist)
		wishlistGroup.DELETE("/:id", handlers.RemoveFromWishlist)
	}

	// Payment routes
	paymentGroup := r.Group("/payments")
	paymentGroup.Use(middlewares.AuthMiddleware())
	{
		paymentGroup.POST("", handlers.ProcessPayment)
		paymentGroup.GET("/:order_id", handlers.GetPaymentStatus)
	}

	// Checkout route
	r.POST("/checkout", middlewares.AuthMiddleware(), handlers.Checkout)
}
