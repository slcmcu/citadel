package docker

import (
	"fmt"
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
func New(namespace, dockerUrl string, data *citadel.ServiceData) (citadel.Service, error) {
	client, err := dockerclient.NewDockerClient(dockerUrl)
	if err != nil {
		return nil, err
	}

	return &Service{
		data:      data,
		docker:    client,
		namespace: namespace,
	}, nil
}

func (d *Service) Data() *citadel.ServiceData {
	return d.data
}

func (d *Service) List(t *citadel.Task) ([]*citadel.ServiceData, error) {
	containers, err := d.docker.ListContainers(false)
	if err != nil {
		return nil, err
	}
	out := []*citadel.ServiceData{}

	for _, c := range containers {
		service, err := d.containerToService(c.Id, c.Names[0])
		if err != nil {
			return nil, err
		} else if service != nil {
			out = append(out, service)
		}
	}

	return out, nil
}

func (d *Service) Run(t *citadel.Task) (*citadel.RunResult, error) {
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

func (d *Service) Stop(t *citadel.Task) (*citadel.StopResult, error) {
	if t.Service.Type == "docker" {
		return nil, fmt.Errorf("cannot stop docker service")
	}

	if err := d.docker.StopContainer(t.Service.Name, 5); err != nil {
		return nil, err
	}

	// FIXME:
	return nil, nil
}

func (d *Service) containerToService(id, name string) (*citadel.ServiceData, error) {
	c, err := d.docker.InspectContainer(id)
	if err != nil {
		return nil, err
	}

	return &citadel.ServiceData{
		Type:   c.Image,
		Name:   name,
		Cpus:   stringToCpus(c.Config.Cpuset),
		Memory: c.Config.Memory,
		Addr:   c.NetworkSettings.IpAddress,
	}, nil
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
