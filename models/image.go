package models

import "gorm.io/gorm"

type Image struct {
  gorm.Model
  Image string `json:"image"`
  ImageURL string `json:"imageURL"`
  Alt string `json:"alt"`
}
