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

func NewEtcdRepository(machines []string) Repository {
	r := &etcdRepository{
		client: etcd.NewClient(machines),
	}
	r.client.SyncCluster()

	return r
}

func (e *etcdRepository) SaveHost(h *citadel.Host) error {
	body, err := e.marshal(h)
	if err != nil {
		return err
	}
	if _, err := e.client.Set(path.Join("/citadel/hosts", h.Name), body, 0); err != nil {
		return err
	}
	return nil
}

func (e *etcdRepository) DeleteHost(h *citadel.Host) error {
	_, err := e.client.Delete(path.Join("/citadel/hosts", h.Name), true)
	return err
}

func (e *etcdRepository) FetchHost(name string) (*citadel.Host, error) {
	resp, err := e.client.Get(path.Join("/citadel/hosts", name), false, false)
	if err != nil {
		return nil, err
	}
	var h *citadel.Host
	if err := e.unmarshal(resp.Node.Value, &h); err != nil {
		return nil, err
	}
	return h, nil
}

func (e *etcdRepository) FetchHosts() ([]*citadel.Host, error) {
	hosts := []*citadel.Host{}
	resp, err := e.client.Get("/citadel/hosts", true, true)
	if err != nil {
		if isNotFoundErr(err) {
			return hosts, nil
		}
		return nil, err
	}
	for _, n := range resp.Node.Nodes {
		var h *citadel.Host
		if err := e.unmarshal(n.Value, &h); err != nil {
			return nil, err
		}
		hosts = append(hosts, h)
	}
	return hosts, nil
}

func (e *etcdRepository) FetchContainerGroup() ([]*citadel.ContainerGroup, error) {
	images := []*citadel.ContainerGroup{}
	resp, err := e.client.Get("/citadel/containers", true, true)
	if err != nil {
		if isNotFoundErr(err) {
			return images, nil
		}
		return nil, err
	}

	for _, n := range resp.Node.Nodes {
		i := &citadel.ContainerGroup{
			Name:      n.Key,
			Instances: n.Nodes.Len(),
			Status:    "healthy",
		}
		images = append(images, i)
	}
	return images, nil
}

func (e *etcdRepository) marshal(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (e *etcdRepository) unmarshal(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

func isNotFoundErr(err error) bool {
	return strings.Contains(err.Error(), "Key not found")
}
