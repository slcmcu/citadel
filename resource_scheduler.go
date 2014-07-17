package citadel

type ResourceScheduler struct {
	registry Registry
}

func NewResourceScheduler(r Registry) Scheduler {
	return &ResourceScheduler{
		registry: r,
	}
}

func (m *ResourceScheduler) Schedule(app *Application, hosts []*Host) ([]*Host, error) {
	var (
		accepted = []*Host{}
		cpus     = app.totalCpus()
		memory   = app.totalMemory()
	)

	for _, h := range hosts {
		// fast path is to make sure that the host can run this
		// if that passes then we need to make sure that it has
		// enough capacity with it's current load of containers
		if h.Memory >= memory && h.Cpus >= cpus {
			reservedCpus, reservedMemory, err := m.getTotalReservation(h)
			if err != nil {
				return nil, err
			}

			if (h.Cpus-reservedCpus) >= cpus && (h.Memory-reservedMemory) >= memory {
				accepted = append(accepted, h)
			}
		}
	}

	return accepted, nil
}

func (m *ResourceScheduler) getTotalReservation(h *Host) (int, int, error) {
	containers, err := m.registry.FetchContainers(h.ID)
	if err != nil {
		return 0, 0, err
	}

	var cpus, memory int

	for _, c := range containers {
		cpus += len(c.Config.Cpus)
		memory += c.Config.Memory
	}

	return cpus, memory, nil
}
