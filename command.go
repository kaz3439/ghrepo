package main

import (
	"bufio"
	"code.google.com/p/goauth2/oauth"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
	"github.com/skratchdot/open-golang/open"
	"github.com/wsxiaoys/terminal/color"
	"os"
	"strconv"
	"time"
)


var createFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "organization, O",
		Usage: "",
	},
	cli.StringFlag{
		Name:  "description, D",
		Usage: "",
	},
	cli.StringFlag{
		Name:  "homepage, H",
		Usage: "",
	},
	cli.BoolFlag{
		Name:  "private, P",
		Usage: "",
	},
	cli.BoolTFlag{
		Name:  "issue, I",
		Usage: "",
	},
	cli.BoolTFlag{
		Name:  "wiki, W",
		Usage: "",
	},
	cli.BoolTFlag{
		Name:  "download, X",
		Usage: "",
	},
	cli.IntFlag{
		Name:  "teamid, T",
		Usage: "",
	},
	cli.BoolFlag{
		Name:  "details, d",
		Usage: "",
	},
}

const (
  flagOrg = "organization"
  flagDesc = "description"
  flagHP = "homepage"
  flagPrivate = "private"
  flagIssue = "issue"
  flagWiki = "wiki"
  flagTeamID = "teamid"
  flagDownload = "download"
  flagName = "name"
  flagDetail = "details"
)

var editFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "name, N",
		Usage: "",
	},
	cli.StringFlag{
		Name:  "description, D",
		Usage: "",
	},
	cli.StringFlag{
		Name:  "homepage, H",
		Usage: "",
	},
	cli.BoolTFlag{
		Name:  "issue, I",
		Usage: "",
	},
	cli.BoolTFlag{
		Name:  "wiki, W",
		Usage: "",
	},
	cli.BoolTFlag{
		Name:  "download, X",
		Usage: "",
	},
	cli.BoolFlag{
		Name:  "details, d",
		Usage: "",
	},
}

var Commands = []cli.Command{
	commandCreate,
	commandOpen,
	commandEdit,
	commandToken,
}

var tokenSubCommands = []cli.Command{
	commandSetToken,
	commandResetToken,
}

var commandCreate = cli.Command{
	Name:        "create",
	ShortName:   "c",
	Usage:       "create repository",
	Description: "",
	Flags:       createFlags,
	Action:      doCreate,
}

var commandOpen = cli.Command{
	Name:        "open",
	ShortName:   "o",
	Usage:       "open repository on browser",
	Description: "",
	Action: doOpen,
}

var commandEdit = cli.Command{
	Name:        "edit",
	ShortName:   "e",
	Usage:       "edit repository",
	Description: "",
	Flags:       editFlags,
	Action:      doEdit,
}

var commandToken = cli.Command{
	Name:        "token",
	ShortName:   "t",
	Usage:       "token manager",
	Description: "",
	Subcommands: tokenSubCommands,
}

var commandSetToken = cli.Command{
	Name:   "set",
	Usage:  "",
	Action: doSetToken,
}

var commandResetToken = cli.Command{
	Name:   "reset",
	Usage:  "",
	Action: doResetToken,
}

func doSetToken(c *cli.Context) {
	token := c.Args().First()
	if len(token) == 0 {
		return
	}

	var configuration = NewConfiguration(token)
	if err := configuration.Persist(); err != nil {
		fmt.Println(err)
	}
}

func doResetToken(c *cli.Context) {
	var configuration = NewConfiguration("")
	if err := configuration.Persist(); err != nil {
		fmt.Println(err)
	}
}

