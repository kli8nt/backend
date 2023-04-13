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

// check if hook exists

func HookExists(client *github.Client, owner string, repo string, url string) (bool, error) {

	hooks, _, err := client.Repositories.ListHooks(context.Background(), owner, repo, nil)
	if err != nil {
		return false, err
	}

	for _, hook := range hooks {
		if hook.Config["url"] == url {
			return true, nil
		}
	}

	return false, nil
}
