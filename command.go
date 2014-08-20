package main

import (
	"bufio"
	"code.google.com/p/goauth2/oauth"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
	"code.google.com/p/gcfg"
	"os"
	"strconv"
	"regexp"
	"github.com/skratchdot/open-golang/open"
	"path/filepath"
)


var createFlags = []cli.Flag{
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
	prompt := c.Bool("details")
	if description := getRepositryField("description", c.String("description"), prompt).(string); description != "" {
		newRepository.Description = &description
	}
	if homepage := getRepositryField("homepage", c.String("homepage"), prompt).(string); homepage != "" {
		newRepository.Homepage = &homepage
	}
	if teamid := getRepositryField("teamid", c.Int("teamid"), prompt).(int); teamid != 0 {
		newRepository.TeamID = &teamid
	}
	private := getRepositryField("private", c.Bool("private"), prompt).(bool)
	newRepository.Private = &private
	issue := getRepositryField("issue", c.Bool("issue"), prompt).(bool)
	newRepository.HasIssues = &issue
	wiki := getRepositryField("wiki", c.Bool("wiki"), prompt).(bool)
	newRepository.HasWiki = &wiki
	download := getRepositryField("download", c.Bool("downloads"), prompt).(bool)
	newRepository.HasDownloads = &download

	// create repository
	client := newGithubClient(configuration)
	repositry, _, createErr := client.Repositories.Create("", &newRepository)
	if createErr != nil {
		fmt.Println(createErr)
		return
	}
	output := "\n\n" +
		"=========================\n" +
		"                         \n" +
		"* We are sccessful in Creating a repository! Push an existing repository from the command line\n" +
		"                         \n" +
		"git remote add origin %s \n" +
		"git push -u origin master\n" +
		"                         \n" +
		"=========================\n" +
		"\n\n"
	fmt.Printf(output, *repositry.GitURL)
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
	prompt := c.Bool("details")
	if name := getRepositryField("name", c.String("name"), prompt).(string); name != "" {
		repository.Name = &name
	}
	if description := getRepositryField("description", c.String("description"), prompt).(string); description != "" {
		repository.Description = &description
	}
	if homepage := getRepositryField("homepage", c.String("homepage"), prompt).(string); homepage != "" {
		repository.Homepage = &homepage
	}
	if issue := getRepositryField("issue", c.String("issue"), prompt).(bool); issue != *repository.HasIssues {
		repository.HasIssues = &issue
	}
	if wiki := getRepositryField("wiki", c.String("wiki"), prompt).(bool); wiki != *repository.HasWiki {
		repository.HasWiki = &wiki
	}
	if download := getRepositryField("download", c.String("download"), prompt).(bool); download != *repository.HasDownloads {
		repository.HasDownloads = &download
	}

	// edit repository
	edittedRepository, _, editErr := client.Repositories.Edit(owner, repo, repository)
	if editErr != nil {
		fmt.Println(editErr)
		return
	}
	output := "\n\n" +
		"=========================\n" +
		"                         \n" +
		"* We are sccessful in Editting the repository!\n" +
		"name:        %s          \n" +
		"description: %s          \n" +
		"homepage:    %s          \n" +
		"issue:       %s          \n" +
		"wiki:        %s          \n" +
		"download:    %s          \n" +
		"=========================\n" +
		"\n\n"
	fmt.Printf(output, edittedRepository.Name,
		edittedRepository.Description,
		edittedRepository.Homepage,
		edittedRepository.HasIssues,
		edittedRepository.HasWiki,
		edittedRepository.HasDownloads)

}

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

func doOpen(c *cli.Context) {
  dir, pathErr := filepath.Abs(filepath.Dir(os.Args[0]))
  if pathErr != nil {
    fmt.Println(pathErr)
    return
  }

  gitconfigPath := fmt.Sprintf("%s/.git/config", dir)
  finfo, configFileErr := os.Stat(gitconfigPath)
  if configFileErr != nil || finfo.IsDir() {
    fmt.Println(".git/config not found")
    return
  }

  gitconfig := &GitConfig{}
  err := gcfg.ReadFileInto(gitconfig, gitconfigPath)
  if err != nil {
    fmt.Println(err)
  }

  reg, _ := regexp.Compile("git@github.com:(.+).git")
  matches := reg.FindStringSubmatch(gitconfig.Remote["origin"].Url)
  githuburl := fmt.Sprintf("https://github.com/%s/tree/master", matches[1])
  openErr := open.Run(githuburl)
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
