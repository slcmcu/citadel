package citadel

import (
	"strconv"
	"strings"

	"github.com/samalba/dockerclient"
)

type (
	ImageInfo struct {
		Name string
		Tag  string
	}
)

func parsePortInformation(info *dockerclient.ContainerInfo, c *Container) error {
	for pp, b := range info.NetworkSettings.Ports {
		parts := strings.Split(pp, "/")
		rawPort, proto := parts[0], parts[1]

		for _, binding := range b {
			port, err := strconv.Atoi(binding.HostPort)
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
	}

	return nil
}

func FromDockerContainer(id, image string, engine *Engine) (*Container, error) {
	info, err := engine.client.InspectContainer(id)
	if err != nil {
		return nil, err
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

	container := &Container{
		ID:     id,
		Engine: engine,
		Image: &Image{
			Name:        image,
			Cpus:        float64(info.Config.CpuShares) / 100.0 * engine.Cpus,
			Memory:      float64(info.Config.Memory / 1024 / 1024),
			Environment: env,
			Hostname:    info.Config.Hostname,
			Domainname:  info.Config.Domainname,
			Type:        cType,
			Labels:      labels,
		},
	}

	if err := parsePortInformation(info, container); err != nil {
		return nil, err
	}

	return container, nil
}

func parseImageName(name string) *ImageInfo {
	imageInfo := &ImageInfo{
		Name: name,
		Tag:  "latest",
	}
	img := strings.Split(name, ":")
	if len(img) == 2 {
		imageInfo.Name = img[0]
		imageInfo.Tag = img[1]
	}
	return imageInfo
}
