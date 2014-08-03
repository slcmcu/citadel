package citadel

import (
	"fmt"
	"strings"

	"github.com/samalba/dockerclient"
)

// ValidateImage ensures that the required fields are set on the container
func ValidateImage(c *Image) error {
	switch {
	case c.Cpus == 0:
		return fmt.Errorf("container cannot have cpus equal to 0")
	case c.Memory == 0:
		return fmt.Errorf("container cannot have memory equal to 0")
	case c.Image == "":
		return fmt.Errorf("container must have an image")
	case c.Type == "":
		return fmt.Errorf("container must have a type")
	}

	return nil
}

func fromDockerContainer(container *dockerclient.Container, engine *Engine) (*Container, error) {
	info, err := engine.client.InspectContainer(container.Id)
	if err != nil {
		return nil, err
	}

	var ports []*Port
	for _, port := range container.Ports {
		p := &Port{
			Proto:     port.Type,
			Port:      port.PublicPort,
			ImagePort: port.PrivatePort,
		}
		ports = append(ports, p)
	}

	var (
		cType  = ""
		labels = []string{}
		env    = make(map[string]string)
	)

	for _, e := range info.Config.Env {
		vals := strings.Split(e, "=")
		k, v := vals[0], vals[1]

		switch k {
		case "_citadel_type":
			cType = v
		case "_citadel_labels":
			labels = strings.Split(v, ",")
		case "HOME", "DEBIAN_FRONTEND", "PATH":
			continue
		default:
			env[k] = v
		}
	}

	return &Container{
		ID:     container.Id,
		Engine: engine,
		Ports:  ports,
		Image: &Image{
			Image:       container.Image,
			Cpus:        float64(info.Config.CpuShares) / 100.0 * engine.Cpus,
			Memory:      float64(info.Config.Memory / 1024 / 1024),
			Environment: env,
			Hostname:    info.Config.Hostname,
			Domainname:  info.Config.Domainname,
			Type:        cType,
			Labels:      labels,
		},
	}, nil
}
