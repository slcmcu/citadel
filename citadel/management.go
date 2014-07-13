package main

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/citadel/citadel"
	"github.com/codegangsta/cli"
	"github.com/coreos/go-etcd/etcd"
	"github.com/gorilla/mux"
)

var (
	registry *etcd.Client
)

var managementCommand = cli.Command{
	Name:   "management",
	Usage:  "run the management ui for the cluster",
	Action: managementAction,
	Flags: []cli.Flag{
		cli.StringFlag{"assets", "management", "assests for the web ui"},
		cli.StringFlag{"addr", ":3002", "address for the web ui to listen on"},
	},
}

func managementAction(context *cli.Context) {
	var (
		// we cannot use the mux router with static files, there is a bug where it does
		// not handle the full directory
		globalMux = http.NewServeMux()
		apiRouter = mux.NewRouter()
	)

	registry = etcd.NewClient([]string{"http://127.0.0.1:4001"})

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
	hosts, err := fetchHosts()
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

func fetchHosts() ([]*citadel.Host, error) {
	hosts := []*citadel.Host{}

	resp, err := registry.Get("/citadel/hosts", true, true)
	if err != nil {
		return nil, err
	}

	for _, n := range resp.Node.Nodes {
		var host *citadel.Host
		if err := json.Unmarshal([]byte(n.Value), &host); err != nil {
			return nil, err
		}

		hosts = append(hosts, host)
	}

	return hosts, nil
}

func fetchContainers() ([]*citadel.Container, error) {
	hosts, err := fetchHosts()
	if err != nil {
		return nil, err
	}

	out := []*citadel.Container{}
	for _, h := range hosts {
		resp, err := registry.Get(filepath.Join("/citadel", h.ID, "containers"), true, true)
		if err != nil {
			return nil, err
		}

		for _, node := range resp.Node.Nodes {
			var container *citadel.Container
			if err := json.Unmarshal([]byte(node.Value), &container); err != nil {
				return nil, err
			}

			out = append(out, container)
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
