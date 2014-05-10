package repository

import (
	"citadelapp.io/citadel"
)

type Repository interface {
	RegisterSlave(uuid string, r *citadel.Slave, ttl int) error
	UpdateSlave(uuid string, ttl int) error

	FetchSlaves() ([]*citadel.Slave, error)
	FetchSlave(uuid string) (*citadel.Slave, error)
	RemoveSlave(uuid string) error

	SaveContainer(string, *citadel.Container) error
	FetchContainers(string) (citadel.Containers, error)

	FetchConfig() (*citadel.Config, error)
}
