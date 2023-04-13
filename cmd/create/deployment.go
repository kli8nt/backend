package create

import (
	"context"

	"github.com/google/go-github/github"
)

// CreateDeployment creates a new deployment

func CreateDeployment(client *github.Client, owner string, repo string, ref string, environment string, description string) (*github.Deployment, error) {

	deploymentRequest := &github.DeploymentRequest{
		Ref:         &ref,
		Environment: &environment,
		Description: &description,
	}

	deployment, _, err := client.Repositories.CreateDeployment(context.Background(), owner, repo, deploymentRequest)
	if err != nil {
		return nil, err
	}

	return deployment, nil
}
