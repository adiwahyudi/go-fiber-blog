package usecase

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go-blog/internal/entity"
	"go-blog/internal/helper"
	"go-blog/internal/model"
	"go-blog/internal/model/converter"
	"go-blog/internal/repository"
	"gorm.io/gorm"
	"strings"
)

type PostUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	PostRepository *repository.PostRepository
	TagRepository  *repository.TagRepository
	UserRepository *repository.UserRepository
}

func NewPostUseCase(
	db *gorm.DB, logger *logrus.Logger, validate *validator.Validate, postRepository *repository.PostRepository, tagRepository *repository.TagRepository, userRepository *repository.UserRepository,
) *PostUseCase {
	return &PostUseCase{
		DB:             db,
		Log:            logger,
		Validate:       validate,
		PostRepository: postRepository,
		TagRepository:  tagRepository,
		UserRepository: userRepository,
	}
}

func (c *PostUseCase) Create(ctx context.Context, userId string, request *model.CreatePostRequest) (*model.PostResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	var tags []*entity.Tag
	for _, item := range request.Tags {
		if item.ID == 0 && len(item.Name) < 0 {
			item.Name = "something"
		}
		if err := c.Validate.Struct(item); err != nil {
			c.Log.Warnf("Invalid tags request on body: %+v", err)
			return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid tags request on body")
		}
		if item.ID == 0 {
			newTagName := strings.TrimSpace(item.Name)
			slug := helper.GenerateSlug(newTagName)
			newTag := &entity.Tag{
				Name: newTagName,
				Slug: slug,
			}

			if err := c.TagRepository.Create(tx, newTag); err != nil {
				c.Log.Warnf("Failed create tag to database : %+v", err)
				return nil, fiber.ErrInternalServerError
			}

			tags = append(tags, newTag)
		} else {
			existTag := new(entity.Tag)
			if err := c.TagRepository.FindById(tx, existTag, item.ID); err != nil {
				c.Log.Warnf("Failed to find tag with ID %v : %+v", item.ID, err)
				return nil, fiber.ErrInternalServerError
			}
			tags = append(tags, existTag)
		}
	}

	newTitle := strings.TrimSpace(request.Title)
	slug := helper.GenerateSlug(newTitle)

	post := &entity.Post{
		Title:   newTitle,
		Slug:    slug,
		Content: request.Content,
		UserID:  userId,
		Tags:    tags,
	}

	if err := c.PostRepository.Create(tx, post); err != nil {
		c.Log.Warnf("Failed to create posts: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, userId); err != nil {
		c.Log.Warnf("Failed to get user data: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	post.User = entity.User{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.PostToResponse(post), nil
}

func (c *PostUseCase) List(ctx context.Context, request *model.SearchPostRequest) ([]model.PostResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	posts, total, err := c.PostRepository.Find(tx, request)
	if err != nil {
		c.Log.Warnf("Failed to get posts : %+v", err)
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	response := make([]model.PostResponse, len(posts))
	for i, post := range posts {
		response[i] = *converter.PostToResponse(&post)
	}

	return response, total, nil
}

func (c *PostUseCase) GetBySlug(ctx context.Context, slug string) (*model.PostResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	post := new(entity.Post)
	if err := c.PostRepository.FindBySlug(tx, post, slug); err != nil {
		c.Log.Warnf("Failed to find post by slug '%s': %+v", slug, err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.PostToResponse(post), nil
}
