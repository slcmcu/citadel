package repository

import (
	"testing"

	"github.com/citadel/citadel"
)

func TestAddTask(t *testing.T) {
	repo := getRepository()
	args := make(map[string]interface{})
	args["foo"] = "bar"
	task := &citadel.Task{
		Command: "test-command",
		Host:    "test-host",
		Type:    "test",
		Args:    args,
	}
	if err := repo.AddTask(task); err != nil {
		t.Error(err)
	}
	teardown()
}

func TestFetchTasks(t *testing.T) {
	repo := getRepository()
	args := make(map[string]interface{})
	args["foo"] = "bar"
	task := &citadel.Task{
		Command: "test-command",
		Host:    "test-host",
		Type:    "test",
		Args:    args,
	}
	if err := repo.AddTask(task); err != nil {
		t.Error(err)
	}
	tasks, err := repo.FetchTasks()
	if err != nil {
		t.Error(err)
	}
	if len(tasks) == 0 {
		t.Errorf("Expected 1 task; received 0")
	}
}

func TestDeleteTask(t *testing.T) {
	repo := getRepository()
	args := make(map[string]interface{})
	args["foo"] = "bar"
	task := &citadel.Task{
		Command: "test-command",
		Host:    "test-host",
		Type:    "test",
		Args:    args,
	}
	if err := repo.AddTask(task); err != nil {
		t.Error(err)
	}
	tasks, err := repo.FetchTasks()
	if err != nil {
		t.Error(err)
	}
	if len(tasks) == 0 {
		t.Errorf("Expected 1 task; received 0")
	}
	tsk := tasks[0]
	// delete
	if err := repo.DeleteTask(*tsk.ID); err != nil {
		t.Error(err)
	}
	nTasks, err := repo.FetchTasks()
	if err != nil {
		t.Error(err)
	}
	if len(nTasks) == 0 {
		t.Errorf("Expected 0 tasks; received %s", len(nTasks))
	}
}