func doCreate(c *cli.Context) {
	name := c.Args().First()
	if len(name) == 0 {
		return
	}

	// check configuration
	configuration, err := OpenConfiguration()
	if err != nil {
		fmt.Println(err)
		return
	}
	if configuration.GithubToken == "" {
		configuration.GithubToken = promptPersonalGithubToken()
		configuration.Persist()
	}

	// set repository attributes
	newRepository := github.Repository{Name: &name}
	prompt := c.Bool(flagDetail)
	if description := getRepositryField(flagDesc, c.String(flagDesc), prompt).(string); description != "" {
		newRepository.Description = &description
	}
	if homepage := getRepositryField(flagHP, c.String(flagHP), prompt).(string); homepage != "" {
		newRepository.Homepage = &homepage
	}
	if teamid := getRepositryField(flagTeamID, c.Int(flagTeamID), prompt).(int); teamid != 0 {
		newRepository.TeamID = &teamid
	}
	private := getRepositryField(flagPrivate, c.Bool(flagPrivate), prompt).(bool)
	newRepository.Private = &private
	issue := getRepositryField(flagIssue, c.Bool(flagIssue), prompt).(bool)
	newRepository.HasIssues = &issue
	wiki := getRepositryField(flagWiki, c.Bool(flagWiki), prompt).(bool)
	newRepository.HasWiki = &wiki
	download := getRepositryField(flagDownload, c.Bool(flagDownload), prompt).(bool)
	newRepository.HasDownloads = &download

	// create repository
	org := c.String(flagOrg)
	client := newGithubClient(configuration)
	networkError := make(chan error)
	resultRepository := make(chan *github.Repository)
	go func () {
	  repository, _, createErr := client.Repositories.Create(org, &newRepository)
	  if createErr != nil {
	    networkError <- createErr
	  }
	  resultRepository <- repository
	}()

	var repository *github.Repository

	loop:
	  for {
	    select {
	    case createErr := <-networkError:
		fmt.Printf("\n\n")
		color.Printf("@{r}!!! Error Occured !!!")
		fmt.Printf("\n\n")
		fmt.Println(createErr)
		fmt.Printf("\n\n")
		break loop
	      case repository = <-resultRepository:
		break loop
	      case <-time.After(time.Minute):
		fmt.Println(" @{r} !!! Timeout !!! ")
		break loop
	      default:
		time.Sleep(time.Second/2)
		fmt.Printf(".")
	    }
	  }

	if repository == nil {
	  return
	}

	output := "\n\n" +
		"=========================\n" +
		"                         \n" +
		"@{g}* We are sccessful in Creating a repository! Push an existing repository from the command line@{|}\n" +
		"                         \n" +
		"git remote add origin %s \n" +
		"git push -u origin master\n" +
		"                         \n" +
		"=========================\n" +
		"\n\n"
	color.Printf(output, *repository.GitURL)
}

