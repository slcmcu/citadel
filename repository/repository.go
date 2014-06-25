package repository

import (
	"fmt"

	"github.com/coreos/go-etcd/etcd"
)

type Repository struct {
	client    *etcd.Client
	namespace string
}

// New returns a new repository connected the the database
func New(etcdMachines []string, namespace string) *Repository {
	c := etcd.NewClient(etcdMachines)

	repo := &Repository{
		client:    c,
		namespace: fmt.Sprintf("/%s", namespace),
	}
	repo.init()
	return repo
}

// init initializes etcd with top level keyspaces
func (r *Repository) init() {
	keys := []string{
		r.namespace,
		fmt.Sprintf("%s/hosts", r.namespace),
		fmt.Sprintf("%s/tasks", r.namespace),
	}
	for _, k := range keys {
		r.client.CreateDir(k, 0)
	}

}
