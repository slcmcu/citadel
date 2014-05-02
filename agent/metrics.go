package main

import (
	"citadelapp.io/citadel"
	"github.com/cloudfoundry/gosigar"
	rethink "github.com/dancannon/gorethink"
)

const (
	HOST_METRICS_TABLE = "host_metrics"
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
	du := sigar.FileSystemUsage{}
	dirPath := "/"
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
func PushHostMetrics() {
	load1, load5, load15, err := getLoadAverage()
	if err != nil {
		log.Errorf("Unable to get host load: %s", err)
		return
	}
	memUsage, err := getMemoryUsage()
	if err != nil {
		log.Error("Unable to get memory usage: %s", err)
		return
	}
	diskUsage, err := getDiskUsage()
	if err != nil {
		log.Error("Unable to get disk usage: %s", err)
		return
	}
	// send to db
	session, err := newRethinkSession()
	if err != nil {
		log.Error("Error connection to RethinkDB: %s", err)
		return
	}
	defer session.Close()
	load := map[string]float64{
		"1":  load1,
		"5":  load5,
		"15": load15,
	}
	metric := citadel.HostMetric{
		Load:   load,
		Memory: memUsage,
		Disks:  diskUsage,
	}
	if _, err = rethink.Table(HOST_METRICS_TABLE).Insert(metric).Run(session); err != nil {
		log.Error("Error pushing host metrics to RethinkDB: %s", err)
		return
	}
}
