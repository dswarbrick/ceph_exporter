package osd

import (
	"log/slog"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/dswarbrick/ceph_exporter/asok"
)

// OSDCollector implements the prometheus.Collector interface.
type OSDCollector struct {
	logger *slog.Logger
}

func NewOSDCollector(logger *slog.Logger) OSDCollector {
	c := OSDCollector{logger: logger}
	return c
}

func (c OSDCollector) Collect(ch chan<- prometheus.Metric) {
	matches, err := filepath.Glob(filepath.Join(*asokDir, "ceph-osd.*.asok"))
	if err != nil {
		c.logger.Error("Could not find OSD admin sockets", "err", err)
		return
	}

	for _, p := range matches {
		client := asok.NewAdminSocketClient(p)
		c.logger.Debug("OSD admin socket ping", "path", p, "response", client.Ping())
	}
}

func (n OSDCollector) Describe(ch chan<- *prometheus.Desc) {
}
