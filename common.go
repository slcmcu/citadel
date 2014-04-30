package citadel

type (
	Disk struct {
		Name       string `json:"name" gorethink:"name"`
		TotalSpace int    `json:"total_space" gorethink:"total_space"` // in bytes
	}
	Host struct {
		Name      string  `json:"name" gorethink:"name"`
		IPAddress string  `json:"ip_address" gorethink:"ip_address"`
		Cpus      int     `json:"cpus" gorethink:"cpus"`
		Memory    int     `json:"memory" gorethink:"memory"`
		Disks     []*Disk `json:"disks" gorethink:"disks"`
	}
)
