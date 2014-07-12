package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/citadel/citadel"
	"github.com/codegangsta/cli"
	"github.com/samalba/dockerclient"
)

var hostCommand = cli.Command{
	Name:   "host",
	Usage:  "run the host and connect it to the cluster",
	Action: hostAction,
	Flags: []cli.Flag{
		cli.StringFlag{"config", "host.toml", "config for the host"},
	},
}

func hostAction(context *cli.Context) {
	config, err := loadConfig(context.String("config"))
	if err != nil {
		logger.WithField("error", err).Fatal("load config")
	}

	host, err := citadel.NewHost(config.ID, config.Cpus, config.Memory, config.Labels, getClient(config), logger)
	if err != nil {
		logger.WithField("error", err).Fatal("create host")
	}

	server := citadel.NewServer(host)
	go waitForInterrupt()

	if err := http.ListenAndServe(config.Addr, server); err != nil {
		logger.WithField("error", err).Fatal("listen and serve")
	}
}

func waitForInterrupt() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	for _ = range sigChan {
		os.Exit(0)
	}
}

func getClient(config *Config) *dockerclient.DockerClient {
	client, err := dockerclient.NewDockerClient(config.Docker)
	if err != nil {
		logger.WithField("error", err).Fatal("unable to connect to docker")
	}

	return client
}
