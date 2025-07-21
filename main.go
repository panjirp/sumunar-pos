package main

import (
	"fmt"
	"sumunar-pos-core/config"
	"sumunar-pos-core/db"
	"sumunar-pos-core/pkg/validator"

	"github.com/labstack/echo/v4"
)

func main() {
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

	// Load env and config
	config.LoadConfig()

	// connect to db
	db.Connect()

	// setup echo
	e := echo.New()

	// validator
	e.Validator = validator.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello from "+config.Cfg.AppName)
	})

	addr := fmt.Sprintf(":%s", config.Cfg.Port)
	e.Logger.Fatal(e.Start(addr))
}
