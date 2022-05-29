package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Album struct {
	gorm.Model
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Cover       string         `json:"cover"`
	CoverURL    string         `json:"coverURL"`
	Slug        string         `json:"slug"`
	RootImages  pq.StringArray `gorm:"type:text[]" json:"root_images"`
	Images      pq.StringArray `gorm:"type:text[]" json:"images"`
}
