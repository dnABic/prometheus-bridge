package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/storage/remote"
)

var collector = initializeMetrics()

func ReceiveMetrics(o interface{}) ContextHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		collector.HttpRequestCount.WithLabelValues("recieve").Inc()

		reqBuf, err := ioutil.ReadAll(snappy.NewReader(r.Body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		stream, ok := MessagingStream(ctx)
		if !ok {
			http.Error(w, "Messaging stream is not associated with the request!, This is probably a bug!", http.StatusInternalServerError)
			return
		}
		err = stream.Publish(reqBuf, o)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
	}
}

func SendMetrics(o interface{}) func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		collector.HttpRequestCount.WithLabelValues("send").Inc()
		stream, ok := MessagingStream(ctx)
		if !ok {
			http.Error(w, "Messaging stream is not associated with the request!, This is probably a bug!", http.StatusInternalServerError)
			return
		}
		msgs, err := stream.Consume(o)

		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		for _, msg := range msgs {
			var req remote.WriteRequest
			if err := proto.Unmarshal(msg, &req); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			for _, ts := range req.Timeseries {
				m := make(model.Metric, len(ts.Labels))
				for _, l := range ts.Labels {
					m[model.LabelName(l.Name)] = model.LabelValue(l.Value)
				}

				for _, s := range ts.Samples {
					fmt.Fprintf(w, "%s %f %d\n", m, s.Value, s.TimestampMs)
				}
			}
		}
	}
}
