package main

import (
	"time"

	"citadelapp.io/citadel/repository"
	"citadelapp.io/citadel/slave"
)

func register(s *slave.Slave, ttl int, repo repository.Repository) error {
	if err := repo.RegisterSlave(s.ID, s, ttl); err != nil {
		return err
	}
	go heartbeat(s.ID, ttl, repo)
	return nil
}

func heartbeat(uuid string, ttl int, repo repository.Repository) {
	for _ = range time.Tick(time.Duration(ttl-2) * time.Second) {
		for i := 0; i < 5; i++ {
			err := repo.UpdateSlave(uuid, ttl)
			if err == nil {
				continue
			}
			logger.WithField("error", err).Error("updating ttl")
			time.Sleep(500 * time.Second)
		}
	}
}
