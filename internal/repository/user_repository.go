package repository

import (
	"github.com/sirupsen/logrus"
	"go-blog/internal/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (r *Repository[T]) FindByEmail(tx *gorm.DB, user *entity.User, email string) error {
	return tx.Where("email = ?", email).Take(user).Error
}
func (r *Repository[T]) FindByUsername(tx *gorm.DB, user *entity.User, username string) error {
	return tx.Where("username = ?", username).Take(user).Error
}
