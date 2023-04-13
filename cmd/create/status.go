package create

import (
	"context"

	"github.com/google/go-github/github"
)

// CreateDeploymentStatus creates a new deployment status

func CreateDeploymentStatus(client *github.Client, owner string, repo string, deploymentID int64, state string, description string, logURL string) (*github.DeploymentStatus, error) {

	deploymentStatusRequest := &github.DeploymentStatusRequest{
		State:       &state,
		Description: &description,
		//Environment: &environment,
		LogURL: &logURL,
	}

	deploymentStatus, _, err := client.Repositories.CreateDeploymentStatus(context.Background(), owner, repo, deploymentID, deploymentStatusRequest)
	if err != nil {
		return nil, err
	}

	return deploymentStatus, nil
}

// UpdateDeploymentStatus updates a deployment status

func UpdateDeploymentStatus(client *github.Client, owner string, repo string, deploymentID int64, deploymentStatusID int64, state string, description string, logURL string) (*github.DeploymentStatus, error) {

	deploymentStatusRequest := &github.DeploymentStatusRequest{
		State:       &state,
		Description: &description,
		//Environment: &environment,
		LogURL: &logURL,
	}

	deploymentStatus, _, err := client.Repositories.CreateDeploymentStatus(context.Background(), owner, repo, deploymentID, deploymentStatusRequest)
	if err != nil {
		return nil, err
	}

	return deploymentStatus, nil
}

// GetDeploymentStatus gets a deployment status

func GetDeploymentStatus(client *github.Client, owner string, repo string, deploymentID int64, deploymentStatusID int64) (*github.DeploymentStatus, error) {

	deploymentStatus, _, err := client.Repositories.GetDeploymentStatus(context.Background(), owner, repo, deploymentID, deploymentStatusID)
	if err != nil {
		return nil, err
	}

	return deploymentStatus, nil
}

// ListDeploymentStatuses lists all deployment statuses

func ListDeploymentStatuses(client *github.Client, owner string, repo string, deploymentID int64) ([]*github.DeploymentStatus, error) {

	deploymentStatuses, _, err := client.Repositories.ListDeploymentStatuses(context.Background(), owner, repo, deploymentID, nil)
	if err != nil {
		return nil, err
	}

	return deploymentStatuses, nil
}
