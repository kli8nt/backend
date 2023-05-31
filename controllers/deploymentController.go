package controllers

import (
	"github.com/adamlahbib/gitaz/initializers"
	"github.com/adamlahbib/gitaz/models"
)

func AddDeployment(deployment models.Deployment) {
	initializers.DB.Create(&deployment)
}

func FetchDeploymentByRepoName(repoName string) models.Deployment {
	var deployment models.Deployment
	initializers.DB.First(&deployment, "repo_name = ?", repoName)
	return deployment
}
