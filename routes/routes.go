// routes/routes.go

package routes

import (
	"sumunar-pos-core/internal/auth"
	"sumunar-pos-core/internal/order"
	"sumunar-pos-core/internal/product"
	"sumunar-pos-core/internal/productservice"
	"sumunar-pos-core/internal/servicetype"
	"sumunar-pos-core/internal/store"
	"sumunar-pos-core/internal/user"
	"sumunar-pos-core/middleware"

	"github.com/labstack/echo/v4"

	_ "sumunar-pos-core/docs" // penting!

	echoSwagger "github.com/swaggo/echo-swagger"
)

func RegisterRoutes(e *echo.Echo, authHandler *auth.Handler, userHandler *user.Handler,
	storeHandler *store.Handler, serviceHandler *servicetype.Handler, productHandler *product.Handler,
	productServiceHandler *productservice.Handler, orderHandler *order.Handler) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// /api group (root for protected and nested routes)
	api := e.Group("/api")

	// Public Auth routes
	auth := api.Group("/auth")
	auth.POST("/login", authHandler.Login)
	auth.POST("/register", authHandler.Register)
	auth.POST("/google-login", authHandler.GoogleLogin)

	authprotected := auth.Group("/protected")
	authprotected.Use(middleware.JWTAuthMiddleware)
	authprotected.POST("/refresh", authHandler.RefreshToken)
	authprotected.POST("/logout", authHandler.Logout)

	// Protected routes
	api.Use(middleware.JWTAuthMiddleware)

	// Users
	users := api.Group("/users")
	users.GET("", userHandler.ListUsers)
	users.GET("/:id", userHandler.GetUser)
	users.POST("", userHandler.CreateUser)

	// Stores (only for admin/owner)
	stores := api.Group("/stores", middleware.RequireRoles("admin", "owner"))
	stores.POST("", storeHandler.Create)
	stores.GET("", storeHandler.FindAll)
	stores.GET("/:id", storeHandler.FindByID)
	stores.PUT("/:id", storeHandler.Update)
	stores.DELETE("/:id", storeHandler.Delete)
	stores.POST("/:id/logo", storeHandler.UploadLogo, middleware.ValidateImageFile)

	// Service Type (only for admin/owner)
	services := api.Group("/servicetypes", middleware.RequireRoles("admin", "owner"))
	services.POST("", serviceHandler.Create)
	services.GET("", serviceHandler.FindAll)
	services.GET("/:id", serviceHandler.FindByID)
	services.PUT("/:id", serviceHandler.Update)
	services.DELETE("/:id", serviceHandler.Delete)

	// Products (only for admin/owner)
	products := api.Group("/products", middleware.RequireRoles("admin", "owner"))
	products.POST("", productHandler.Create)
	products.GET("", productHandler.FindAll)
	products.GET("/:id", productHandler.FindByID)
	products.PUT("/:id", productHandler.Update)
	products.DELETE("/:id", productHandler.Delete)

	// Product Service (only for admin/owner)
	productservice := api.Group("/productservice", middleware.RequireRoles("admin", "owner"))
	productservice.POST("", productServiceHandler.Create)
	productservice.GET("", productServiceHandler.FindAll)
	productservice.GET("/:id", productServiceHandler.FindByID)
	productservice.PUT("/:id", productServiceHandler.Update)
	productservice.DELETE("/:id", productServiceHandler.Delete)

	// Order (only for admin/owner)
	order := api.Group("/order", middleware.RequireRoles("admin", "owner"))
	order.POST("", orderHandler.Create)
	order.GET("", orderHandler.FindAll)
	// order.GET("/:id", orderHandler.FindByID)
	order.PUT("/:id", orderHandler.Update)
	order.DELETE("/:id", orderHandler.Delete)
}
