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

	prompt := c.Bool("detail")
	newRepository := github.Repository{Name: &name}
	if c.String("description") != "" {
		newRepository.Description = getRepositryField(c.String("description"), prompt).(*string)
	}
	if c.String("homepage") != "" {
		newRepository.Homepage = getRepositryField(c.String("homepage"), prompt).(*string)
	}
	if c.Int("teamid") != 0 {
		newRepository.TeamID = getRepositryField(c.Int("teamid"), prompt).(*int)
	}

	newRepository.Private = getRepositryField(c.Bool("private"), prompt).(bool)
	newRepository.HasIssues = getRepositryField(c.Bool("issue"), prompt).(bool)
	newRepository.HasWiki = getRepositryField(c.Bool("wiki"), prompt).(bool)
	newRepository.HasDownloads = getRepositryField(c.Bool("downloads"), prompt).(bool)

	fmt.Printf("%#v", newRepository)
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

	//client := setClient(configuration)
	//repositry, _, createErr := client.Repositories.Create("", &newRepository)
	//if createErr != nil {
	//	return
	//}
	//fmt.Printf("%#v", repositry)
}

func promptPersonalGithubToken() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("your personal github token: ")
	scanner.Scan()
	githubToken := scanner.Text()
	return githubToken
}

func setClient(configuration *Configuration) *github.Client {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: configuration.GithubToken},
	}

	client := github.NewClient(t.Client())
	return client
}

func getRepositryField(field interface{}, prompt bool) interface{} {
	if prompt == false {
		return field
	}

	scan := bufio.NewScanner(os.Stdin)
	defaultValue := ""
	switch field.(type) {
	case string:
		defaultValue = field.(string)
	case bool:
		defaultValue = strconv.FormatBool(field.(bool))
	case int:
		defaultValue = strconv.FormatInt(field.(int64), 10)
	}

	fmt.Printf("description [%s] : ", defaultValue)
	scan.Scan()
	field = scan.Text()
	return field
}
