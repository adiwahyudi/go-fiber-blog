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

type TagUseCase struct {
	DB            *gorm.DB
	Log           *logrus.Logger
	Validate      *validator.Validate
	TagRepository *repository.TagRepository
}

func NewTagUseCase(
	db *gorm.DB, logger *logrus.Logger, validate *validator.Validate, tagRepository *repository.TagRepository,
) *TagUseCase {
	return &TagUseCase{
		DB:            db,
		Log:           logger,
		Validate:      validate,
		TagRepository: tagRepository,
	}
}

func (c *TagUseCase) Create(ctx context.Context, request *model.CreateTagResponse) (*model.TagResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	if request.ID == 0 {
		c.Log.Warnf("Tag already exist")
		return nil, fiber.NewError(fiber.StatusBadRequest, "Tag already exist")
	}

	newName := strings.TrimSpace(request.Name)
	slug := helper.GenerateSlug(newName)

	tag := &entity.Tag{
		Name: newName,
		Slug: slug,
	}

	if err = c.TagRepository.Create(tx, tag); err != nil {
		c.Log.Warnf("Failed create tag to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.TagToCreateResponse(tag), nil
}
