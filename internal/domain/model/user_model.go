package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"not null;unique;default:null"`
	Username string `gorm:"not null;uniqueIndex"`
	Name     string `gorm:"not null;default:null"`
	Password string `gorm:"not null;default:null"`
}
