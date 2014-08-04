package main

import (
	"log"

	"github.com/citadel/citadel"
)

func main() {
	engines := []*citadel.Engine{}

	cluster, err := citadel.NewCluster(engines...)
	if err != nil {
		log.Fatal(err)
	}

	events, err := cluster.Events()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for e := range events {
			log.Println(e)
		}
	}()

	containers, err := cluster.ListContainers()
	if err != nil {
		log.Fatal(err)
	}

	c1 := containers[0]

	if err := cluster.Kill(c1, 9); err != nil {
		log.Fatal(err)
	}

	if err := cluster.Remove(c1); err != nil {
		log.Fatal(err)
	}

	image := &citadel.Image{
		Name:   "crosbymichael/redis",
		Memory: 256,
		Cpus:   0.4,
	}

	container, err := cluster.Start("service", image)
	if err != nil {
		log.Fatal(err)
	}
}
