package main

import (
	"log"

	"github.com/citadel/citadel"
	"github.com/citadel/citadel/cluster"
	"github.com/citadel/citadel/scheduler"
)

func main() {
	engines := []*citadel.Engine{}

	c, err := cluster.New(scheduler.NewResourceManager(), engines...)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	image := &citadel.Image{
		Name:   "crosbymichael/redis",
		Memory: 256,
		Cpus:   0.4,
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
