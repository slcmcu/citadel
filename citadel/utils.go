package main

import (
	"strings"

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
