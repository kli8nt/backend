package controllers

import (
	"github.com/adamlahbib/gitaz/initializers"
	"github.com/adamlahbib/gitaz/models"
)

func AddDeployment(deployment models.Deployment) {
	initializers.DB.Create(&deployment)
}
