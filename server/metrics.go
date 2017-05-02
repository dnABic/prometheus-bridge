package server

import (
	"github.com/prometheus/client_golang/prometheus"
)

type metrics struct {
	HttpRequestCount *prometheus.CounterVec
}

func initializeMetrics() *metrics {
	m := &metrics{
		prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_request_count",
				Help: "Http connection to the server",
			},
			[]string{"uri"},
		),
	}

	prometheus.MustRegister(m.HttpRequestCount)

	return m
}
