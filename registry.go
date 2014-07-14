package citadel

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/coreos/go-etcd/etcd"
)

type Registry interface {
	SaveHost(*Host) error
	FetchHost(string) (*Host, error)
	FetchHosts() ([]*Host, error)
	DeleteHost(string) error

	SaveApplication(*Application) error
	FetchApplication(string) (*Application, error)
	FetchApplications() ([]*Application, error)
	DeleteApplication(string) error

	SaveContainer(string, *Container) error
	DeleteContainer(string, string) error
	FetchContainer(string, string) (*Container, error)
	FetchContainers(string) ([]*Container, error)
}

type registry struct {
	client *etcd.Client
}

func NewRegistry(machines []string) Registry {
	return &registry{
		client: etcd.NewClient(machines),
	}
}

func (r *registry) SaveApplication(a *Application) error {
	data, err := json.Marshal(a)
	if err != nil {
		return err
	}

	_, err = r.client.Set(filepath.Join("/citadel/applications", a.ID), string(data), 0)
	return err
}

func (r *registry) DeleteApplication(id string) error {
	_, err := r.client.Delete(filepath.Join("/citadel/applications", id), false)
	return err
}

func (r *registry) FetchApplications() ([]*Application, error) {
	out := []*Application{}

	resp, err := r.client.Get("/citadel/applications", true, true)
	if err != nil {
		if isNotFound(err) {
			return out, nil
		}

		return nil, err
	}

	for _, node := range resp.Node.Nodes {
		var app *Application
		if err := json.Unmarshal([]byte(node.Value), &app); err != nil {
			return nil, err
		}

		out = append(out, app)
	}

	return out, nil
}

func (r *registry) FetchApplication(id string) (*Application, error) {
	resp, err := r.client.Get(filepath.Join("/citadel/applications", id), false, false)
	if err != nil {
		return nil, err
	}

	var app *Application
	if err := json.Unmarshal([]byte(resp.Node.Value), &app); err != nil {
		return nil, err
	}

	return app, nil
}

func (r *registry) FetchContainers(id string) ([]*Container, error) {
	out := []*Container{}

	resp, err := r.client.Get(filepath.Join("/citadel", id, "containers"), true, true)
	if err != nil {
		if isNotFound(err) {
			return out, nil
		}
		return nil, err
	}

	for _, node := range resp.Node.Nodes {
		var container *Container
		if err := json.Unmarshal([]byte(node.Value), &container); err != nil {
			return nil, err
		}

		out = append(out, container)
	}

	return out, nil
}

func (r *registry) FetchContainer(hostID, id string) (*Container, error) {
	resp, err := r.client.Get(filepath.Join("/citadel", hostID, "containers", id), false, false)
	if err != nil {
		return nil, err
	}

	var container *Container
	if err := json.Unmarshal([]byte(resp.Node.Value), &container); err != nil {
		return nil, err
	}

	return container, nil
}

func (r *registry) SaveContainer(hostID string, c *Container) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	if _, err := r.client.Set(filepath.Join("/citadel", hostID, "containers", c.ID), string(data), 0); err != nil {
		return err
	}

	return nil
}

func (r *registry) DeleteContainer(hostID, id string) error {
	_, err := r.client.Delete(filepath.Join("/citadel", hostID, "containers", id), false)
	return err
}

func (r *registry) SaveHost(h *Host) error {
	data, err := json.Marshal(h)
	if err != nil {
		return err
	}

	if _, err := r.client.Set(filepath.Join("/citadel/hosts", h.ID), string(data), 0); err != nil {
		return err
	}

	return nil
}

func (r *registry) FetchHosts() ([]*Host, error) {
	hosts := []*Host{}

	resp, err := r.client.Get("/citadel/hosts", true, true)
	if err != nil {
		if isNotFound(err) {
			return hosts, nil
		}
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

func (r *registry) FetchHost(id string) (*Host, error) {
	resp, err := r.client.Get(filepath.Join("/citadel/hosts", id), false, false)
	if err != nil {
		return nil, err
	}

	var host *Host
	if err := json.Unmarshal([]byte(resp.Node.Value), &host); err != nil {
		return nil, err
	}

	return host, nil
}

func (r *registry) DeleteHost(id string) error {
	_, err := r.client.Delete(filepath.Join("/citadel/hosts", id), false)
	return err
}

func isNotFound(err error) bool {
	return strings.Contains(err.Error(), "Key not found")
}
