package repository

import (
	"citadelapp.io/citadel"
	"encoding/json"
	"fmt"

	uuid "github.com/nu7hatch/gouuid"
)

func (r *Repository) AddTask(t *citadel.Task) error {
	// assign new random uuid if not present
	if t.ID == nil {
		uid, err := uuid.NewV4()
		if err != nil {
			return err
		}
		id := uid.String()
		t.ID = &id
	}

	d, err := json.Marshal(t)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s/tasks/%s", r.namespace, *t.ID)
	if _, err := r.client.Set(key, string(d), 0); err != nil {
		return err
	}
	return nil
}

func (r *Repository) FetchTasks() ([]*citadel.Task, error) {
	key := fmt.Sprintf("%s/tasks/", r.namespace)
	results, err := r.client.Get(key, false, true)
	if err != nil {
		return nil, err
	}

	var tasks []*citadel.Task
	for _, n := range results.Node.Nodes {
		var task *citadel.Task
		if err := json.Unmarshal([]byte(n.Value), &task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *Repository) DeleteTask(id string) error {
	key := fmt.Sprintf("%s/tasks/%s", r.namespace, id)
	if _, err := r.client.RawDelete(key, true, true); err != nil {
		return err
	}

	return nil
}
