package main

import (
	"net/http"

	alog "github.com/apex/log"
	"github.com/magicst0ne/rackserver_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)


var (
	rootLoggerCtx *alog.Entry

	sc = &SafeConfig{
		C: &Config{},
	}

	configFile = kingpin.Flag(
		"config.file",
		"Path to configuration file.",
	).Default("config.yml").String()
	listenAddress = kingpin.Flag(
		"web.listen-address",
		"Address to listen on for web interface and telemetry.",
	).Default(":9610").String()
)

func init() {
	rootLoggerCtx = alog.WithFields(alog.Fields{
		"app": "rackserver_exporter",
	})
}

// define new http handleer
func metricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var (
			hostConfig *HostConfig
			err        error
			ok         bool
			group      []string
		)

		registry := prometheus.NewRegistry()

		target := r.URL.Query().Get("target")
		if target == "" {
			http.Error(w, "'target' parameter must be specified", 400)
			return
		}

		group, ok = r.URL.Query()["group"]

		targetLoggerCtx := rootLoggerCtx.WithFields(alog.Fields{
        	"target": target,
        	"group": group,
    	})

		if ok && len(group[0]) >= 1 {
			// Trying to get hostConfig from group.
			if hostConfig, err = sc.HostConfigForGroup(group[0]); err != nil {
				targetLoggerCtx.WithError(err).Error("error getting credentials")
				return
			}
		}

		targetLoggerCtx.Info(hostConfig.Username)
		targetLoggerCtx.Info(hostConfig.Password)

		collector := collector.NewRedfishCollector(target, hostConfig.Username, hostConfig.Password, hostConfig.BasicAuth, targetLoggerCtx)
		registry.MustRegister(collector)
		gatherers := prometheus.Gatherers{
			prometheus.DefaultGatherer,
			registry,
		}
		// Delegate http serving to Prometheus client library, which will call collector.Collect.
		h := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)

	}
}

func main() {

	log.AddFlags(kingpin.CommandLine)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	err := sc.ReloadConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/redfish", metricsHandler())
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
            <head>
            <title>Rackserver Exporter</title>
            </head>
			<body>
            <h1>Rackserver Exporter</h1>
            <form action="/redfish">
            <label>Target:</label> <input type="text" name="target" placeholder="X.X.X.X" value="172.17.100.144"><br>
            <label>Group:</label> <input type="text" name="group" placeholder="group" value="dell"><br>
            <input type="submit" value="Submit">
						</form>
						<p><a href="/metrics">Local metrics</a></p>
            </body>
            </html>`))
	})

	log.Info("app started. listening on ", *listenAddress)
	err = http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}