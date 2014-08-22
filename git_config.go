package main

import (
  "os"
  "fmt"
  "code.google.com/p/gcfg"
  "path/filepath"
  "errors"
  "regexp"
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

const (
  githubURL = "https://github.com"
  sshURLRegex = "git@github.com:(.+).git"
  svnURLRegex = "(https://github.com/.+).git"
)

func (gitconfig *GitConfig) RemoteURL(remoteName string) (string, error) {
  if remoteIsURLPath(remoteName) {
    return fmt.Sprintf("%s/%s", githubURL, remoteName), nil
  }

  remote := gitconfig.Remote[remoteName]
  if remote == nil {
    return "", errors.New("invalid remote name")
  }

  remoteUrl := remote.Url
  sshRegex := regexp.MustCompile(sshURLRegex)
  svnRegex := regexp.MustCompile(svnURLRegex)
  var githubUrl string
  if  sshMatches := sshRegex.FindStringSubmatch(remoteUrl); len(sshMatches) == 2 { //ssh
    githubUrl = fmt.Sprintf("%s/%s", githubURL, sshMatches[1])
  } else if svnMatches := svnRegex.FindStringSubmatch(remoteUrl); len(svnMatches) == 2  { //svn
    githubUrl = svnMatches[1]
  } else { //https
    githubUrl = remoteUrl
  }
  return githubUrl, nil
}

func (gitconfig *GitConfig) RemoteURLFromCurrentBranch() (string, error) {
    var tree string
    for key := range gitconfig.Branch {
      tree = key
      break
    }
    if gitconfig.Branch[tree] == nil {
      return "", errors.New("invalid branch name")
    }
    return gitconfig.Branch[tree].Remote, nil
}

func remoteIsURLPath(remote string) bool {
  if regexp.MustCompile("^[^/]+/[^/]+$").MatchString(remote) {
    return true
  } else {
    return false
  }
}
