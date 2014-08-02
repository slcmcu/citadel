package citadel

import (
	"io/ioutil"
	"log"
	"testing"
)

var testLogger = log.New(ioutil.Discard, "", log.LstdFlags)

func TestScheduleHighMemory(t *testing.T) {
	resources := []*Engine{
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
		s = newEngineManger(testLogger)
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
	resources := []*Engine{
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
		s = newEngineManger(testLogger)
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
