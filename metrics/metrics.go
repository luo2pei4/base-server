package metrics

import "github.com/prometheus/client_golang/prometheus"

type ServerCollector struct {
	desc *prometheus.Desc
}

func init() {
	prometheus.MustRegister()
}

func (c *ServerCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

func (c *ServerCollector) Collect(ch chan<- *prometheus.Metric) {}
