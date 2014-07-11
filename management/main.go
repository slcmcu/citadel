package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/citadel/citadel"
	"github.com/citadel/citadel/repository"
	"github.com/gorilla/mux"
)

var (
	assets       string
	etcdMachines string
	listenAddr   string
	repo         *repository.Repository
	logger       = logrus.New()
)

func init() {
	flag.StringVar(&assets, "assets", "management", "path the the http assets")
	flag.StringVar(&etcdMachines, "etcd-machines", "http://127.0.0.1:4001", "comma separated list of etcd machines")
	flag.StringVar(&listenAddr, "listenAddr", ":3002", "management listen address")

	flag.Parse()
}

func httpError(w http.ResponseWriter, err error) {
	logger.Error(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func getHosts(w http.ResponseWriter, r *http.Request) {
	hosts, err := repo.FetchHosts()
	if err != nil {
		httpError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(hosts); err != nil {
		httpError(w, err)
		return
	}
}

func getContainers(w http.ResponseWriter, r *http.Request) {
	containers, err := repo.FetchContainers()
	if err != nil {
		httpError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(containers); err != nil {
		httpError(w, err)
		return
	}
}

func postTasks(w http.ResponseWriter, r *http.Request) {
	var t *citadel.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		httpError(w, err)
		return
	}

	if err := repo.AddTask(t); err != nil {
		httpError(w, err)
		return
	}
}

func main() {
	var (
		err error
		// we cannot use the mux router with static files, there is a bug where it does
		// not handle the full directory
		globalMux = http.NewServeMux()
		apiRouter = mux.NewRouter()
	)

	machines := strings.Split(etcdMachines, ",")
	if err != nil {
		logger.WithField("error", err).Fatal("unable to parse etcd machines")
	}

	repo = repository.New(machines, "citadel")
	apiRouter.HandleFunc("/api/hosts", getHosts)
	apiRouter.HandleFunc("/api/containers", getContainers).Methods("GET")
	apiRouter.HandleFunc("/api/tasks", postTasks).Methods("POST")

	globalMux.Handle("/api/", apiRouter)
	globalMux.Handle("/", http.FileServer(http.Dir(assets)))

	if err := http.ListenAndServe(listenAddr, globalMux); err != nil {
		logger.WithField("error", err).Fatal("serve management ui")
	}
}
