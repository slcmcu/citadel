package citadel

import (
	"sort"
)

type Cpus []*Cpu

type Cpu struct {
	ID     int
	UsedBy int
}

func (c Cpus) Len() int {
	return len(c)
}

func (c Cpus) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c Cpus) Less(i, j int) bool {
	var (
		ci = c[i]
		cj = c[j]
	)
	return ci.UsedBy < cj.UsedBy
}

func SortCpusByCount(cpus Cpus) {
	sort.Sort(cpus)
}
