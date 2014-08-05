package scheduler

import (
	"math/rand"
	"strings"
	"time"

	"github.com/citadel/citadel"
)

type HostScheduler struct {
}

func (h *HostScheduler) Schedule(c *citadel.Image, e *citadel.Engine) (bool, error) {
	if len(c.Labels) > 0 && h.validHost(e, c.Labels) {
		return true, nil
	}
	return false, nil
}

func (h *HostScheduler) validHost(e *citadel.Engine, labels []string) bool {
	// randomize lookup
	rand.Seed(time.Now().UTC().UnixNano())
	i := rand.Intn(len(labels) - 1)
	label := labels[i]
	parts := strings.Split(label, "host:")
	if len(parts) != 2 {
		return false
	}
	host := parts[1]
	if e.ID == host {
		return true
	}
	return false
}
