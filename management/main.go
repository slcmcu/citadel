package main

import (
	"encoding/json"
	"flag"
	"net/http"

	"citadelapp.io/citadel/repository"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

var (
	assets   string
	repoAddr string

	repo *repository.Repository

	logger = logrus.New()
)

func init() {
	flag.StringVar(&assets, "assets", "management", "path the the http assets")
	flag.StringVar(&repoAddr, "repository", "127.0.0.1:28015", "repository address")

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

func main() {
	var (
		err error
		// we cannot use the mux router with static files, there is a bug where it does
		// not handle the full directory
		globalMux = http.NewServeMux()
		apiRouter = mux.NewRouter()
	)

	if repo, err = repository.New(repoAddr); err != nil {
		logger.WithField("error", err).Fatal("cannot connect to repository")
	}

	apiRouter.HandleFunc("/api/hosts", getHosts)
	apiRouter.HandleFunc("/api/containers", getContainers)

	globalMux.Handle("/api/", apiRouter)
	globalMux.Handle("/", http.FileServer(http.Dir(assets)))

	if err := http.ListenAndServe(":3002", globalMux); err != nil {
		logger.WithField("error", err).Fatal("serve management ui")
	}
}
