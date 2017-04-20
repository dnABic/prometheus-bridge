package main

import "flag"

type Options struct {
	Port      *int
	AMQPUri   *string
	QueueName *string
	Count     *uint
}

func GetOptions() Options {
	o := Options{
		flag.Int("port", 9091, "Port to start server"),
		flag.String("amqp-uri", "amqp://guest:guest@localhost:5672/", "AMQP Connection string"),
		flag.String("queue-name", "metrics", "AMQP queue name for metrics"),
		flag.Uint("count", 50, "How many messsages should be fetched at once when delivering metrics."),
	}

	flag.Parse()

	return o
}
