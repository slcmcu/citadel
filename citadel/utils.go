package main

import (
	"strings"

	"citadelapp.io/citadel"
	"citadelapp.io/citadel/handler"
	"citadelapp.io/citadel/repository"

	"github.com/codegangsta/cli"
)

func parseRawCommand(context *cli.Context) (string, string) {
	parts := strings.SplitN(context.Args().First(), ":", 2)
	switch len(parts) {
	case 1:
		return parts[0], ""
	case 2:
		return parts[0], parts[1]
	default:
		logger.Fatalf("invalid command format %s", context.Args().First())
	}
	return "", ""
}

func newService(context *cli.Context) (citadel.Service, error) {
	repo := repository.NewEtcdRepository(machines, false)

	service, err := repo.FetchService(context.String("service"))
	if err != nil {
		return nil, err
	}

	return handler.NewClient(service), nil
}
