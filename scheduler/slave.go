package main

import (
	"fmt"

	"citadelapp.io/citadel"
	"citadelapp.io/citadel/scheduler/slave"
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func slaveMain(context *cli.Context) {
	s, err := slave.New(uuid, logger, docker)
	if err != nil {
		return err
	}

	sub, err := nc.Subscribe(fmt.Sprintf("execute.%s", uuid), func(c *citadel.Container) {
		state, err := s.Execute(c)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"error": err,
				"uuid":  uuid,
			}).Error("executing container")
		}
	})
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

}
