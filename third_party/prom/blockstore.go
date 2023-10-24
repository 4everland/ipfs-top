package prom

import "github.com/prometheus/client_golang/prometheus"

type BlockStoreMetrics struct {
	cacheHits prometheus.Counter
	hits      prometheus.Counter
	requests  prometheus.Counter
}

func NewBlockStoreMetrics() *BlockStoreMetrics {
	m := &BlockStoreMetrics{
		cacheHits: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: BlockStoreSys,
			Name:      "cache_hits_total",
			Help:      "Total number of blockstore disk cache hits.",
		}),
		hits: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: BlockStoreSys,
			Name:      "hits_total",
			Help:      "Total number of blockstore hits.",
		}),
		requests: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: BlockStoreSys,
			Name:      "blockstore_requests_total",
			Help:      "Total number of blockstore requests.",
		}),
	}

	prometheus.MustRegister(m.requests, m.cacheHits, m.hits)
	return m
}

func (m *BlockStoreMetrics) IncrRequestHits() {
	m.requests.Inc()
}

func (m *BlockStoreMetrics) IncrBlockStoreHits() {
	m.hits.Inc()
	m.IncrRequestHits()
}

func (m *BlockStoreMetrics) IncrCacheHits() {
	m.cacheHits.Inc()
	m.IncrBlockStoreHits()
	m.IncrRequestHits()
}
