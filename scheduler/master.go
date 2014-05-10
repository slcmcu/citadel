package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"citadelapp.io/citadel"
	"citadelapp.io/citadel/repository"
	"citadelapp.io/citadel/scheduler/master"
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func registerMaster(m *master.Master, ttl int, repo repository.Repository) {
	if err := repo.RegisterMaster(&m.Master, ttl); err != nil {
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

	http.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var task *citadel.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			logger.WithField("error", err).Error("decoding task")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		logger.WithFields(logrus.Fields{
			"instances": task.Instances,
			"image":     task.Container.Image,
			"cpus":      task.Container.Cpus,
			"memory":    task.Container.Memory,
		}).Info("scheduling task")

		slaves, err := m.Schedule(task, repo)
		if err != nil {
			logger.WithField("error", err).Error("cannot schedule task")
			if err == master.ErrNoValidOffers {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		var scheduleError error
		for _, s := range slaves {
			var reply *citadel.Container
			if err := nc.Request(fmt.Sprintf("execute.%s", s.ID), task.Container, &reply, timeout); err != nil {
				logger.WithField("error", err).Error("cannot publish task")
				if scheduleError == nil {
					scheduleError = err
				}
			}

			logger.WithFields(logrus.Fields{
				"slave":        s.ID,
				"container_id": reply.ID,
			}).Info("scheduled task")
		}
		if scheduleError != nil {
			http.Error(w, scheduleError.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/pull", func(w http.ResponseWriter, r *http.Request) {
		image := r.URL.Query().Get("image")
		logger.WithField("image", image).Info("pulling image")

		if err := nc.Publish("slaves.pull", image); err != nil {
			logger.WithField("error", err).Error("pull image")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	if err := http.ListenAndServe(m.Addr, nil); err != nil {
		logger.WithField("error", err).Fatal("serve master")
	}
}
