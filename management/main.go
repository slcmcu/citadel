package main

import (
	"flag"

	"citadelapp.io/citadel"
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

	m.Group("/api", func(r martini.Router) {
		r.Get("/hosts", func(rdr render.Render) {
			data, err := repo.FetchHosts()
			if err != nil {
				log.Fatal(err)
			}
			rdr.JSON(200, data)
		})

		r.Get("/hosts/memory", func(rdr render.Render) {

		})
	})

	m.Use(martini.Static("."))
	m.Run()
}
