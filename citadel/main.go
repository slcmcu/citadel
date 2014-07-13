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
	app.Author = "citadel team"

	app.Commands = []cli.Command{
		hostCommand,
		managementCommand,
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}
