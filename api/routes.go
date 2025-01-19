package api

import (
	"github.com/geoo115/Ecommerce/api/handlers"
	"github.com/geoo115/Ecommerce/api/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/signup", handlers.Signup)
	r.POST("/login", handlers.Login)
	r.POST("/logout", middlewares.AuthMiddleware(), handlers.Logout)

	r.GET("/admin/reports/sales", middlewares.AdminMiddleware(), handlers.SalesReport)
	r.GET("/admin/reports/inventory", middlewares.AdminMiddleware(), handlers.InventoryReport)

	//product categories routes
	r.GET("/categories", handlers.ListCategories)
	r.POST("/categories", middlewares.AdminMiddleware(), handlers.AddCategory)
	r.DELETE("/categories/:id", middlewares.AdminMiddleware(), handlers.DeleteCategory)

	// Product routes
	r.GET("/products", handlers.ListProducts)
	r.GET("/product/:id", handlers.GetProduct)
	r.POST("/product", middlewares.AdminMiddleware(), middlewares.ValidateProduct(), handlers.AddProduct)
	r.PUT("/product/:id", middlewares.AdminMiddleware(), handlers.EditProduct)
	r.DELETE("/product/:id", middlewares.AdminMiddleware(), handlers.DeleteProduct)
	r.GET("/products/search", handlers.SearchProducts)

	// Order routes
	r.POST("/orders", middlewares.AuthMiddleware(), handlers.PlaceOrder)
	r.GET("/orders", middlewares.AuthMiddleware(), handlers.ListOrders)
	r.GET("/orders/:id", middlewares.AuthMiddleware(), handlers.GetOrder)
	r.PUT("/orders/:id/cancel", middlewares.AuthMiddleware(), handlers.CancelOrder)

	// Cart routes
	r.POST("/cart", middlewares.AuthMiddleware(), handlers.AddToCart)
	r.GET("/cart", middlewares.AuthMiddleware(), handlers.ListCart)
	r.DELETE("/cart/:id", middlewares.AuthMiddleware(), handlers.RemoveFromCart)

	// Address routes
	r.POST("/address", middlewares.AuthMiddleware(), handlers.AddAddress)
	r.PUT("/address/:id", middlewares.AuthMiddleware(), handlers.EditAddress)
	r.DELETE("/address/:id", middlewares.AuthMiddleware(), handlers.DeleteAddress)

	//review routes
	r.POST("/reviews", middlewares.AuthMiddleware(), handlers.AddReview)
	r.GET("/reviews/:product_id", handlers.ListReviews)

	//wishlist routes
	r.GET("/wishlist", middlewares.AuthMiddleware(), handlers.ListWishlist)
	r.POST("/wishlist", middlewares.AuthMiddleware(), handlers.AddToWishlist)
	r.DELETE("/wishlist/:id", middlewares.AuthMiddleware(), handlers.RemoveFromWishlist)

	// Payment routes
	r.POST("/payments", middlewares.AuthMiddleware(), handlers.ProcessPayment)
	r.GET("/payments/:order_id", middlewares.AuthMiddleware(), handlers.GetPaymentStatus)
	r.POST("/checkout", middlewares.AuthMiddleware(), handlers.Checkout)
}
