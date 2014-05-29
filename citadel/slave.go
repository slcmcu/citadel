package main

import (
	"net/http"
	"os"
	"strings"

	"citadelapp.io/citadel"
	"citadelapp.io/citadel/docker"
	"citadelapp.io/citadel/handler"
	"citadelapp.io/citadel/repository"
	"citadelapp.io/citadel/slave"
	"github.com/codegangsta/cli"
)

var slaveCommand = cli.Command{
	Name:        "slave",
	Action:      slaveAction,
	Description: "slave service that manages other services",
	Flags: []cli.Flag{
		cli.StringFlag{"namespace", "stackbrew", "docker namespace"},
		cli.StringFlag{"type", "docker", "service type"},
		cli.StringFlag{"addr", "127.0.0.1:3001", "address of the service"},
		cli.IntFlag{"memory", 0, "memory amount of the service"},
		cli.IntFlag{"cpus", 1, "number of cpus for the service"},
	},
}

func slaveAction(context *cli.Context) {
	var (
		service citadel.Service

		repo = repository.NewEtcdRepository(machines, false)
		data = &citadel.ServiceData{
			Name:   context.GlobalString("service"),
			Memory: context.Int("memory"),
			Cpus:   context.Int("cpus"),
			Addr:   context.String("addr"),
			Type:   context.String("type"),
		}
	)

	switch data.Type {
	case "docker":
		url := strings.Replace(os.Getenv("DOCKER_HOST"), "tcp", "http", -1)
		daemon, err := docker.New(context.String("namespace"), url, data)
		if err != nil {
			logger.Fatal(err)
		}

		service, err = slave.New(data, daemon, repo)
		if err != nil {
			logger.Fatal(err)
		}
	default:
		logger.Fatalf("unknown slave type %s", data.Type)
	}

	// FIXME: register with the cluster

	server := handler.New(service, logger)
	if err := http.ListenAndServe(data.Addr, server); err != nil {
		logger.Fatal(err)
	}
}
