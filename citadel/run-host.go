package main

import (
	"net/http"

	"citadelapp.io/citadel"
	"citadelapp.io/citadel/repository"
	"citadelapp.io/citadel/utils"
	"github.com/codegangsta/cli"
	"github.com/dancannon/gorethink"
	"github.com/samalba/dockerclient"
)

var runHostCommand = cli.Command{
	Name:   "run-host",
	Usage:  "run the host and connect it to the cluster",
	Action: runHostAction,
	Flags: []cli.Flag{
		cli.StringFlag{"region", "", "region where the host is running"},
		cli.StringFlag{"addr", "", "external ip address for the host"},
		cli.StringFlag{"docker", "unix:///var/run/docker.sock", "docker remote ip address"},
		cli.IntFlag{"cpus", -1, "number of cpus available to the host"},
		cli.IntFlag{"memory", -1, "number of mb of memory available to the host"},
	},
}

func runHostAction(context *cli.Context) {
	var (
		cpus   = context.Int("cpus")
		memory = context.Int("memory")
		addr   = context.String("addr")
		region = context.String("region")
	)

	id, err := utils.GetMachineID()
	if err != nil {
		logger.WithField("error", err).Fatal("unable to read machine id")
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
		ID:     id,
		Memory: memory,
		Cpus:   cpus,
		Addr:   addr,
		Region: region,
	}

	if err := r.SaveHost(host); err != nil {
		logger.WithField("error", err).Fatal("unable to save host")
	}
	defer r.DeleteHost(id)

	client, err := dockerclient.NewDockerClient(context.String("docker"))
	if err != nil {
		logger.WithField("error", err).Fatal("unable to connect to docker")
	}

	if err := loadContainers(id, r, client); err != nil {
		logger.WithField("error", err).Fatal("unable to load containers")
	}

	if err := http.ListenAndServe(":8787", nil); err != nil {
		logger.WithField("error", err).Fatal("unable to listen on http")
	}
}

func loadContainers(hostId string, r *repository.Repository, client *dockerclient.DockerClient) error {
	sesson := r.Session()

	// delete all containers for this host and recreate them
	if _, err := gorethink.Table("containers").Filter(func(row gorethink.RqlTerm) interface{} {
		return row.Field("host_id").Eq(hostId)
	}).Delete().Run(sesson); err != nil {
		return err
	}

	containers, err := client.ListContainers(true)
	if err != nil {
		return err
	}

	for _, c := range containers {
		full, err := client.InspectContainer(c.Id)
		if err != nil {
			return err
		}

		cc := &citadel.Container{
			ID:     full.Id,
			Image:  utils.RemoveTag(c.Image),
			HostID: hostId,
			Cpus:   full.Config.CpuShares, // FIXME: not the right place, this is cpuset
		}

		if full.Config.Memory > 0 {
			cc.Memory = full.Config.Memory / 1024 / 1024
		}

		if full.State.Running {
			cc.State.Status = citadel.Running
		} else {
			cc.State.Status = citadel.Stopped
		}
		cc.State.ExitCode = full.State.ExitCode

		if err := r.SaveContainer(cc); err != nil {
			return err
		}
	}

	return nil
}
