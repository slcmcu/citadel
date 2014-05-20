package repository

import "citadelapp.io/citadel"

type Repository interface {
	FetchServices(string) ([]*citadel.Service, error)
	FetchService(string) (*citadel.Service, error)

	SaveService(string, *citadel.Service) error

	FetchConfig() (*citadel.Config, error) // global cluster configuration
}
