package server

import (
	"context"
	"io/ioutil"
	"net/http"
	"prometheus-amqp-bridge/messaging"

	"github.com/golang/snappy"
)

type ServerHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request)

func HandleWithContext(ctx context.Context, f ServerHandler) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f(ctx, w, r)
	})
}

func ReceiveMetrics(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	reqBuf, err := ioutil.ReadAll(snappy.NewReader(r.Body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stream, ok := MessagingStream(ctx)
	if !ok {
		http.Error(w, "Messaging stream is not associated with the request!, This is probably a bug!", http.StatusBadRequest)
	}
	err = stream.Publish(reqBuf, &messaging.RabbitMQPublishSettings{"metrics"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
