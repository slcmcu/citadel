package scheduler

import (
	"reflect"

	"citadelapp.io/citadel/repository"
	"github.com/stevedonovan/luar"
)

var (
	lua     = luar.Init()
	intType = reflect.TypeOf(int(0))
)

type Resource struct {
	Cpus       int     `json:"cpus,omitempty"`
	CpuProfile string  `json:"cpu_profile,omitempty"`
	Memory     float64 `json:"memory,omitempty"`
	Disk       float64 `json:"disk,omitempty"`
}

// Weight returns the current weight for resources avaliable
// on the host
func Weight(resource *Resource) (int, error) {
	f := luar.NewLuaObjectFromName(lua, "GetWeight")
	r, err := f.Callf([]reflect.Type{intType}, resource)
	if err != nil {
		return -1, err
	}
	return r[0].(int), nil
}

func Init(repo repository.Repository) error {
	plugin, err := repo.FetchPlugin()
	if err != nil {
		return err
	}
	return loadPlugin(plugin)
}

func loadPlugin(plugin string) error {
	if err := lua.DoString(plugin); err != nil {
		return err
	}
	return nil
}
