package main

import (
  "github.com/google/go-github/github"
  "code.google.com/p/goauth2/oauth"
)

type Client struct {
  client *github.Client
}

func NewClient(configuration *Configuration) *Client {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: configuration.GithubToken},
	}

	client := &Client{github.NewClient(t.Client())}
	return client
}

func (client *Client) CreateRepository(org string, newRepository *github.Repository) (*github.Repository, error ){
  repository, _, createErr := client.client.Repositories.Create(org, newRepository)
  return repository, createErr
}

func (client *Client) GetRepository(owner string, repo string) (*github.Repository, error ) {
    repository, _, err := client.client.Repositories.Get(owner, repo)
    return repository, err
}

func (client *Client) EditRepository(owner string, repo string, repository *github.Repository) (*github.Repository, error ) {
    edittedRepository, _, editErr := client.client.Repositories.Edit(owner, repo, repository)
    return edittedRepository, editErr
}
