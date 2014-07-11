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

	app.Flags = []cli.Flag{
		cli.StringFlag{"etcd-machines", "http://127.0.0.1:4001", "comma separated list of etcd hosts"},
	}

	app.Commands = []cli.Command{
		hostCommand,
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}
