package main

import (
	"net/http"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	versioncollector "github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/common/promslog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

func main() {
	var (
		metricsPath = kingpin.Flag(
			"web.telemetry-path", "Path under which to expose metrics.",
		).Default("/metrics").String()
		toolkitFlags = kingpinflag.AddFlags(kingpin.CommandLine, ":9654")
	)

	promslogConfig := &promslog.Config{}
	flag.AddFlags(kingpin.CommandLine, promslogConfig)
	kingpin.Version(version.Print("ceph_exporter"))
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promslog.New(promslogConfig)

	prometheus.MustRegister(versioncollector.NewCollector("ceph_exporter"))

	http.Handle(*metricsPath, promhttp.Handler())

	if *metricsPath != "/" {
		landingConfig := web.LandingConfig{
			Name:        "Ceph Exporter",
			Description: "Prometheus Ceph Exporter",
			HeaderColor: "#f0424d",
			Version:     version.Info(),
			Links: []web.LandingLinks{
				{
					Address: *metricsPath,
					Text:    "Metrics",
				},
			},
		}
		landingPage, err := web.NewLandingPage(landingConfig)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
		http.Handle("/", landingPage)
	}

	if err := web.ListenAndServe(&http.Server{}, toolkitFlags, logger); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
