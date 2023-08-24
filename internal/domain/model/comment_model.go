package model

import "time"

type Comment struct {
	ID        uint   `gorm:"primaryKey"`
	Commenter string `gorm:"not null;type:varchar(255)"`
	PostID    uint   `gorm:"not null"`
	Content   string `gorm:"type:text;type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time

	User User `gorm:"foreignKey:Commenter;references:Username"`
	Post Post `gorm:"foreignKey:PostID; references:ID;constraint:OnDelete:CASCADE"`
}
