package scheduler

import (
	"fmt"
	"strings"

	"github.com/citadel/citadel"
	"github.com/samalba/dockerclient"
)

// ImageScheduler only returns engines that already have the image pulled
// locally on disk for docker to use
type ImageScheduler struct {
}

func (i *ImageScheduler) Schedule(c *citadel.Image, e *citadel.Engine) (bool, error) {
	fullImage := c.Name

	if !strings.Contains(fullImage, ":") {
		fullImage = fmt.Sprintf("%s:latest", fullImage)
	}

	images, err := e.client.ListImages()
	if err != nil {
		return false, err
	}

	if i.containsImage(fullImage, images) {
		return true, nil
	}

	return false, nil
}

func (i *ImageScheduler) containsImage(requested string, images []*dockerclient.Image) bool {
	for _, img := range images {
		for _, tag := range img.RepoTags {
			if requested == tag {
				return true
			}
		}
	}

	return false
}
