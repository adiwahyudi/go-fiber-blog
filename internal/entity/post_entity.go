package entity

import (
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID          uint           `gorm:"primaryKey;not null"`
	Title       string         `gorm:"type:varchar(100);not null"`
	Slug        string         `gorm:"type:varchar(255);not null"`
	Content     string         `gorm:"type:longtext;not null"`
	UserID      string         `gorm:"type:varchar(36)"`
	Tags        []*Tag         `gorm:"many2many:post_tags"`
	User        User           `gorm:"foreignKey:UserID;references:ID"`
	PublishedAt *time.Time     `gorm:"TIMESTAMP NULL"`
	CreatedAt   *time.Time     `gorm:"autoCreateTime"`
	UpdatedAt   *time.Time     `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
