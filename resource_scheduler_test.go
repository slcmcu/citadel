package citadel

import "testing"

func TestScheduleHighMemory(t *testing.T) {
	var (
		r = newMockRegistry()
		s = NewResourceScheduler(r)
	)

	hosts := []*Host{
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

	app := &Application{
		ID: "1",
		Containers: []*Config{
			{
				Cpus:   0.1,
				Memory: 1248,
			},
		},
	}

	h1, err := s.Schedule(app, hosts)
	if err != nil {
		t.Fatal(err)
	}

	if h1.ID != "2" {
		t.Fatal("not host id 2")
	}
}

// this test should schedule the task on the smallest host because the
// utilization is better than running the task on a large host
func TestScheduleBestUtilization(t *testing.T) {
	var (
		r = newMockRegistry()
		s = NewResourceScheduler(r)
	)

	hosts := []*Host{
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

	app := &Application{
		ID: "1",
		Containers: []*Config{
			{
				Cpus:   0.1,
				Memory: 92,
			},
		},
	}

	h1, err := s.Schedule(app, hosts)
	if err != nil {
		t.Fatal(err)
	}

	if h1.ID != "1" {
		t.Fatal("not host id 1")
	}

}
