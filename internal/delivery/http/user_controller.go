package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go-blog/internal/delivery/http/middleware"
	"go-blog/internal/model"
	"go-blog/internal/usecase"
)

type UserController struct {
	Log     *logrus.Logger
	UseCase *usecase.UserUseCase
}

func NewUserController(logger *logrus.Logger, useCase *usecase.UserUseCase) *UserController {
	return &UserController{
		Log:     logger,
		UseCase: useCase,
	}
}

// RegisterUser godoc
// @Tags Auth
// @Summary Register a new user
// @Description API for register a new user.
// @ID register-user
// @Router /api/auth/register [post]
// @Accept json
// @Param register body model.RegisterUserRequest true "Request"
// @Produce json
// @Success 201
func (c *UserController) RegisterUser(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)

	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})

}

// Login godoc
// @Tags Auth
// @Summary Login to services
// @Description API for auth login.
// @ID login-user
// @Router /api/auth/login [post]
// @Accept json
// @Param register body model.LoginUserRequest true "Request"
// @Produce json
// @Success 201
func (c *UserController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Login(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

// Update godoc
// @Tags Users
// @Summary Update User
// @Description API for update user that currently logged in.
// @ID update-current-user
// @Security Bearer
// @Router /api/users [patch]
// @Accept json
// @Param register body model.UpdateUserRequest true "Request"
// @Produce json
// @Success 201
func (c *UserController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := new(model.UpdateUserRequest)
	request.ID = auth.ID

	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	if err := c.UseCase.Update(ctx.UserContext(), request); err != nil {
		c.Log.WithError(err).Error("error updating user")
		return err
	}

	return ctx.JSON(model.WebResponse[string]{Data: "Successfully update user"})

}

// Delete godoc
// @Tags Users
// @Summary Delete user by ID.
// @Description API for delete user by ID, restricted for admin only.
// @ID delete-user
// @Security Bearer
// @Router /api/users/{userId} [delete]
// @Accept json
// @Param userId path string true "User ID"
// @Produce json
// @Success 201
func (c *UserController) Delete(ctx *fiber.Ctx) error {
	userId := ctx.Params("userId")

	if err := c.UseCase.Delete(ctx.UserContext(), userId); err != nil {
		c.Log.WithError(err).Error("error updating user")
		return err
	}

	return ctx.JSON(model.WebResponse[string]{Data: "Successfully delete user"})

}
