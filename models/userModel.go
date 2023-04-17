package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Name     string `json:"name"`
	Token    string `json:"token"`
}