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
	cli.BoolFlag{
		Name:  "issue, I",
		Usage: "",
	},
	cli.BoolFlag{
		Name:  "wiki, W",
		Usage: "",
	},
	cli.BoolFlag{
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

	newRepository := github.Repository{Name: &name}
	if description := c.String("description"); description != "" {
		newRepository.Description = &description
	}
	if homepage := c.String("homepage"); homepage != "" {
		newRepository.Homepage = &homepage
	}
	if private := c.Bool("private"); private != false {
		newRepository.Private = &private
	}
	if issue := c.Bool("issue"); issue != false {
		newRepository.HasIssues = &issue
	}
	if wiki := c.Bool("wiki"); wiki != false {
		newRepository.HasWiki = &wiki
	}
	if downloads := c.Bool("downloads"); downloads != false {
		newRepository.HasDownloads = &downloads
	}
	if teamid := c.Int("teamid"); teamid != 0 {
		newRepository.TeamID = &teamid
	}

	// Repositry structure doesn't have the followings
	//autoinit := c.Int("autoinit")
	//gitignore := c.String("gitignore")
	//license := c.String("license")

	scan := bufio.NewScanner(os.Stdin)
	description := ""
	if newRepository.Description != nil {
		description = *newRepository.Description
	}
	fmt.Printf("description [%s] : ", description)
	scan.Scan()
	description = scan.Text()

	homepage := ""
	if newRepository.Homepage != nil {
		homepage = *newRepository.Homepage
	}
	fmt.Printf("homepage [%s] : ", homepage)
	scan.Scan()
	homepage = scan.Text()

	private := false
	if newRepository.Private != nil {
		private = *newRepository.Private
	}
	fmt.Printf("private [%s] : ", strconv.FormatBool(private))
	scan.Scan()
	private, _ = strconv.ParseBool(scan.Text())
	//if detailFlag := c.Bool("detail"); detailFlag == true {
	//	if newRepository.description != nil {
	//		scanner.Scan()
	//		githubToken := scanner.Text()
	//	}
	//}

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
