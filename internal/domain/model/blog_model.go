package model

import "time"

type Blog struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null;type:varchar(255)"`
	Description string `gorm:"type:text"`
	Owner       string `gorm:"uniqueIndex;not null;type:varchar(255)"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Posts []Post
	User  User `gorm:"foreignKey:Owner;references:Username"`
}
