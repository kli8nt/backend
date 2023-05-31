package controllers

import (
	"fmt"

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

// fetch all deployments by username

func FetchDeploymentsByUsername(username string) []models.Deployment {
	var deployments []models.Deployment
	initializers.DB.Find(&deployments)
	fmt.Println(deployments)
	// initializers.DB.Find(&deployments, "username = ?", username)
	return deployments
}
