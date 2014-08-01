package citadel

import "fmt"

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
