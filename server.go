// Copyright 2017 The Mobility House GmbH
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/prometheus/common/model"

	"github.com/prometheus/prometheus/storage/remote"

	"prometheus-amqp-bridge/messaging"
	"prometheus-amqp-bridge/server"
)

func main() {
	stream := &messaging.RabbitMQStream{}
	stream.Connect("amqp://guest:guest@localhost:5672/", messaging.Options{})
	defer stream.Close()

	ctx := server.NewContext(context.Background(), stream)

	http.HandleFunc("/receive", server.HandleWithContext(ctx, server.ReceiveMetrics))

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		msg, err := stream.Consume(&messaging.RabbitMQConsumeSettings{QueueName: "metrics"})

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
	})

	fmt.Println("Starting server on port 9091")
	http.ListenAndServe(":9091", nil)
}
