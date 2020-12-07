package metric

import (
	"fmt"
	"os"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/vontikov/stoa/internal/cluster"
	"github.com/vontikov/stoa/internal/util"
)

const namespace string = "github_com_vontikov"

var (
	// Info contains the app information.
	Info prometheus.Gauge

	once sync.Once
)

type leaderCollector struct {
	c cluster.Cluster
	d *prometheus.Desc
}

func newLeaderCollector(app string, c cluster.Cluster) *leaderCollector {
	return &leaderCollector{
		c: c,
		d: prometheus.NewDesc(
			fmt.Sprintf("%s_%s_leader", namespace, app),
			"Leadership status",
			nil, nil),
	}
}

func (c *leaderCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.d
}

func (c *leaderCollector) Collect(ch chan<- prometheus.Metric) {
	var v float64
	if c.c.IsLeader() {
		v = 1.0
	}
	ch <- prometheus.MustNewConstMetric(c.d, prometheus.GaugeValue, v)
}

func Init(app, version string, cluster cluster.Cluster) {
	once.Do(func() {
		hostname, _ := os.Hostname()

		Info = promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: app,
			Name:      "info",
			Help:      "Application info",
			ConstLabels: prometheus.Labels{
				"version":  version,
				"hostname": util.HostHostname(hostname),
			},
		})

		prometheus.MustRegister(newLeaderCollector(app, cluster))
	})
}
