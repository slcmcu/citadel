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
		cli.StringSliceFlag{"etcd-machines", &cli.StringSlice{"http://127.0.0.1:4001"}, "etcd hosts"},
	}

	app.Commands = []cli.Command{
		appCommand,
		containerCommand,
		hostCommand,
		hostsCommand,
		managementCommand,
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}
