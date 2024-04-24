package metrics

import (
	"github.com/luo2pei4/base-server/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

type serviceCollector struct {
	desc *prometheus.Desc
}

func init() {
	// prometheus.MustRegister(collectors.NewBuildInfoCollector())
	// prometheus.MustRegister(collectors.NewGoCollector(
	// 	collectors.WithGoCollectorRuntimeMetrics(collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/.*")}),
	// ))
	prometheus.MustRegister(newServiceCollector())
}

func newServiceCollector() *serviceCollector {
	return &serviceCollector{
		desc: prometheus.NewDesc("service_stats", "statistics exposed by base-server service", nil, nil),
	}
}

func (s *serviceCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- s.desc
}

func (s serviceCollector) Collect(ch chan<- prometheus.Metric) {
	httpMetricsPrometheus(ch)
}

func httpMetricsPrometheus(ch chan<- prometheus.Metric) {
	for api, value := range middleware.GlobalHTTPStats.TotalRequests.Load() {
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName("http", "requests", "total"),
				"Total number of http requests",
				[]string{"api"}, nil),
			prometheus.CounterValue,
			float64(value),
			api,
		)
	}
	for api, value := range middleware.GlobalHTTPStats.TotalErrors.Load() {
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName("http", "errors", "total"),
				"Total number of error requests",
				[]string{"api"}, nil),
			prometheus.CounterValue,
			float64(value),
			api,
		)
	}
	for api, value := range middleware.GlobalHTTPStats.TotalCanceled.Load() {
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName("http", "canceled", "total"),
				"Total number of canceled requests",
				[]string{"api"}, nil),
			prometheus.CounterValue,
			float64(value),
			api,
		)
	}
}
