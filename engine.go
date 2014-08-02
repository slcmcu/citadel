package citadel

import "github.com/samalba/dockerclient"

type containers []*Container

type Engine struct {
	ID     string   `json:"id,omitempty"`
	Addr   string   `json:"addr,omitempty"`
	Cpus   float64  `json:"cpus,omitempty"`
	Memory float64  `json:"memory,omitempty"`
	Labels []string `json:"labels,omitempty"`

	client     *dockerclient.DockerClient
	containers containers
}

func (d *Engine) SetClient(client *dockerclient.DockerClient) {
	d.client = client
}

func (d *Engine) Containers() containers {
	return d.containers
}

func (d *Engine) loadContainers() error {
	d.containers = containers{}

	c, err := d.client.ListContainers(false)
	if err != nil {
		return err
	}

	for _, ci := range c {
		cc, err := asCitadelContainer(&ci, d)
		if err != nil {
			return err
		}

		d.containers = append(d.containers, cc)
	}

	return nil
}

func (c containers) totalCpuAndMemory() (cpu float64, mem float64) {
	for _, ci := range c {
		cpu += ci.Cpus
		mem += ci.Memory
	}

	return cpu, mem
}
