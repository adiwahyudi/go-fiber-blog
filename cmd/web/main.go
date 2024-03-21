package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	_ "go-blog/docs"
	"go-blog/internal/config"
)

// @title                       Go-blog Backend
// @description                 Go-blog Backend with Go.
// @termsOfService              http://swagger.io/terms/
// @contact.name                I Wayan Adi Wahyudi
// @contact.email               adi.wahyudi14@gmail.com
// @license.name                Apache 2.0
// @license.url                 http://www.apache.org/licenses/LICENSE-2.0.html
// @schemes                     http https
// @securityDefinitions.apikey  Bearer
// @in                          header
// @name                        Authorization
// @securityDefinitions.basic   BasicAuth
// @description                 "Type 'Bearer {TOKEN}' to correctly set the API Key"
// @BasePath                    /
func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	app := config.NewFiber(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator()

	// Fiber
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     viperConfig.GetString("CORS_ALLOW_ORIGINS"),
		AllowCredentials: true,
	}))

	// swag init -g ./cmd/web/main.go
	app.Get("/swagger/*", swagger.HandlerDefault)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
	})

	webPort := viperConfig.GetInt("APP_PORT")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
