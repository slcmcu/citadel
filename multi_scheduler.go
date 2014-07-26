package citadel

type MultiScheduler struct {
	schedulers []Scheduler
}

func NewMultiScheduler(s ...Scheduler) Scheduler {
	return &MultiScheduler{
		schedulers: s,
	}
}

func (m *MultiScheduler) Schedule(resources []*Resource, c *Container) ([]*Resource, error) {
	var (
		err  error
		prev = resources
	)

	for _, s := range m.schedulers {
		if prev, err = s.Schedule(prev, c); err != nil {
			return nil, err
		}
	}

	return prev, nil
}
