package model

import "time"

type Post struct {
	ID      uint   `gorm:"primaryKey"`
	Title   string `gorm:"not null;type:varchar(255)"`
	Slug    string `gorm:"not null;index;type:varchar(510)"`
	Content string `gorm:"not null;type:text"`
	// Author  string `gorm:"not null;type:varchar(255)"`
	BlogID uint `gorm:"not null;index"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Blog Blog `gorm:"foreignKey:BlogID; references:ID"`
}
