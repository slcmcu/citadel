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
* database
    * stores metrics (influxdb)
    * stores runtime data (etcd)
    * aggregates host information (etcd)
    * lock server for the cluster (etcd)


### Metrics
Host metrics are collected and stored in `select * from metrics.hosts.b8f6b1166755` where `b8f6b1166755` is the host's unique name.


**setup config**
```bash
curl -s http://dev.crosbymichael.com:4001/v2/keys/citadel/config -XPUT -d value='{"poll_interval":30, "influx_host":"d
ev.crosbymichael.com:8086", "influx_user": "citadel", "influx_password":"koye", "influx_database":"citadel", "natsd":"nats://dev.crosbymichael.com:4222"}'
```
