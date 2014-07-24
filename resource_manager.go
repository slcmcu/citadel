package citadel

// ResourceManager is responsible for managing the resources of the cluster
type ResourceManager struct {
	registry Registry
}

type placement struct {
	r     *Resource
	Score float64
}

func newResourceManger(registry Registry) *ResourceManager {
	return &ResourceManager{
		registry: registry,
	}
}

// PlaceContainer uses the provided resources to make a decision on which resource the container
// should run based on best utilization of the resources.
func (r *ResourceManager) PlaceContainer(resources []*Resource, c *Container) (*placement, error) {
	return nil, nil
}
