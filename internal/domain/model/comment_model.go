package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Commenter string `gorm:"not null"`
	PostID    uint   `gorm:"not null"`
	Content   string `gorm:"type:text"`

	User User `gorm:"foreignKey:Commenter; references:Username"`
	Post Post `gorm:"foreignKey:PostID; references:ID"`
}
