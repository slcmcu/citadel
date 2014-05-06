package main

import (
	"flag"
	"fmt"
	"strings"

	"citadelapp.io/citadel/metrics"
	"citadelapp.io/citadel/repository"
	"github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

var (
	machines string
	log      = logrus.New()
)

func init() {
	flag.StringVar(&machines, "machines", "127.0.0.1:4001", "Comma separated list of etcd machines")
	flag.Parse()
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())

	etcdMachines := strings.Split(machines, ",")
	repo := repository.NewEtcdRepository(etcdMachines)
	conf, err := repo.FetchConfig()
	if err != nil {
		log.Fatalf("Unable to get config from etcd: %s", err)
	}
	store, err := metrics.NewStore(conf)
	if err != nil {
		log.Fatal(err)
	}

	m.Group("/api", func(r martini.Router) {
		r.Get("/hosts", func(rdr render.Render) {
			data, err := repo.FetchHosts()
			if err != nil {
				log.Fatal(err)
			}
			rdr.JSON(200, data)
		})

		r.Get("/containers", func(rdr render.Render) {
			data, err := repo.FetchContainerGroup()
			if err != nil {
				log.Fatal(err)
			}
			rdr.JSON(200, data)
		})

		r.Get("/hosts/:name/metrics", func(params martini.Params, rdr render.Render) {
			table := fmt.Sprintf("metrics.host.%s", params["name"])
			data, err := store.Fetch(fmt.Sprintf("select * from %s group by time(5m) where time > now() -12h limit 200", table))
			if err != nil {
				log.Fatal(err)
			}
			rdr.JSON(200, data)
		})
	})

	m.Use(martini.Static("."))
	m.Run()
}
