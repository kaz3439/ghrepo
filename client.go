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

func (client *Client) CreateRepository(org string, newRepository *github.Repository) (<-chan *github.Repository, <-chan error) {
  networkErr := make(chan error)
  resultRepository := make(chan *github.Repository)
  go func () {
    repository, _, createErr := client.client.Repositories.Create(org, newRepository)
    if createErr != nil {
      networkErr <- createErr
      return
    }
    resultRepository <- repository
  }()
  return resultRepository, networkErr
}

func (client *Client) GetRepository(owner string, repo string) (<-chan *github.Repository, <-chan error ) {
    networkErr := make(chan error)
    resultRepository := make(chan *github.Repository)
    go func() {
      repository, _, err := client.client.Repositories.Get(owner, repo)
      if err != nil {
	networkErr <- err
	return
      }
      resultRepository <- repository
    }()
    return resultRepository, networkErr
}

func (client *Client) EditRepository(owner string, repo string, repository *github.Repository) (<-chan *github.Repository, <-chan error) {
    networkErr := make(chan error)
    resultRepository := make(chan *github.Repository)
    go func () {
      edittedRepository, _, editErr := client.client.Repositories.Edit(owner, repo, repository)
      if editErr != nil {
	networkErr <- editErr
	return
      }
      resultRepository <- edittedRepository
    }()
    return resultRepository, networkErr
}
