## Citadel
Citadel provides a service oriented approach to cluster management.  Everything is a service and every service has a set of standard endpoints.

### Keyspace

```bash
/citadel/config

/citadel/services
/citadel/services/<name>
/citadel/services/<name>/config
/citadel/services/<name>/services
/citadel/services/<name>/services/<name>
```

### Goals
* Provide a higher level of abstraction on containers by aggragating them based on their image
* Not a service discovery method
* Schedule tasks with tight resourse limits and monitor resources on the cluster
* Host independent, all you need is docker

### Metrics
Host metrics are collected and stored in `select * from metrics.hosts.b8f6b1166755` where `b8f6b1166755` is the host's unique name.


**setup config**
```bash
curl -s http://dev.crosbymichael.com:4001/v2/keys/citadel/config -XPUT -d value='{"poll_interval":30, "influx_host":"dev.crosbymichael.com:8086", "influx_user": "citadel", "influx_password":"koye", "influx_database":"citadel", "natsd":"nats://dev.crosbymichael.com:4222", "ttl": 30, "master_timeout", "10s"}'
```

