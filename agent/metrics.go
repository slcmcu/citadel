package main

import (
	"fmt"

	"citadelapp.io/citadel"
	"citadelapp.io/citadel/metrics"
	"github.com/cloudfoundry/gosigar"
)

func getLoadAverage() (float64, float64, float64, error) {
	load := sigar.LoadAverage{}
	if err := load.Get(); err != nil {
		return -1, -1, -1, err
	}
	return load.One, load.Five, load.Fifteen, nil
}

func getMemoryUsage() (*metrics.Memory, error) {
	mem := sigar.Mem{}
	if err := mem.Get(); err != nil {
		return nil, err
	}
	metric := &metrics.Memory{
		Total: float64(mem.Total),
		Used:  float64(mem.Free),
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
			Path:      dirPath,
			Total:     float64(du.Total),
			Used:      float64(du.Used),
			Free:      float64(du.Free),
			Files:     float64(du.Files),
			Available: float64(du.Avail),
		},
	}
	return disks, nil
}

func getCpuMetrics() (*metrics.Cpu, error) {
	c := sigar.Cpu{}
	if err := c.Get(); err != nil {
		return nil, err
	}
	metric := &metrics.Cpu{
		Nice: float64(c.Nice),
		User: float64(c.User),
		Sys:  float64(c.Sys),
		Wait: float64(c.Wait),
	}
	return metric, nil
}

func pushHostMetrics(host *citadel.Host, store metrics.Store) error {
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

	m := &metrics.Metric{
		Cpu:    cpu,
		Memory: mem,
		Load1:  load1,
		Load5:  load5,
		Load15: load15,
	}
	return store.Save(fmt.Sprintf("metrics.host.%s", host.Name), m)
}
