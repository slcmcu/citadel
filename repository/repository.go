package repository

import "citadelapp.io/citadel"

type Repository interface {
	FetchServices(string) ([]*citadel.ServiceData, error)
	FetchService(string) (*citadel.ServiceData, error)

	SaveService(string, *citadel.ServiceData) error

	FetchConfig() (*citadel.Config, error) // global cluster configuration
}
