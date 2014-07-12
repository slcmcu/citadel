package citadel

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/citadel/citadel/utils"
	"github.com/samalba/dockerclient"
)

func containerFromDocker(h *Host, info *dockerclient.ContainerInfo) (*Container, error) {
	c := &Container{
		ID:     strings.TrimPrefix(info.Name, "/"),
		Image:  info.Image,
		HostID: h.ID,
		Cpus:   utils.CpusetTOI(info.Config.Cpuset),
	}

	// if cpuset is not specified then the container has all the cpus on the host
	if len(c.Cpus) == 0 {
		for i := 0; i < h.Cpus; i++ {
			c.Cpus = append(c.Cpus, i)
		}
	}

	if info.Config.Memory > 0 {
		c.Memory = info.Config.Memory / 1024 / 1024
	}

	if info.State.Running {
		c.State.Status = Running
	} else {
		c.State.Status = Stopped
	}

	c.State.ExitCode = info.State.ExitCode

	if err := copyPorts(info, c); err != nil {
		return nil, err
	}

	return c, nil
}

func copyPorts(info *dockerclient.ContainerInfo, c *Container) error {
	if info.HostConfig != nil && info.HostConfig.PortBindings != nil {
		for cp, bindings := range info.HostConfig.PortBindings {
			var (
				container int
				proto     string
			)

			if _, err := fmt.Sscanf(cp, "%d/%s", &container, &proto); err != nil {
				return err
			}

			for _, b := range bindings {
				hostPort, err := strconv.Atoi(b.HostPort)
				if err != nil {
					return err
				}

				c.Ports = append(c.Ports, &Port{
					Proto:     proto,
					Container: container,
					Host:      hostPort,
				})
			}
		}
	}

	return nil
}
