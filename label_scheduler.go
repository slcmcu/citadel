package citadel

type LabelScheduler struct {
}

func (l *LabelScheduler) Schedule(t *Transaction) error {
	if len(t.Container.Labels) == 0 {
		return nil
	}

	accpeted := []*Docker{}

	for _, r := range t.GetEngines() {
		if l.contains(r, t.Container.Labels) {
			accpeted = append(accpeted, r)
		}
	}

	t.Reduce(accpeted)

	return nil
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
