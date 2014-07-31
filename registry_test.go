package citadel

import "fmt"

type mockRegistry struct {
	resources map[string]*reservation
}

type reservation struct {
	r      *Docker
	cpus   float64
	memory float64
}

func newMockRegistry(rs []*Docker) *mockRegistry {
	m := &mockRegistry{
		resources: make(map[string]*reservation, len(rs)),
	}

	for _, r := range rs {
		m.resources[r.ID] = &reservation{r: r}
	}

	return m
}

func (r *mockRegistry) FetchDockers() ([]*Docker, error) {
	out := []*Docker{}

	for _, rs := range r.resources {
		out = append(out, rs.r)
	}

	return out, nil
}

func (r *mockRegistry) GetTotalReservations(id string) (float64, float64, error) {
	rs := r.resources[id]
	if rs == nil {
		return 0, 0, fmt.Errorf("no resource for id %s", id)
	}

	return rs.cpus, rs.memory, nil
}

func (r *mockRegistry) PlaceReservation(id string, c *Container) error {
	re := r.resources[id]
	if re == nil {
		return fmt.Errorf("no resource for id %s", id)
	}

	re.cpus += c.Cpus
	re.memory += c.Memory

	return nil
}

func (r *mockRegistry) ReleaseReservation(id string, c *Container) error {
	re := r.resources[id]
	if re == nil {
		return fmt.Errorf("no resource for id %s", id)
	}

	re.cpus -= c.Cpus
	re.memory -= c.Memory

	return nil
}

func (r *mockRegistry) DeleteDocker(id string) error {
	delete(r.resources, id)

	return nil
}

func (r *mockRegistry) SaveDocker(rs *Docker) error {
	r.resources[rs.ID] = &reservation{r: rs}

	return nil
}

func (r *mockRegistry) Close() error {
	return nil
}
