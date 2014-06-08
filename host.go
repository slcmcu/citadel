package citadel

// Host represents a host machine on the cluster
// running docker containers
type Host struct {
	// ID is a unique identifier for a host
	ID string `json:"id,omitempty" gorethink:"id,omitempty"`
	// Region is the region/datacenter that the host is running in
	Region string `json:"region,omitempty", gorethink:"region,omitempty"`
	// Addr is the ip and port to connect to a specific host
	Addr string `json:"addr,omitempty" gorethink:"addr,omitempty"`
	// Cpus is the number of cpus the host has available
	Cpus int `json:"cpus,omitempty" gorethink:"cpus,omitempty"`
	// Memory is the amount of memory in mb the host has available
	Memory int `json:"memory,omitempty" gorethink:"memory,omitempty"`
}
