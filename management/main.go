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

		r.Get("/hosts/memory", func(rdr render.Render) {
			rows, err := rethink.Table("host_metric").WithFields("timestamp", "memory").OrderBy(rethink.Desc("timestamp")).Group(func(d rethink.RqlTerm) rethink.RqlTerm {
				return d.Field("timestamp").Minutes().Mod(2).Eq(0)
			}).Ungroup().Filter(func(d rethink.RqlTerm) rethink.RqlTerm {
				return d.Field("group")
			}).Map(func(d rethink.RqlTerm) rethink.RqlTerm {
				return d.Field("reduction")
			}).Nth(0).Map(func(d rethink.RqlTerm) rethink.RqlTerm {
				return d.Field("memory").Field("used").Div(d.Field("memory").Field("total")).Mul(100)
			}).Limit(300).Run(session)

			if err != nil {
				log.Fatal(err)
			}
			type stat struct {
				Key   int     `json:"key"`
				Value float64 `json:"value"`
			}
			var (
				i    int
				data = []*stat{}
			)
			for rows.Next() {
				i++
				var v float64
				if err := rows.Scan(&v); err != nil {
					log.Fatal(err)
				}
				data = append(data, &stat{i, v})
			}
			rdr.JSON(200, data)
		})

	})

	m.Use(martini.Static("."))

	m.Run()
}
