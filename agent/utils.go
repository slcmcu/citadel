package main

import (
	"crypto/sha1"
	"encoding/hex"
	"net"
	"runtime"
	"strings"

	"citadelapp.io/citadel"
)

// getAgentName gets the agent name based upon the first available mac address
func getAgentName() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	name := "unknown"
	for _, iface := range ifaces {
		n := iface.HardwareAddr.String()
		if n != "" {
			name = strings.Replace(n, ":", "", -1)
			break
		}
	}
	return name, nil
}

func generateHostId(name string) string {
	h := sha1.New()
	h.Write([]byte(name))
	return hex.EncodeToString(h.Sum(nil))
}

func getHostInfo(name string) (*citadel.Host, error) {
	cpus := runtime.NumCPU()
	memUsage, err := getMemoryUsage()
	if err != nil {
		return nil, err
	}

	diskUsage, err := getDiskUsage()
	if err != nil {
		return nil, err
	}

	return &citadel.Host{
		Name:        name,
		IPAddress:   listen,
		Cpus:        cpus,
		TotalMemory: memUsage.Total,
		Disks:       diskUsage,
	}, nil
}
