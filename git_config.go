package main

import (
  "os"
  "fmt"
  "code.google.com/p/gcfg"
  "path/filepath"
  "errors"
)

const gitConfigRelativePath = ".git/config"

type GitConfig struct {
  Core struct {
    Repositoryformatversion int
    Filemode bool
    Bare bool
    Logallrefupdates bool
    Ignorecase bool
    Precomposeunicode bool
  }
  Remote map[string]*struct {
    Url string
    Fetch string
  }
  Branch map[string]*struct {
    Remote string
    Merge string
  }
}

func NewGitConfig() (*GitConfig, error) {
  currentDirectory, pathErr := filepath.Abs(filepath.Dir(os.Args[0]))
  if pathErr != nil {
    return nil, pathErr
  }
  gitconfigPath := fmt.Sprintf("%s/%s", currentDirectory, gitConfigRelativePath)
  finfo, configFileErr := os.Stat(gitconfigPath)
  if configFileErr != nil || finfo.IsDir() {
    fileNotFoundError := errors.New("configuration file not found")
    return nil, fileNotFoundError
  }

  gitconfig := &GitConfig{}
  err := gcfg.ReadFileInto(gitconfig, gitconfigPath)
  if err != nil {
    return nil, err
  }
  return gitconfig, nil
}
