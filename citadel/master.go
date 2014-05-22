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
	Name:   "master",
	Action: masterAction,
	Flags: []cli.Flag{
		cli.StringFlag{"type", "", "service type"},
		cli.StringFlag{"addr", "127.0.0.1:3000", "address of the service"},
		cli.IntFlag{"memory", 0, "memory amount of the service"},
		cli.IntFlag{"cpus", 1, "number of cpus for the service"},
	},
}

func masterAction(context *cli.Context) {
	var (
		data = &citadel.ServiceData{
			Name:   context.String("service"),
			Memory: context.Int("memory"),
			Cpus:   context.Int("cpus"),
			Addr:   context.String("addr"),
		}

		server = handler.New(master.New(data, repository.NewEtcdRepository(machines, false)), logger)
	)

	// FIXME: register master with the cluster

	if err := http.ListenAndServe(data.Addr, server); err != nil {
		logger.Fatal(err)
	}
}
