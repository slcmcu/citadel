package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"citadelapp.io/citadel/repository"
	"github.com/codegangsta/cli"
)

var hostCommand = cli.Command{
	Name:   "hosts",
	Usage:  "view host information in the cluster",
	Action: hostAction,
}

func hostAction(context *cli.Context) {
	machines := strings.Split(context.GlobalString("etcd-machines"), ",")
	r := repository.New(machines, "citadel")

	hosts, err := r.FetchHosts()
	if err != nil {
		logger.WithField("error", err).Fatal("unable to fetch all hosts")
	}

	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	fmt.Fprint(w, "ID\tREGION\tADDR\tCPUS\tMEMORY\n")

	for _, h := range hosts {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\n", h.ID, h.Region, h.Addr, h.Cpus, h.Memory)
	}

	w.Flush()
}
