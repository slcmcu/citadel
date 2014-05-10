package repository

import (
	"encoding/json"
	"path"
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

// RegisterSlave registers the uuid and slave information into the key
// /citadel/slaves/<uuid> and /citadel/slaves/<uuid>/config with the ttl
func (e *etcdRepository) RegisterSlave(uuid string, slave *citadel.Slave, ttl int) error {
	data, err := e.marshal(slave)
	if err != nil {
		return err
	}
	if _, err := e.client.CreateDir(path.Join("/citadel/slaves", uuid), uint64(ttl)); err != nil {
		return err
	}
	if _, err := e.client.Set(path.Join("/citadel/slaves", uuid, "config"), data, 0); err != nil {
		return err
	}
	return nil
}

func (e *etcdRepository) UpdateSlave(uuid string, ttl int) error {
	if _, err := e.client.UpdateDir(path.Join("/citadel/slaves", uuid), uint64(ttl)); err != nil {
		return err
	}
	return nil
}

func (e *etcdRepository) RemoveSlave(uuid string) error {
	_, err := e.client.Delete(path.Join("/citadel/slaves", uuid), true)
	return err
}

func (e *etcdRepository) FetchSlave(uuid string) (*citadel.Slave, error) {
	resp, err := e.client.Get(path.Join("/citadel/slaves", uuid, "config"), false, false)
	if err != nil {
		return nil, err
	}
	var s *citadel.Slave
	if err := e.unmarshal(resp.Node.Value, &s); err != nil {
		return nil, err
	}
	return s, nil
}

func (e *etcdRepository) FetchSlaves() ([]*citadel.Slave, error) {
	slaves := []*citadel.Slave{}
	resp, err := e.client.Get("/citadel/slaves", true, true)
	if err != nil {
		if isNotFoundErr(err) {
			return slaves, nil
		}
		return nil, err
	}
	for _, n := range resp.Node.Nodes {
		if n.Dir {
			for _, sdir := range n.Nodes {
				if !sdir.Dir {
					var s *citadel.Slave
					if err := e.unmarshal(n.Value, &s); err != nil {
						return nil, err
					}
					slaves = append(slaves, s)
				}
			}
		}
	}
	return slaves, nil
}

func (e *etcdRepository) SaveContainer(uuid string, c *citadel.Container) error {
	data, err := e.marshal(c)
	if err != nil {
		return err
	}
	_, err = e.client.Set(path.Join("/citadel/slaves", uuid, "containers", c.ID), data, 0)
	return err
}

func (e *etcdRepository) FetchContainers(uuid string) (citadel.Containers, error) {
	containers := citadel.Containers{}
	resp, err := e.client.Get(path.Join("/citadel/slaves", uuid, "containers"), false, true)
	if err != nil {
		if isNotFoundErr(err) {
			return containers, nil
		}
		return nil, err
	}

	for _, n := range resp.Node.Nodes {
		var c *citadel.Container
		if err := e.unmarshal(n.Value, &c); err != nil {
			return nil, err
		}
		containers[c.ID] = c
	}
	return containers, nil
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
