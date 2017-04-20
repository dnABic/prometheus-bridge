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
	"log"
	"net/http"

	"prometheus-bridge/messaging"
	"prometheus-bridge/server"
)

func main() {
	o := GetOptions()

	log.Println("Trying to connect to AMQP: ", *o.AMQPUri)

	stream := messaging.NewRabbitMQStream(*o.AMQPUri, messaging.Options{})
	defer stream.Close()

	ctx := server.NewContext(context.Background(), stream)

	http.HandleFunc("/receive", server.HandleWithContext(ctx, server.ReceiveMetrics(&messaging.RabbitMQPublishSettings{*o.QueueName})))
	http.HandleFunc("/metrics", server.HandleWithContext(ctx, server.SendMetrics(&messaging.RabbitMQConsumeSettings{*o.QueueName, *o.Count})))

	fmt.Printf("Starting server on port %d\n", *o.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", *o.Port), nil)
}
