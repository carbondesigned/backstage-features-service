package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title   string   `json:"title"`
	Excerpt string   `json:"excerpt"`
	Body    string   `json:"body"`
	Tags    []string `gorm:"many2many:post_tags;" json:"tags"`
	Views   int      `json:"views"`
	Author  string   `json:"author"`
	Cover   string   `json:"cover"`
	Likes   int      `json:"likes"`
	Slug    string   `json:"slug"`
}
