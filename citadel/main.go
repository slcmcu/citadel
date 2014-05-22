package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var (
	logger   = logrus.New()
	machines = []string{os.Getenv("ETCD_MACHINES")}
)

func main() {
	app := cli.NewApp()
	app.Name = "citadel"
	app.Version = "0.1"
	app.Author = "@crosbymichael"
	app.Email = "michael@crosbymichael.com"
	app.Action = viewAction

	app.Flags = []cli.Flag{
		cli.StringFlag{"service", "/master", "service endpoint to hit"},
	}

	app.Commands = []cli.Command{
		newCommand,
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}
