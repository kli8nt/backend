package create

import (
	"context"
	"log"

	"github.com/cloudflare/cloudflare-go"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// new cloudflare client

func NewCloudflareClient(token string, email string) *cloudflare.API {
	client, err := cloudflare.New(token, email)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

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
