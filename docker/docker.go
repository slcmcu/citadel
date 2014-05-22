package docker

import (
	"fmt"
	"path/filepath"
	"strings"

	"citadelapp.io/citadel"
	"github.com/samalba/dockerclient"
)

type Service struct {
	data      *citadel.ServiceData
	namespace string
	docker    *dockerclient.DockerClient
}

// New returns a new Service interacting with a docker daemon.  The namespace
// is the image namespace required to create new containers
func New(namespace string, data *citadel.ServiceData, docker *dockerclient.DockerClient) citadel.Service {

	return &Service{
		data:      data,
		docker:    docker,
		namespace: namespace,
	}
}

func (d *Service) Data() *citadel.ServiceData {
	return d.data
}

func (d *Service) List(t *citadel.Task) ([]citadel.Service, error) {
	containers, err := d.docker.ListContainers(false)
	if err != nil {
		return nil, err
	}
	out := []citadel.Service{}

	for _, c := range containers {
		var service citadel.Service

		switch t.Service.Type {
		case "docker":
			service, err = d.containerToService(c.Id, c.Names[0])
		default:
			_, name := filepath.Split(c.Image)
			if name == t.Service.Type {
				service, err = d.containerToService(c.Id, c.Names[0])
			}
		}

		if err != nil {
			return nil, err
		} else if service != nil {
			out = append(out, service)
		}
	}

	return out, nil
}

func (d *Service) Run(t *citadel.Task) (interface{}, error) {
	if t.Service.Type == "docker" {
		return nil, fmt.Errorf("cannot start docker service")
	}

	config := &dockerclient.ContainerConfig{
		Image:  fmt.Sprintf("%s/%s:latest", d.namespace, t.Service.Type),
		Memory: t.Service.Memory,
		Cpuset: cpusToString(t.Service.Cpus),
	}

	id, err := d.docker.CreateContainer(config)
	if err != nil {
		return nil, err
	}

	if err := d.docker.StartContainer(id); err != nil {
		return nil, err
	}

	// FIXME: return a report of the run with state information
	return nil, nil
}

func (d *Service) Stop(t *citadel.Task) (interface{}, error) {
	if t.Service.Type == "docker" {
		return nil, fmt.Errorf("cannot stop docker service")
	}

	if err := d.docker.StopContainer(t.Service.Name, 5); err != nil {
		return nil, err
	}

	// FIXME:
	return nil, nil
}

func (d *Service) containerToService(id, name string) (citadel.Service, error) {
	return NewContainer(id, name, d.docker)
}

func cpusToString(cpus int) string {
	out := []string{}
	for i := 0; i < cpus; i++ {
		out = append(out, fmt.Sprint(i))
	}
	return strings.Join(out, ",")
}

func stringToCpus(cpuset string) int {
	return len(strings.Split(cpuset, ","))
}
