package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"citadelapp.io/citadel"
	"citadelapp.io/citadel/repository"
	"citadelapp.io/citadel/slave"
	"citadelapp.io/citadel/utils"
	"github.com/Sirupsen/logrus"
	"github.com/apcera/nats"
	"github.com/codegangsta/cli"
	"github.com/samalba/dockerclient"
)

func register(s *slave.Slave, ttl int, repo repository.Repository) error {
	if err := repo.RegisterSlave(s.ID, s, ttl); err != nil {
		return err
	}
	go heartbeat(s.ID, ttl, repo)
	return nil
}

func heartbeat(uuid string, ttl int, repo repository.Repository) {
	for _ = range time.Tick(time.Duration(ttl-2) * time.Second) {
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

func execute(s *slave.Slave, c *citadel.Container, repo repository.Repository, nc *nats.EncodedConn) {
	if err := s.Execute(c); err != nil {
		logger.WithFields(logrus.Fields{
			"error": err,
			"uuid":  s.ID,
		}).Error("executing container")
		return
	}

	if err := repo.SaveContainer(s.ID, c); err != nil {
		logger.WithFields(logrus.Fields{
			"error":     err,
			"uuid":      s.ID,
			"container": c.ID,
		}).Error("saving container")
	}
	nc.Publish("containers.start", c)
}

func eventHandler(event *dockerclient.Event, args ...interface{}) {
	var (
		s    = args[0].(*slave.Slave)
		repo = args[1].(repository.Repository)
	)

	switch event.Status {
	case "die":
		if err := s.RemoveContainer(event.Id); err != nil {
			logger.WithFields(logrus.Fields{
				"error": err,
				"event": event.Status,
				"id":    event.Id,
			}).Error("cannot remove container from slave")
		}
		if err := repo.RemoveContainer(s.ID, event.Id); err != nil {
			logger.WithFields(logrus.Fields{
				"error": err,
				"event": event.Status,
				"id":    event.Id,
			}).Error("cannot remove container")
		}
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

	s, err := slave.New(uuid, docker)
	if err != nil {
		logger.WithField("error", err).Fatal("unable to initialize slave")
	}

	if err := register(s, conf.TTL, repo); err != nil {
		logger.WithField("error", err).Fatal("register slave")
	}
	nc.Publish("slaves.joining", uuid)

	docker.StartMonitorEvents(eventHandler, s, repo)

	execSub, err := nc.Subscribe(fmt.Sprintf("execute.%s", uuid), func(msg *nats.Msg) {
		var c *citadel.Container
		if err := json.Unmarshal(msg.Data, &c); err != nil {
			logger.WithField("error", err).Error("unmarshal container from message")
			return
		}
		logger.WithField("image", c.Image).Info("executing")
		execute(s, c, repo, nc)
		if err := nc.Publish(msg.Reply, c); err != nil {
			logger.WithField("error", err).Error("sending response")
		}
	})
	if err != nil {
		logger.WithField("error", err).Fatal("subscribe")
	}
	defer execSub.Unsubscribe()

	pullSub, err := nc.Subscribe("slaves.pull", func(image string) {
		logger.WithField("image", image).Info("pulling")
		if err := s.PullImage(image); err != nil {
			logger.WithField("error", err).Error("pull image")
		}
	})
	if err != nil {
		logger.WithField("error", err).Fatal("subscribe")
	}
	defer pullSub.Unsubscribe()

	for s := range sig {
		nc.Publish("slaves.leaving", uuid)
		repo.RemoveSlave(uuid)

		logger.WithField("signal", s.String()).Info("exiting")
		return
	}
}
