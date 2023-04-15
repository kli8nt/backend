package create

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

// set the website for the repo

func SetWebsite(client *github.Client, owner string, repo string, website string) (*github.Repository, error) {
	ctx := context.Background()

	Repository := &github.Repository{
		Homepage: &website,
	}

	_, _, err := client.Repositories.Edit(ctx, owner, repo, Repository)
	if err != nil {
		return nil, err
	}

	return Repository, nil

}

// CreateCheckRun creates a new check run for the latest commit in the given repository and branch, and returns the created check run.
func CreateCheckRun(client *github.Client, owner, repo, branch string) (*github.CheckRun, error) {
	ctx := context.Background()

	// Get the SHA of the latest commit for the branch
	ref, _, err := client.Git.GetRef(ctx, owner, repo, fmt.Sprintf("refs/heads/%s", branch))
	if err != nil {
		return nil, fmt.Errorf("failed to get branch reference: %v", err)
	}
	sha := ref.GetObject().GetSHA()

	st := "in_progress"

	check := &github.CreateCheckRunOptions{
		Name:    "My Check",
		HeadSHA: sha,
		Status:  &st,
	}
	newCheck, _, err := client.Checks.CreateCheckRun(ctx, owner, repo, *check)
	if err != nil {
		return nil, fmt.Errorf("failed to create check run: %v", err)
	}
	return newCheck, nil
}

// UpdateCheckRunStatus updates the status of the check run with the given ID to the given status and returns any error that occurred.
func UpdateCheckRunStatus(client *github.Client, owner, repo string, checkID int64, status string) error {
	ctx := context.Background()

	options := &github.UpdateCheckRunOptions{
		Status: &status,
	}
	_, _, err := client.Checks.UpdateCheckRun(ctx, owner, repo, checkID, *options)
	if err != nil {
		return fmt.Errorf("failed to update check run status: %v", err)
	}
	return nil
}
