Prometheus Bridge
=================

What is this
-------------

This simple service bridges two Prometheus servers in a way that one can push metrics to the other one. It uses Prometheus `remote_write` and standard `metrics` collector.

It uses RabbitMQ as a transfer medium, so after remote push from upstream Prometheus service the message is stored in RabbitMQ and then upon metric requests from the downstream Prometheus server theses messages will be fetched from RabbitMQ and will be exposed.

How to use it
-------------
1. Setup your upstream Prometheus instance to send metrics the Bridge using `remote_write` config
2. Setup your downstream Prometheus instance to scrape metrics from the Bridge
3. Make sure your RabbitMQ instance is configured and accessible.
4. Start the Bridge service.

Build
-----

Usual `go get ...` can build this package. 
If you're on MAC and want to build the Docker image you should build the application for Linux 64bit. Try this:

`GOOS=linux CGO_ENABLED=0 go build`

Additionally a `Dockerfile` is also provided: `docker image build -t bridge:edge .` this expects the Linux binary in the current folder under `prometheus-bridge` name. The above `go build` command should do the job.


Integration Tests
-----------------

Integration tests are a simple `docker-compose` file for now. It launches an `etcd` instance, a `prometheus` instance, a `rabbitmq` instance and finally the `prometheus-bridge`. Then it waits until everything is up and running then it runs a curl command against the bridge `/metrics` endpoint. Which should return `200` response. The payload is ignored for now but it's something that can be enhanced.

To run integration tests you have to build a Docker image for the project with `mobilityhouse/prometheus-bridge:edge`. Follow the instructions below:

```
prometheus-bridge $ GOOS=linux CGO_ENABLED=0 go build
prometheus-bridge $ docker build -t mobilityhouse/prometheus-bridge:edge .
prometheus-bridge $ cd integration
prometheus-bridge/integraion $ docker-compose up --build
```

If everything goes well you should see the final `curl` command with some metrics executed successfully.

Release Cycle
-------------

Every commit/merge to master will be release under pre-release in [github releases](https://github.com/mobilityhouse/prometheus-bridge/releases) with `latest` tag and git tags will be release under release.

Docker images
-------------

Docker images are available in [Docker hub](https://hub.docker.com/r/mobilityhouse/prometheus-bridge/) and push using build pipeline. The convention is that every image with exact version tag like `v0.1` is a release candidate and intermediate build are versions are `v0.1-3-gb003338` which means:

The latest version before this build is: `v0.1`
The number of commits from the latest release to this commit is: `3`
And the git SHA for this build is: `gb003338`

Build pipeline
--------------

Build pipeline is in the private Jenkins instance using `Jenkinsfile` in the root folder.

Contribution
------------

Via Pull Requests and Github Issues

License
-------

Copyright belongs to The Mobility House under Apache License 2.0 - Please refere to LICENSE file for more details.
