package handler

import (
	"encoding/json"
	"net/http"

	"citadelapp.io/citadel"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

type Handler struct {
	router  *mux.Router
	service citadel.Service
	logger  *logrus.Logger
}

func New(service citadel.Service, logger *logrus.Logger) http.Handler {
	r := mux.NewRouter()

	h := &Handler{
		router:  r,
		service: service,
		logger:  logger,
	}

	r.HandleFunc("/", h.listhandler).Methods("POST")
	r.HandleFunc("/run", h.runhandler).Methods("POST")
	r.HandleFunc("/stop", h.stophandler).Methods("POST")

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) listhandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("list")

	task, err := h.unmarshalTask(r)
	if err != nil {
		h.httpError(w, err)
		return
	}

	response, err := h.service.List(task)
	if err != nil {
		h.httpError(w, err)
		return
	}

	h.marshal(w, response)
}

func (h *Handler) runhandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("run")

	task, err := h.unmarshalTask(r)
	if err != nil {
		h.httpError(w, err)
		return
	}

	response, err := h.service.Run(task)
	if err != nil {
		h.httpError(w, err)
		return
	}

	h.marshal(w, response)
}

func (h *Handler) stophandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("stop")

	task, err := h.unmarshalTask(r)
	if err != nil {
		h.httpError(w, err)
		return
	}

	response, err := h.service.Stop(task)
	if err != nil {
		h.httpError(w, err)
		return
	}

	h.marshal(w, response)
}

func (h *Handler) unmarshalTask(r *http.Request) (*citadel.Task, error) {
	var t *citadel.Task
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, err
	}
	return t, nil
}

func (h *Handler) marshal(w http.ResponseWriter, v interface{}) {
	w.Header().Add("content-type", "application/json")

	if err := json.NewEncoder(w).Encode(v); err != nil {
		h.httpError(w, err)
	}
}

func (h *Handler) httpError(w http.ResponseWriter, err error) {
	h.logger.Error(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
