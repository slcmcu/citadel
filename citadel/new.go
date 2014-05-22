package main

import (
	"path/filepath"

	"citadelapp.io/citadel"
	"github.com/codegangsta/cli"
)

var newCommand = cli.Command{
	Name:   "new",
	Action: newAction,
	Flags: []cli.Flag{
		cli.StringFlag{"type", "", "service type"},
		cli.StringFlag{"addr", "", "address of the service"},
		cli.IntFlag{"memory", 0, "memory amount of the service"},
		cli.IntFlag{"cpus", 1, "number of cpus for the service"},
	},
}

func newAction(context *cli.Context) {
	var (
		fullPath = context.Args().First()
		memory   = context.Int("memory")
		cpus     = context.Int("cpus")
		addr     = context.String("addr")
		tpe      = context.String("type")
	)

	service, err := newService(context)
	if err != nil {
		logger.Fatal(err)
	}

	switch {
	case fullPath == "":
		logger.Fatal("name connot be empty")
	case addr == "":
		logger.Fatal("addr cannot be empty")
	case memory == 0:
		logger.Fatal("memory cannot be 0")
	}

	_, name := filepath.Split(fullPath)
	data := &citadel.ServiceData{
		Name:   name,
		Cpus:   cpus,
		Addr:   addr,
		Memory: memory,
		Type:   tpe,
	}

	rundata, err := service.Run(&citadel.Task{
		Name:      name,
		Service:   data,
		Instances: 1,
	})

	if err != nil {
		logger.Fatal(err)
	}

	// FIXME: print out the rundata in a nice report
	logger.Println(rundata)
}
