package metrics

import (
	"fmt"

	"citadelapp.io/citadel"
	"github.com/influxdb/influxdb-go"
)

type Store interface {
	Save(table string, m *Metric) error
	Fetch(query string) ([]*Metric, error)
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

func (i *influxStore) Fetch(query string) ([]*Metric, error) {
	resp, err := i.client.Query(query)
	if err != nil {
		return nil, err
	}

	out := []*Metric{}
	for _, s := range resp {
		for _, p := range s.Points {
			m := &Metric{
				Memory: &Memory{},
				Cpu:    &Cpu{},
			}
			for j, c := range s.Columns {
				switch c {
				case "load_1":
					m.Load1 = p[j].(float64)
				case "load_5":
					m.Load5 = p[j].(float64)
				case "load_15":
					m.Load15 = p[j].(float64)
				case "cpu_nice":
					m.Cpu.Nice = p[j].(float64)
				case "cpu_sys":
					m.Cpu.Sys = p[j].(float64)
				case "cpu_wait":
					m.Cpu.Wait = p[j].(float64)
				case "cpu_user":
					m.Cpu.User = p[j].(float64)
				case "memory_used":
					m.Memory.Used = p[j].(float64)
				case "memory_total":
					m.Memory.Total = p[j].(float64)
				case "time":
					m.Time = p[j].(float64)
				case "sequence_number":
					// ignore
				default:
					return nil, fmt.Errorf("unknown column %s from table %s", c, s.Name)
				}
			}
			out = append(out, m)
		}
	}
	return out, nil
}
