package main

import (
	"flag"
	"net/http"
	"strings"
	"time"

	"citadelapp.io/citadel"
	"citadelapp.io/citadel/metrics"
	"citadelapp.io/citadel/repository"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

var (
	listen   string
	machines string
	log      = logrus.New()
)

func init() {
	flag.StringVar(&listen, "listen", "", "Listen address")
	flag.StringVar(&machines, "machines", "127.0.0.1:4001", "Comma separated list of etcd machines")
	flag.Parse()
}

func collectMetrics(store metrics.Store, host *citadel.Host, conf *citadel.Config) {
	for _ = range time.Tick(time.Duration(conf.PullInterval) * time.Second) {
		if err := pushHostMetrics(host, store); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	if listen == "" {
		log.Fatal("You must specify a listen address")
	}
	var (
		etcdMachines = strings.Split(machines, ",")
		repo         = repository.NewEtcdRepository(etcdMachines)
	)

	conf, err := repo.FetchConfig()
	if err != nil {
		log.Fatal(err)
	}
	agentName, err := getAgentName()
	if err != nil {
		log.Fatal(err)
	}
	host, err := getHostInfo(agentName)
	if err != nil {
		log.Fatal(err)
	}

	// save to repo
	if err := repo.SaveHost(host); err != nil {
		log.Fatalf("Unable to save host: %s", err)
	}

	log.WithFields(logrus.Fields{
		"cpus":      host.Cpus,
		"memory":    host.TotalMemory,
		"diskspace": host.Disks,
	}).Debug("Host Info")

	store, err := metrics.NewStore(conf)
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(logrus.Fields{
		"nodename": agentName,
		"address":  listen,
	}).Info("Citadel Agent")

	go collectMetrics(store, host, conf)

	r := mux.NewRouter()
	if err := http.ListenAndServe(listen, r); err != nil {
		log.Fatal(err)
	}
}
