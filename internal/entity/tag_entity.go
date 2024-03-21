package entity

import (
	"time"
)

type Tag struct {
	ID        uint       `gorm:"primaryKey;not null"`
	Name      string     `gorm:"unique;type:varchar(100);not null" `
	Slug      string     `gorm:"unique;type:varchar(255);not null"`
	Posts     []Post     `gorm:"many2many:post_tags"`
	CreatedAt *time.Time `gorm:"autoCreateTime"`
}
