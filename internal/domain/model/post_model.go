package model

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title     string `gorm:"not null"`
	TitleSlug string `gorm:"not null; index"`
	Content   string `gorm:"not null; type:text"`
	Author    string `gorm:"not null"`
	BlogID    uint   `gorm:"not null; index"`

	Blog Blog `gorm:"foreignKey:BlogID; references:ID"`
}
