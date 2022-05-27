package models

import (
	pq "github.com/lib/pq"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title   string         `json:"title"`
	Excerpt string         `json:"excerpt"`
	Body    string         `json:"body"`
	Tags    pq.StringArray `gorm:"type:text[]" json:"tags"`
	Views   int            `json:"views"`
	Author  string         `json:"author"`
	Cover   string         `json:"cover"`
	Likes   int            `json:"likes"`
	Slug    string         `json:"slug"`
}
