# Bastion
This is a reference implementation of the Citadel Docker scheduler.

# Usage
Create a custom config using `bastion.conf.sample` as an example.  If you use TLS Docker hosts, the use the below to specify the certificate paths.

## Start Docker with TLS enabled
Place the sample certs in `/certs`.  Add the following to your Docker config and restart the daemon:

`--tls --tlscert --tlskey --tlscacert=/certs/ca.pem --tlscert=/certs/server-cert.pem --tlskey=/certs/server-key.pem -H unix:///var/run/docker.sock -H tcp://0.0.0.0:2375 --tlsverify`

## Run Bastion:
There is a pre-built Docker image available for testing.  It comes bundled with the example certs.  This example shows bind mounting an external config file into the container.

`docker run -it -p 8080:8080 -v bastion.conf:/etc/bastion.conf ehazlett/bastion -conf /etc/bastion.conf`

Create the following `go-demo.json`:

```
{
    "name": "bastion-demo",
    "image": "ehazlett/go-demo",
    "hostname": "bastion-demo.example.com",
    "domain": "example.com",
    "cpus": "0.2",
    "memory": 256,
    "type": "service",
    "labels": ["us-east-1"],
    "environment": {
        "FOO": "bar"
    }

}
```

Then use `curl` to start the container:

`curl -d @go-demo.json <bastion-host-ip:8080>/`

For example, if you are running bastion local:

`curl -d @go-demo.json http://127.0.0.1:8080/`

Bastion will pull the image and then start the container.  Bastion will return the error if one occurs otherwise it will return a `201 Created` on success (no content).
