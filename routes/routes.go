// routes/routes.go

package routes

import (
	"sumunar-pos-core/internal/auth"
	"sumunar-pos-core/internal/user"
	"sumunar-pos-core/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, authHandler *auth.Handler, userHandler *user.Handler) {
	// Public routes
	auth := e.Group("/auth")
	auth.POST("/login", authHandler.Login)
	auth.POST("/register", authHandler.Register)
	auth.POST("/google-login", authHandler.GoogleLogin)

	// Protected routes
	api := e.Group("/api")
	api.Use(middleware.JWTAuthMiddleware)
	api.GET("/users", userHandler.ListUsers)
	api.GET("/users/:id", userHandler.GetUser)
	api.POST("/users", userHandler.CreateUser)

	api.POST("/refresh", authHandler.RefreshToken)
	api.POST("/logout", authHandler.Logout)

	storeGroup := e.Group("/stores")
	storeGroup.Use(middleware.JWTAuthMiddleware)
	storeGroup.Use(middleware.RequireRoles("admin", "owner"))
}
