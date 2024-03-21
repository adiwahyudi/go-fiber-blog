package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go-blog/internal/delivery/http"
	"go-blog/internal/delivery/http/middleware"
)

type RouteConfig struct {
	App            fiber.Router
	UserController *http.UserController
	PostController *http.PostController
	AuthMiddleware *middleware.Middleware
	Config         *viper.Viper
}

func (c *RouteConfig) Setup() {
	c.SetupViewsRoutes()
	c.SetupGuestRoutes()
	c.SetupProtectedRoutes()
}
func (c *RouteConfig) SetupViewsRoutes() {
	c.App.Get("/home", func(ctx *fiber.Ctx) error {
		return ctx.Render("home", fiber.Map{
			"Name": "Adi Wahyudi",
		})
	})
}
func (c *RouteConfig) SetupGuestRoutes() {
	// Auth
	auth := c.App.Group("/auth")
	auth.Post("/login", c.UserController.Login)
	auth.Post("/register", c.UserController.RegisterUser)

	// Post
	c.App.Get("/posts", c.PostController.List)
	c.App.Get("/posts/:username", c.PostController.ListByUser)
	c.App.Get("/post/:slug", c.PostController.FindBySlug)
}

func (c *RouteConfig) SetupProtectedRoutes() {

	// Users
	users := c.App.Group("/users")
	users.Patch("", c.UserController.Update)
	users.Patch("/:userId", c.UserController.Delete)

	// Post
	posts := c.App.Group("/posts")
	posts.Post("", c.PostController.CreatePost)
}
