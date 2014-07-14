package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/citadel/citadel"
	"github.com/codegangsta/cli"
)

var containerCommand = cli.Command{
	Name:   "containers",
	Usage:  "view containers",
	Action: containerAction,
}

func containerAction(context *cli.Context) {
	registry = citadel.NewRegistry(context.GlobalStringSlice("etcd-machines"))

	hosts, err := registry.FetchHosts()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	fmt.Fprint(w, "NAME\tHOST\tAPP\tTYPE\tSTATUS\n")

	for _, h := range hosts {
		containers, err := registry.FetchContainers(h.ID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		for _, c := range containers {
			// lets not show the group containers here for a better UI
			if c.Config.Type != citadel.Group {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", c.Name, c.HostID, c.ApplicationID, c.Config.Type, c.State.Status)
			}
		}
	}

	w.Flush()
}
