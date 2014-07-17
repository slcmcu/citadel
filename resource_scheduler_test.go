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
				Cpus:   []int{0},
				Memory: 1248,
			},
		},
	}

	accepted, err := s.Schedule(app, hosts)
	if err != nil {
		t.Fatal(err)
	}

	if len(accepted) != 1 {
		t.Fatal("no hosts found")
	}

	h1 := accepted[0]

	if h1.ID != "2" {
		t.Fatal("not host id 2")
	}
}
