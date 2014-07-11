package repository

import (
	"encoding/json"
	"fmt"

	"github.com/citadel/citadel"
)

func (r *Repository) SaveContainer(c *citadel.Container) error {
	d, err := json.Marshal(c)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s/hosts/%s/containers/%s", r.namespace, c.HostID, c.ID)
	if _, err := r.client.Set(key, string(d), 0); err != nil {
		return err
	}
	return nil
}

func (r *Repository) FetchContainers() ([]*citadel.Container, error) {
	// find all hosts and then find all containers for each host
	hosts, err := r.FetchHosts()
	if err != nil {
		return nil, err
	}

	// find containers
	var containers []*citadel.Container
	for _, h := range hosts {
		cKey := fmt.Sprintf("%s/hosts/%s/containers", r.namespace, h.ID)
		cResults, err := r.client.Get(cKey, false, true)
		if err != nil {
			return nil, err
		}
		for _, cr := range cResults.Node.Nodes {
			var c *citadel.Container
			if err := json.Unmarshal([]byte(cr.Value), &c); err != nil {
				return nil, err
			}
			containers = append(containers, c)
		}
	}

	return containers, nil
}

func (r *Repository) DeleteContainer(hostId string, id string) error {
	// remove host container
	key := fmt.Sprintf("%s/hosts/%s/containers/%s", r.namespace, hostId, id)
	if _, err := r.client.Delete(key, true); err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteHostContainers(hostId string) error {
	// remove host containers
	key := fmt.Sprintf("%s/hosts/%s/containers/", r.namespace, hostId)
	if _, err := r.client.RawDelete(key, true, true); err != nil {
		return err
	}

	return nil
}
