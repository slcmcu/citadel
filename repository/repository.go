package repository

import (
	"citadelapp.io/citadel"
)

type Repository interface {
	RegisterSlave(uuid string, r *citadel.Resource, ttl int) error
	UpdateSlave(uuid string, ttl int) error
	SaveState(uuid string, state *citadel.State) error

	SaveHost(*citadel.Host) error
	DeleteHost(*citadel.Host) error
	FetchHost(name string) (*citadel.Host, error)
	FetchHosts() ([]*citadel.Host, error)

	FetchConfig() (*citadel.Config, error)

	FetchContainerGroup() ([]*citadel.ContainerGroup, error)
}
