package repository

import (
	"testing"

	"github.com/citadel/citadel"
)

func TestAddHost(t *testing.T) {
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
	teardown()
}

func TestFetchHosts(t *testing.T) {
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
	hosts, err := repo.FetchHosts()
	if err != nil {
		t.Error(err)
	}
	if len(hosts) == 0 {
		t.Errorf("Expected 1 host; received 0")
	}
	teardown()
}

func TestDeleteHost(t *testing.T) {
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
	hosts, err := repo.FetchHosts()
	if err != nil {
		t.Error(err)
	}
	if len(hosts) == 0 {
		t.Errorf("Expected 1 host; received 0")
	}
	// delete
	if err := repo.DeleteHost(h.ID); err != nil {
		t.Error(err)
	}
	nhosts, err := repo.FetchHosts()
	if err != nil {
		t.Error(err)
	}
	if len(nhosts) != 0 {
		t.Errorf("Expected 0 hosts; received %s", len(nhosts))
	}
	teardown()
}
