package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"

	"citadelapp.io/citadel/metrics"
	"citadelapp.io/citadel/repository"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

var (
	assets string

	store metrics.Store
	repo  repository.Repository

	machines = []string{os.Getenv("ETCD_MACHINES")}
	logger   = logrus.New()
)

func init() {
	flag.StringVar(&assets, "assets", "management", "Path the the http assets")
	flag.Parse()
}

func getservices(w http.ResponseWriter, r *http.Request) {
	services, err := repo.FetchServices("/")
	if err != nil {
		httpError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(services); err != nil {
		httpError(w, err)
		return
	}
}

func getservice(w http.ResponseWriter, r *http.Request) {

}

func httpError(w http.ResponseWriter, err error) {
	logger.Error(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func main() {
	var (
		err error
		// we cannot use the mux router with static files, there is a bug where it does
		// not handle the full directory
		globalMux = http.NewServeMux()
		apiRouter = mux.NewRouter()
	)

	repo = repository.NewEtcdRepository(machines, false)

	conf, err := repo.FetchConfig()
	if err != nil {
		logger.WithField("error", err).Fatal("fetch config")
	}

	if store, err = metrics.NewStore(conf); err != nil {
		logger.WithField("error", err).Fatal("new metrics store")
	}

	apiRouter.HandleFunc("/api/services", getservices)
	apiRouter.HandleFunc("/api/services/{name:.*}", getservice)

	globalMux.Handle("/api/", apiRouter)
	globalMux.Handle("/", http.FileServer(http.Dir(assets)))

	if err := http.ListenAndServe(":3002", globalMux); err != nil {
		logger.WithField("error", err).Fatal("serve management ui")
	}
}
