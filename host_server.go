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

	// host information endpoints
	s.r.HandleFunc("/host", s.hostHandler).Methods("GET")

	s.r.HandleFunc("/stop/{id:.*}", s.stopHandler).Methods("POST")
	s.r.HandleFunc("/run/{id:.*}", s.runHandler).Methods("POST")

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

func (s *Server) runHandler(w http.ResponseWriter, r *http.Request) {
	id := getId(r)

	// TODO: this needs to return some basic information
	if err := s.host.RunContainer(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) stopHandler(w http.ResponseWriter, r *http.Request) {
	id := getId(r)

	if err := s.host.StopContainer(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) registerHandler(w http.ResponseWriter, r *http.Request) {

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

func getId(r *http.Request) string {
	return mux.Vars(r)["id"]
}
