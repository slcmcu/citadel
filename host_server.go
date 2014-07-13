package citadel

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	host *Host
	r    *mux.Router
}

func NewServer(h *Host) *Server {
	s := &Server{
		host: h,
		r:    mux.NewRouter(),
	}

	s.r.HandleFunc("/stop", s.stopHandler).Methods("POST")
	s.r.HandleFunc("/run", s.runHandler).Methods("POST")
	s.r.HandleFunc("/host", s.hostHandler).Methods("GET")
	s.r.HandleFunc("/containers/{id:.*}", s.getContainer).Methods("GET")
	s.r.HandleFunc("/containers", s.listHandler).Methods("GET")

	return s
}

func (s *Server) Close() error {
	return s.host.Close()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}

func (s *Server) hostHandler(w http.ResponseWriter, r *http.Request) {
	s.marshal(w, s.host)
}

func (s *Server) listHandler(w http.ResponseWriter, r *http.Request) {
	containers, err := s.host.Containers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.marshal(w, containers)
}

func (s *Server) getContainer(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	c, err := s.host.Container(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	s.marshal(w, c)
}

func (s *Server) runHandler(w http.ResponseWriter, r *http.Request) {
	var container *Container
	if err := s.unmarshal(r, &container); err != nil {
		// TODO: this could be a bad content type error
		// need to pick between the two
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.host.RunContainer(container); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.marshal(w, container)
}

func (s *Server) stopHandler(w http.ResponseWriter, r *http.Request) {
	var container *Container
	if err := s.unmarshal(r, &container); err != nil {
		// TODO: this could be a bad content type error
		// need to pick between the two
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.host.StopContainer(container); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.marshal(w, container)
}

func (s *Server) unmarshal(r *http.Request, v interface{}) error {
	defer r.Body.Close()

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return err
	}

	if mediaType != "application/json" {
		return fmt.Errorf("invalid content type, expect application/json")
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func (s *Server) marshal(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(v); err != nil {
		s.host.logger.WithField("error", err).Error("encode json")
	}
}
