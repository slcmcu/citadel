package main

import (
	"strings"
)

type Vol struct{}

func parseVolumes(volumes string) ([]string, map[string]struct{}) {
	binds := []string{}
	vols := make(map[string]struct{})
	v := strings.Split(volumes, " ")
	for _, vol := range v {
		// bind
		if strings.Contains(vol, ":") {
			binds = append(binds, vol)
		} else {
			vols[vol] = Vol{}

		}
	}
	return binds, vols
}
