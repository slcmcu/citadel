package citadel

import "github.com/samalba/dockerclient"

type Docker struct {
	ID     string   `json:"id,omitempty"`
	Addr   string   `json:"addr,omitempty"`
	Cpus   float64  `json:"cpus,omitempty"`
	Memory float64  `json:"memory,omitempty"`
	Labels []string `json:"labels,omitempty"`

	client *dockerclient.DockerClient
}

func (d *Docker) SetClient(client *dockerclient.DockerClient) {
	d.client = client
}

func (d *Docker) GetCpuAndMemoryReservation() (cpu float64, mem float64, err error) {
	containers, err := d.client.ListContainers(false)
	if err != nil {
		return 0, 0, err
	}

	for _, ci := range containers {
		info, err := d.client.InspectContainer(ci.Id)
		if err != nil {
			return 0, 0, err
		}

		cpu += float64(info.Config.CpuShares) / 100.0 * d.Cpus
		mem += float64(info.Config.Memory / 1024 / 1024)
	}

	return cpu, mem, nil
}
