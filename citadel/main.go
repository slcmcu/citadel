package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/citadel/citadel"
	"github.com/codegangsta/cli"
)

var (
	logger   = logrus.New()
	registry citadel.Registry
)

func main() {
	app := cli.NewApp()
	app.Usage = "mangage your docker containers across hosts"
	app.Name = "citadel"
	app.Version = "0.1"
	app.Author = "citadel team"

	app.Before = func(context *cli.Context) error {
		registry = citadel.NewRegistry(context.GlobalStringSlice("etcd-machines"))

		return nil
	}

	app.Flags = []cli.Flag{
		cli.StringSliceFlag{"etcd-machines", &cli.StringSlice{"http://127.0.0.1:4001"}, "etcd hosts"},
	}

	app.Commands = []cli.Command{
		appCommand,
		deleteCommand,
		startCommand,
		stopCommand,
		loadCommand,
		liveCommand,
		containerCommand,
		hostCommand,
		hostsCommand,
		managementCommand,
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}
