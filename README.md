## Citadel
Citadel provides a service oriented approach to cluster management.  Everything is a service and every service has a set of standard endpoints.

### Endpoints
* POST /
    * return a list of all sub services
* POST /run
    * run a new service as the child of the current
* POST /stop
    * stop the child service of the current

### Keyspace

```bash
/citadel/config

/citadel/services
/citadel/services/<name>
/citadel/services/<name>/config
/citadel/services/<name>/services
/citadel/services/<name>/services/<name>
```


### Running


```bash
export ETCD_MACHINES='http://dev.crosbymichael.com:4001'

# start a master process binding to 127.0.0.1:3001
citadel master --cpus 4 --memory 8000

# query the master
citadel

# query the master for a service
citadel /master

# start a slave connecting to a docker
export DOCKER_HOST=tcp://192.168.56.101:4243

#            slave name
citadel --service local slave --cpus 2 --memory 1024

# query the master for the slave
citadel /master

# see the containers on the slave
citadel /master/local

```



### End Goals
* Provide a higher level of abstraction on containers by aggragating them based on their image
* Not a service discovery method
* Schedule tasks with tight resourse limits and monitor resources on the cluster
* Host independent, all you need is docker

