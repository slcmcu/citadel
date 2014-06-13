package repository

import (
	"citadelapp.io/citadel"
	"github.com/dancannon/gorethink"
)

// SaveHost will save a single host with the unique id
func (r *Repository) SaveHost(h *citadel.Host) error {
	if _, err := gorethink.Table("hosts").Insert(h).RunWrite(r.session); err != nil {
		return err
	}

	return nil
}

// FetchHosts returns all the hosts in the cluseter
func (r *Repository) FetchHosts() ([]*citadel.Host, error) {
	results, err := gorethink.Table("hosts").Run(r.session)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	var hosts []*citadel.Host
	if err := results.ScanAll(&hosts); err != nil {
		return nil, err
	}

	return hosts, nil
}

// DeleteHost deletes a specific host by the host id
func (r *Repository) DeleteHost(id string) error {
	// remove host containers
	if _, err := gorethink.Table("containers").Filter(gorethink.Row.Field("host_id").Eq(id)).Delete().Run(r.session); err != nil {
		return err
	}
	// remove host
	if _, err := gorethink.Table("hosts").Get(id).Delete().Run(r.session); err != nil {
		return err
	}

	return nil
}
