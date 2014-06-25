package repository

import (
	"encoding/json"
	"fmt"
	"strings"

	"citadelapp.io/citadel"
)

// SaveHost will save a single host with the unique id
func (r *Repository) SaveHost(h *citadel.Host) error {
	d, err := json.Marshal(h)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s/hosts/%s/info", r.namespace, h.ID)
	if _, err := r.client.Set(key, string(d), 0); err != nil {
		return err
	}
	return nil
}

// FetchHosts returns all the hosts in the cluseter
func (r *Repository) FetchHosts() ([]*citadel.Host, error) {
	key := fmt.Sprintf("%s/hosts/", r.namespace)
	results, err := r.client.Get(key, false, true)
	if err != nil {
		return nil, err
	}

	var hosts []*citadel.Host
	for _, x := range results.Node.Nodes {
		hn, err := r.client.Get(x.Key, false, true)
		if err != nil {
			return nil, err
		}
		for _, h := range hn.Node.Nodes {
			if strings.Contains(h.Key, "info") {
				var host *citadel.Host
				if err := json.Unmarshal([]byte(h.Value), &host); err != nil {
					return nil, err
				}
				hosts = append(hosts, host)
			}
		}
	}

	return hosts, nil
}

// DeleteHost deletes a specific host by the host id
func (r *Repository) DeleteHost(id string) error {
	// remove host containers
	key := fmt.Sprintf("%s/hosts/%s", r.namespace, id)
	if _, err := r.client.RawDelete(key, true, true); err != nil {
		return err
	}

	return nil
}
