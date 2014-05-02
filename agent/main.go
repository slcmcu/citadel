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
	"github.com/Sirupsen/logrus"
	rethink "github.com/dancannon/gorethink"
)

const (
	HOST_METRICS_INTERVAL = 5
	HOST_TABLE            = "host"
)

var (
	listenAddress string
	rethinkDbHost string

	log = logrus.New()
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

// getHostFilter returns an RqlTerm filtered on the specified host
func getHostFilter(name string) rethink.RqlTerm {
	return rethink.Table(HOST_TABLE).Filter(map[string]string{"name": name})
}

func initHostInfo(name string) (*citadel.Host, error) {
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
		IPAddress: listenAddress,
		Cpus:      cpus,
		Memory:    memUsage,
		Disks:     diskUsage,
	}

	session, err := citadel.NewRethinkSession(rethinkDbHost)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	row, err := getHostFilter(name).RunRow(session)
	if row.IsNil() {
		if _, err := rethink.Table(HOST_TABLE).Insert(hostInfo).RunWrite(session); err != nil {
			return nil, err
		}
	} else {
		if _, err := getHostFilter(name).Update(hostInfo).Run(session); err != nil {
			return nil, err
		}
	}

	log.WithFields(logrus.Fields{
		"cpus":      cpus,
		"memory":    memUsage,
		"diskspace": diskUsage,
	}).Debug("Initializing host info")

	return hostInfo, nil
}

func init() {
	flag.StringVar(&listenAddress, "l", "", "Listen address")
	flag.StringVar(&rethinkDbHost, "rethink-host", "127.0.0.1:28015", "RethinkDB Address")

	flag.Parse()
}

func main() {
	if listenAddress == "" {
		log.Fatal("You must specify a listen address")
	}

	agentName := getAgentName()

	host, err := initHostInfo(agentName)
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(logrus.Fields{
		"nodename": agentName,
		"address":  listenAddress,
	}).Info("Citadel Agent")

	log.WithFields(logrus.Fields{
		"host": rethinkDbHost,
	}).Debug("Connecting to RethinkDB")

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

	hostMetricsTick := time.Tick(HOST_METRICS_INTERVAL * time.Second)

main:
	for {
		select {
		case <-hostMetricsTick:
			if err := pushHostMetrics(host); err != nil {
				log.Fatal(err)
			}
		case <-sig:
			break main
		}
	}
	log.Info("Shutting down Citadel")
}
