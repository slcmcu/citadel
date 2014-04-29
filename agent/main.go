package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"citadelapp.io/citadel"
	"github.com/Sirupsen/logrus"
	rethink "github.com/dancannon/gorethink"
)

var (
	listenAddress string
	rethinkDbHost string
	rethinkDbPort int
	rethinkDbName string

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

// getMemoryTotal gets the total available memory in bytes
func getMemoryTotal() int {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		log.Fatal("Unable to get memory info: %s", err)
	}
	defer f.Close()
	total := 0
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) == 3 {
			if strings.Index(fields[0], "MemTotal") == 0 {
				t, err := strconv.Atoi(fields[1])
				if err != nil {
					log.Fatal("Unable to parse memory total: %s", err)
				}
				total = t * 1024 // convert to bytes
				break
			}
		}
	}
	return total
}

// getDiskTotal returns the total disk space in bytes
func getDiskTotal() uint64 {
	var stat syscall.Statfs_t
	syscall.Statfs("/", &stat)
	return stat.Bavail * uint64(stat.Bsize)
}

func initHostInfo(name string) error {
	var (
		cpus      = runtime.NumCPU()
		memTotal  = getMemoryTotal()
		diskTotal = getDiskTotal()
	)
	disks := []*citadel.Disk{
		{
			Name:       "/",
			TotalSpace: int(getDiskTotal()),
		},
	}

	hostInfo := citadel.Host{
		Name:      name,
		IPAddress: listenAddress,
		Cpus:      cpus,
		Memory:    memTotal,
		Disks:     disks,
	}

	session, err := newRethinkSession()
	if err != nil {
		return err
	}
	defer session.Close()

	if _, err := rethink.Table("hosts").Insert(hostInfo).RunWrite(session); err != nil {
		return err
	}

	log.WithFields(logrus.Fields{
		"cpus":      cpus,
		"memory":    memTotal,
		"diskspace": diskTotal,
	}).Debug("Initializing host info")

	return nil
}

func newRethinkSession() (*rethink.Session, error) {
	return rethink.Connect(rethink.ConnectOpts{
		Address:     fmt.Sprintf("%s:%d", rethinkDbHost, rethinkDbPort),
		Database:    rethinkDbName,
		MaxIdle:     10,
		IdleTimeout: time.Second * 10,
	})
}

func init() {
	flag.StringVar(&listenAddress, "l", "", "Listen address")

	flag.StringVar(&rethinkDbHost, "rethink-host", "127.0.0.1", "RethinkDB Host")
	flag.IntVar(&rethinkDbPort, "rethink-port", 28015, "RethinkDB Port")
	flag.StringVar(&rethinkDbName, "rethink-name", "citadel", "RethinkDB Name")
}

func main() {
	flag.Parse()

	if listenAddress == "" {
		log.Fatal("You must specify a listen address")
	}

	agentName := getAgentName()

	if err := initHostInfo(agentName); err != nil {
		log.Fatal(err)
	}

	log.WithFields(logrus.Fields{
		"nodename": agentName,
		"address":  listenAddress,
	}).Info("Citadel Agent")

	log.WithFields(logrus.Fields{
		"host": rethinkDbHost,
		"port": rethinkDbPort,
		"name": rethinkDbName,
	}).Debug("Connecting to RethinkDB")

}
