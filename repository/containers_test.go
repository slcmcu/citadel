package repository

import (
	"testing"
	"time"

	"citadelapp.io/citadel"
)

func TestAddContainer(t *testing.T) {
	repo := getRepository()
	state := citadel.State{
		Status:    "created",
		StartedAt: time.Now(),
		ExitedAt:  time.Now(),
		ExitCode:  0,
	}
	c := &citadel.Container{
		ID:     "abc123",
		Image:  "foo",
		HostID: "test-host",
		Cpus:   1,
		Memory: 64,
		State:  state,
	}
	if err := repo.SaveContainer(c); err != nil {
		t.Error(err)
	}
	teardown()
}

func TestFetchContainers(t *testing.T) {
	repo := getRepository()
	h := &citadel.Host{
		ID:     "test-host",
		Region: "test",
		Addr:   "none",
		Cpus:   1,
		Memory: 512,
	}
	if err := repo.SaveHost(h); err != nil {
		t.Error(err)
	}
	state := citadel.State{
		Status:    "created",
		StartedAt: time.Now(),
		ExitedAt:  time.Now(),
		ExitCode:  0,
	}
	c := &citadel.Container{
		ID:     "abc123",
		Image:  "foo",
		HostID: "test-host",
		Cpus:   1,
		Memory: 64,
		State:  state,
	}
	if err := repo.SaveContainer(c); err != nil {
		t.Error(err)
	}

	containers, err := repo.FetchContainers()
	if err != nil {
		t.Error(err)
	}

	if len(containers) == 0 {
		t.Errorf("Expected 1 container; received 0")
	}
	teardown()
}

func TestDeleteContainer(t *testing.T) {
	repo := getRepository()
	h := &citadel.Host{
		ID:     "test-host",
		Region: "test",
		Addr:   "none",
		Cpus:   1,
		Memory: 512,
	}
	if err := repo.SaveHost(h); err != nil {
		t.Error(err)
	}
	state := citadel.State{
		Status:    "created",
		StartedAt: time.Now(),
		ExitedAt:  time.Now(),
		ExitCode:  0,
	}
	c := &citadel.Container{
		ID:     "abc123",
		Image:  "foo",
		HostID: "test-host",
		Cpus:   1,
		Memory: 64,
		State:  state,
	}
	if err := repo.SaveContainer(c); err != nil {
		t.Error(err)
	}

	containers, err := repo.FetchContainers()
	if err != nil {
		t.Error(err)
	}

	if len(containers) == 0 {
		t.Errorf("Expected 1 container; received 0")
	}
	// delete
	if err := repo.DeleteContainer(h.ID, c.ID); err != nil {
		t.Error(err)
	}

	nContainers, err := repo.FetchContainers()
	if err != nil {
		t.Error(err)
	}

	if len(nContainers) != 0 {
		t.Errorf("Expected 0 containers; received %s", len(nContainers))
	}
	teardown()
}
