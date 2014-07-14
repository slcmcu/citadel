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
		Flags: []cli.Flag{
			cli.StringSliceFlag{"hosts", &cli.StringSlice{}, "hosts to load the app on"},
		},
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
		Flags: []cli.Flag{
			cli.StringSliceFlag{"hosts", &cli.StringSlice{}, "hosts to load the app on"},
		},
	}

	stopCommand = cli.Command{
		Name:   "stop",
		Usage:  "stop an application",
		Action: stopAction,
		Flags: []cli.Flag{
			cli.StringSliceFlag{"hosts", &cli.StringSlice{}, "hosts to load the app on"},
		},
	}
)

func deleteAction(context *cli.Context) {
	app, err := loadApp(context)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, hn := range context.StringSlice("hosts") {
		if err := runTrans(app, hn, "DELETE", "app"); err != nil {
			logger.WithField("error", err).Errorf("deleting app from %s", hn)
		}
	}

	if err := registry.DeleteApplication(app.ID); err != nil {
		logger.WithField("error", err).Fatal("delete application from registry")
	}
}

func loadAction(context *cli.Context) {
	app, err := loadApp(context)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := registry.SaveApplication(app); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, hn := range context.StringSlice("hosts") {
		if err := runTrans(app, hn, "POST", "app"); err != nil {
			logger.WithField("error", err).Errorf("cannot load app on %s", hn)
		}
	}
}

func startAction(context *cli.Context) {
	app, err := loadApp(context)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, hn := range context.StringSlice("hosts") {
		if err := runTrans(app, hn, "POST", "run"); err != nil {
			logger.WithField("error", err).Errorf("cannot run container on %s", hn)
		}
	}
}

func stopAction(context *cli.Context) {
	app, err := loadApp(context)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, hn := range context.StringSlice("hosts") {
		if err := runTrans(app, hn, "POST", "stop"); err != nil {
			logger.WithField("error", err).Errorf("cannot stop container on %s", hn)
		}
	}
}

func runTrans(app *citadel.Application, hostName, method, endpoint string) error {
	var (
		c = http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	)

	host, err := registry.FetchHost(hostName)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://%s/%s/%s", host.Addr, endpoint, app.ID)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	var tran *citadel.Transaction
	if err := json.NewDecoder(resp.Body).Decode(&tran); err != nil {
		return err
	}

	compiled, err := template.New("transaction").Parse(tranactionTemplate)
	if err != nil {
		return err
	}

	if err := compiled.Execute(os.Stdout, tran); err != nil {
		return err
	}

	return nil
}

func loadApp(context *cli.Context) (*citadel.Application, error) {
	f, err := os.Open(context.Args().First())
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
