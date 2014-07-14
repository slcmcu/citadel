package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/citadel/citadel"
	"github.com/codegangsta/cli"
)

const tranactionTemplate = `
TRANSACTION: {{.ID}}
TYPE: {{.Type}}
CONTAINERS:
{{range $container := .Containers}}
    ID: {{$container.ID}} 
    NAME: {{$container.Name}}
{{end}}
`

var (
	deleteCommand = cli.Command{
		Name:   "delete",
		Usage:  "delete an application and all containers",
		Action: deleteAction,
	}

	loadCommand = cli.Command{
		Name:   "load",
		Usage:  "load an application",
		Action: loadAction,
	}

	startCommand = cli.Command{
		Name:   "start",
		Usage:  "start an application",
		Action: startAction,
	}

	stopCommand = cli.Command{
		Name:   "stop",
		Usage:  "stop an application",
		Action: stopAction,
	}
)

func deleteAction(context *cli.Context) {

}

func loadAction(context *cli.Context) {

}

func startAction(context *cli.Context) {
	runTrans(context, "run")
}

func stopAction(context *cli.Context) {
	runTrans(context, "stop")
}

func runTrans(context *cli.Context, endpoint string) {
	var (
		appFile  = context.Args().Get(0)
		hostName = context.Args().Get(1)
		c        = http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	)

	app, err := loadApp(appFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	host, err := registry.FetchHost(hostName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/%s/%s", host.Addr, endpoint, app.ID), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	resp, err := c.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var tran *citadel.Transaction
	if err := json.NewDecoder(resp.Body).Decode(&tran); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	compiled, err := template.New("transaction").Parse(tranactionTemplate)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := compiled.Execute(os.Stdout, tran); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func loadApp(p string) (*citadel.Application, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var app *citadel.Application
	if err := json.NewDecoder(f).Decode(&app); err != nil {
		return nil, err
	}

	return app, nil
}
