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
	"citadelapp.io/citadel/repository"
	"github.com/Sirupsen/logrus"
	"github.com/influxdb/influxdb-go"
)

var (
	configPath string
	log        = logrus.New()
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

func initHostInfo(name string, conf *citadel.Config) (*citadel.Host, error) {
	cpus := runtime.NumCPU()
	memUsage, err := getMemoryUsage()
	if err != nil {
		return nil, err
	}

	diskUsage, err := getDiskUsage()
	if err != nil {
		return nil, err
	}

	hostInfo := &citadel.Host{
		Name:      name,
		IPAddress: conf.Listen,
		Cpus:      cpus,
		Memory:    memUsage,
		Disks:     diskUsage,
	}

	repo := repository.NewEtcdRepository(conf.Machines)
	if err := repo.SaveHost(hostInfo); err != nil {
		return nil, err
	}

	log.WithFields(logrus.Fields{
		"cpus":      cpus,
		"memory":    memUsage,
		"diskspace": diskUsage,
	}).Debug("Initializing host info")

	return hostInfo, nil
}

func init() {
	flag.StringVar(&configPath, "config", "config.toml", "path to the configuration file")
	flag.Parse()
}

func main() {
	conf, err := citadel.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	agentName := getAgentName()

	host, err := initHostInfo(agentName, conf)
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(logrus.Fields{
		"nodename": agentName,
		"address":  conf.Listen,
	}).Info("Citadel Agent")

	metrics, err := influxdb.NewClient(&influxdb.ClientConfig{
		Database: conf.InfluxDatabase,
		Host:     conf.InfluxHost,
		Username: conf.InfluxUser,
		Password: conf.InfluxPassword,
	})
	if err != nil {
		log.Fatal(err)
	}

	var (
		sig             = make(chan os.Signal)
		hostMetricsTick = time.Tick(time.Duration(conf.PullInterval) * time.Second)
	)
	signal.Notify(sig, os.Interrupt)

	for {
		select {
		case <-hostMetricsTick:
			if err := pushHostMetrics(host, metrics); err != nil {
				log.Fatal(err)
			}
		case <-sig:
			log.Info("Shutting down Citadel")
			return
		}
	}
}
