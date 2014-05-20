package main

import "github.com/codegangsta/cli"

func viewAction(context *cli.Context) {
	name, command := parseRawCommand(context)
	logger.Println(name, command)
}
