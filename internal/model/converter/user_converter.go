package converter

import (
	"go-blog/internal/entity"
	"go-blog/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: &user.CreatedAt,
		UpdatedAt: &user.UpdatedAt,
	}
}

func UserToTokenResponse(token string) *model.UserResponse {
	return &model.UserResponse{
		Token: token,
	}
}
