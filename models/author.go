package models

import "gorm.io/gorm"

type Author struct {
	gorm.Model

	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Image    string `json:"image"`
	Admin    bool   `json:"admin"`
}
