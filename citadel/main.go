package main

import (
	"os"

	"citadelapp.io/citadel"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var (
	logger  = logrus.New()
	service *citadel.Service
)

func newLocalService() *citadel.Service {
	nilCommand := &citadel.NilCommand{}

	return &citadel.Service{
		Name: "cli",
		Commands: map[string]citadel.Command{
			"list":  nilCommand,
			"start": nilCommand,
			"stop":  nilCommand,
		},
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "citadel"
	app.Version = "0.1"
	app.Author = "@crosbymichael"
	app.Email = "michael@crosbymichael.com"
	app.Action = viewAction

	app.Commands = []cli.Command{
		newCommand,
	}

	service = newLocalService()

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}
