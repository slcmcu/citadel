package citadel

import "github.com/samalba/dockerclient"

// SetContainerCpus will take a setting like 1.2 and set the correct cpu shares and
// cpu set on the config
func SetContainerCpus(d *Docker, c *Container, config *dockerclient.ContainerConfig) {
	config.CpuShares = int(c.Cpus * 100.0 / d.Cpus)
}
