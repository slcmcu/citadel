package main

import (
	"encoding/json"
	"net/http"

	"github.com/citadel/citadel"
	"github.com/codegangsta/cli"
	"github.com/gorilla/mux"
)

var (
	registry citadel.Registry
)

var managementCommand = cli.Command{
	Name:   "management",
	Usage:  "run the management ui for the cluster",
	Action: managementAction,
	Flags: []cli.Flag{
		cli.StringFlag{"assets", "assets", "assests for the web ui"},
		cli.StringFlag{"addr", ":3002", "address for the web ui to listen on"},
		cli.StringSliceFlag{"etcd-machines", &cli.StringSlice{"http://127.0.0.1:4001"}, "etcd hosts"},
	},
}

func managementAction(context *cli.Context) {
	var (
		// we cannot use the mux router with static files, there is a bug where it does
		// not handle the full directory
		globalMux = http.NewServeMux()
		apiRouter = mux.NewRouter()
	)

	registry = citadel.NewRegistry(context.StringSlice("etcd-machines"))

	apiRouter.HandleFunc("/api/hosts", getHosts)
	apiRouter.HandleFunc("/api/containers", getContainers).Methods("GET")

	globalMux.Handle("/api/", apiRouter)
	globalMux.Handle("/", http.FileServer(http.Dir(context.String("assets"))))

	if err := http.ListenAndServe(context.String("addr"), globalMux); err != nil {
		logger.WithField("error", err).Fatal("serve management ui")
	}
}

func httpError(w http.ResponseWriter, err error) {
	logger.Error(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func getHosts(w http.ResponseWriter, r *http.Request) {
	hosts, err := registry.FetchHosts()
	if err != nil {
		httpError(w, err)
		return
	}

	marshal(w, hosts)
}

func getContainers(w http.ResponseWriter, r *http.Request) {
	containers, err := fetchContainers()
	if err != nil {
		httpError(w, err)
		return
	}

	marshal(w, containers)
}

func fetchContainers() ([]interface{}, error) {
	hosts, err := registry.FetchHosts()
	if err != nil {
		return nil, err
	}

	out := []interface{}{}

	for _, h := range hosts {
		containers, err := registry.FetchContainers(h)
		if err != nil {
			return nil, err
		}

		for _, c := range containers {
			out = append(out, struct {
				*citadel.Container
				Host string `json:"host,omitempty"`
			}{
				Container: c,
				Host:      h.ID,
			})
		}
	}

	return out, nil
}

func marshal(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(v); err != nil {
		logger.WithField("error", err).Error("encode json")
	}
}
