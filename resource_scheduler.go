package citadel

type ResourceScheduler struct {
	registry Registry
}

func NewResourceScheduler(r Registry) Scheduler {
	return &ResourceScheduler{
		registry: r,
	}
}

func (m *ResourceScheduler) Schedule(app *Application, hosts []*Host) (*Host, error) {
	var (
		scores = []*score{}
		cpus   = app.totalCpus()
		memory = app.totalMemory()
	)

	for _, h := range hosts {
		// fast path is to make sure that the host can run this
		// if that passes then we need to make sure that it has
		// enough capacity with it's current load of containers
		if h.Memory >= memory && float64(h.Cpus) >= cpus {
			reservedCpus, reservedMemory, err := m.getTotalReservation(h)
			if err != nil {
				return nil, err
			}

			var (
				cpuScore    = (float64(reservedCpus) / float64(h.Cpus)) * 100.0
				memoryScore = (float64(reservedMemory) / float64(h.Memory)) * 100.0
				totalScore  = ((cpuScore + memoryScore) / 200.0) * 100.0
			)

			if totalScore <= 100.0 {
				scores = append(scores, &score{h: h, score: totalScore})
			}
		}
	}

	if len(scores) == 0 {
		return nil, ErrNoValidHost
	}

	sortScores(scores)

	return scores[0].h, nil
}

func (m *ResourceScheduler) getTotalReservation(h *Host) (float64, int, error) {
	containers, err := m.registry.FetchContainers(h.ID)
	if err != nil {
		return 0, 0, err
	}

	var (
		cpus   float64
		memory int
	)

	for _, c := range containers {
		cpus += c.Config.Cpus
		memory += c.Config.Memory
	}

	return cpus, memory, nil
}
