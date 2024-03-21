package usecase

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go-blog/internal/delivery/http/middleware"
	"go-blog/internal/entity"
	"go-blog/internal/model"
	"go-blog/internal/model/converter"
	"go-blog/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
	Middleware     *middleware.Middleware
}

func NewUserUseCase(
	db *gorm.DB, logger *logrus.Logger, validate *validator.Validate, userRepository *repository.UserRepository, mddlwr *middleware.Middleware,
) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            logger,
		Validate:       validate,
		UserRepository: userRepository,
		Middleware:     mddlwr,
	}
}

func (c *UserUseCase) Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to generate bcrype hash : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	user := &entity.User{
		ID:       uuid.New().String(),
		Name:     request.Name,
		Username: request.Username,
		Email:    request.Email,
		Password: string(password),
	}
	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.Warnf("Failed create user to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body  : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindByEmail(tx, user, request.Email); err != nil {
		c.Log.Warnf("Failed find user by email : %+v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password.")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.Warnf("Failed to compare user password with bcrype hash : %+v", err)
		return nil, fiber.ErrUnauthorized
	}
	auth := &model.Auth{
		ID: user.ID,
	}
	token, err := c.Middleware.GenerateToken(auth)
	if err != nil {
		c.Log.Warnf("Failed to create JWT token : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.UserResponse{Token: token}, nil
}

func (c *UserUseCase) Update(ctx context.Context, request *model.UpdateUserRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body  : %+v", err)
		return fiber.ErrBadRequest
	}

	if err := c.UserRepository.FindById(tx, &entity.User{}, request.ID); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return fiber.ErrNotFound
	}

	updatedUser := new(entity.User)
	updatedUser.Name = request.Name
	updatedUser.Email = request.Email
	if request.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			c.Log.Warnf("Failed to generate bcrype hash : %+v", err)
			return fiber.ErrInternalServerError
		}
		updatedUser.Password = string(password)
	}

	if err := c.UserRepository.Updates(tx, updatedUser, request.ID); err != nil {
		c.Log.Warnf("Failed save user : %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return fiber.ErrInternalServerError
	}

	return nil
}

func (c *UserUseCase) Delete(ctx context.Context, id string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Var(id, "uuid4"); err != nil {
		c.Log.Warnf("Invalid user id : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user id")
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, id); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return fiber.ErrNotFound
	}

	if err := c.UserRepository.Delete(tx, user); err != nil {
		c.Log.Warnf("Failed delete user : %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return fiber.ErrInternalServerError
	}

	return nil
}
