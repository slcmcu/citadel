package repository

import (
	"citadelapp.io/citadel"
	"citadelapp.io/citadel/master"
	"citadelapp.io/citadel/slave"
)

type Repository interface {
	RegisterSlave(uuid string, s *slave.Slave, ttl int) error
	UpdateSlave(uuid string, ttl int) error

	FetchSlaves() ([]*slave.Slave, error)
	FetchSlave(uuid string) (*slave.Slave, error)
	RemoveSlave(uuid string) error

	SaveContainer(string, *citadel.Container) error
	FetchContainers(string) (citadel.Containers, error)
	RemoveContainer(string, string) error

	RegisterMaster(m *master.Master, ttl int) error
	UpdateMaster(ttl int) error

	FetchConfig() (*citadel.Config, error)
}
