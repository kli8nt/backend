package create

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// CommonOptions is the common create command options.
type CommonOptions struct {
	Owner       string
	Repo        string
	Token       string
	Description string
}

// new github client

func NewClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}
