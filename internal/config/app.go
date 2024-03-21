package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go-blog/internal/delivery/http"
	"go-blog/internal/delivery/http/middleware"
	"go-blog/internal/delivery/http/route"
	"go-blog/internal/repository"
	"go-blog/internal/usecase"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	Config   *viper.Viper
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
}

func Bootstrap(config *BootstrapConfig) {
	// setup middleware
	authMiddleware := middleware.NewMiddleware(config.Config, config.Log)

	// Setup repository
	userRepository := repository.NewUserRepository(config.Log)
	postRepository := repository.NewPostRepository(config.Log)
	tagRepository := repository.NewTagRepository(config.Log)

	// Setup use case
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, authMiddleware)
	postUseCase := usecase.NewPostUseCase(config.DB, config.Log, config.Validate, postRepository, tagRepository, userRepository)
	_ = usecase.NewTagUseCase(config.DB, config.Log, config.Validate, tagRepository)

	// Setup controller
	userController := http.NewUserController(config.Log, userUseCase)
	postController := http.NewPostController(config.Log, postUseCase)

	// Testing route
	config.App.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Pong! 🎉")
	})

	// setup route
	routeConfig := route.RouteConfig{
		App:            config.App,
		UserController: userController,
		PostController: postController,
		AuthMiddleware: authMiddleware,
		Config:         config.Config,
	}

	routeConfig.Setup()
}
