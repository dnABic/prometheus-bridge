package messaging

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type metrics struct {
	BrokerConnections     *prometheus.CounterVec
	PublishMessageSummary prometheus.Histogram
	ConsumeMessageSummary prometheus.Histogram
}

const (
	ConnectionFailed    = "failed"
	ConnectionSucceeded = "succeeded"
)

func measureElapsedTime(s time.Time, h prometheus.Histogram) {
	e := time.Since(s)

	h.Observe(e.Seconds())
}

func initializeMetrics() *metrics {
	m := &metrics{
		prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "broker_connections_count",
				Help: "Counts connections established to the broker.",
			},
			[]string{"broker", "status"},
		),

		prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name: "broker_publised_messages_duration",
				Help: "Latency of publishing messages to the broker",
			},
		),
		prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name: "broker_consumed_messages_duration",
				Help: "Latency of consumed messages from the broker",
			},
		),
	}

	prometheus.MustRegister(m.BrokerConnections)
	prometheus.MustRegister(m.PublishMessageSummary)
	prometheus.MustRegister(m.ConsumeMessageSummary)

	return m
}
