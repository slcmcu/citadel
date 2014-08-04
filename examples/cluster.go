package main

import (
	"log"

	"github.com/citadel/citadel"
	"github.com/citadel/citadel/cluster"
	"github.com/citadel/citadel/scheduler"
)

func main() {
	boot2docker := &citadel.Engine{
		ID:     "boot2docker",
		Addr:   "http://192.168.56.101:2375",
		Memory: 2048,
		Cpus:   4,
		Labels: []string{"local"},
	}

	if err := boot2docker.Connect(nil); err != nil {
		log.Fatal(err)
	}

	c, err := cluster.New(scheduler.NewResourceManager(), boot2docker)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	if err := c.RegisterScheduler("service", &scheduler.LabelScheduler{}); err != nil {
		log.Fatal(err)
	}

	image := &citadel.Image{
		Name:   "crosbymichael/redis",
		Memory: 256,
		Cpus:   0.4,
		Type:   "service",
	}

	container, err := c.Start(image)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v\n", container)

	containers, err := c.ListContainers()
	if err != nil {
		log.Fatal(err)
	}

	c1 := containers[0]

	if err := c.Kill(c1, 9); err != nil {
		log.Fatal(err)
	}

	if err := c.Remove(c1); err != nil {
		log.Fatal(err)
	}
}
