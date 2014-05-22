package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"path"

	"citadelapp.io/citadel"
)

type ClientService struct {
	data   *citadel.ServiceData
	client *http.Client
}

func NewClient(data *citadel.ServiceData) citadel.Service {
	return &ClientService{
		data:   data,
		client: &http.Client{},
	}
}

func (c *ClientService) Data() *citadel.ServiceData {
	return c.data
}

func (c *ClientService) List(t *citadel.Task) ([]*citadel.ServiceData, error) {
	r, err := c.newRequest(t, "/")
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out []*citadel.ServiceData
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	return out, nil
}

func (c *ClientService) Run(t *citadel.Task) (*citadel.RunResult, error) {
	r, err := c.newRequest(t, "/run")
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result *citadel.RunResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *ClientService) Stop(t *citadel.Task) (*citadel.StopResult, error) {
	r, err := c.newRequest(t, "/stop")
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result *citadel.StopResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *ClientService) newRequest(task *citadel.Task, p string) (*http.Request, error) {
	var data io.Reader

	if task != nil {
		buf := bytes.NewBuffer(nil)
		if err := json.NewEncoder(buf).Encode(task); err != nil {
			return nil, err
		}
		data = buf
	}

	return http.NewRequest("POST", path.Join(c.data.Addr, p), data)
}
