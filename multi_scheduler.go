package citadel

type MultiScheduler struct {
	schedulers []Scheduler
}

func NewMultiScheduler(s ...Scheduler) Scheduler {
	return &MultiScheduler{
		schedulers: s,
	}
}

func (m *MultiScheduler) Schedule(t *Transaction) error {
	for _, s := range m.schedulers {
		if err := s.Schedule(t); err != nil {
			return err
		}
	}

	return nil
}
