package metrics

import (
	"net/http"
	"regexp"

	"github.com/luo2pei4/base-server/internal/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type serviceMetricsCollector struct {
	desc *prometheus.Desc
}

var serviceCollector *serviceMetricsCollector

func init() {
	serviceCollector = newServiceMetricsCollector()
}

func PrometheusHandler() http.Handler {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		serviceCollector,
		collectors.NewBuildInfoCollector(),
		collectors.NewGoCollector(
			collectors.WithGoCollectorRuntimeMetrics(collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/.*")}),
		),
	)
	gatherers := prometheus.Gatherers{
		registry,
	}
	return promhttp.InstrumentMetricHandler(
		registry,
		promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{
			ErrorHandling: promhttp.ContinueOnError,
		}))
}

func newServiceMetricsCollector() *serviceMetricsCollector {
	return &serviceMetricsCollector{
		desc: prometheus.NewDesc("service_stats", "statistics exposed by base-server service", nil, nil),
	}
}

func (s *serviceMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- s.desc
}

func (s serviceMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	httpMetricsPrometheus(ch)
}

func httpMetricsPrometheus(ch chan<- prometheus.Metric) {
	for api, value := range middleware.GlobalHTTPStats.TotalRequests.Load() {
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName("service_stats", "http", "requests_total"),
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
				prometheus.BuildFQName("service_stats", "http", "errors_total"),
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
				prometheus.BuildFQName("service_stats", "http", "canceled_total"),
				"Total number of canceled requests",
				[]string{"api"}, nil),
			prometheus.CounterValue,
			float64(value),
			api,
		)
	}
}
