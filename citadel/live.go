package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/codegangsta/cli"
)

var liveCommand = cli.Command{
	Name:   "live",
	Usage:  "live keeps group containers going",
	Action: liveAction,
	Flags: []cli.Flag{
		cli.StringSliceFlag{"volumes", &cli.StringSlice{}, "modify volume perms"},
	},
}

func liveAction(context *cli.Context) {
	for _, v := range context.StringSlice("volumes") {
		parts := strings.Split(v, ":")

		uid, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		gid, err := strconv.Atoi(parts[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := os.Chown(parts[0], uid, gid); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	for _ = range sigChan {
		os.Exit(0)
	}
}
