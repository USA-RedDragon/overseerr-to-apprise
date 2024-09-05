package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	exampleCounter prometheus.Counter
}

func NewMetrics() *Metrics {
	metrics := &Metrics{
		exampleCounter: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "example_counter",
			Help: "An example counter",
		}),
	}
	metrics.register()
	return metrics
}

func (m *Metrics) register() {
	prometheus.MustRegister(m.exampleCounter)
}

func (m *Metrics) IncrementExampleCounter() {
	m.exampleCounter.Inc()
}
