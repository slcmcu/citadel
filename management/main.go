package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"citadelapp.io/citadel/metrics"
	"citadelapp.io/citadel/repository"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

var (
	machines string
	assets   string

	store metrics.Store
	repo  repository.Repository
	log   = logrus.New()
)

func init() {
	flag.StringVar(&machines, "machines", "127.0.0.1:4001", "Comma separated list of etcd machines")
	flag.StringVar(&assets, "assets", "management", "Path the the http assets")
	flag.Parse()
}

func getContainers(w http.ResponseWriter, r *http.Request) {
	data, err := repo.FetchContainerGroup()
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Fatal(err)
	}
}

func getHosts(w http.ResponseWriter, r *http.Request) {
	data, err := repo.FetchHosts()
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Fatal(err)
	}
}

func getHostMetrics(w http.ResponseWriter, r *http.Request) {
	var (
		params = mux.Vars(r)
		table  = fmt.Sprintf("metrics.host.%s", params["name"])
	)

	data, err := store.Fetch(fmt.Sprintf("select * from %s group by time(5m) where time > now() -12h limit 200", table))
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var (
		err          error
		globalMux    = http.NewServeMux()
		apiRouter    = mux.NewRouter()
		etcdMachines = strings.Split(machines, ",")
	)

	repo = repository.NewEtcdRepository(etcdMachines)
	conf, err := repo.FetchConfig()
	if err != nil {
		log.Fatal(err)
	}
	if store, err = metrics.NewStore(conf); err != nil {
		log.Fatal(err)
	}

	apiRouter.HandleFunc("/api/containers", getContainers)
	apiRouter.HandleFunc("/api/hosts", getHosts)
	apiRouter.HandleFunc("/api/hosts/{name:.*}/metrics", getHostMetrics)

	globalMux.Handle("/api/", apiRouter)
	globalMux.Handle("/", http.FileServer(http.Dir(assets)))

	if err := http.ListenAndServe(":3000", globalMux); err != nil {
		log.Fatal(err)
	}
}
