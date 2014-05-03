package main

import (
	"flag"
	"fmt"

	"citadelapp.io/citadel"
	"citadelapp.io/citadel/metrics"
	"citadelapp.io/citadel/repository"
	"github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

var (
	configPath string
	log        = logrus.New()
)

func init() {
	flag.StringVar(&configPath, "config", "config.toml", "path to the configuration file")
	flag.Parse()
}

func main() {
	conf, err := citadel.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	m := martini.Classic()
	m.Use(render.Renderer())

	repo := repository.NewEtcdRepository(conf.Machines)
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

		r.Get("/hosts/:name/metrics", func(params martini.Params, rdr render.Render) {
			table := fmt.Sprintf("metrics.hosts.%s", params["name"])
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
