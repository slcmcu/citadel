package citadel

import "log"

// DockerManager is responsible for managing the resources of the cluster
type DockerManager struct {
	logger *log.Logger
}

func newDockerManger(logger *log.Logger) *DockerManager {
	return &DockerManager{
		logger: logger,
	}
}

// PlaceContainer uses the provided resources to make a decision on which resource the container
// should run based on best utilization of the resources.
func (r *DockerManager) PlaceContainer(resources []*Docker, c *Container) (*Docker, error) {
	var (
		scores = []*score{}
	)

	for _, re := range resources {
		if re.Memory < c.Memory || re.Cpus < c.Cpus {
			continue
		}

		cpus, memory := re.ReservedCpus, re.ReservedMemory

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
