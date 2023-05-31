package models

import "gorm.io/gorm"

type Repo struct {
	gorm.Model
	OwnerID uint
	Owner   User   `gorm:"foreignKey:OwnerID;"`
	Name    string `json:"name"`
	RepoUrl    string `json:"repo_url"`
	Branch  string `json:"branch"`
}
