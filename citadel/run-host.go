package main

import (
	"net/http"

	"citadelapp.io/citadel"
	"citadelapp.io/citadel/repository"
	"citadelapp.io/citadel/utils"
	"github.com/codegangsta/cli"
)

var runHostCommand = cli.Command{
	Name:   "run-host",
	Usage:  "run the host and connect it to the cluster",
	Action: runHostAction,
	Flags: []cli.Flag{
		cli.StringFlag{"region", "", "region where the host is running"},
		cli.StringFlag{"addr", "", "external ip address for the host"},
		cli.IntFlag{"cpus", -1, "number of cpus available to the host"},
		cli.IntFlag{"memory", -1, "number of mb of memory available to the host"},
	},
}

func runHostAction(context *cli.Context) {
	var (
		cpus   = context.Int("cpus")
		memory = context.Int("memory")
		addr   = context.String("addr")
		region = context.String("region")
	)

	id, err := utils.GetMachineID()
	if err != nil {
		logger.WithField("error", err).Fatal("unable to read machine id")
	}

	switch {
	case cpus < 1:
		logger.Fatal("cpus must have a value")
	case memory < 1:
		logger.Fatal("memory must have a value")
	case addr == "":
		logger.Fatal("addr must have a value")
	case region == "":
		logger.Fatal("region must have a value")
	}

	r, err := repository.New(context.GlobalString("repository"))
	if err != nil {
		logger.WithField("error", err).Fatal("unable to connect to repository")
	}
	defer r.Close()

	host := &citadel.Host{
		ID:     id,
		Memory: memory,
		Cpus:   cpus,
		Addr:   addr,
		Region: region,
	}

	if err := r.SaveHost(host); err != nil {
		logger.WithField("error", err).Fatal("unable to save host")
	}
	defer r.DeleteHost(id)

	if err := http.ListenAndServe(addr, nil); err != nil {
	}
}
