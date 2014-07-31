package citadel

import "log"

// ResourceManager is responsible for managing the resources of the cluster
type ResourceManager struct {
	registry Registry
	logger   *log.Logger
}

func newResourceManger(registry Registry, logger *log.Logger) *ResourceManager {
	return &ResourceManager{
		logger:   logger,
		registry: registry,
	}
}

// PlaceContainer uses the provided resources to make a decision on which resource the container
// should run based on best utilization of the resources.
func (r *ResourceManager) PlaceContainer(resources []*Resource, c *Container) (*Resource, error) {
	var (
		scores = []*score{}
	)

	for _, re := range resources {
		if re.Memory < c.Memory || re.Cpus < c.Cpus {
			continue
		}

		cpus, memory, err := r.registry.GetTotalReservations(re.ID)
		if err != nil {
			return nil, err
		}

		var (
			cpuScore    = ((cpus + c.Cpus) / re.Cpus) * 100.0
			memoryScore = ((memory + c.Memory) / re.Memory) * 100.0
			total       = ((cpuScore + memoryScore) / 200.0) * 100.0
		)

		log.Printf("resource=%s score=%f\n", re.ID, total)
		if total <= 100.0 {
			scores = append(scores, &score{r: re, score: total})
		}
	}

	if len(scores) == 0 {
		return nil, ErrUnableToSchedule
	}

	sortScores(scores)

	return scores[0].r, nil
}
