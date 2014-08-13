package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "ghrepo - CLI tool for Github Reposity"
	app.Usage = ""
	app.Commands = Commands
	app.Run(os.Args)
}
