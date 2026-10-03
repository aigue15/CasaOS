package github

import (
	"github.com/google/go-github/v92/github"
)

func GetGithubClient() (*github.Client, error) {
	return github.NewClient()
}
