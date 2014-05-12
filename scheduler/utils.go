package main

import (
	"citadelapp.io/citadel"
	"citadelapp.io/citadel/repository"
	"citadelapp.io/citadel/utils"
	"github.com/Sirupsen/logrus"
	"github.com/apcera/nats"
	"github.com/codegangsta/cli"
	"github.com/samalba/dockerclient"
)

func getUUID() string {
	uuid, err := utils.GetUUID()
	if err != nil {
		logger.WithField("error", err).Fatal("unable to generate uuid")
	}
	return uuid
}

func getRepositoryAndConfig(context *cli.Context) (repository.Repository, *citadel.Config) {
	repo := repository.NewEtcdRepository(machines.Value(), false)
	conf, err := repo.FetchConfig()
	if err != nil {
		logger.WithField("error", err).Fatalln("fetch config", machines.Value())
	}
	return repo, conf
}

func getDocker(context *cli.Context) *dockerclient.DockerClient {
	docker, err := dockerclient.NewDockerClient(context.String("docker"))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error":  err,
			"docker": context.String("docker"),
		}).Fatal("unable to connect to a docker endpoint")
	}
	return docker
}

func getNats(conf *citadel.Config) *nats.EncodedConn {
	opts := nats.DefaultOptions
	opts.Servers = conf.Natsd

	nc, err := opts.Connect()
	if err != nil {
		logger.WithField("error", err).Fatal("natsd connect")
	}
	c, err := nats.NewEncodedConn(nc, "json")
	if err != nil {
		logger.WithField("error", err).Fatal("natsd encoded conn")
	}
	return c
}
