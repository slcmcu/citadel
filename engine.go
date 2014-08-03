package citadel

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/samalba/dockerclient"
)

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

func (e *Engine) Run(c *Container) error {
	var (
		err    error
		env    = []string{}
		client = e.client
		i      = c.Image
	)
	c.Engine = e

	for k, v := range i.Environment {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	env = append(env,
		fmt.Sprintf("_citadel_type=%s", i.Type),
		fmt.Sprintf("_citadel_labels=%s", strings.Join(i.Labels, ",")),
	)

	config := &dockerclient.ContainerConfig{
		Hostname:     i.Hostname,
		Domainname:   i.Domainname,
		Image:        i.Name,
		Cmd:          i.Args,
		Memory:       int(i.Memory) * 1024 * 1024,
		Env:          env,
		CpuShares:    int(i.Cpus * 100.0 / e.Cpus),
		ExposedPorts: make(map[string]struct{}),
	}

	hostConfig := &dockerclient.HostConfig{
		PublishAllPorts: len(i.BindPorts) == 0,
		PortBindings:    make(map[string][]dockerclient.PortBinding),
	}

	for _, b := range i.BindPorts {
		key := fmt.Sprintf("%d/%s", b.Port, b.Proto)
		config.ExposedPorts[key] = struct{}{}

		hostConfig.PortBindings[key] = []dockerclient.PortBinding{
			{
				HostPort: fmt.Sprint(b.Port),
			},
		}
	}

retry:
	if c.ID, err = client.CreateContainer(config, ""); err != nil {
		if err != dockerclient.ErrNotFound {
			return err
		}

		if err := client.PullImage(i.Name, "latest"); err != nil {
			return err
		}

		goto retry
	}

	if err := client.StartContainer(c.ID, hostConfig); err != nil {
		return err
	}

	return e.updatePortInformation(c)
}

func (e *Engine) ListImages() ([]string, error) {
	images, err := e.client.ListImages()
	if err != nil {
		return nil, err
	}

	out := []string{}

	for _, i := range images {
		for _, t := range i.RepoTags {
			out = append(out, t)
		}
	}

	return out, nil
}

func (e *Engine) updatePortInformation(c *Container) error {
	info, err := e.client.InspectContainer(c.ID)
	if err != nil {
		return err
	}

	for pp, b := range info.NetworkSettings.Ports {
		parts := strings.Split(pp, "/")
		rawPort, proto := parts[0], parts[1]

		port, err := strconv.Atoi(b[0].HostPort)
		if err != nil {
			return err
		}

		containerPort, err := strconv.Atoi(rawPort)
		if err != nil {
			return err
		}

		c.Ports = append(c.Ports, &Port{
			Proto:         proto,
			Port:          port,
			ContainerPort: containerPort,
		})
	}

	return nil
}

func (d *Engine) loadContainers() error {
	d.containers = containers{}

	c, err := d.client.ListContainers(false)
	if err != nil {
		return err
	}

	for _, ci := range c {
		cc, err := fromDockerContainer(&ci, d)
		if err != nil {
			return err
		}

		d.containers = append(d.containers, cc)
	}

	return nil
}

func (c containers) totalCpuAndMemory() (cpu float64, mem float64) {
	for _, ci := range c {
		cpu += ci.Image.Cpus
		mem += ci.Image.Memory
	}

	return cpu, mem
}
