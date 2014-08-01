package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/citadel/citadel"
	"github.com/gorilla/mux"
	"github.com/samalba/dockerclient"
)

var (
	configPath     string
	config         *Config
	clusterManager *citadel.ClusterManager

	logger = log.New(os.Stderr, "[bastion] ", log.LstdFlags)
)

func init() {
	flag.StringVar(&configPath, "conf", "", "config file")
	flag.Parse()
}

func destroy(w http.ResponseWriter, r *http.Request) {
	var container *citadel.Container
	if err := json.NewDecoder(r.Body).Decode(&container); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := destroyContainer(container); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func receive(w http.ResponseWriter, r *http.Request) {
	var container *citadel.Container
	if err := json.NewDecoder(r.Body).Decode(&container); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := runContainer(container); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func engines(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	if err := json.NewEncoder(w).Encode(config.Engines); err != nil {
		logger.Println(err)
	}
}

func runContainer(container *citadel.Container) error {
	docker, err := clusterManager.ScheduleContainer(container)
	if err != nil {
		return err
	}

	logger.Printf("using host %s (%s)\n", docker.ID, docker.Addr)

	// TODO: error check on run instead of pulling every time?
	if err := docker.Client.PullImage(container.Image, ""); err != nil {
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

	containerId, err := docker.Client.CreateContainer(containerConfig, container.Name)
	if err != nil {
		return err
	}

	if err := docker.Client.StartContainer(containerId, hostConfig); err != nil {
		return err
	}

	logger.Printf("launched %s (%s) on %s\n", container.Name, containerId[:5], docker.ID)

	return nil
}

func destroyContainer(container *citadel.Container) error {
	engines := clusterManager.Engines()
	for _, engine := range engines {
		_, err := engine.Client.InspectContainer(container.Name)
		if err != nil {
			logger.Printf("no container found on host %s name=%s\n", engine.ID, container.Name)
			continue
		}
		if err := engine.Client.KillContainer(container.Name); err != nil {
			return err
		}
		if err := engine.Client.RemoveContainer(container.Name); err != nil {
			return err
		}
	}
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

	for _, d := range config.Engines {
		if err := setDockerClient(d, tlsConfig); err != nil {
			logger.Fatal(err)
		}
	}

	clusterManager = citadel.NewClusterManager(config.Engines, logger)
	clusterManager.RegisterScheduler("service", &citadel.LabelScheduler{})

	r := mux.NewRouter()
	r.HandleFunc("/run", receive).Methods("POST")
	r.HandleFunc("/destroy", destroy).Methods("POST")

	logger.Printf("bastion listening on %s\n", config.ListenAddr)
	if err := http.ListenAndServe(config.ListenAddr, r); err != nil {
		logger.Fatal(err)
	}
}
