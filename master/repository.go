package master

import (
	"citadelapp.io/citadel"
	"citadelapp.io/citadel/slave"
)

type Repository interface {
	FetchSlaves() ([]*slave.Slave, error)
	FetchContainers(string) (citadel.Containers, error)
}
