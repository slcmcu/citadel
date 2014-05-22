package citadel

type Service interface {
	Data() *ServiceData

	List(t *Task) ([]*ServiceData, error)
	Run(t *Task) (interface{}, error)
	Stop(t *Task) (interface{}, error)
}

// ServiceData represents a resource in a cluster that can be consumed
type ServiceData struct {
	Name   string `json:"name,omitempty"`
	Type   string `json:"type,omitempty"`
	Addr   string `json:"addr,omitempty"`   // address to the service, can include port number
	Cpus   int    `json:"cpus,omitempty"`   // number of cpus assiged to the service
	Memory int    `json:"memory,omitempty"` // memory limit in mb
}
