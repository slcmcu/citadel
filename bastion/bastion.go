package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/samalba/dockerclient"
)

var (
	SSL_CERT    string
	SSL_KEY     string
	CA_CERT     string
	LISTEN_ADDR string
	HOSTS_STR   string
	HOSTS       []string
)

type (
	Container struct {
		Name        string            `json:"name"`
		Image       string            `json:"image"`
		Cpus        float64           `json:"cpus,string"`
		Memory      int               `json:"memory"`
		Type        string            `json:"type"`
		Environment map[string]string `json:"environment"`
		Hostname    string            `json:"hostname"`
		Domainname  string            `json:"domain"`
	}
)

func init() {
	flag.StringVar(&HOSTS_STR, "hosts", "", "Docker hosts - comma separated (i.e. https://1.2.3.4:2375)")
	flag.StringVar(&SSL_CERT, "ssl-cert", "", "SSL Certificate for Docker Hosts")
	flag.StringVar(&SSL_KEY, "ssl-key", "", "SSL Certificate Key for Docker Hosts")
	flag.StringVar(&CA_CERT, "ca-cert", "", "SSL CA Certificate Key for Docker Hosts")
	flag.StringVar(&LISTEN_ADDR, "listen", ":8080", "Listen address")
	flag.Parse()
	if HOSTS_STR != "" {
		HOSTS = strings.Split(HOSTS_STR, ",")
	}
}

func main() {
	http.HandleFunc("/", receive)
	log.Fatal(http.ListenAndServe(LISTEN_ADDR, nil))
}

func receive(w http.ResponseWriter, r *http.Request) {
	var container Container
	if err := json.NewDecoder(r.Body).Decode(&container); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := runContainer(&container); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func getHost() (string, error) {
	// TODO: use citadel for selection
	// check for available hosts
	if len(HOSTS) == 0 {
		return "", errors.New("no hosts available for execution")
	}
	host := HOSTS[rand.Intn(len(HOSTS))]
	return host, nil
}

func runContainer(container *Container) error {
	// TLS config
	var tlsConfig tls.Config
	tlsConfig.InsecureSkipVerify = true
	certPool := x509.NewCertPool()
	file, err := ioutil.ReadFile(CA_CERT)
	if err != nil {
		log.Printf("error reading ca cert %s: %s", CA_CERT, err)
		return err
	}
	certPool.AppendCertsFromPEM(file)
	tlsConfig.RootCAs = certPool
	_, errCert := os.Stat(SSL_CERT)
	_, errKey := os.Stat(SSL_KEY)
	if errCert == nil && errKey == nil {
		cert, err := tls.LoadX509KeyPair(SSL_CERT, SSL_KEY)
		if err != nil {
			log.Printf("error loading X509 key: %s", err)
			return err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	host, err := getHost()
	if err != nil {
		log.Printf("error getting a host for container: %s", err)
		return err
	}
	log.Printf("using host %s\n", host)
	docker, err := dockerclient.NewDockerClient(host, &tlsConfig)
	if err != nil {
		log.Printf("unable to connect to docker daemon: %s", err)
		return err
	}
	// TODO: error check on run instead of pulling every time?
	if err := docker.PullImage(container.Image, ""); err != nil {
		log.Printf("unable to pull image: %s", err)
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
		Memory:     container.Memory * 1048576,
		Env:        env,
	}
	// TODO: allow to be customized?
	hostConfig := &dockerclient.HostConfig{
		PublishAllPorts: true,
	}
	containerId, err := docker.CreateContainer(containerConfig, container.Name)
	if err != nil {
		log.Printf("error creating container: %s", err)
		return err
	}
	if err := docker.StartContainer(containerId, hostConfig); err != nil {
		log.Printf("error starting container: %s", err)
		return err
	}
	log.Printf("launched %s (%s) on %s\n", container.Name, containerId[:5], host)
	return nil
}
