package repository

import (
	"path"
	"path/filepath"
	"strings"

	"citadelapp.io/citadel"
	"github.com/coreos/go-etcd/etcd"
)

type Repository struct {
	client *etcd.Client
}

func New(machines []string, sync bool) *Repository {
	r := &Repository{
		client: etcd.NewClient(machines),
	}
	if sync {
		r.client.SyncCluster()
	}
	return r
}

// name == /local
// translatted to /citadel/services/local/services
func (e *Repository) FetchServices(name string) ([]*citadel.ServiceData, error) {
	out := []*citadel.ServiceData{}

	resp, err := e.client.Get(path.Join("/citadel/services", buildServiceName(name, "services")), true, true)
	if err != nil {
		if isNotFoundErr(err) {
			return out, nil
		}
		return nil, err
	}

	if err := e.createServiceData(resp.Node.Nodes, &out); err != nil {
		return nil, err
	}

	return out, nil
}

func (e *Repository) createServiceData(nodes etcd.Nodes, out *[]*citadel.ServiceData) error {
	for _, n := range nodes {
		for _, sdir := range n.Nodes {
			_, name := filepath.Split(sdir.Key)

			switch name {
			case "config":
				var s *citadel.ServiceData
				if err := e.unmarshal(sdir.Value, &s); err != nil {
					return err
				}
				s.Name = strings.TrimSuffix(strings.TrimPrefix(sdir.Key, "/citadel/services"), "/config")
				*out = append(*out, s)
			case "services":
				if err := e.createServiceData(sdir.Nodes, out); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (e *Repository) FetchService(name string) (*citadel.ServiceData, error) {
	resp, err := e.client.Get(path.Join("/citadel/services", buildServiceName(name, "config")), true, true)
	if err != nil {
		return nil, err
	}

	var s *citadel.ServiceData
	if err := e.unmarshal(resp.Node.Value, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// name == local/redis
// translatted to /citadel/services/local/services/redis/config
func (e *Repository) SaveService(name string, s *citadel.ServiceData) error {
	data, err := e.marshal(s)
	if err != nil {
		return err
	}

	_, err = e.client.Set(path.Join("/citadel/services", buildServiceName(name, "config")), data, 0)
	return err
}
