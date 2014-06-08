package repository

import (
	"citadelapp.io/citadel"
	"github.com/dancannon/gorethink"
)

type Repository struct {
	session *gorethink.Session
}

func New(addr string) (*Repository, error) {
	s, err := gorethink.Connect(gorethink.ConnectOpts{
		Address:  addr,
		Database: "citadel",
	})
	if err != nil {
		return nil, err
	}

	return &Repository{
		session: s,
	}, nil
}

// Closes closes the underlying session to the datastore
func (r *Repository) Close() error {
	return r.session.Close()
}

// SaveHost will save a single host with the unique id
func (r *Repository) SaveHost(h *citadel.Host) error {
	if _, err := gorethink.Table("hosts").Insert(h).RunWrite(r.session); err != nil {
		return err
	}

	return nil
}

// FetchHosts returns all the hosts in the cluseter
func (r *Repository) FetchHosts() ([]*citadel.Host, error) {
	results, err := gorethink.Table("hosts").Run(r.session)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	var hosts []*citadel.Host
	if err := results.ScanAll(&hosts); err != nil {
		return nil, err
	}

	return hosts, nil
}
