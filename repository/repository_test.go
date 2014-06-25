package repository

import (
	"os"
)

func getRepository() *Repository {
	host := []string{os.Getenv("ETCD_HOST")}
	repo := New(host, "citadel-testing")
	return repo
}

func teardown() {
	repo := getRepository()
	repo.client.Delete(repo.namespace, true)
}
