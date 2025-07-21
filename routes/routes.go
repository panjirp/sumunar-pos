// routes/routes.go

package routes

import (
	"sumunar-pos-core/internal/auth"
	"sumunar-pos-core/internal/store"
	"sumunar-pos-core/internal/user"
	"sumunar-pos-core/middleware"

	"github.com/labstack/echo/v4"

	_ "sumunar-pos-core/docs" // penting!

	echoSwagger "github.com/swaggo/echo-swagger"
)

func RegisterRoutes(e *echo.Echo, authHandler *auth.Handler, userHandler *user.Handler, storeHandler *store.Handler) {
	// Swagger (can be outside /api if you want it public)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// /api group (root for protected and nested routes)
	api := e.Group("/api")

	// Public Auth routes
	auth := api.Group("/auth")
	auth.POST("/login", authHandler.Login)
	auth.POST("/register", authHandler.Register)
	auth.POST("/google-login", authHandler.GoogleLogin)

	// Token routes (could be public or protected depending on design)
	api.POST("/refresh", authHandler.RefreshToken)
	api.POST("/logout", authHandler.Logout)

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
}
