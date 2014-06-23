package repository

import (
	"github.com/dancannon/gorethink"
)

type Repository struct {
	session *gorethink.Session
}

// New returns a new repository connected the the database
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

func (r *Repository) Session() *gorethink.Session {
	return r.session
}
