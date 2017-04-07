package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"prometheus-amqp-bridge/messaging"

	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/storage/remote"
)

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

func SendMetrics(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	stream, ok := MessagingStream(ctx)
	if !ok {
		http.Error(w, "Messaging stream is not associated with the request!, This is probably a bug!", http.StatusBadRequest)
	}
	msg, err := stream.Consume(&messaging.RabbitMQConsumeSettings{"metrics"})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req remote.WriteRequest
	if err := proto.Unmarshal(msg, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
