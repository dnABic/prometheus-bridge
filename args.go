package main

import "flag"

type Options struct {
	Port      *int
	AMQPUri   *string
	QueueName *string
}

func GetOptions() Options {
	o := Options{
		flag.Int("port", 9091, "Port to start server"),
		flag.String("amqp-uri", "amqp://guest:guest@localhost:5672/", "AMQP Connection string"),
		flag.String("metrics", "metrics", "AMQP queue name for metrics"),
	}

	flag.Parse()

	return o
}
