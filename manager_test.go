package citadel

import (
	"io/ioutil"
	"log"
	"testing"
)

type testScheduler struct {
}

func (t *testScheduler) Schedule(rs []*Docker, c *Container) ([]*Docker, error) {
	return rs, nil
}

var defaultLogger = log.New(ioutil.Discard, "[citadel.tests] ", log.LstdFlags)

func TestNewManager(t *testing.T) {
	m := NewClusterManager(nil, defaultLogger)

	if m == nil {
		t.Fatal("NewClusterManager should not return a nil manager")
	}
}

func TestNoSchedulerForType(t *testing.T) {
	m := NewClusterManager(nil, defaultLogger)

	_, err := m.ScheduleContainer(&Container{
		Type: "service",
	})

	if err == nil {
		t.Fatal("expected error to not be nil")
	}

	if err != ErrNoSchedulerForType {
		t.Fatalf("expected error type ErrNoSchedulerForType but received %q", err)
	}
}

func TestSchedulerRegister(t *testing.T) {
	var (
		m = NewClusterManager(nil, defaultLogger)
		s = &testScheduler{}
	)

	if err := m.RegisterScheduler("service", s); err != nil {
		t.Fatal(err)
	}
}

func TestSchedulerRegisterWithConflict(t *testing.T) {
	var (
		m = NewClusterManager(nil, defaultLogger)
		s = &testScheduler{}
	)

	if err := m.RegisterScheduler("service", s); err != nil {
		t.Fatal(err)
	}

	err := m.RegisterScheduler("service", s)
	if err == nil {
		t.Fatal("expected error to not be nil")
	}

	if err != ErrSchedulerExists {
		t.Fatalf("expected error ErrSchedulerExists but received %q", err)
	}
}
