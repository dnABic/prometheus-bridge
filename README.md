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

Usual `go get ...` can build this package. Additionally a `Dockerfile` is also provided: `docker image build -t bridge:edge .`

Contribution
------------

Via Pull Requests and Github Issues

License
-------

Copyright belongs to The Mobility House under Apache License 2.0 
