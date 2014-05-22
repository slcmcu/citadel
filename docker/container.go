package docker

import (
	"citadelapp.io/citadel"
	"github.com/samalba/dockerclient"
)

type ContainerService struct {
	id     string
	client *dockerclient.DockerClient
	data   *citadel.ServiceData
}

func NewContainer(id, name string, d *dockerclient.DockerClient) (citadel.Service, error) {
	c, err := d.InspectContainer(id)
	if err != nil {
		return nil, err
	}

	data := &citadel.ServiceData{
		Type:   c.Image,
		Name:   name,
		Cpus:   stringToCpus(c.Config.Cpuset),
		Memory: c.Config.Memory,
		Addr:   c.NetworkSettings.IpAddress,
	}

	return &ContainerService{
		id:     id,
		client: d,
		data:   data,
	}, nil
}

func (c *ContainerService) Data() *citadel.ServiceData {
	return c.data
}

func (c *ContainerService) List(t *citadel.Task) ([]citadel.Service, error) {
	panic("not implemented")
}

func (c *ContainerService) Run(t *citadel.Task) (interface{}, error) {
	panic("not implemented")
}

func (c *ContainerService) Stop(t *citadel.Task) (interface{}, error) {
	panic("not implemented")
}
