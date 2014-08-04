package citadel

import (
	"strings"

	"github.com/samalba/dockerclient"
)

func fromDockerContainer(container *dockerclient.Container, engine *Engine) (*Container, error) {
	info, err := engine.client.InspectContainer(container.Id)
	if err != nil {
		return nil, err
	}

	var ports []*Port
	for _, port := range container.Ports {
		p := &Port{
			Proto:         port.Type,
			Port:          port.PublicPort,
			ContainerPort: port.PrivatePort,
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
			Name:        container.Image,
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
