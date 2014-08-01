package citadel

type LabelScheduler struct {
}

func (l *LabelScheduler) Schedule(engines []*Docker, c *Container) ([]*Docker, error) {
	if len(c.Labels) == 0 {
		return engines, nil
	}

	out := []*Docker{}

	for _, r := range engines {
		if l.contains(r, c.Labels) {
			out = append(out, r)
		}
	}

	return out, nil
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
