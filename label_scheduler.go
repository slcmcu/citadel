package citadel

type LabelScheduler struct {
}

func (l *LabelScheduler) Schedule(resources []*Resource, c *Container) ([]*Resource, error) {
	if len(c.Labels) == 0 {
		return resources, nil
	}

	out := []*Resource{}

	for _, r := range resources {
		if l.contains(r, c.Labels) {
			out = append(out, r)
		}
	}

	return out, nil
}

func (l *LabelScheduler) contains(r *Resource, constraints []string) bool {
	for _, c := range constraints {
		if !l.resourceContains(r, c) {
			return false
		}
	}

	return true
}

func (l *LabelScheduler) resourceContains(r *Resource, c string) bool {
	for _, l := range r.Labels {
		if l == c {
			return true
		}
	}

	return false
}
