package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"citadelapp.io/citadel"
	"citadelapp.io/citadel/master"
	"citadelapp.io/citadel/repository"
	"github.com/Sirupsen/logrus"
	"github.com/apcera/nats"
	"github.com/gorilla/mux"
)

type masterHandler struct {
	m    *master.Master
	nc   *nats.EncodedConn
	repo repository.Repository

	router *mux.Router
}

func newMasterHandler(m *master.Master, nc *nats.EncodedConn, repo repository.Repository) http.Handler {
	r := mux.NewRouter()

	h := &masterHandler{
		router: r,
		m:      m,
		nc:     nc,
		repo:   repo,
	}
	r.HandleFunc("/run", h.runHander)
	r.HandleFunc("/pull", h.pullHandler)

	return h
}

func (h *masterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *masterHandler) runHander(w http.ResponseWriter, r *http.Request) {
	var (
		task *citadel.Task
		m    = h.m
		nc   = h.nc
		repo = h.repo
	)

	defer r.Body.Close()

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
		if err := nc.Request(fmt.Sprintf("execute.%s", s.ID), task.Container, &reply, m.Timeout); err != nil {
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
}

func (h *masterHandler) pullHandler(w http.ResponseWriter, r *http.Request) {
	var (
		image = r.URL.Query().Get("image")
		nc    = h.nc
	)
	logger.WithField("image", image).Info("pulling image")

	if err := nc.Publish("slaves.pull", image); err != nil {
		logger.WithField("error", err).Error("pull image")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
