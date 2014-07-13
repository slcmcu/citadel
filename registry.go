package citadel

import (
	"encoding/json"
	"path/filepath"

	"github.com/coreos/go-etcd/etcd"
)

type Registry struct {
	client *etcd.Client
}

func NewRegistry(machines []string) *Registry {
	return &Registry{
		client: etcd.NewClient(machines),
	}
}

func (r *Registry) FetchContainers(h *Host) ([]*Container, error) {
	resp, err := r.client.Get(filepath.Join("/citadel", h.ID, "containers"), true, true)
	if err != nil {
		return nil, err
	}

	out := []*Container{}
	for _, node := range resp.Node.Nodes {
		var container *Container
		if err := json.Unmarshal([]byte(node.Value), &container); err != nil {
			return nil, err
		}

		out = append(out, container)
	}

	return out, nil
}

func (r *Registry) FetchContainer(h *Host, id string) (*Container, error) {
	resp, err := r.client.Get(filepath.Join("/citadel", h.ID, "containers", id), false, false)
	if err != nil {
		return nil, err
	}

	var container *Container
	if err := json.Unmarshal([]byte(resp.Node.Value), &container); err != nil {
		return nil, err
	}

	return container, nil
}

func (r *Registry) SaveContainer(h *Host, c *Container) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	if _, err := r.client.Set(filepath.Join("/citadel", h.ID, "containers", c.ID), string(data), 0); err != nil {
		return err
	}

	return nil
}

func (r *Registry) DeleteContainer(h *Host, c *Container) error {
	_, err := r.client.Delete(filepath.Join("/citadel", h.ID, "containers", c.ID), false)
	return err
}

func (r *Registry) SaveHost(h *Host) error {
	data, err := json.Marshal(h)
	if err != nil {
		return err
	}

	if _, err := r.client.Set(filepath.Join("/citadel/hosts", h.ID), string(data), 0); err != nil {
		return err
	}

	return nil
}

func (r *Registry) FetchHosts() ([]*Host, error) {
	hosts := []*Host{}

	resp, err := r.client.Get("/citadel/hosts", true, true)
	if err != nil {
		return nil, err
	}

	for _, n := range resp.Node.Nodes {
		var host *Host
		if err := json.Unmarshal([]byte(n.Value), &host); err != nil {
			return nil, err
		}

		hosts = append(hosts, host)
	}

	return hosts, nil
}

func (r *Registry) DeleteHost(h *Host) error {
	_, err := r.client.Delete(filepath.Join("/citadel/hosts", h.ID), false)
	return err
}
