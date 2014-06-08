package repository

import "citadelapp.io/citadel"

func (e *Repository) FetchConfig() (*citadel.Config, error) {
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
