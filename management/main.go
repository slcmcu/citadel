package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"citadelapp.io/citadel"
	"citadelapp.io/citadel/metrics"
	"citadelapp.io/citadel/repository"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

var (
	machines string
	assets   string

	store  metrics.Store
	repo   repository.Repository
	logger = logrus.New()
)

type Node struct {
	ID     string  `json:"id,omitempty"`
	Addr   string  `json:"addr,omitempty"`
	Type   string  `json:"type,omitempty"`
	Cpus   int     `json:"cpus,omitempty"`
	Memory float64 `json:"memory,omitempty"`
}

func init() {
	flag.StringVar(&machines, "machines", "127.0.0.1:4001", "Comma separated list of etcd machines")
	flag.StringVar(&assets, "assets", "management", "Path the the http assets")
	flag.Parse()
}

func getContainers(w http.ResponseWriter, r *http.Request) {
	slaves, err := repo.FetchSlaves()
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	containers := []*citadel.Container{}
	for _, s := range slaves {
		cs, err := repo.FetchContainers(s.ID)
		if err != nil {
			logger.WithField("error", err).Error("fetch containers for slave")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, c := range cs {
			containers = append(containers, c)
		}
	}

	if err := json.NewEncoder(w).Encode(containers); err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getNodes(w http.ResponseWriter, r *http.Request) {
	master, err := repo.FetchMaster()
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slaves, err := repo.FetchSlaves()
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nodes := []*Node{
		{Type: "master", ID: master.ID, Addr: master.Addr},
	}

	for _, s := range slaves {
		nodes = append(nodes, &Node{
			Type:   "slave",
			ID:     s.ID,
			Addr:   s.IP,
			Cpus:   s.Cpus,
			Memory: s.Memory,
		})
	}

	if err := json.NewEncoder(w).Encode(nodes); err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getNodeMetrics(w http.ResponseWriter, r *http.Request) {
	var (
		params = mux.Vars(r)
		table  = fmt.Sprintf("metrics.host.%s", params["name"])
	)

	data, err := store.Fetch(fmt.Sprintf("select * from %s group by time(5m) where time > now() -12h limit 200", table))
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	var (
		err error
		// we cannot use the mux router with static files, there is a bug where it does
		// not handle the full directory
		globalMux    = http.NewServeMux()
		apiRouter    = mux.NewRouter()
		etcdMachines = strings.Split(machines, ",")
	)

	repo = repository.NewEtcdRepository(etcdMachines, false)
	conf, err := repo.FetchConfig()
	if err != nil {
		logger.WithField("error", err).Fatal("fetch config")
	}
	if store, err = metrics.NewStore(conf); err != nil {
		logger.WithField("error", err).Fatal("new metrics store")
	}

	apiRouter.HandleFunc("/api/containers", getContainers)
	apiRouter.HandleFunc("/api/nodes", getNodes)
	apiRouter.HandleFunc("/api/nodes/{name:.*}/metrics", getNodeMetrics)

	globalMux.Handle("/api/", apiRouter)
	globalMux.Handle("/", http.FileServer(http.Dir(assets)))

	if err := http.ListenAndServe(":3000", globalMux); err != nil {
		logger.WithField("error", err).Fatal("serve management ui")
	}
}
