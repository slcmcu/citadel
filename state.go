package citadel

import "net"

// State contains information about the running state of a container
type State struct {
	Slave     string     `json:"slave,omitempty"`     // Save UUID
	Container *Container `json:"container,omitempty"` // Container
	IP        net.IP     `json:"ip,omitempty"`        // IP of the running container
}

type States map[string]*State

func (s States) Len() int {
	return len(s)
}

func (s States) Cpus() int {
	i := 0
	for _, c := range s {
		i += c.Container.Cpus
	}
	return i
}

func (s States) Memory() float64 {
	i := 0.0
	for _, c := range s {
		i += c.Container.Memory
	}
	return i
}

func (s States) Volumes() []string {
	out := []string{}
	for _, c := range s {
		for k := range c.Container.Volumes {
			out = append(out, k)
		}
	}
	return out
}
