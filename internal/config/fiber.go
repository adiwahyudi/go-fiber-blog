package config

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/spf13/viper"
	"go-blog/version"
)

func NewFiber(config *viper.Viper) *fiber.App {
	engine := html.New("./views", ".html")
	var app = fiber.New(fiber.Config{
		AppName:      fmt.Sprintf("%s v%s", config.GetString("APP_NAME"), version.Version),
		ErrorHandler: NewErrorHandler(),
		Prefork:      config.GetBool("APP_PREFORK"),
		Views:        engine,
	})

	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		return ctx.Status(code).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}
}
