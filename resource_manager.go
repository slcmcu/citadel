package citadel

import "log"

// ResourceManager is responsible for managing the engines of the cluster
type ResourceManager struct {
	logger *log.Logger
}

func newEngineManger(logger *log.Logger) *ResourceManager {
	return &ResourceManager{
		logger: logger,
	}
}

// PlaceImage uses the provided engines to make a decision on which resource the container
// should run based on best utilization of the engines.
func (r *ResourceManager) PlaceContainer(c *Container, engines []*Engine) (*Engine, error) {
	scores := []*score{}

	for _, re := range engines {
		if re.Memory < c.Image.Memory || re.Cpus < c.Image.Cpus {
			continue
		}

		cpus, memory := re.containers.totalCpuAndMemory()

		var (
			cpuScore    = ((cpus + c.Image.Cpus) / re.Cpus) * 100.0
			memoryScore = ((memory + c.Image.Memory) / re.Memory) * 100.0
			total       = ((cpuScore + memoryScore) / 200.0) * 100.0
		)

		r.logger.Printf("resource=%s score=%f\n", re.ID, total)

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
