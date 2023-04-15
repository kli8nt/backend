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

func CreateStatusCheck(client *github.Client, owner string, repo string, branch string, state string, description string, targetURL string, gcontext string) (*github.RepoStatus, error) {

	ctx := context.Background()

	commit, _, err := client.Repositories.GetCommit(ctx, owner, repo, branch)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit for branch %s: %v", branch, err)
	}
	sha := commit.GetSHA()

	status := &github.RepoStatus{
		State:       &state,
		Description: &description,
		TargetURL:   &targetURL,
		Context:     &gcontext,
	}

	repoStatus, _, err := client.Repositories.CreateStatus(ctx, owner, repo, sha, status)
	if err != nil {
		return nil, err
	}

	return repoStatus, nil
}
