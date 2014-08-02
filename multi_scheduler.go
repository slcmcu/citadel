package citadel

type MultiScheduler struct {
	schedulers []Scheduler
}

func NewMultiScheduler(s ...Scheduler) Scheduler {
	return &MultiScheduler{
		schedulers: s,
	}
}

func (m *MultiScheduler) Schedule(c *Container, e *Engine) (bool, error) {
	for _, s := range m.schedulers {
		canrun, err := s.Schedule(c, e)
		if err != nil {
			return false, err
		}

		if !canrun {
			return false, nil
		}
	}

	return true, nil
}
