package main

import (
	"net/http"

	"citadelapp.io/citadel"
	"citadelapp.io/citadel/handler"
	"citadelapp.io/citadel/master"
	"citadelapp.io/citadel/repository"
	"github.com/codegangsta/cli"
)

var masterCommand = cli.Command{
	Name:        "master",
	Action:      masterAction,
	Description: "master service to interace with slaves in the cluseter",
	Flags: []cli.Flag{
		cli.StringFlag{"addr", "127.0.0.1:3000", "address of the service"},
		cli.IntFlag{"memory", 0, "memory amount of the service"},
		cli.IntFlag{"cpus", 1, "number of cpus for the service"},
	},
}

func masterAction(context *cli.Context) {
	var (
		data = &citadel.ServiceData{
			Name:   context.GlobalString("service"),
			Memory: context.Int("memory"),
			Cpus:   context.Int("cpus"),
			Addr:   context.String("addr"),
			Type:   "master",
		}

		server = handler.New(master.New(data, repository.NewEtcdRepository(machines, false)), logger)
	)

	// FIXME: register master with the cluster

	if err := http.ListenAndServe(data.Addr, server); err != nil {
		logger.Fatal(err)
	}
}
