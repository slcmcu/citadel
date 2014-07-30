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
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/citadel/citadel"
	"github.com/citadel/citadel/redis"
	"github.com/samalba/dockerclient"
)

var (
	CONFIG_FILE    string
	config         *Config
	logger         = logrus.New()
	clusterManager *citadel.ClusterManager
)

type (
	Config struct {
		SSLCertificate string             `json:"ssl-cert,omitempty"`
		SSLKey         string             `json:"ssl-key,omitempty"`
		CACertificate  string             `json:"ca-cert,omitempty"`
		RedisAddr      string             `json:"redis-addr,omitempty"`
		RedisPass      string             `json:"redis-pass,omitempty"`
		ListenAddr     string             `json:"listen-addr,omitempty"`
		Hosts          []citadel.Resource `json:"hosts,omitempty"`
	}
)

func init() {
	flag.StringVar(&CONFIG_FILE, "conf", "", "config file")
	flag.Parse()
	// load config
	data, err := ioutil.ReadFile(CONFIG_FILE)
	if err != nil {
		logger.Fatalf("unable to open config: %s", err)
	}
	var cfg Config
	if err := json.Unmarshal([]byte(data), &cfg); err != nil {
		logger.Fatalf("unable to parse config: %s", err)
	}
	config = &cfg
}

func main() {
	registry := redis.NewRedisRegistry(config.RedisAddr, config.RedisPass)
	// load host resources
	for _, host := range config.Hosts {
		registry.SaveResource(&host)
	}
	defaultLogger := log.New(ioutil.Discard, "[bastion] ", log.LstdFlags)
	clusterManager = citadel.NewClusterManager(registry, defaultLogger)
	labelScheduler := citadel.LabelScheduler{}
	clusterManager.RegisterScheduler("service", &labelScheduler)

	logger.Infof("bastion listening on %s", config.ListenAddr)
	http.HandleFunc("/", receive)
	logger.Fatal(http.ListenAndServe(config.ListenAddr, nil))
}

func getTLSConfig() (*tls.Config, error) {
	// TLS config
	var tlsConfig tls.Config
	tlsConfig.InsecureSkipVerify = true
	certPool := x509.NewCertPool()
	file, err := ioutil.ReadFile(config.CACertificate)
	if err != nil {
		logger.Errorf("error reading ca cert %s: %s", config.CACertificate, err)
		return nil, err
	}
	certPool.AppendCertsFromPEM(file)
	tlsConfig.RootCAs = certPool
	_, errCert := os.Stat(config.SSLCertificate)
	_, errKey := os.Stat(config.SSLKey)
	if errCert == nil && errKey == nil {
		cert, err := tls.LoadX509KeyPair(config.SSLCertificate, config.SSLKey)
		if err != nil {
			logger.Errorf("error loading X509 key: %s", err)
			return &tlsConfig, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	return &tlsConfig, nil
}

func receive(w http.ResponseWriter, r *http.Request) {
	var container citadel.Container
	if err := json.NewDecoder(r.Body).Decode(&container); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	switch r.Method {
	case "POST":
		if err := runContainer(&container); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		fmt.Fprintf(w, "bastion")
	}
	w.WriteHeader(http.StatusCreated)
}

func getDockerClient(host string) (*dockerclient.DockerClient, error) {
	tlsConfig, err := getTLSConfig()
	if err != nil {
		logger.Errorf("unable to get TLS config: %s", err)
		return nil, err
	}
	if err != nil {
		logger.Errorf("error getting a host for container: %s", err)
		return nil, err
	}
	docker, err := dockerclient.NewDockerClient(host, tlsConfig)
	if err != nil {
		logger.Errorf("unable to connect to docker daemon: %s", err)
		return nil, err
	}
	return docker, nil
}

func runContainer(container *citadel.Container) error {
	// schedule
	resource, err := clusterManager.ScheduleContainer(container)
	if err != nil {
		logger.Errorf("error scheduling container: %s", err)
		return err
	}
	logger.Errorf("using host %s (%s)\n", resource.ID, resource.Addr)
	docker, err := getDockerClient(resource.Addr)
	if err != nil {
		logger.Errorf("error getting docker client: %s", err)
		return err
	}
	// TODO: error check on run instead of pulling every time?
	if err := docker.PullImage(container.Image, ""); err != nil {
		logger.Errorf("unable to pull image: %s", err)
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
		logger.Errorf("error creating container: %s", err)
		return err
	}
	if err := docker.StartContainer(containerId, hostConfig); err != nil {
		logger.Errorf("error starting container: %s", err)
		return err
	}
	logger.Errorf("launched %s (%s) on %s\n", container.Name, containerId[:5], resource.ID)
	return nil
}
