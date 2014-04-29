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

	"citadelapp.io/citadel/common"
	"github.com/Sirupsen/logrus"
	rethink "github.com/dancannon/gorethink"
)

var (
	listenAddress string
	listenPort    int
	rethinkDbHost string
	rethinkDbPort int
	rethinkDbName string
	enableDebug   bool
	log           = logrus.New()
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

// getIP detects the first available non-local address
func getIP() string {
	// attempt to detect non-local ip if none is specified via the flag
	if listenAddress != "0.0.0.0" && strings.Index(listenAddress, "127") != 0 {
		return listenAddress
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalf("Unable to get network addresses: %s", err)
	}
	ip := ""
	for _, addr := range addrs {
		a := addr.String()
		if strings.Index(a, "127") != 0 {
			i := strings.Split(a, "/")
			if len(i) == 0 {
				log.Fatal("Error parsing IP")
			}
			ip = i[0]
			break
		}
	}
	return ip
}

func initHostInfo(name string) {
	cpus := runtime.NumCPU()
	ip := getIP()
	memTotal := getMemoryTotal()
	diskTotal := getDiskTotal()
	// host info
	disks := []*common.Disk{}
	disk := common.Disk{
		Name:       "/",
		TotalSpace: int(getDiskTotal()),
	}
	disks = append(disks, &disk)
	hostInfo := common.Host{
		Name:      name,
		IPAddress: getIP(),
		Cpus:      cpus,
		Memory:    memTotal,
		Disks:     disks,
	}
	session, err := newRethinkSession()
	if err != nil {
		log.Fatalf("Error connecting to RethinkDB: %s", err)
	}
	tbl := rethink.Table("hosts")
	if _, err := tbl.Insert(hostInfo).RunWrite(session); err != nil {
		log.Fatalf("Error pushing host info to rethink: %s", err)
	}
	log.WithFields(logrus.Fields{
		"cpus":      cpus,
		"ip":        ip,
		"memory":    memTotal,
		"diskspace": diskTotal,
	}).Debug("Initializing host info")
}

func newRethinkSession() (*rethink.Session, error) {
	// get rethink pool
	session, err := rethink.Connect(rethink.ConnectOpts{
		Address:     fmt.Sprintf("%s:%d", rethinkDbHost, rethinkDbPort),
		Database:    rethinkDbName,
		MaxIdle:     10,
		IdleTimeout: time.Second * 10,
	})
	return session, err
}

func init() {
	flag.StringVar(&listenAddress, "l", "", "Listen address")
	flag.IntVar(&listenPort, "p", 3001, "Listen port")
	flag.StringVar(&rethinkDbHost, "rethink-host", "127.0.0.1", "RethinkDB Host")
	flag.IntVar(&rethinkDbPort, "rethink-port", 28015, "RethinkDB Port")
	flag.StringVar(&rethinkDbName, "rethink-name", "citadel", "RethinkDB Name")
	flag.BoolVar(&enableDebug, "debug", false, "Enable debug logging")
}

func main() {
	flag.Parse()
	if listenAddress == "" {
		log.Error("You must specify a listen address")
		return
	}
	if enableDebug {
		log.Level = logrus.Debug
	}
	agentName := getAgentName()
	// add host info
	log.WithFields(logrus.Fields{
		"nodename": agentName,
		"address":  listenAddress,
	}).Info("Citadel Agent")
	log.WithFields(logrus.Fields{
		"host": rethinkDbHost,
		"port": rethinkDbPort,
		"name": rethinkDbName,
	}).Debug("Connecting to RethinkDB")
	initHostInfo(agentName)
}
