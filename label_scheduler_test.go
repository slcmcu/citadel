package citadel

import "testing"

func TestScheduleLabel(t *testing.T) {
	var (
		r = newMockRegistry()
		s = NewLabelScheduler(r)
	)

	hosts := []*Host{
		{
			ID:     "1",
			Labels: []string{"hostname:docker.com", "local", "rack:2"},
		},
		{
			ID: "2",
		},
	}

	app := &Application{
		Constraints: []string{"hostname:docker.com"},
	}

	h1, err := s.Schedule(app, hosts)
	if err != nil {
		t.Fatal(err)
	}

	if h1.ID != "1" {
		t.Fatalf("expected app to be scheduled on 1 but was sent to %s", h1.ID)
	}
}
