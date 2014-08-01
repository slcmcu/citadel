package citadel

import (
	"fmt"
	"strings"

	"github.com/samalba/dockerclient"
)

// ImageScheduler only returns engines that already have the image pulled
// locally on disk for docker to use
type ImageScheduler struct {
}

func (i *ImageScheduler) Schedule(engines []*Docker, c *Container) ([]*Docker, error) {
	out := []*Docker{}

	fullImage := c.Image
	if !strings.Contains(fullImage, ":") {
		fullImage = fmt.Sprintf("%s:latest", fullImage)
	}

	for _, e := range engines {
		images, err := e.Client.ListImages()
		if err != nil {
			return nil, err
		}

		if i.containsImage(fullImage, images) {
			out = append(out, e)
		}
	}

	return out, nil
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
