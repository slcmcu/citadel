package metrics

import (
	"citadelapp.io/citadel"
	"github.com/influxdb/influxdb-go"
)

type Store interface {
	Save(table string, m *Metric) error
}

type influxStore struct {
	client *influxdb.Client
}

func NewStore(conf *citadel.Config) (Store, error) {
	client, err := influxdb.NewClient(&influxdb.ClientConfig{
		Database: conf.InfluxDatabase,
		Host:     conf.InfluxHost,
		Username: conf.InfluxUser,
		Password: conf.InfluxPassword,
	})
	if err != nil {
		return nil, err
	}
	return &influxStore{client: client}, nil
}

func (i *influxStore) Save(table string, m *Metric) error {
	s := &influxdb.Series{
		Name: table,
		Columns: []string{"load_1", "load_5", "load_15",
			"cpu_nice", "cpu_sys", "cpu_wait", "cpu_user",
			"memory_used", "memory_total"},
		Points: [][]interface{}{
			[]interface{}{m.Load1, m.Load5, m.Load15, m.Cpu.Nice, m.Cpu.Sys, m.Cpu.Wait, m.Cpu.User, m.Memory.Used, m.Memory.Total},
		},
	}
	return i.client.WriteSeries([]*influxdb.Series{s})
}
