package scheduler

type CpuProfile float64

const (
	DISK   = 1.0
	CPUS   = 1.0
	MEMORY = 0.8

	High   CpuProfile = 1.0
	Medium            = 0.75
	Low               = 0.5
)

type Resource struct {
	Cpus       int        `json:"cpus,omitempty"`
	CpuProfile CpuProfile `json:"cpu_profile,omitempty"`
	Memory     float64    `json:"memory,omitempty"`
	Disk       float64    `json:"disk,omitempty"`
}

// Weight returns the current weight for resources avaliable
// on the host
//
// current: currently used resources on the host
// requested: the requested resources to run a container
//
// problem: are these hard resource requirements or are they soft limits?
// container needs up to 2gb but it's currently only using 1gb so can we
// allocate that extra 1gb to another app?
func Weight(current, requested *Resource) int {
	weight := (current.Disk - requested.Disk) * DISK
	if requested.Cpus > 0 {
		weight = float64((current.Cpus - requested.Cpus)) * CPUS
	}
	weight = (current.Memory - requested.Memory) * MEMORY
	weight *= float64(requested.CpuProfile)

	return int(weight)
}
