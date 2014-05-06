package main

import (
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"

	"citadelapp.io/citadel"
	"citadelapp.io/citadel/metrics"
	"citadelapp.io/citadel/repository"
	"github.com/Sirupsen/logrus"
)

var (
	listen   string
	machines string
	log      = logrus.New()
)

// getAgentName gets the agent name based upon the first available mac address
func getAgentName() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("Unable to detect agent name: %s", err)
	}
	name := "unknown"
	for _, iface := range ifaces {
		n := iface.HardwareAddr.String()
		if n != "" {
			name = strings.Replace(n, ":", "", -1)
			break
		}
	}
	return name
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

	host := &citadel.Host{
		Name:        name,
		IPAddress:   listen,
		Cpus:        cpus,
		TotalMemory: memUsage.Total,
		Disks:       diskUsage,
	}

	return host, nil
}

func init() {
	flag.StringVar(&listen, "listen", "", "Listen address")
	flag.StringVar(&machines, "machines", "127.0.0.1:4001", "Comma separated list of etcd machines")
	flag.Parse()
}

func main() {
	if listen == "" {
		log.Fatal("You must specify a listen address")
	}
	etcdMachines := strings.Split(machines, ",")
	repo := repository.NewEtcdRepository(etcdMachines)

	conf, err := repo.FetchConfig()
	if err != nil {
		log.Fatal(err)
	}

	agentName := getAgentName()

	host, err := getHostInfo(agentName)
	if err != nil {
		log.Fatal(err)
	}

	// save to repo
	if err := repo.SaveHost(host); err != nil {
		log.Fatalf("Unable to save host: %s", err)
	}
	log.WithFields(logrus.Fields{
		"cpus":      host.Cpus,
		"memory":    host.TotalMemory,
		"diskspace": host.Disks,
	}).Debug("Host Info")

	store, err := metrics.NewStore(conf)
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(logrus.Fields{
		"nodename": agentName,
		"address":  listen,
	}).Info("Citadel Agent")

	var (
		sig             = make(chan os.Signal)
		hostMetricsTick = time.Tick(time.Duration(conf.PullInterval) * time.Second)
	)
	signal.Notify(sig, os.Interrupt)

	for {
		select {
		case <-hostMetricsTick:
			if err := pushHostMetrics(host, store); err != nil {
				log.Fatal(err)
			}
		case <-sig:
			log.Info("Shutting down Citadel")
			return
		}
	}
}
