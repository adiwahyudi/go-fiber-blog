package entity

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        string `gorm:"primaryKey;not null;type:varchar(36)"`
	Name      string `gorm:"not null;type:varchar(100)"`
	Username  string `gorm:"not null;type:varchar(30);unique"`
	Email     string `gorm:"not null;type:varchar(255);unique"`
	Password  string `gorm:"not null;type:varchar(255)"`
	Posts     []Post `gorm:"foreignKey:UserID;references:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
