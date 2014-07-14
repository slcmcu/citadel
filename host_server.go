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

	s.r.HandleFunc("/app/{id:.*}", s.loadHandler).Methods("POST")
	s.r.HandleFunc("/app/{id:.*}", s.deleteHandler).Methods("DELETE")

	// /run runs the givin application's containers on the host
	s.r.HandleFunc("/run/{id:.*}", s.runHandler).Methods("POST")

	// /stop stops the givin application's containers running on the host
	s.r.HandleFunc("/stop/{id:.*}", s.stopHandler).Methods("POST")

	return s
}

func (s *Server) Close() error {
	return s.host.Close()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}

func (s *Server) runHandler(w http.ResponseWriter, r *http.Request) {
	id := getId(r)

	tran := s.host.RunContainer(id)
	if tran.Err != nil {
		w.Header().Set("Status", fmt.Sprint(http.StatusInternalServerError))
	}

	s.marshal(w, tran)
}

func (s *Server) stopHandler(w http.ResponseWriter, r *http.Request) {
	id := getId(r)

	tran := s.host.StopContainer(id)
	if tran.Err != nil {
		w.Header().Set("Status", fmt.Sprint(http.StatusInternalServerError))
	}

	s.marshal(w, tran)
}

func (s *Server) loadHandler(w http.ResponseWriter, r *http.Request) {
	id := getId(r)

	tran := s.host.Load(id)
	if tran.Err != nil {
		w.Header().Set("Status", fmt.Sprint(http.StatusInternalServerError))
	}

	s.marshal(w, tran)
}

func (s *Server) deleteHandler(w http.ResponseWriter, r *http.Request) {
	id := getId(r)

	tran := s.host.Delete(id)
	if tran.Err != nil {
		w.Header().Set("Status", fmt.Sprint(http.StatusInternalServerError))
	}

	s.marshal(w, tran)
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
