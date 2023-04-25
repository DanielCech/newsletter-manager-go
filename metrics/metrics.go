package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	namespace = ""
	subsystem = ""
)

// FQName builds fully-qualified metrics name using namespace, subsystem and name.
func FQName(name string) string {
	return prometheus.BuildFQName(namespace, subsystem, name)
}

// MustRegister registers collectors.
func MustRegister(collectors ...prometheus.Collector) {
	prometheus.MustRegister(collectors...)
}

// CounterOpts creates new counter opts using namespace, subsystem, name and labels.
func CounterOpts(name string, labels prometheus.Labels) prometheus.CounterOpts {
	return prometheus.CounterOpts{
		Namespace:   namespace,
		Subsystem:   subsystem,
		Name:        name,
		ConstLabels: labels,
	}
}

// GaugeOpts creates new gauge opts using namespace, subsystem, name and labels.
func GaugeOpts(name string, labels prometheus.Labels) prometheus.GaugeOpts {
	return prometheus.GaugeOpts{
		Namespace:   namespace,
		Subsystem:   subsystem,
		Name:        name,
		ConstLabels: labels,
	}
}

// HistogramOpts creates new histogram opts using namespace, subsystem, name and labels.
func HistogramOpts(name string, labels prometheus.Labels) prometheus.HistogramOpts {
	return prometheus.HistogramOpts{
		Namespace:   namespace,
		Subsystem:   subsystem,
		Name:        name,
		ConstLabels: labels,
	}
}

// SummaryOpts creates new summary opts using namespace, subsystem, name and labels.
func SummaryOpts(name string, labels prometheus.Labels) prometheus.SummaryOpts {
	return prometheus.SummaryOpts{
		Namespace:   namespace,
		Subsystem:   subsystem,
		Name:        name,
		ConstLabels: labels,
	}
}
