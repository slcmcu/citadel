package citadel

type LabelScheduler struct {
	registry Registry
}

func NewLabelScheduler(r Registry) Scheduler {
	return &LabelScheduler{
		registry: r,
	}
}

func (l *LabelScheduler) Schedule(app *Application, hosts []*Host) (*Host, error) {
	if len(app.Constraints) == 0 {
		return hosts[0], nil
	}

	accepted := []*Host{}

	for _, h := range hosts {
		if l.contains(h, app.Constraints) {
			accepted = append(accepted, h)
		}
	}

	if len(accepted) == 0 {
		return nil, ErrNoValidHost
	}

	return accepted[0], nil
}

func (l *LabelScheduler) contains(h *Host, constraints []string) bool {
	for _, c := range constraints {
		if !l.hostContains(h, c) {
			return false
		}
	}

	return true
}

func (l *LabelScheduler) hostContains(h *Host, c string) bool {
	for _, l := range h.Labels {
		if l == c {
			return true
		}
	}

	return false
}
