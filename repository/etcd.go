package repository

import (
	"encoding/json"
	"path"
	"strings"

	"citadelapp.io/citadel"
	"citadelapp.io/citadel/master"
	"citadelapp.io/citadel/slave"
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
func (e *etcdRepository) RegisterSlave(uuid string, s *slave.Slave, ttl int) error {
	data, err := e.marshal(s)
	if err != nil {
		return err
	}
	if _, err := e.client.CreateDir(path.Join("/citadel/slaves", uuid), uint64(ttl)); err != nil {
		return err
	}
	_, err = e.client.Set(path.Join("/citadel/slaves", uuid, "config"), data, 0)
	return err
}

func (e *etcdRepository) UpdateSlave(uuid string, ttl int) error {
	_, err := e.client.UpdateDir(path.Join("/citadel/slaves", uuid), uint64(ttl))
	return err
}

func (e *etcdRepository) RemoveSlave(uuid string) error {
	_, err := e.client.Delete(path.Join("/citadel/slaves", uuid), true)
	return err
}

func (e *etcdRepository) FetchSlave(uuid string) (*slave.Slave, error) {
	resp, err := e.client.Get(path.Join("/citadel/slaves", uuid, "config"), false, false)
	if err != nil {
		return nil, err
	}
	var s *slave.Slave
	if err := e.unmarshal(resp.Node.Value, &s); err != nil {
		return nil, err
	}
	return s, nil
}

func (e *etcdRepository) FetchSlaves() ([]*slave.Slave, error) {
	slaves := []*slave.Slave{}
	resp, err := e.client.Get("/citadel/slaves", true, true)
	if err != nil {
		if isNotFoundErr(err) {
			return slaves, nil
		}
		return nil, err
	}

	for _, n := range resp.Node.Nodes {
		for _, sdir := range n.Nodes {
			if sdir.Key == path.Join(n.Key, "config") {
				var s *slave.Slave
				if err := e.unmarshal(sdir.Value, &s); err != nil {
					return nil, err
				}
				slaves = append(slaves, s)
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

func (e *etcdRepository) RemoveContainer(uuid, id string) error {
	_, err := e.client.Delete(path.Join("/citadel/slaves", uuid, "containers", id), true)
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

func (e *etcdRepository) RegisterMaster(m *master.Master, ttl int) error {
	data, err := e.marshal(m)
	if err != nil {
		return err
	}
	_, err = e.client.Create(path.Join("/citadel/master"), data, uint64(ttl))
	return err
}

func (e *etcdRepository) UpdateMaster(ttl int) error {
	_, err := e.client.UpdateDir(path.Join("/citadel/master"), uint64(ttl))
	return err
}

func (e *etcdRepository) FetchMaster() (*master.Master, error) {
	resp, err := e.client.Get("/citadel/master", false, false)
	if err != nil {
		return nil, err
	}

	var m *master.Master
	if err := e.unmarshal(resp.Node.Value, &m); err != nil {
		return nil, err
	}
	return m, nil
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
