package citadel

import (
	"fmt"
	"strings"

	"github.com/samalba/dockerclient"
)

// Container is a docker container running on a specific engine
type Container struct {
	// Name is the unique container name
	Name string `json:"name,omitempty"`

	// Image is the docker image to base the container off of
	Image string `json:"image,omitempty"`

	// Cpus is the number of cpu resources to give to the container
	Cpus float64 `json:"cpus,omitempty"`

	// Memory is the amount of memory in MB for the container
	Memory float64 `json:"memory,omitempty"`

	// Envionrment is the environment vars to set on the container
	Environment map[string]string `json:"environment,omitempty"`

	// Hostname is the host name to set for the container
	Hostname string `json:"hostname,omitempty"`

	// Domainname is the domain name to set for the container
	Domainname string `json:"domain,omitempty"`

	// Args are cli arguments to pass to the image
	Args []string `json:"args,omitempty"`

	// Type is the container type, often service, batch, etc...
	Type string `json:"type,omitempty"`

	// Labels are matched with constraints on the engines
	Labels []string `json:"labels,omitempty"`

	// UserData is user defined data that is passed to the container
	UserData map[string][]string `json:"user_data,omitempty"`

	// Placement is the host information from the cluster
	Placement *Placement `json:"placement,omitempty"`
}

func (c *Container) Run(e *Engine) error {
	var (
		env    = []string{}
		client = e.client
	)

	for k, v := range c.Environment {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	env = append(env,
		fmt.Sprintf("_citadel_type=%s", c.Type),
		fmt.Sprintf("_citadel_labels=%s", strings.Join(c.Labels, ",")),
	)

	config := &dockerclient.ContainerConfig{
		Hostname:   c.Hostname,
		Domainname: c.Domainname,
		Image:      c.Image,
		Cmd:        c.Args,
		Memory:     int(c.Memory) * 1024 * 1024,
		Env:        env,
		CpuShares:  int(c.Cpus * 100.0 / e.Cpus),
	}

	hostConfig := &dockerclient.HostConfig{
		PublishAllPorts: true,
	}

retry:
	if _, err := client.CreateContainer(config, c.Name); err != nil {
		if err != dockerclient.ErrNotFound {
			return err
		}

		if err := client.PullImage(c.Image, "latest"); err != nil {
			return err
		}

		goto retry
	}

	return client.StartContainer(c.Name, hostConfig)
}

func (c *Container) Kill(e *Engine) error {
	return e.client.KillContainer(c.Name)
}
