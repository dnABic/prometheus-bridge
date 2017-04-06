FROM golang

ADD vendor /usr/src/go/src/server/vendor
ADD server.go /usr/src/go/src/server/server.go
RUN export GOPATH=/usr/src/go \
    && cd /usr/src/go/src/server \
    && go build \
    && cp server /usr/sbin/

EXPOSE 9091

CMD ["server"]
