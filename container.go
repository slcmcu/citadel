package citadel

import (
	"fmt"
	"strings"

	"github.com/samalba/dockerclient"
)

// Container is a docker container running on a specific engine
type Container struct {
	// ID is the container's id
	// This is set by citadel
	ID string `json:"id,omitempty"`

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

	// BindPorts ensures that the container has exclusive access to the specified ports
	BindPorts []*Port `json:"bind_ports,omitempty"`

	// UserData is user defined data that is passed to the container
	UserData map[string][]string `json:"user_data,omitempty"`

	// Placement is the host information from the cluster
	// This is set by citadel
	Placement *Placement `json:"placement,omitempty"`
}

func (c *Container) Run(e *Engine) (string, error) {
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
		Hostname:     c.Hostname,
		Domainname:   c.Domainname,
		Image:        c.Image,
		Cmd:          c.Args,
		Memory:       int(c.Memory) * 1024 * 1024,
		Env:          env,
		CpuShares:    int(c.Cpus * 100.0 / e.Cpus),
		ExposedPorts: make(map[string]struct{}),
	}

	hostConfig := &dockerclient.HostConfig{
		PublishAllPorts: len(c.BindPorts) == 0,
		PortBindings:    make(map[string][]dockerclient.PortBinding),
	}

	for _, b := range c.BindPorts {
		key := fmt.Sprintf("%d/%s", b.Port, b.Proto)
		config.ExposedPorts[key] = struct{}{}

		hostConfig.PortBindings[key] = []dockerclient.PortBinding{
			{
				HostPort: fmt.Sprint(b.Port),
			},
		}
	}

retry:
	cc, err := client.CreateContainer(config, "")
	if err != nil {
		if err != dockerclient.ErrNotFound {
			return "", err
		}

		if err := client.PullImage(c.Image, "latest"); err != nil {
			return "", err
		}

		goto retry
	}

	return cc, client.StartContainer(cc, hostConfig)
}
