package main

import (
	"bufio"
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
	"os"
	"strconv"
)

type Configuration struct {
	GithubToken string
}

func (configuration *Configuration) Persist() error {
	bdata, marshallErr := json.Marshal(configuration)
	if marshallErr != nil {
		return marshallErr
	}
	os.Remove("config.json")
	file, createErr := os.Create("config.json")
	if createErr != nil {
		return createErr
	}
	file.WriteString(string(bdata))
	return nil
}

func NewConfiguration(token string) *Configuration {
	return &Configuration{token}
}

func OpenConfiguration() (*Configuration, error) {
	if _, err := os.Stat("config.json"); err != nil {
		return nil, err
	}

	file, openErr := os.Open("config.json")
	if openErr != nil {
		return nil, openErr
	}

	var configuration = &Configuration{}
	decoder := json.NewDecoder(file)
	decodeErr := decoder.Decode(&configuration)
	return configuration, decodeErr
}

var createFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "description, D",
		Usage: "",
	},
	cli.StringFlag{
		Name:  "homepage, H",
		Usage: "",
	},

	cli.StringFlag{
		Name:  "gitignore, G",
		Usage: "",
	},
	cli.StringFlag{
		Name:  "license, L",
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
	cli.BoolFlag{
		Name:  "autoinit, A",
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
	Action: func(c *cli.Context) {
	},
}

var commandEdit = cli.Command{
	Name:        "edit",
	ShortName:   "e",
	Usage:       "edit repository",
	Description: "",
	Action: func(c *cli.Context) {
	},
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

	prompt := c.Bool("details")
	newRepository := github.Repository{Name: &name}

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

	// Repositry structure doesn't have the followings
	//autoinit := c.Int("autoinit")
	//gitignore := c.String("gitignore")
	//license := c.String("license")

	configuration, err := OpenConfiguration()
	if err != nil {
		fmt.Println(err)
	}

	if configuration.GithubToken == "" {
		configuration.GithubToken = promptPersonalGithubToken()
		configuration.Persist()
	}

	client := newGithubClient(configuration)
	repositry, _, createErr := client.Repositories.Create("", &newRepository)
	if createErr != nil {
		return
	}
	fmt.Printf("%#v", repositry)
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
