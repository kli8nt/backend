package models

import "gorm.io/gorm"

type History struct {
	gorm.Model
	DeploymentID uint
	Deployment   Deployment `gorm:"foreignKey:DeploymentID;"`
	LogsURL      string     `json:"logs_url"`
	Status       string     `json:"status"`
}
