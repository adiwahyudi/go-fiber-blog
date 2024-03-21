package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go-blog/internal/delivery/http/middleware"
	"go-blog/internal/model"
	"go-blog/internal/usecase"
	"math"
	"strings"
)

type PostController struct {
	Log     *logrus.Logger
	UseCase *usecase.PostUseCase
}

func NewPostController(logger *logrus.Logger, useCase *usecase.PostUseCase) *PostController {
	return &PostController{
		Log:     logger,
		UseCase: useCase,
	}
}

// CreatePost godoc
// @Tags Posts
// @Summary Create a post.
// @Description API create a post.
// @Security Bearer
// @ID create-post
// @Router /api/posts [post]
// @Param _ body model.CreatePostRequest true "Request create post"
// @Accept json
// @Produce json
// @Success 201
func (c *PostController) CreatePost(ctx *fiber.Ctx) error {
	user := middleware.GetUser(ctx)

	request := new(model.CreatePostRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx.UserContext(), user.ID, request)
	if err != nil {
		c.Log.Warnf("Failed to create post : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.PostResponse]{Data: response})

}

// List godoc
// @Tags Posts
// @Summary Get all posts.
// @Description API get all posts.
// @ID get-posts
// @Router /api/posts [get]
// @Param title query string false "Title"
// @Param sort query string false "Sort"
// @Param tags query []string false "Tags"
// @Param page query int false "Page Number" default(1)
// @Param size query int false "Size" default(10)
// @Accept json
// @Produce json
// @Success 200
func (c *PostController) List(ctx *fiber.Ctx) error {
	request := &model.SearchPostRequest{
		Title: ctx.Query("title", ""),
		Tags:  strings.Split(ctx.Query("tags"), ","),
		Sort:  ctx.Query("sort", ""),
		Paginate: model.Pagination{
			Page: ctx.QueryInt("page", 1),
			Size: ctx.QueryInt("size", 10),
		},
	}
	c.Log.Infof("tags: %+v", request.Tags)
	response, total, err := c.UseCase.List(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to load posts: %+v", err)
		return err
	}
	paging := &model.PageMetadata{
		Page:      request.Paginate.Page,
		Size:      request.Paginate.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Paginate.Size))),
	}
	return ctx.JSON(model.WebResponse[[]model.PostResponse]{Data: response, Paging: paging})
}

// ListByUser godoc
// @Tags Posts
// @Summary Get all posts by specific user.
// @Description API get all posts by specific user.
// @ID get-posts-by-user
// @Router /api/posts/{username} [get]
// @Param username path string true "Username"
// @Param title query string false "Title"
// @Param sort query string false "Sort"
// @Param page query int false "Page Number" default(1)
// @Param size query int false "Size" default(10)
// @Accept json
// @Produce json
// @Success 200
func (c *PostController) ListByUser(ctx *fiber.Ctx) error {
	request := &model.SearchPostRequest{
		Username: ctx.Params("username", ""),
		Sort:     ctx.Query("sort", ""),
		Title:    ctx.Query("title", ""),
		//Tag:      ctx.Query("tag", ""),
		Paginate: model.Pagination{
			Page: ctx.QueryInt("page", 1),
			Size: ctx.QueryInt("size", 10),
		},
	}
	response, total, err := c.UseCase.List(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to load posts: %+v", err)
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Paginate.Page,
		Size:      request.Paginate.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Paginate.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.PostResponse]{Data: response, Paging: paging})
}

// FindBySlug godoc
// @Tags Posts
// @Summary Get posts by post slug.
// @Description API get post detail by slug.
// @ID posts-by-slug
// @Router /api/post/{slug} [get]
// @Accept json
// @Param slug path string true "Post Slug"
// @Produce json
// @Success 200
func (c *PostController) FindBySlug(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")

	post, err := c.UseCase.GetBySlug(ctx.UserContext(), slug)
	if err != nil {
		c.Log.Warnf("Failed to load post : %+v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.PostResponse]{Data: post})
}
