package citadel

type LabelScheduler struct {
}

func (l *LabelScheduler) Schedule(c *Container, e *Docker) (bool, error) {
	if len(c.Labels) == 0 || l.contains(e, c.Labels) {
		return true, nil
	}

	return false, nil
}

func (l *LabelScheduler) contains(r *Docker, constraints []string) bool {
	for _, c := range constraints {
		if !l.resourceContains(r, c) {
			return false
		}
	}

	return true
}

func (l *LabelScheduler) resourceContains(r *Docker, c string) bool {
	for _, l := range r.Labels {
		if l == c {
			return true
		}
	}

	return false
}
