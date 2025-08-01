package main

import (
	"fmt"
	"sumunar-pos-core/config"
	"sumunar-pos-core/db"
	docs "sumunar-pos-core/docs"
	"sumunar-pos-core/internal/auth"
	"sumunar-pos-core/internal/customer"
	"sumunar-pos-core/internal/order"
	"sumunar-pos-core/internal/product"
	"sumunar-pos-core/internal/productservice"
	"sumunar-pos-core/internal/servicetype"
	"sumunar-pos-core/internal/store"
	"sumunar-pos-core/internal/user"
	"sumunar-pos-core/internal/userstore"
	"sumunar-pos-core/pkg/validator"
	"sumunar-pos-core/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Swagger docs metadata
	// @title Sumunar POS API
	// @version 1.0
	// @description This is a POS backend for Sumunar.
	// @contact.name Panji Rachmatullah
	// @contact.email panji@example.com
	// @host localhost:8080
	// @BasePath /api/v1
	// @securityDefinitions.apikey BearerAuth
	// @in header
	// @name Authorization

	// Load configuration
	config.LoadConfig()

	// Connect to DB
	db.Connect()
	dbConn := db.GetTxBeginner() // dbConn is db.DBTX

	// Setup Echo
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:8000"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	e.Validator = validator.New()

	// Setup Swagger docs config
	docs.SwaggerInfo.Title = "Sumunar POS API"
	docs.SwaggerInfo.Description = "API documentation for Sumunar POS"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"

	// Simple health check
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello from "+config.Cfg.AppName)
	})

	// ==== Init Repositories ====
	userRepo := user.NewUserRepository(dbConn)
	userStoreRepo := userstore.NewUserStoreRepository(dbConn)
	refreshTokenRepo := auth.NewRefreshTokenRepo(dbConn)
	storeRepo := store.NewStoreRepository(dbConn)
	serviceTypeRepo := servicetype.NewServiceRepository(dbConn)
	productRepo := product.NewProductRepository(dbConn)
	productServiceRepo := productservice.NewProductServiceRepository(dbConn)
	orderRepo := order.NewOrderRepository(dbConn)
	customerRepo := customer.NewCustomerRepository(dbConn)

	// ==== Init Services ====
	authService := auth.NewService(userRepo, refreshTokenRepo)
	userService := user.NewService(userRepo)
	userStoreService := userstore.NewService(userStoreRepo)
	storeService := store.NewService(storeRepo, userStoreService, dbConn)
	serviceTypeService := servicetype.NewService(serviceTypeRepo, dbConn)
	productService := product.NewService(productRepo, dbConn)
	productServiceService := productservice.NewService(productServiceRepo, dbConn)
	orderService := order.NewService(orderRepo, productServiceRepo, customerRepo, dbConn)

	// ==== Init Handlers ====
	authHandler := auth.NewHandler(authService)
	userHandler := user.NewHandler(userService)
	storeHandler := store.NewHandler(storeService)
	serviceTypeHandler := servicetype.NewHandler(serviceTypeService)
	productHandler := product.NewHandler(productService)
	productServiceHandler := productservice.NewHandler(productServiceService)
	orderHandler := order.NewHandler(orderService)

	// ==== Register Routes ====
	routes.RegisterRoutes(
		e,
		authHandler,
		userHandler,
		storeHandler,
		serviceTypeHandler,
		productHandler,
		productServiceHandler,
		orderHandler,
	)

	// Start server
	addr := fmt.Sprintf(":%s", config.Cfg.Port)
	e.Logger.Fatal(e.Start(addr))
}
