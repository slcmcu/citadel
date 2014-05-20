package dockercommand

import (
	"fmt"
	"path/filepath"
	"strings"

	"citadelapp.io/citadel"
	"github.com/samalba/dockerclient"
)

type DockerCommand struct {
	namespace string
	docker    *dockerclient.DockerClient
}

// New returns a new DockerCommand interacting with a docker daemon.  The namespace
// is the image namespace required to create new containers
func New(namespace string, docker *dockerclient.DockerClient) *DockerCommand {
	return &DockerCommand{
		docker:    docker,
		namespace: namespace,
	}
}

func (d *DockerCommand) List(s *citadel.Service) ([]*citadel.Service, error) {
	containers, err := d.docker.ListContainers(false)
	if err != nil {
		return nil, err
	}
	out := []*citadel.Service{}

	for _, c := range containers {
		var service *citadel.Service

		switch s.Type {
		case "docker":
			service, err = d.containerToService(c.Id, c.Names[0])
		default:
			_, name := filepath.Split(c.Image)
			if name == s.Type {
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

func (d *DockerCommand) Start(s *citadel.Service) error {
	if s.Type == "docker" {
		return fmt.Errorf("cannot start docker service")
	}

	config := &dockerclient.ContainerConfig{
		Image:  fmt.Sprintf("%s/%s:latest", d.namespace, s.Type),
		Memory: s.Memory,
		Cpuset: cpusToString(s.Cpus),
	}

	id, err := d.docker.CreateContainer(config)
	if err != nil {
		return err
	}
	return d.docker.StartContainer(id)
}

func (d *DockerCommand) Stop(s *citadel.Service) error {
	if s.Type == "docker" {
		return fmt.Errorf("cannot stop docker service")
	}
	return d.docker.StopContainer(s.Name, 5)
}

func (d *DockerCommand) containerToService(id, name string) (*citadel.Service, error) {
	data, err := d.docker.InspectContainer(id)
	if err != nil {
		return nil, err
	}

	return &citadel.Service{
		Type:   data.Image,
		Name:   name,
		Cpus:   stringToCpus(data.Config.Cpuset),
		Memory: data.Config.Memory,
		Addr:   data.NetworkSettings.IpAddress,
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
