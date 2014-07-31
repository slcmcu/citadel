package docker

import (
	"crypto/tls"

	"github.com/citadel/citadel"
	"github.com/samalba/dockerclient"
)

type Registry struct {
	resources map[string]*client
}

type client struct {
	r      *citadel.Docker
	client *dockerclient.DockerClient
}

func New(resources []*citadel.Docker, config *tls.Config) (citadel.Registry, error) {
	r := &Registry{
		resources: make(map[string]*client, len(resources)),
	}

	for _, rr := range resources {
		c, err := dockerclient.NewDockerClient(rr.Addr, config)
		if err != nil {
			return nil, err
		}

		r.resources[rr.ID] = &client{r: rr, client: c}
	}

	return r, nil
}

func (r *Registry) FetchDockers() ([]*citadel.Docker, error) {
	out := []*citadel.Docker{}

	for _, c := range r.resources {
		containers, err := c.client.ListContainers(false)
		if err != nil {
			return nil, err
		}

		for _, ci := range containers {
			info, err := c.client.InspectContainer(ci.Id)
			if err != nil {
				return nil, err
			}

			c.r.ReservedCpus += float64(info.Config.CpuShares)
			c.r.ReservedMemory += float64(info.Config.Memory / 1024 / 1024)
		}

		out = append(out, c.r)
	}

	return out, nil
}

func (r *Registry) Close() error {
	return nil
}
