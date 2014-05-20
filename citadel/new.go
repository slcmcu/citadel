package main

import "github.com/codegangsta/cli"

var newCommand = cli.Command{
	Name:   "new",
	Action: newAction,
	Flags: []cli.Flag{
		cli.StringFlag{"addr", "", "address of the service"},
		cli.Float64Flag{"memory", 0, "memory amount of the service"},
		cli.IntFlag{"cpus", 1, "number of cpus for the service"},
	},
}

func newAction(context *cli.Context) {
	var (
		name   = context.Args().First()
		memory = context.Float64("memory")
		cpus   = context.Int("cpus")
		addr   = context.String("addr")
	)

	if name == "" {
		logger.Fatal("name connot be empty")
	}

	logger.Println(name, memory, cpus, addr)
}
