package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/citadel/citadel"
	"github.com/codegangsta/cli"
)

var hostsCommand = cli.Command{
	Name:   "hosts",
	Usage:  "display host information",
	Action: hostsAction,
}

func hostsAction(context *cli.Context) {
	registry = citadel.NewRegistry(context.GlobalStringSlice("etcd-machines"))

	hosts, err := registry.FetchHosts()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	fmt.Fprint(w, "ID\tLABELS\tADDR\tCPUS\tMEMORY\n")

	for _, h := range hosts {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\n", h.ID, strings.Join(h.Labels, ","), h.Addr, h.Cpus, h.Memory)
	}

	w.Flush()
}
