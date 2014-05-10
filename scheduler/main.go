package main

import (
	"log"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var (
	logger   = logrus.New()
	machines = &cli.StringSlice{}
)

func main() {
	app := cli.NewApp()
	app.Name = "citadel scheduler"
	app.Version = "0.1"
	app.Author = "@crosbymichael"
	app.Email = "michael@crosbymichael.com"

	app.Flags = []cli.Flag{
		cli.StringSliceFlag{"etcd", machines, "etcd machines to connect to"},
	}
	app.Commands = []cli.Command{
		{
			Name:        "slave",
			Description: "run as a slave in the cluster",
			Action:      slaveMain,
			Flags: []cli.Flag{
				cli.StringFlag{"docker", "unix:///var/run/docker.sock", "docker endpoint"},
			},
		},
		{
			Name:        "master",
			Description: "run as the master in the cluster",
			Action:      masterMain,
			Flags: []cli.Flag{
				cli.StringFlag{"addr", "", "http address to reach the master"},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
