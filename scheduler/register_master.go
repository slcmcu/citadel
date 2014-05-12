package main

import (
	"time"

	"citadelapp.io/citadel/master"
	"citadelapp.io/citadel/repository"
)

func registerMaster(m *master.Master, ttl int, repo repository.Repository) {
	if err := repo.RegisterMaster(m, ttl); err != nil {
		logger.WithField("error", err).Fatal("register master")
	}
	go masterHeartbeat(repo, ttl)
}

func masterHeartbeat(repo repository.Repository, ttl int) {
	for _ = range time.Tick(time.Duration(ttl-2) * time.Second) {
		for i := 0; i < 5; i++ {
			err := repo.UpdateMaster(ttl)
			if err == nil {
				continue
			}
			logger.WithField("error", err).Error("updating ttl")
			time.Sleep(500 * time.Second)
		}
	}
}
