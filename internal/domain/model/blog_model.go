package model

import "gorm.io/gorm"

type Blog struct {
	gorm.Model
	Name  string `gorm:"not null"`
	Owner string `gorm:"uniqueIndex; not null"`

	User User `gorm:"foreignKey:Owner;references:Username"`
}
