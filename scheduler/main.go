package main

import (
	"log"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var (
	logger = logrus.New()
)

func masterMain(context *cli.Context) {

}

func main() {
	app := cli.NewApp()
	app.Name = "citadel scheduler"
	app.Version = "0.1"
	app.Author = "@crosbymichael"
	app.Email = "michael@crosbymichael.com"
	app.Commands = []cli.Command{
		{
			Name:        "slave",
			Description: "run as a slave in the cluster",
			Action:      slaveMain,
		},
		{
			Name:        "master",
			Description: "run as the master in the cluster",
			Action:      masterMain,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
