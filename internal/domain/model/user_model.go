package model

import "time"

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"not null;unique;default:null;type:varchar(255)"`
	Username string `gorm:"not null;uniqueIndex;type:varchar(255)"`
	Name     string `gorm:"not null;default:null;type:varchar(255)"`
	Password []byte `gorm:"not null;default:null;type:bytea"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
