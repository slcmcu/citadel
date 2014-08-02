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

func (i *ImageScheduler) Schedule(t *Transaction) error {
	var (
		accpeted  = []*Docker{}
		fullImage = t.Container.Image
	)

	if !strings.Contains(fullImage, ":") {
		fullImage = fmt.Sprintf("%s:latest", fullImage)
	}

	for _, e := range t.GetEngines() {
		images, err := e.client.ListImages()
		if err != nil {
			return err
		}

		if i.containsImage(fullImage, images) {
			accpeted = append(accpeted, e)
		}
	}

	t.Reduce(accpeted)

	return nil
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
