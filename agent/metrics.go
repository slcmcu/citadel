package main

import (
	"citadelapp.io/citadel"
	"github.com/cloudfoundry/gosigar"
	"github.com/influxdb/influxdb-go"
)

func getLoadAverage() (float64, float64, float64, error) {
	load := sigar.LoadAverage{}
	if err := load.Get(); err != nil {
		return -1, -1, -1, err
	}
	return load.One, load.Five, load.Fifteen, nil
}

func getMemoryUsage() (*citadel.MemoryUsageMetric, error) {
	mem := sigar.Mem{}
	if err := mem.Get(); err != nil {
		return nil, err
	}
	metric := &citadel.MemoryUsageMetric{
		Total: mem.Total,
		Used:  mem.Free,
	}
	return metric, nil
}

func getDiskUsage() ([]*citadel.Disk, error) {
	var (
		du      = sigar.FileSystemUsage{}
		dirPath = "/"
	)
	if err := du.Get(dirPath); err != nil {
		return nil, err
	}
	disks := []*citadel.Disk{
		{
			Total:     du.Total,
			Used:      du.Used,
			Free:      du.Free,
			Files:     du.Files,
			Available: du.Avail,
			Path:      dirPath,
		},
	}
	return disks, nil
}

func getCpuMetrics() (*citadel.CpuMetric, error) {
	c := sigar.Cpu{}
	if err := c.Get(); err != nil {
		return nil, err
	}
	metric := &citadel.CpuMetric{
		Nice: c.Nice,
		User: c.User,
		Sys:  c.Sys,
		Wait: c.Wait,
	}
	return metric, nil
}

func pushHostMetrics(host *citadel.Host, client *influxdb.Client) error {
	load1, load5, load15, err := getLoadAverage()
	if err != nil {
		return err
	}
	mem, err := getMemoryUsage()
	if err != nil {
		return err
	}
	cpu, err := getCpuMetrics()
	if err != nil {
		return err
	}

	s := &influxdb.Series{
		Name: "metrics.hosts." + host.Name,
		Columns: []string{"load_1", "load_5", "load_15",
			"cpu_nice", "cpu_sys", "cpu_wait", "cpu_user",
			"memory_used", "memory_total"},
		Points: [][]interface{}{
			[]interface{}{load1, load5, load15, cpu.Nice, cpu.Sys, cpu.Wait, cpu.User, mem.Used, mem.Total},
		},
	}
	return client.WriteSeries([]*influxdb.Series{s})
}
