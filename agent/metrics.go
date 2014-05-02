package main

import (
	"time"

	"citadelapp.io/citadel"
	"github.com/cloudfoundry/gosigar"
	rethink "github.com/dancannon/gorethink"
)

const (
	HOST_METRICS_TABLE = "host_metric"
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
		Free:  mem.Free,
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

func pushHostMetrics(host *citadel.Host) error {
	load1, load5, load15, err := getLoadAverage()
	if err != nil {
		return err
	}
	memUsage, err := getMemoryUsage()
	if err != nil {
		return err
	}
	diskUsage, err := getDiskUsage()
	if err != nil {
		return err
	}
	cpu, err := getCpuMetrics()
	if err != nil {
		return err
	}
	session, err := citadel.NewRethinkSession(rethinkDbHost)
	if err != nil {
		return err
	}
	defer session.Close()

	load := map[string]float64{
		"1":  load1,
		"5":  load5,
		"15": load15,
	}
	metric := citadel.HostMetric{
		Name:      host.Name,
		Load:      load,
		Memory:    memUsage,
		Disks:     diskUsage,
		Cpu:       cpu,
		Timestamp: time.Now(),
	}
	if _, err = rethink.Table(HOST_METRICS_TABLE).Insert(metric).Run(session); err != nil {
		return err
	}
	return nil
}
