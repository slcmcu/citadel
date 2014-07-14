package citadel

import (
	"fmt"
	"strconv"

	"github.com/samalba/dockerclient"
)

func createPorts(info *dockerclient.ContainerInfo) ([]*Port, error) {
	out := []*Port{}

	if info.HostConfig != nil && info.HostConfig.PortBindings != nil {
		for cp, bindings := range info.HostConfig.PortBindings {
			var (
				container int
				proto     string
			)

			if _, err := fmt.Sscanf(cp, "%d/%s", &container, &proto); err != nil {
				return nil, err
			}

			for _, b := range bindings {
				hostPort, err := strconv.Atoi(b.HostPort)
				if err != nil {
					return nil, err
				}

				out = append(out, &Port{
					Proto:     proto,
					Container: container,
					Host:      hostPort,
				})
			}
		}
	}

	return out, nil
}
