package repository

import (
	"path"

	"citadelapp.io/citadel"
)

func (e *Repository) SaveHost(h *citadel.Host, ttl int) error {
	data, err := e.marshal(h)
	if err != nil {
		return err
	}

	_, err := e.client.Set(path.Join("/citadel/hosts", h.ID), data, ttl)
	return err
}

func (e *Repository) FetchHost(id string) (*citadel.Host, error) {

}

func (e *Repository) FetchHosts() ([]*citadel.Host, error) {

}

func (e *Repository) UpdateHostTTL(h *citadel.Host, ttl int) error {

}
