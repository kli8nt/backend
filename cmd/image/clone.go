package image

import (
	"fmt"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func CloneRepo(url, token string, dirname) error {
	// Clone the repository into the current directory
	_, err := git.PlainClone("./"+dirname, false, &git.CloneOptions{
		URL:      url,
		Auth:     &http.BasicAuth{Username: "token", Password: token},
		Progress: os.Stdout,
	})

	if err != nil {
		return err
	}

	fmt.Println("Repository cloned successfully!")
	return nil
}
