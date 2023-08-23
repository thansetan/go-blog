package model

import "time"

type List struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null;type:varchar(255)"`
	Slug        string `gorm:"not null;index;type:varchar(510)"`
	Description string `gorm:"type:text"`
	Owner       string `gorm:"not null;type:varchar(255)"`

	CreatedAt time.Time
	UpdatedAt time.Time

	User      User   `gorm:"foreignKey:Owner; references:Username"`
	Posts     []Post `gorm:"many2many:list_posts"`
	ListPosts []ListPost
}

type ListPost struct {
	ListID uint `gorm:"uniqueIndex:idx_list_id_post_id"`
	PostID uint `gorm:"uniqueIndex:idx_list_id_post_id"`

	List List `gorm:"foreignKey:ListID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Post Post `gorm:"foreignKey:PostID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ListPost) TableName() string {
	return "list_posts"
}
