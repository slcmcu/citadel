package citadel

import (
	"fmt"

	"github.com/samalba/dockerclient"
)

// SetContainerCpus will take a setting like 1.2 and set the correct cpu shares and
// cpu set on the config
func SetContainerCpus(d *Docker, c *Container, config *dockerclient.ContainerConfig) {
	config.CpuShares = int(c.Cpus * 100.0 / d.Cpus)
}

// ValidateContainer ensures that the required fields are set on the container
func ValidateContainer(c *Container) error {
	switch {
	case c.Cpus == 0:
		return fmt.Errorf("container cannot have cpus equal to 0")
	case c.Memory == 0:
		return fmt.Errorf("container cannot have memory equal to 0")
	case c.Image == "":
		return fmt.Errorf("container must have an image")
	case c.Name == "":
		return fmt.Errorf("container must have a name")
	case c.Type == "":
		return fmt.Errorf("container must have a type")
	}

	return nil
}
