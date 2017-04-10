FROM golang

ADD vendor /usr/src/go/src/prometheus-bridge/vendor
ADD messaging /usr/src/go/src/prometheus-bridge/messaging
ADD server.go /usr/src/go/src/prometheus-bridge/server.go
ADD args.go /usr/src/go/src/prometheus-bridge/args.go
ADD server /usr/src/go/src/prometheus-bridge/server
RUN export GOPATH=/usr/src/go \
    && cd /usr/src/go/src/prometheus-bridge \
    && go build \
    && cp prometheus-bridge /usr/sbin/server

EXPOSE 9091

ENTRYPOINT ["server"]
