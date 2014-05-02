package main

import (
	"flag"

	"citadelapp.io/citadel"
	"github.com/Sirupsen/logrus"
	rethink "github.com/dancannon/gorethink"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

var (
	addr string
	log  = logrus.New()
)

func init() {
	flag.StringVar(&addr, "l", "", "rethinkdb address")
	flag.Parse()
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())

	session, err := citadel.NewRethinkSession(addr)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	m.Group("/api", func(r martini.Router) {
		r.Get("/hosts", func(rdr render.Render) {
			data := []*citadel.Host{}
			rows, err := rethink.Table("host").Run(session)
			if err != nil {
				log.Fatal(err)
			}
			for rows.Next() {
				var host *citadel.Host
				if err := rows.Scan(&host); err != nil {
					log.Fatal(err)
				}
				data = append(data, host)
			}
			rdr.JSON(200, data)
		})
	})

	m.Use(martini.Static("."))

	m.Run()
}
