package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"citadelapp.io/citadel"
	"citadelapp.io/citadel/repository"
	"citadelapp.io/citadel/utils"
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/dancannon/gorethink"
	"github.com/samalba/dockerclient"
)

type (
	HostEngine struct {
		client     *dockerclient.DockerClient
		repository *repository.Repository
		id         string
		listenAddr string
	}
)

var runHostCommand = cli.Command{
	Name:   "run-host",
	Usage:  "run the host and connect it to the cluster",
	Action: runHostAction,
	Flags: []cli.Flag{
		cli.StringFlag{"host-id", "", "specify host id (default: detected)"},
		cli.StringFlag{"region", "", "region where the host is running"},
		cli.StringFlag{"addr", "", "external ip address for the host"},
		cli.StringFlag{"docker", "unix:///var/run/docker.sock", "docker remote ip address"},
		cli.IntFlag{"cpus", -1, "number of cpus available to the host"},
		cli.IntFlag{"memory", -1, "number of mb of memory available to the host"},
		cli.StringFlag{"listen, l", ":8787", "listen address"},
	},
}

func runHostAction(context *cli.Context) {
	var (
		cpus       = context.Int("cpus")
		memory     = context.Int("memory")
		addr       = context.String("addr")
		region     = context.String("region")
		hostId     = context.String("host-id")
		listenAddr = context.String("listen")
	)
	if hostId == "" {
		id, err := utils.GetMachineID()
		if err != nil {
			logger.WithField("error", err).Fatal("unable to read machine id")
		}
		hostId = id
	}

	switch {
	case cpus < 1:
		logger.Fatal("cpus must have a value")
	case memory < 1:
		logger.Fatal("memory must have a value")
	case addr == "":
		logger.Fatal("addr must have a value")
	case region == "":
		logger.Fatal("region must have a value")
	}

	r, err := repository.New(context.GlobalString("repository"))
	if err != nil {
		logger.WithField("error", err).Fatal("unable to connect to repository")
	}
	defer r.Close()

	host := &citadel.Host{
		ID:     hostId,
		Memory: memory,
		Cpus:   cpus,
		Addr:   addr,
		Region: region,
	}

	if err := r.SaveHost(host); err != nil {
		logger.WithField("error", err).Fatal("unable to save host")
	}
	defer r.DeleteHost(hostId)

	client, err := dockerclient.NewDockerClient(context.String("docker"))
	if err != nil {
		logger.WithField("error", err).Fatal("unable to connect to docker")
	}

	hostEngine := &HostEngine{
		client:     client,
		repository: r,
		id:         hostId,
		listenAddr: listenAddr,
	}
	// start
	go hostEngine.run()
	// watch for operations
	go hostEngine.watch()
	// handle stop signal
	hostEngine.waitForInterrupt()
}

func (eng *HostEngine) waitForInterrupt() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	for _ = range sigChan {
		// stop engine
		eng.stop()
		os.Exit(0)
	}
}

func (eng *HostEngine) run() {
	logger.Info("Starting Citadel")
	if err := eng.loadContainers(); err != nil {
		logger.WithField("error", err).Fatal("unable to load containers")
	}

	// listen for events
	eng.client.StartMonitorEvents(eng.dockerEventHandler)

	if err := http.ListenAndServe(eng.listenAddr, nil); err != nil {
		logger.WithField("error", err).Fatal("unable to listen on http")
	}
}

func (eng *HostEngine) stop() {
	logger.Info("Stopping")
	// remove host from repository
	eng.repository.DeleteHost(eng.id)
}

func (eng *HostEngine) loadContainers() error {
	sesson := eng.repository.Session()

	// delete all containers for this host and recreate them
	if _, err := gorethink.Table("containers").Filter(func(row gorethink.RqlTerm) interface{} {
		return row.Field("host_id").Eq(eng.id)
	}).Delete().Run(sesson); err != nil {
		return err
	}

	containers, err := eng.client.ListContainers(true)
	if err != nil {
		return err
	}

	for _, c := range containers {
		cc, err := eng.generateContainerInfo(c)
		if err != nil {
			return err
		}
		if err := eng.repository.SaveContainer(cc); err != nil {
			return err
		}
	}

	return nil
}

func (eng *HostEngine) generateContainerInfo(cnt interface{}) (*citadel.Container, error) {
	c := cnt.(dockerclient.Container)
	info, err := eng.client.InspectContainer(c.Id)
	if err != nil {
		return nil, err
	}
	cc := &citadel.Container{
		ID:     info.Id,
		Image:  utils.CleanImageName(c.Image),
		HostID: eng.id,
		Cpus:   info.Config.CpuShares, // FIXME: not the right place, this is cpuset
	}

	if info.Config.Memory > 0 {
		cc.Memory = info.Config.Memory / 1024 / 1024
	}

	if info.State.Running {
		cc.State.Status = citadel.Running
	} else {
		cc.State.Status = citadel.Stopped
	}
	cc.State.ExitCode = info.State.ExitCode
	return cc, nil
}

func (eng *HostEngine) dockerEventHandler(event *dockerclient.Event, args ...interface{}) {
	switch event.Status {
	case "destroy":
		// remove container from repository
		if err := eng.repository.DeleteContainer(event.Id); err != nil {
			logger.Warnf("Unable to remove container from repository: %s", err)
			return
		}
	default:
		// reload containers into repository
		// when adding a single container, the Container struct is not
		// returned but instead ContainerInfo.  to keep the same
		// generateContainerInfo for a citadel container, i simply
		// re-run the loadContainers.  this can probably be improved.
		eng.loadContainers()
	}
}

