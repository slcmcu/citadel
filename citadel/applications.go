package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/codegangsta/cli"
)

var appCommand = cli.Command{
	Name:   "apps",
	Usage:  "view and mangage applications",
	Action: appAction,
}

func appAction(context *cli.Context) {
	apps, err := registry.FetchApplications()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	fmt.Fprint(w, "ID\tCONTAINERS\tPORTS\tCPUS\tMEMORY\n")

	for _, a := range apps {
		var (
			ports  = ""
			cpus   float64
			memory int
		)

		for _, p := range a.Ports {
			if p.Host != 0 {
				ports += fmt.Sprintf("%d ", p.Host)
			}
		}

		for _, c := range a.Containers {
			cpus += c.Cpus
			memory += c.Memory
		}

		fmt.Fprintf(w, "%s\t%d\t%s\t%f\t%d\n", a.ID, len(a.Containers), ports, cpus, memory)
	}

	w.Flush()
}
