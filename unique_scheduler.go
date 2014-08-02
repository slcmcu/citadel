package citadel

// UniqueScheduler only returns engines that do not have the image running
type UniqueScheduler struct {
}

func (u *UniqueScheduler) Schedule(c *Container, e *Engine) (bool, error) {
	if u.hasContainer(c, e.containers) {
		return false, nil
	}
	return true, nil
}

func (u *UniqueScheduler) hasContainer(container *Container, containers []*Container) bool {
	for _, c := range containers {
		if c.Name == container.Name {
			return true
		}
	}
	return false
}
