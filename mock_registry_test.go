package citadel

import "errors"

var (
	mockErrNotFound = errors.New("not found")
)

type mockRegistry struct {
	hosts      map[string]*Host
	apps       map[string]*Application
	containers map[string]*Container
}

func newMockRegistry() *mockRegistry {
	return &mockRegistry{
		hosts:      make(map[string]*Host),
		apps:       make(map[string]*Application),
		containers: make(map[string]*Container),
	}
}

func (r *mockRegistry) SaveApplication(a *Application) error {
	r.apps[a.ID] = a

	return nil
}

func (r *mockRegistry) DeleteApplication(id string) error {
	delete(r.apps, id)

	return nil
}

func (r *mockRegistry) FetchApplications() ([]*Application, error) {
	out := []*Application{}

	for _, v := range r.apps {
		out = append(out, v)
	}

	return out, nil
}

func (r *mockRegistry) FetchApplication(id string) (*Application, error) {
	a, exists := r.apps[id]

	if !exists {
		return nil, mockErrNotFound
	}

	return a, nil
}

func (r *mockRegistry) FetchContainers(id string) ([]*Container, error) {
	out := []*Container{}

	for _, v := range r.containers {
		out = append(out, v)
	}

	return out, nil
}

func (r *mockRegistry) FetchContainer(hostID, id string) (*Container, error) {
	a, exists := r.containers[id]

	if !exists {
		return nil, mockErrNotFound
	}

	return a, nil
}

func (r *mockRegistry) SaveContainer(hostID string, c *Container) error {
	r.containers[c.ID] = c

	return nil
}

func (r *mockRegistry) DeleteContainer(hostID, id string) error {
	return nil
}

func (r *mockRegistry) SaveHost(h *Host) error {
	r.hosts[h.ID] = h

	return nil
}

func (r *mockRegistry) FetchHosts() ([]*Host, error) {
	out := []*Host{}

	for _, v := range r.hosts {
		out = append(out, v)
	}

	return out, nil
}

func (r *mockRegistry) FetchHost(id string) (*Host, error) {
	a, exists := r.hosts[id]

	if !exists {
		return nil, mockErrNotFound
	}

	return a, nil
}

func (r *mockRegistry) DeleteHost(id string) error {
	delete(r.hosts, id)

	return nil
}
