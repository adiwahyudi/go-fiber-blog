package repository

import (
	"github.com/sirupsen/logrus"
	"go-blog/internal/entity"
	"gorm.io/gorm"
)

type TagRepository struct {
	Repository[entity.Tag]
	Log *logrus.Logger
}

func NewTagRepository(log *logrus.Logger) *TagRepository {
	return &TagRepository{
		Log: log,
	}
}

func (r *TagRepository) FindByName(db *gorm.DB, tag *entity.Tag, name string) error {
	return db.Take(tag, "name = ?", name).Error
}
