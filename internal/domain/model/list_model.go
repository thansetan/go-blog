package model

import "gorm.io/gorm"

type List struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Slug        string `gorm:"not null;index"`
	Description string `gorm:"type:text"`
	Owner       string `gorm:"not null"`

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
