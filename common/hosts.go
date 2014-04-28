package common

type (
	Disk struct {
		Name       string
		TotalSpace int // in bytes
	}
	Host struct {
		Name      string
		IPAddress string
		Cpus      int
		Memory    int
		Disks     []*Disk
	}
)
