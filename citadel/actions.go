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
		Usage:  "load an application but does not load on hosts",
		Action: loadAction,
		Flags: []cli.Flag{
			cli.StringSliceFlag{"hosts", &cli.StringSlice{}, "hosts to load the app on"},
		},
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
	var (
		appFile  = context.Args().Get(0)
		hostName = context.Args().Get(1)
	)

	app := runTrans(appFile, hostName, "DELETE", "app")

	if err := registry.DeleteApplication(app.ID); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func loadAction(context *cli.Context) {
	var (
		appFile = context.Args().Get(0)
	)

	app, err := loadApp(appFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := registry.SaveApplication(app); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, hn := range context.StringSlice("hosts") {
		runTrans(appFile, hn, "POST", "app")
	}
}

func startAction(context *cli.Context) {
	var (
		appFile  = context.Args().Get(0)
		hostName = context.Args().Get(1)
	)

	runTrans(appFile, hostName, "POST", "run")
}

func stopAction(context *cli.Context) {
	var (
		appFile  = context.Args().Get(0)
		hostName = context.Args().Get(1)
	)

	runTrans(appFile, hostName, "POST", "stop")
}

func runTrans(appFile, hostName, method, endpoint string) *citadel.Application {
	var (
		c = http.Client{
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

	url := fmt.Sprintf("https://%s/%s/%s", host.Addr, endpoint, app.ID)

	req, err := http.NewRequest(method, url, nil)
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

	return app
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
