package create

import (
	"context"

	"github.com/google/go-github/github"
)

// create github hook to listen for events on a repo

func CreateHook(client *github.Client, owner string, repo string, url string, events []string) (*github.Hook, error) {

	hookRequest := &github.Hook{
		Events: events,
		Config: map[string]interface{}{
			"url":          url,
			"content_type": "json",
		},
	}

	hook, _, err := client.Repositories.CreateHook(context.Background(), owner, repo, hookRequest)
	if err != nil {
		return nil, err
	}

	return hook, nil
}
