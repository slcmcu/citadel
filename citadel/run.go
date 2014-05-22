package main

import (
	"citadelapp.io/citadel"
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var runCommand = cli.Command{
	Name:   "run",
	Action: runAction,
	Flags: []cli.Flag{
		cli.StringFlag{"type", "", "service type"},
		cli.IntFlag{"memory", 0, "memory amount of the service"},
		cli.IntFlag{"cpus", 1, "number of cpus for the service"},
		cli.IntFlag{"instances", 1, "number of instances to run"},
	},
}

func runAction(context *cli.Context) {
	var (
		name      = context.Args().First()
		memory    = context.Int("memory")
		cpus      = context.Int("cpus")
		tpe       = context.String("type")
		instances = context.Int("instances")
	)

	parent, err := newService(context)
	if err != nil {
		logger.Fatal(err)
	}

	switch {
	case name == "":
		logger.Fatal("name connot be empty")
	case memory == 0:
		logger.Fatal("memory cannot be 0")
	case cpus == 0:
		logger.Fatal("cpus cannot be 0")
	case tpe == "":
		logger.Fatal("type cannot be empty")
	}

	data := &citadel.ServiceData{
		Name:   name,
		Cpus:   cpus,
		Memory: memory,
		Type:   tpe,
	}

	task := &citadel.Task{
		Name:      name,
		Service:   data,
		Instances: instances,
	}

	logger.WithFields(logrus.Fields{
		"instaces": instances,
		"name":     name,
		"type":     tpe,
	}).Info("running task")

	rundata, err := parent.Run(task)
	if err != nil {
		logger.Fatal(err)
	}

	// FIXME: print out the rundata in a nice report
	logger.Println(rundata)
}
