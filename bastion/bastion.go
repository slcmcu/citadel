package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/citadel/citadel"
	"github.com/citadel/citadel/docker"
	"github.com/samalba/dockerclient"
)

var (
	configPath     string
	config         *Config
	clusterManager *citadel.ClusterManager

	logger = log.New(os.Stderr, "[bastion] ", log.LstdFlags)
)

type Config struct {
	SSLCertificate string              `json:"ssl-cert,omitempty"`
	SSLKey         string              `json:"ssl-key,omitempty"`
	CACertificate  string              `json:"ca-cert,omitempty"`
	RedisAddr      string              `json:"redis-addr,omitempty"`
	RedisPass      string              `json:"redis-pass,omitempty"`
	ListenAddr     string              `json:"listen-addr,omitempty"`
	Hosts          []*citadel.Resource `json:"hosts,omitempty"`
}

func init() {
	flag.StringVar(&configPath, "conf", "", "config file")
	flag.Parse()
}

func loadConfig() error {
	f, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(&config)
}

func getTLSConfig() (*tls.Config, error) {
	// TLS config
	var tlsConfig tls.Config
	tlsConfig.InsecureSkipVerify = true
	certPool := x509.NewCertPool()

	file, err := ioutil.ReadFile(config.CACertificate)
	if err != nil {
		return nil, err
	}

	certPool.AppendCertsFromPEM(file)
	tlsConfig.RootCAs = certPool
	_, errCert := os.Stat(config.SSLCertificate)
	_, errKey := os.Stat(config.SSLKey)
	if errCert == nil && errKey == nil {
		cert, err := tls.LoadX509KeyPair(config.SSLCertificate, config.SSLKey)
		if err != nil {
			return &tlsConfig, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return &tlsConfig, nil
}

func receive(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var container citadel.Container
		if err := json.NewDecoder(r.Body).Decode(&container); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := runContainer(&container); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func getDockerClient(resource *citadel.Resource) (*dockerclient.DockerClient, error) {
	var tlsConfig *tls.Config
	u, err := url.Parse(resource.Addr)
	if err != nil {
		return nil, err
	}

	// only load tls config if using https
	if u.Scheme == "https" {
		tlsCfg, err := getTLSConfig()
		if err != nil {
			return nil, err
		}
		tlsConfig = tlsCfg
	}

	return dockerclient.NewDockerClient(resource.Addr, tlsConfig)
}

func runContainer(container *citadel.Container) error {
	resource, err := clusterManager.ScheduleContainer(container)
	if err != nil {
		return err
	}

	logger.Printf("using host %s (%s)\n", resource.ID, resource.Addr)

	docker, err := getDockerClient(resource)
	if err != nil {
		return err
	}

	// TODO: error check on run instead of pulling every time?
	if err := docker.PullImage(container.Image, ""); err != nil {
		return err
	}

	// format env
	env := []string{}
	for k, v := range container.Environment {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	containerConfig := &dockerclient.ContainerConfig{
		Hostname:   container.Hostname,
		Domainname: container.Domainname,
		Image:      container.Image,
		Memory:     int(container.Memory) * 1048576,
		Env:        env,
	}

	// TODO: allow to be customized?
	hostConfig := &dockerclient.HostConfig{
		PublishAllPorts: true,
	}

	containerId, err := docker.CreateContainer(containerConfig, container.Name)
	if err != nil {
		return err
	}

	if err := docker.StartContainer(containerId, hostConfig); err != nil {
		return err
	}

	logger.Printf("launched %s (%s) on %s\n", container.Name, containerId[:5], resource.ID)

	return nil
}

func main() {
	if err := loadConfig(); err != nil {
		logger.Fatal(err)
	}

	tlsConfig, err := getTLSConfig()
	if err != nil {
		logger.Fatal(err)
	}

	registry, err := docker.New(config.Hosts, tlsConfig)
	if err != nil {
		logger.Fatal(err)
	}
	defer registry.Close()

	clusterManager = citadel.NewClusterManager(registry, logger)
	clusterManager.RegisterScheduler("service", &citadel.LabelScheduler{})

	http.HandleFunc("/", receive)

	logger.Printf("bastion listening on %s\n", config.ListenAddr)
	if err := http.ListenAndServe(config.ListenAddr, nil); err != nil {
		logger.Fatal(err)
	}
}
