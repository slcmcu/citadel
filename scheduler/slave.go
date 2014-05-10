package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"citadelapp.io/citadel"
	"citadelapp.io/citadel/repository"
	"citadelapp.io/citadel/scheduler/slave"
	"citadelapp.io/citadel/utils"
	"github.com/Sirupsen/logrus"
	"github.com/apcera/nats"
	"github.com/codegangsta/cli"
	"github.com/samalba/dockerclient"
)

func register(s *slave.Slave, ttl int, repo repository.Repository) error {
	if err := repo.RegisterSlave(s.ID, &s.Resource, ttl); err != nil {
		return err
	}
	go heartbeat(s.ID, ttl, repo)
	return nil
}

func heartbeat(uuid string, ttl int, repo repository.Repository) {
	for _ = range time.Tick(time.Duration(ttl) * time.Second) {
		for i := 0; i < 5; i++ {
			err := repo.UpdateSlave(uuid, ttl)
			if err == nil {
				continue
			}
			logger.WithField("error", err).Error("updating ttl")
			time.Sleep(500 * time.Second)
		}
	}
}

func getUUID() string {
	uuid, err := utils.GetUUID()
	if err != nil {
		logger.WithField("error", err).Fatal("unable to generate uuid")
	}
	return uuid
}

func getRepositoryAndConfig(context *cli.Context) (repository.Repository, *citadel.Config) {
	repo := repository.NewEtcdRepository(context.StringSlice("etcd"))
	conf, err := repo.FetchConfig()
	if err != nil {
		logger.WithField("error", err).Fatal("fetch config")
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

func execute(s *slave.Slave, c *citadel.Container, repo repository.Repository) {
	state, err := s.Execute(c)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err,
			"uuid":  s.ID,
		}).Error("executing container")
		return
	}

	if err := repo.SaveState(s.ID, state); err != nil {
		logger.WithFields(logrus.Fields{
			"error":     err,
			"uuid":      s.ID,
			"container": state.ID,
		}).Error("saving contaienr state")
	}
}

func slaveMain(context *cli.Context) {
	var (
		uuid       = getUUID()
		docker     = getDocker(context)
		sig        = make(chan os.Signal, 1)
		repo, conf = getRepositoryAndConfig(context)
		nc         = getNats(conf)
	)
	defer nc.Close()
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	s, err := slave.New(uuid, logger, docker)
	if err != nil {
		logger.WithField("error", err).Fatal("unable to initialize slave")
	}

	if err := register(s, conf.SlaveTTL, repo); err != nil {
		logger.WithField("error", err).Fatal("register slave")
	}
	sub, err := nc.Subscribe(fmt.Sprintf("execute.%s", uuid), func(c *citadel.Container) {
		execute(s, c, repo)
	})
	if err != nil {
		logger.WithField("error", err).Fatal("subscribe")
	}
	defer sub.Unsubscribe()

	for s := range sig {
		logger.WithField("signal", s.String()).Info("exiting")
		return
	}
}
