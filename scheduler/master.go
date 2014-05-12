package main

import (
	"net/http"
	"time"

	"citadelapp.io/citadel/master"
	"github.com/codegangsta/cli"
)

func masterMain(context *cli.Context) {
	var (
		uuid       = getUUID()
		ip         = context.String("addr")
		repo, conf = getRepositoryAndConfig(context)
		nc         = getNats(conf)
	)

	timeout, err := time.ParseDuration(conf.MasterTimeout)
	if err != nil {
		logger.WithField("error", err).Fatal("parse timeout duration")
	}

	m, err := master.New(uuid, ip, timeout)
	if err != nil {
		logger.WithField("error", err).Fatal("initializing master")
	}

	registerMaster(m, conf.TTL, repo)

	h := newMasterHandler(m, nc, repo)
	if err := http.ListenAndServe(m.Addr, h); err != nil {
		logger.WithField("error", err).Fatal("serve master")
	}
}
