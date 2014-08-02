package citadel

import "log"

// ResourceManager is responsible for managing the engines of the cluster
type ResourceManager struct {
	logger *log.Logger
}

func newDockerManger(logger *log.Logger) *ResourceManager {
	return &ResourceManager{
		logger: logger,
	}
}

// PlaceContainer uses the provided engines to make a decision on which resource the container
// should run based on best utilization of the engines.
func (r *ResourceManager) PlaceContainer(t *Transaction) error {
	var (
		scores = []*score{}
	)

	for _, re := range t.engines {
		if re.Memory < t.Container.Memory || re.Cpus < t.Container.Cpus {
			continue
		}

		cpus, memory := re.containers.totalCpuAndMemory()

		var (
			cpuScore    = ((cpus + t.Container.Cpus) / re.Cpus) * 100.0
			memoryScore = ((memory + t.Container.Memory) / re.Memory) * 100.0
			total       = ((cpuScore + memoryScore) / 200.0) * 100.0
		)

		r.logger.Printf("resource=%s score=%f\n", re.ID, total)

		if total <= 100.0 {
			scores = append(scores, &score{r: re, score: total})
		}
	}

	if len(scores) == 0 {
		return ErrUnableToSchedule
	}

	sortScores(scores)

	t.Placement = &Placement{
		Engine: scores[0].r,
	}

	return nil
}
