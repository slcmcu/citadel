package citadel

import (
	"fmt"
	"strings"
)

// UniqueScheduler only returns engines that do not have the image running
type UniqueScheduler struct {
}

func (u *UniqueScheduler) Schedule(c *Image, e *Engine) (bool, error) {
	if u.hasImage(c, e.containers) {
		return false, nil
	}

	return true, nil
}

func (u *UniqueScheduler) hasImage(i *Image, containers []*Container) bool {
	fullImage := i.Name

	if !strings.Contains(fullImage, ":") {
		fullImage = fmt.Sprintf("%s:latest", fullImage)
	}

	for _, c := range containers {
		if c.Image.Name == fullImage {
			return true
		}
	}

	return false
}
