## Citadel
Citadel is the combination of Shipyard and Dockerui to provide a complete view on your cluster of docker instances.  


### Goals
* Provide a higher level of abstraction on containers by aggragating them based on their image
* Not a service discovery method
* Schedule tasks with tight resourse limits and monitor resources on the cluster
* Host independent, all you need is docker


### Development
The only hard dep right now for building the project is to the Go installed.  After that just `go run citadel.go` and navigate to `localhost:3000` in your broswer to begin interacting with the ui.


### Architecture
* api
    * provides the web ui
    * sends messages to agents
* agent
    * runs on each host
    * sends messages to docker
    * collects host metrics 
    * collects container metrics 
* database (rethinkdb and/or etcd/redis)
    * stores metrics 
    * stores runtime data
    * aggregates host information
    * lock server for the cluster