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

	accepted, err := s.Schedule(app, hosts)
	if err != nil {
		t.Fatal(err)
	}

	if len(accepted) != 1 {
		t.Fatalf("scheduler should have found 1 host but found %d", len(hosts))
	}

	h1 := accepted[0]

	if h1.ID != "1" {
		t.Fatalf("expected app to be scheduled on 1 but was sent to %s", h1.ID)
	}
}
