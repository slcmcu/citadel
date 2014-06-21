package repository

import (
	"citadelapp.io/citadel"
	"github.com/dancannon/gorethink"
)

func (r *Repository) AddTask(t *citadel.Task) error {
	if _, err := gorethink.Table("tasks").Insert(t).RunWrite(r.session); err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteTask(id string) error {
	if _, err := gorethink.Table("tasks").Get(id).Delete().Run(r.session); err != nil {
		return err
	}

	return nil
}

func (r *Repository) FetchTasks() ([]*citadel.Task, error) {
	results, err := gorethink.Table("tasks").Run(r.session)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	var tasks []*citadel.Task
	if err := results.ScanAll(&tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}
