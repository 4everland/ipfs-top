package prom

import "github.com/prometheus/client_golang/prometheus"

type ProvideMetrics struct {
	period    prometheus.Gauge
	providers *prometheus.CounterVec
}

func NewProvideMetrics(label ...string) *ProvideMetrics {
	m := &ProvideMetrics{
		period: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: ProviderSys,
			Name:      "period_total",
			Help:      "Total number of periods",
		}),
		providers: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: ProviderSys,
			Name:      "provide_total",
			Help:      "Total number of provide.",
		}, label),
	}

	prometheus.MustRegister(m.period, m.providers)
	return m
}

func (m *ProvideMetrics) Inc() {
	m.period.Inc()
}

func (m *ProvideMetrics) Clean() {
	m.period.Set(0)
}

func (m *ProvideMetrics) Provide(label ...string) {
	m.providers.WithLabelValues(label...).Inc()
}
