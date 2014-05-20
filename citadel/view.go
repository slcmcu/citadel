package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"citadelapp.io/citadel/repository"
	"github.com/codegangsta/cli"
)

func viewAction(context *cli.Context) {
	var (
		repo    = repository.NewEtcdRepository(machines, false)
		name, _ = parseRawCommand(context)
		w       = tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	)

	if name == "" {
		name = "/"
	}

	services, err := repo.FetchServices(name)
	if err != nil {
		logger.Fatal(err)
	}

	fmt.Fprint(w, "NAME\tTYPE\tADDR\tCPUS\tMEMORY\n")
	for _, s := range services {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%f\n", s.Name, s.Type, s.Addr, s.Cpus, s.Memory)
	}

	if err := w.Flush(); err != nil {
		logger.Fatal(err)
	}
}
