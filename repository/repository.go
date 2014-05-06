package repository

import (
	"citadelapp.io/citadel"
)

type Repository interface {
	SaveHost(*citadel.Host) error
	DeleteHost(*citadel.Host) error
	FetchHost(name string) (*citadel.Host, error)
	FetchHosts() ([]*citadel.Host, error)

	FetchConfig() (*citadel.Config, error)

	FetchContainerGroup() ([]*citadel.ContainerGroup, error)

	FetchPlugin() (string, error)
}
