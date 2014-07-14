## Citadel
Manage a multihost docker setup.

### Application
An application is a group of containers that are deployed together on a single host.  You can run multiple instances of an application on multiple hosts or on a single host if the application has a portable configuration.  

When an application is registered on the host an application container will be created with the `<app_name>.group` container name.  This application container will handle the port mapping, a static ip for the application, and volumes.  The actual containers running your code can come and go within the application and ip, port, and data will be preserved with the application container.

#### redis with volumes
```json
{
    "id": "redis",
    "ports": [
        {
            "proto": "tcp",
            "container": 6379
        }
    ],
    "volumes": [
        {
            "path": "/redis" ,
            "uid": 1,
            "gid": 1
        }
    ],
    "containers": [
        {
            "image": "crosbymichael/redis",
            "cpus": [0],
            "memory": 512,
            "type": "service"
        }
    ]
}
```

#### redis master and slave
This sample application configuration will deploy a redis master and slave as an application on any host specified.
```json
{
    "id": "redis_group",
    "ports": [
        {
            "proto": "tcp",
            "host": 6379,
            "container": 6379
        }
    ],
    "containers": [
         {
            "image": "crosbymichael/redis",
            "cpus": [0],
            "memory": 512,
            "type": "service"
         },
         {
            "image": "crosbymichael/redis",
            "cpus": [1],
            "memory": 512,
            "type": "service",
            "args": ["--slaveof", "127.0.0.1", "6379", "--port", "6378"]
         }
    ]
}
```

If a container dies and the type is `service` then it will be automatically restarted.  You don't have to worry about the ip or ports changing because the application container handles that and any data volumes given to the container will be preserved.

### CLI
For now you can use the citadel cli to quickly load applications into the cluster and run them.  If I have a host named `boot2docker` then to run the `redis` application from `samples/redis.jons` then I will do the following:

```bash
# view hosts
citadel hosts 

citadel load --hosts boot2docker redis.json

# view apps in the cluster
citadel apps 

# run the app on the host boot2docker
citadel run --hosts boot2docker redis.json

# look at the containers
citadel containers 
docker ps 
```


### Setup - TODO: WIP
To get citadel running you need a docker host.

```bash
docker run --name etcd -p 4001:4001 -p 7001:7001 coreos/etcd
```

### End Goals
* Provide a higher level of abstraction on containers by agg. them based on their image
* Not a service discovery method
* Schedule tasks with tight resource limits and monitor resources on the cluster
* Host independent, all you need is docker
* libswarm integration to run, stop, rm containers via citadel so you can use the docker cli and existing tools