func (eng *HostEngine) watch() {
	tickerChan := time.NewTicker(time.Millisecond * 2000).C // check for new instances every 2 seconds
	for {
		select {
		case <-tickerChan:
			tasks, err := eng.repository.FetchTasks()
			if err != nil {
				logger.Fatal("unable to fetch queue: %s", err)
			}
			for _, task := range tasks {
				// filter this hosts tasks
				if task.Host == eng.id {
					go eng.taskHandler(task)
				}
			}
		}
	}
}

func (eng *HostEngine) taskHandler(task *citadel.Task) {
	switch task.Command {
	case "run":
		logger.WithFields(logrus.Fields{
			"host": task.Host,
			"args": task.Args,
		}).Info("processing run task")
		eng.runHandler(task)
		return
	case "restart":
		logger.WithFields(logrus.Fields{
			"host": task.Host,
			"args": task.Args,
		}).Info("processing restart task")
		eng.restartHandler(task)
		return
	case "stop":
		logger.WithFields(logrus.Fields{
			"host": task.Host,
			"args": task.Args,
		}).Info("processing stop task")
		eng.stopHandler(task)
		return
	case "destroy":
		logger.WithFields(logrus.Fields{
			"host": task.Host,
			"args": task.Args,
		}).Info("processing destroy task")
		eng.destroyHandler(task)
		return
	default:
		logger.WithFields(logrus.Fields{
			"command": task.Command,
			"args":    task.Args,
		}).Error("unknown task command")
		return
	}
}

func (eng *HostEngine) runHandler(task *citadel.Task) {
	logger.WithFields(logrus.Fields{
		"host":      task.Host,
		"image":     task.Args["image"],
		"cpus":      task.Args["cpus"],
		"memory":    task.Args["memory"],
		"instances": task.Args["instances"],
	}).Info("running container")
	// remove task
	eng.repository.DeleteTask(task.Id)
	instances := int(task.Args["instances"].(float64))
	// run containers
	for i := 0; i < instances; i++ {
		image := task.Args["image"].(string)
		cpus := int(task.Args["cpus"].(float64))
		memory := int(task.Args["memory"].(float64))
		containerConfig := &dockerclient.ContainerConfig{
			Image:     image,
			Memory:    memory * 1048576, // convert to bytes
			CpuShares: cpus,
			Tty:       true,
			OpenStdin: true,
		}
		hostConfig := &dockerclient.HostConfig{
			PublishAllPorts: true,
		}
		// create container
		containerId, err := eng.client.CreateContainer(containerConfig, "")
		if err != nil {
			switch err.Error() {
			case "Not found":
				logger.WithFields(logrus.Fields{
					"host":  task.Host,
					"image": image,
				}).Info("pulling image")
				// missing image; pull
				eng.client.PullImage(image, "latest")
				// containerId is blank if image is missing; create new config
				cId, err := eng.client.CreateContainer(containerConfig, "")
				if err != nil {
					logger.WithFields(logrus.Fields{
						"image": image,
						"err":   err,
					}).Error("error creating container")
					return
				}
				containerId = cId
			default:
				logger.WithFields(logrus.Fields{
					"image": image,
					"err":   err,
				}).Error("error creating container")
				return
			}
		}
		// start container
		if err := eng.client.StartContainer(containerId, hostConfig); err != nil {
			logger.WithFields(logrus.Fields{
				"image": image,
				"err":   err,
			}).Error("error starting container")
			return
		}

		logger.WithFields(logrus.Fields{
			"host":        task.Host,
			"containerId": containerId,
			"image":       image,
		}).Info("started container")
	}
}

func (eng *HostEngine) stopHandler(task *citadel.Task) {
	logger.WithFields(logrus.Fields{
		"host": task.Host,
		"id":   task.Args["containerId"],
	}).Info("stopping container")
	// remove task
	defer eng.repository.DeleteTask(task.Id)
	containerId := task.Args["containerId"].(string)
	if err := eng.client.StopContainer(containerId, 10); err != nil {
		logger.WithFields(logrus.Fields{
			"containerId": containerId,
			"err":         err,
		}).Error("error stopping container")
		return
	}
}

func (eng *HostEngine) restartHandler(task *citadel.Task) {
	logger.WithFields(logrus.Fields{
		"host": task.Host,
		"id":   task.Args["containerId"],
	}).Info("restarting container")
	// remove task
	defer eng.repository.DeleteTask(task.Id)
	containerId := task.Args["containerId"].(string)
	if err := eng.client.RestartContainer(containerId, 10); err != nil {
		logger.WithFields(logrus.Fields{
			"containerId": containerId,
			"err":         err,
		}).Error("error restarting container")
		return
	}
}

func (eng *HostEngine) destroyHandler(task *citadel.Task) {
	logger.WithFields(logrus.Fields{
		"host": task.Host,
		"id":   task.Args["containerId"],
	}).Info("destroying container")
	// remove task
	defer eng.repository.DeleteTask(task.Id)
	containerId := task.Args["containerId"].(string)
	if err := eng.client.KillContainer(containerId); err != nil {
		logger.WithFields(logrus.Fields{
			"containerId": containerId,
			"err":         err,
		}).Error("error killing container")
		return
	}
	if err := eng.client.RemoveContainer(containerId); err != nil {
		logger.WithFields(logrus.Fields{
			"containerId": containerId,
			"err":         err,
		}).Error("error removing container")
		return
	}
}
