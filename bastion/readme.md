# Bastion
This is a reference implementation of the Citadel Docker scheduler.

# Usage
In order to use Bastion, you must setup your Docker remote hosts using TLS.  You can use the example certs in this repository for testing (DO NOT USE IN PRODUCTION) or see https://docs.docker.com/articles/https/ for more information.

## Start Docker with TLS enabled
Place the sample certs in `/certs`.  Add the following to your Docker config and restart the daemon:

`--tls --tlscert --tlskey --tlscacert=/certs/ca.pem --tlscert=/certs/server-cert.pem --tlskey=/certs/server-key.pem -H unix:///var/run/docker.sock -H tcp://0.0.0.0:2375 --tlsverify`

## Run Bastion:
There is a pre-built Docker image available for testing.  It comes bundled with the example certs.

`docker run -it -p 8080:8080 --rm ehazlett/bastion -ca-cert /certs/ca.pem -ssl-cert /certs/client-cert.pem -ssl-key /certs/client-key.pem -hosts <your-hosts>`

For example:

`docker run -it -p 8080:8080 --rm ehazlett/bastion -ca-cert /certs/ca.pem -ssl-cert /certs/client-cert.pem -ssl-key /certs/client-key.pem -hosts https://1.2.3.4:2375`

Create the following `go-demo.json`:

```
{
    "name": "bastion-demo",
    "image": "ehazlett/go-demo",
    "cpus": "0.2",
    "memory": 256,
    "type": "service"
}
```

Then use `curl` to start the container:

`curl -d @go-demo.json <bastion-host-ip:8080>/`

For example, if you are running bastion local:

`curl -d @go-demo.json http://127.0.0.1:8080/`

Bastion will pull the image and then start the container.  Bastion will return the error if one occurs otherwise it will return a `201 Created` on success (no content).
