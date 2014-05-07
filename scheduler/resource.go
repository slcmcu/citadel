package scheduler

import (
	"reflect"

	"citadelapp.io/citadel/repository"
	"github.com/stevedonovan/luar"
)

var (
	acceptFunc *luar.LuaObject
	lua        = luar.Init()
	boolType   = reflect.TypeOf(true)
)

type Resource struct {
	Image      string            `json:"image,omitempty"`
	Cpus       int               `json:"cpus,omitempty"`
	CpuProfile string            `json:"cpu_profile,omitempty"`
	Memory     float64           `json:"memory,omitempty"`
	Disk       float64           `json:"disk,omitempty"`
	Context    map[string]string `json:"context,omitempty"`
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
	acceptFunc = luar.NewLuaObjectFromName(lua, "Accept")
	return nil
}

func Accept(resource *Resource) (bool, error) {
	r, err := acceptFunc.Callf([]reflect.Type{boolType}, resource)
	if err != nil {
		return false, err
	}
	return r[0].(bool), nil
}
