package citadel

import (
	"fmt"
	"testing"
)

type mockRegistry struct {
	resources map[string]*reservation
}

type reservation struct {
	r      *Resource
	cpus   float64
	memory float64
}

func newMockRegistry(rs []*Resource) *mockRegistry {
	m := &mockRegistry{
		resources: make(map[string]*reservation, len(rs)),
	}

	for _, r := range rs {
		m.resources[r.ID] = &reservation{r: r}
	}

	return m
}

func (r *mockRegistry) GetTotalReservations(id string) (float64, float64, error) {
	rs := r.resources[id]
	if rs == nil {
		return 0, 0, fmt.Errorf("no resource for id %s", id)
	}

	return rs.cpus, rs.memory, nil
}

func TestScheduleHighMemory(t *testing.T) {
	resources := []*Resource{
		{
			ID:     "1",
			Cpus:   1,
			Memory: 512,
		},
		{
			ID:     "2",
			Cpus:   2,
			Memory: 2048,
		},
	}

	var (
		r = newMockRegistry(resources)
		s = newResourceManger(r)
	)

	c := &Container{
		Cpus:   0.1,
		Memory: 1248,
	}

	h1, err := s.PlaceContainer(resources, c)
	if err != nil {
		t.Fatal(err)
	}

	if h1.ID != "2" {
		t.Fatalf("not resource id 2 received %s", h1.ID)
	}
}

// this test should schedule the task on the smallest host because the
// utilization is better than running the task on a large host
func TestScheduleBestUtilization(t *testing.T) {
	resources := []*Resource{
		{
			ID:     "1",
			Cpus:   1,
			Memory: 512,
		},
		{
			ID:     "2",
			Cpus:   1,
			Memory: 2048,
		},
	}

	var (
		r = newMockRegistry(resources)
		s = newResourceManger(r)
	)

	c := &Container{
		Cpus:   0.1,
		Memory: 92,
	}

	h1, err := s.PlaceContainer(resources, c)
	if err != nil {
		t.Fatal(err)
	}

	if h1.ID != "1" {
		t.Fatal("not resource id 1")
	}
}
