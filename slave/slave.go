package slave

import "net/http"

type SlaveService struct {
}

func New() *SlaveService {

}

func (m *SlaveService) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, m.mux)
}
