package repository

import (
	"citadelapp.io/citadel"
	"github.com/dancannon/gorethink"
)

func (r *Repository) SaveContainer(c *citadel.Container) error {
	if _, err := gorethink.Table("containers").Insert(c).RunWrite(r.session); err != nil {
		return err
	}

	return nil
}

func (r *Repository) FetchContainers() ([]*citadel.Container, error) {
	results, err := gorethink.Table("containers").Run(r.session)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	var containers []*citadel.Container
	if err := results.ScanAll(&containers); err != nil {
		return nil, err
	}

	return containers, nil
}

func (r *Repository) DeleteContainer(id string) error {
	if _, err := gorethink.Table("containers").Get(id).Delete().Run(r.session); err != nil {
		return err
	}

	return nil
}
