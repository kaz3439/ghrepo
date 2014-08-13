package main

import (
	"bufio"
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
	"os"
)

type Configuration struct {
	GithubToken string
}

var tokenFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "token, t",
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
	Flags:       tokenFlags,
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

	var configuration = Configuration{""}
	configuration.GithubToken = token
	bdata, err := json.Marshal(configuration)
	if err != nil {
		fmt.Println(err)
		return
	}
	os.Remove("config.json")
	file, createErr := os.Create("config.json")
	if createErr != nil {
		return
	}
	file.WriteString(string(bdata))
}

func doResetToken(c *cli.Context) {
	os.Remove("config.json")
	file, createErr := os.Create("config.json")
	if createErr != nil {
		return
	}
	var configuration = Configuration{""}
	bdata, err := json.Marshal(configuration)
	if err != nil {
		fmt.Println(err)
		return
	}
	file.WriteString(string(bdata))
}

func doCreate(c *cli.Context) {
	file, openErr := os.Open("config.json")
	var configuration = Configuration{""}
	if openErr != nil {
		bdata, err := json.Marshal(configuration)
		if err != nil {
			fmt.Println(err)
			return
		}
		file, createErr := os.Create("config.json")
		if createErr != nil {
			return
		}
		file.WriteString(string(bdata))
	} else {
		decoder := json.NewDecoder(file)
		configuration = Configuration{}
		decodeErr := decoder.Decode(&configuration)
		if decodeErr != nil {
			return
		}
	}

	if configuration.GithubToken == "" {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		githubToken := scanner.Text()
		configuration.GithubToken = githubToken
	}

	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: configuration.GithubToken},
	}

	client := github.NewClient(t.Client())
	repos, _, _ := client.Repositories.List("", nil)
	fmt.Printf("%#v", repos)
}
