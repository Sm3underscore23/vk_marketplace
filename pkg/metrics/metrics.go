package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	nameSpaceBalancer = "balancer"
)

var (
	lables = []string{"code", "backend_path"}
)

var (
	ErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: nameSpaceBalancer,
		Subsystem: "http",
		Name:      "errors_total",
	},
		lables,
	)

	InFlightReq = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: nameSpaceBalancer,
		Subsystem: "http",
		Name:      "in_flight_requests_total",
	})

	HealthBackend = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: nameSpaceBalancer,
		Subsystem: "checker",
		Name:      "health_backends_amount",
	})

	NotAvailableBackend = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: nameSpaceBalancer,
		Subsystem: "checker",
		Name:      "not_available_backends_amount",
	})

	ResponseTime = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: nameSpaceBalancer,
		Subsystem: "http",
		Name:      "response_time",
		Objectives: map[float64]float64{
			0.1:  0.05,
			0.5:  0.05,
			0.95: 0.01,
			0.99: 0.005,
			1.0:  0.001,
		},
	},
		lables,
	)
)
