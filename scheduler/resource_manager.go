package scheduler

import (
	"fmt"

	"github.com/citadel/citadel"
)

// ResourceManager is responsible for managing the engines of the cluster
type ResourceManager struct {
}

func NewResourceManager() *ResourceManager {
	return &ResourceManager{}
}

// PlaceImage uses the provided engines to make a decision on which resource the container
// should run based on best utilization of the engines.
func (r *ResourceManager) PlaceContainer(c *citadel.Container, engines []*citadel.Engine) (*citadel.Engine, error) {
	scores := []*score{}

	for _, re := range engines {
		if re.Memory < c.Image.Memory || re.Cpus < c.Image.Cpus {
			continue
		}

		containers, err := re.ListContainers()
		if err != nil {
			return nil, err
		}
		cpus, memory := r.totalCpuAndMemory(containers)

		var (
			cpuScore    = ((cpus + c.Image.Cpus) / re.Cpus) * 100.0
			memoryScore = ((memory + c.Image.Memory) / re.Memory) * 100.0
			total       = ((cpuScore + memoryScore) / 200.0) * 100.0
		)

		if total <= 100.0 {
			scores = append(scores, &score{r: re, score: total})
		}
	}

	if len(scores) == 0 {
		return nil, fmt.Errorf("no resources avaliable to schedule container")
	}

	sortScores(scores)

	return scores[0].r, nil
}

func (r *ResourceManager) totalCpuAndMemory(containers []*citadel.Container) (cpus float64, memory float64) {
	for _, c := range containers {
		cpus += c.Image.Cpus
		memory += c.Image.Memory
	}

	return
}
