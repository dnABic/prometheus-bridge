FROM golang

ADD vendor /usr/src/go/src/prometheus-amqp-bridge/vendor
ADD messaging /usr/src/go/src/prometheus-amqp-bridge/messaging
ADD server.go /usr/src/go/src/prometheus-amqp-bridge/server.go
ADD args.go /usr/src/go/src/prometheus-amqp-bridge/args.go
ADD server /usr/src/go/src/prometheus-amqp-bridge/server
RUN export GOPATH=/usr/src/go \
    && cd /usr/src/go/src/prometheus-amqp-bridge \
    && go build \
    && cp prometheus-amqp-bridge /usr/sbin/server

EXPOSE 9091

ENTRYPOINT ["server"]
