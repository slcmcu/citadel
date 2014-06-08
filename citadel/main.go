package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var (
	logger = logrus.New()
)

func main() {
	app := cli.NewApp()
	app.Name = "citadel"
	app.Version = "0.1"
	app.Author = "@crosbymichael"
	app.Email = "michael@crosbymichael.com"
	//	app.Action = viewAction

	app.Flags = []cli.Flag{
		cli.StringFlag{"repository", "127.0.0.1:28015", "repository to connect to"},
	}

	app.Commands = []cli.Command{
		hostCommand,
		runHostCommand,
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}
