package api

import (
    "github.com/geoo115/Ecommerce/api/handlers"
    "github.com/geoo115/Ecommerce/api/middlewares"
    "github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
    r.POST("/signup", handlers.Signup)
    r.POST("/login", handlers.Login)

    // Product routes
    r.GET("/products", handlers.ListProducts)
    r.POST("/products", handlers.AddProduct)

    // Cart routes
    r.POST("/cart", middlewares.AuthMiddleware(), handlers.AddToCart)
    r.GET("/cart", middlewares.AuthMiddleware(), handlers.ListCart)
    r.DELETE("/cart/:id", middlewares.AuthMiddleware(), handlers.RemoveFromCart)

    // Address routes
    r.POST("/address", middlewares.AuthMiddleware(), handlers.AddAddress)
    r.PUT("/address/:id", middlewares.AuthMiddleware(), handlers.EditAddress)
    r.DELETE("/address/:id", middlewares.AuthMiddleware(), handlers.DeleteAddress)

    r.POST("/checkout", middlewares.AuthMiddleware(), handlers.Checkout)
}