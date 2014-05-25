package repository

import (
	"encoding/json"
	"path"
	"path/filepath"
	"strings"

	"citadelapp.io/citadel"
	"github.com/coreos/go-etcd/etcd"
)

type etcdRepository struct {
	client *etcd.Client
}

func NewEtcdRepository(machines []string, sync bool) Repository {
	r := &etcdRepository{
		client: etcd.NewClient(machines),
	}
	if sync {
		r.client.SyncCluster()
	}
	return r
}

// name == /local
// translatted to /citadel/services/local/services
func (e *etcdRepository) FetchServices(name string) ([]*citadel.ServiceData, error) {
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

func (e *etcdRepository) createServiceData(nodes etcd.Nodes, out *[]*citadel.ServiceData) error {
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

func (e *etcdRepository) FetchService(name string) (*citadel.ServiceData, error) {
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
func (e *etcdRepository) SaveService(name string, s *citadel.ServiceData) error {
	data, err := e.marshal(s)
	if err != nil {
		return err
	}

	_, err = e.client.Set(path.Join("/citadel/services", buildServiceName(name, "config")), data, 0)
	return err
}

func (e *etcdRepository) FetchConfig() (*citadel.Config, error) {
	resp, err := e.client.Get("/citadel/config", false, false)
	if err != nil {
		return nil, err
	}
	var c *citadel.Config
	if err := e.unmarshal(resp.Node.Value, &c); err != nil {
		return nil, err
	}
	return c, nil
}

// marshal encodes the value into a string via the json encoder
func (e *etcdRepository) marshal(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// unmarshal decodes the data using the json decoder into the value v
func (e *etcdRepository) unmarshal(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

// isNotFoundErr returns true if the error is of type Key Not Found
func isNotFoundErr(err error) bool {
	return strings.Contains(err.Error(), "Key not found")
}

func buildServiceName(fullPath, prefix string) string {
	if fullPath == "" || fullPath == "/" {
		return "/"
	}

	dir, name := filepath.Split(fullPath)

	var (
		parts = strings.Split(dir, "/")
		full  = []string{}
	)

	switch len(parts) {
	case 0:
		return path.Join(name, prefix)
	}

	for _, p := range parts {
		if p != "" {
			full = append(full, p, "services")
		}
	}

	return path.Join(append(full, name, prefix)...)
}