func doEdit(c *cli.Context) {
	owner := c.Args().Get(0)
	repo := c.Args().Get(1)
	if len(owner) == 0 || len(repo) == 0 {
		return
	}

	// check configuration
	configuration, err := OpenConfiguration()
	if err != nil {
		fmt.Println(err)
		return
	}
	if configuration.GithubToken == "" {
		configuration.GithubToken = promptPersonalGithubToken()
		configuration.Persist()
	}

	client := newGithubClient(configuration)

	// get repository
	repository, _, err := client.Repositories.Get(owner, repo)
	if err != nil {
		fmt.Println(err)
		return
	}

	// set repository attributes
	repository = &github.Repository{}
	prompt := c.Bool(flagDetail)
	if name := getRepositryField(flagName, c.String(flagName), prompt).(string); name != "" {
		repository.Name = &name
	}
	if description := getRepositryField(flagDesc, c.String(flagDesc), prompt).(string); description != "" {
		repository.Description = &description
	}
	if homepage := getRepositryField(flagHP, c.String(flagHP), prompt).(string); homepage != "" {
		repository.Homepage = &homepage
	}
	if issue := getRepositryField(flagIssue, c.String(flagIssue), prompt).(bool); issue != *repository.HasIssues {
		repository.HasIssues = &issue
	}
	if wiki := getRepositryField(flagWiki, c.String(flagWiki), prompt).(bool); wiki != *repository.HasWiki {
		repository.HasWiki = &wiki
	}
	if download := getRepositryField(flagDownload, c.String(flagDownload), prompt).(bool); download != *repository.HasDownloads {
		repository.HasDownloads = &download
	}

	// edit repository
	connectionErr := make(chan error)
	resultRepository := make(chan *github.Repository)
	go func () {
	  edittedRepository, _, editErr := client.Repositories.Edit(owner, repo, repository)
	  if editErr != nil {
	    connectionErr <- editErr
	    return
	  }
	  resultRepository <- edittedRepository
	}()

	var edittedRepository *github.Repository
	loop:
	  for {
	    select {
	      case editErr := <-connectionErr:
		fmt.Printf("\n\n")
		color.Printf("@{r} !!! Error Occuered !!! ")
		fmt.Printf("\n\n")
		fmt.Println(editErr)
		fmt.Printf("\n\n")
		break loop
	      case edittedRepository = <-resultRepository:
		break loop
	      case <-time.After(time.Minute):
		fmt.Println(" @{r} !!! Timeout !!! ")
		break loop
	      default:
		time.Sleep(time.Second)
		fmt.Printf(".")
	    }
	  }

	if edittedRepository == nil {
	  return
	}

	output := "\n\n" +
		"=========================\n" +
		"                         \n" +
		"@{g}* We are sccessful in Editting the repository!\n@{|}" +
		"name:        %s          \n" +
		"description: %s          \n" +
		"homepage:    %s          \n" +
		"issue:       %s          \n" +
		"wiki:        %s          \n" +
		"download:    %s          \n" +
		"=========================\n" +
		"\n\n"
	color.Printf(output, *edittedRepository.Name,
		*edittedRepository.Description,
		*edittedRepository.Homepage,
		*edittedRepository.HasIssues,
		*edittedRepository.HasWiki,
		*edittedRepository.HasDownloads)

}

func doOpen(c *cli.Context) {
  remote := c.Args().Get(0)
  tree := c.Args().Get(1)

  gitconfig, configErr := NewGitConfig()
  if configErr != nil {
    fmt.Println(configErr)
    return
  }

  if remote == "" {
    result, currentBranchErr := gitconfig.RemoteURLFromCurrentBranch()
    if currentBranchErr != nil {
      fmt.Println(currentBranchErr)
      return
    }
    remote = result
  }

  githubUrl, remoteURLErr := gitconfig.RemoteURL(remote)
  if remoteURLErr != nil {
    fmt.Println(remoteURLErr)
  }

  var openURL string
  if tree != "" {
    openURL = fmt.Sprintf("%s/tree/%s", githubUrl, tree)
  } else {
    openURL = githubUrl
  }
  openErr := open.Run(openURL)
  if openErr != nil {
    fmt.Println(openErr)
  }
}

func promptPersonalGithubToken() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("your personal github token: ")
	scanner.Scan()
	githubToken := scanner.Text()
	return githubToken
}

func newGithubClient(configuration *Configuration) *github.Client {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: configuration.GithubToken},
	}

	client := github.NewClient(t.Client())
	return client
}

func getRepositryField(name string, field interface{}, prompt bool) interface{} {
	if prompt == false {
		return field
	}

	var defaultPrompt string
	switch field.(type) {
	case string:
		defaultPrompt = field.(string)
	case bool:
		defaultPrompt = strconv.FormatBool(field.(bool))
	case int:
		if field.(int) != 0 {
			defaultPrompt = strconv.FormatInt(int64(field.(int)), 10)
		} else {
			defaultPrompt = ""
		}
	default:
		defaultPrompt = ""
	}

	fmt.Printf("%s [%s] : ", name, defaultPrompt)
	scan := bufio.NewScanner(os.Stdin)
	scan.Scan()
	text := scan.Text()
	switch field.(type) {
	case bool:
		if res, err := strconv.ParseBool(text); err != nil {
			field = res
		}
	case int:
		if res, err := strconv.ParseInt(text, 10, 0); err != nil {
			field = int(res)
		}
	default:
		field = text
	}
	return field
}

