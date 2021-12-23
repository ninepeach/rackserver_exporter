package collector

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/magicst0ne/rackserver_exporter/redfish"
	redfishcommon "github.com/magicst0ne/rackserver_exporter/redfish/common"
	"github.com/magicst0ne/rackserver_exporter/redfish/redfishapi"
)

// Metric name parts.
const (
	// Exporter namespace.
	namespace = "rackserver"
	// Subsystem(s).
	exporter = "exporter"
	// Math constant for picoseconds to seconds.
	picoSeconds = 1e12
)

// Metric descriptors.
var (
	totalScrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, exporter, "collector_duration_seconds"),
		"Collector time duration.",
		nil, nil,
	)
)

// RedfishCollector collects redfish metrics. It implements prometheus.Collector.
type RedfishCollector struct {
	redfishClient *redfish.APIClient
	collectors    map[string]prometheus.Collector
	redfishUp     prometheus.Gauge
}

// NewRedfishCollector return RedfishCollector
func NewRedfishCollector(host string, username string, password string, basicauth string, logger *log.Entry) *RedfishCollector {
	var collectors map[string]prometheus.Collector
	collectorLogCtx := logger
	BasicAuth := false
	if basicauth!="" {
		BasicAuth = true
	}
	redfishClient, err := newRedfishClient(host, username, password, BasicAuth)


 	if err != nil {
		collectorLogCtx.WithError(err).Error("error creating redfish client")
	} else {
		chassisCollector := NewChassisCollector(namespace, redfishClient, collectorLogCtx)
		systemCollector := NewSystemCollector(namespace, redfishClient, collectorLogCtx)

		//collectors = map[string]prometheus.Collector{"system": systemCollector}
		collectors = map[string]prometheus.Collector{"chassis": chassisCollector, "system": systemCollector}
	}

	return &RedfishCollector{
		redfishClient: redfishClient,
		collectors:    collectors,
		redfishUp: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "",
				Name:      "up",
				Help:      "redfish up",
			},
		),
	}
}

// Describe implements prometheus.Collector.
func (r *RedfishCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, collector := range r.collectors {
		collector.Describe(ch)
	}

}

// Collect implements prometheus.Collector.
func (r *RedfishCollector) Collect(ch chan<- prometheus.Metric) {

	scrapeTime := time.Now()
	if r.redfishClient != nil {
		defer r.redfishClient.Logout()
		r.redfishUp.Set(1)
		wg := &sync.WaitGroup{}
		wg.Add(len(r.collectors))

		defer wg.Wait()
		for _, collector := range r.collectors {
			go func(collector prometheus.Collector) {
				defer wg.Done()
				collector.Collect(ch)
			}(collector)
		}
	} else {
		r.redfishUp.Set(0)
	}

	ch <- r.redfishUp
	ch <- prometheus.MustNewConstMetric(totalScrapeDurationDesc, prometheus.GaugeValue, time.Since(scrapeTime).Seconds())
}

func newRedfishClient(host string, username string, password string, basicauth bool) (*redfish.APIClient, error) {

	url := fmt.Sprintf("https://%s", host)

	config := redfish.ClientConfig{
		Endpoint: url,
		Username: username,
		Password: password,
		Insecure: true,
		BasicAuth: basicauth,
	}
	redfishClient, err := redfish.Connect(config)
	if err != nil {
		return nil, err
	}
	return redfishClient, nil
}

func parseCommonStatusHealth(status redfishcommon.Health) (float64, bool) {
	if bytes.Equal([]byte(status), []byte("OK")) {
		return float64(1), true
	} else if bytes.Equal([]byte(status), []byte("Warning")) {
		return float64(2), true
	} else if bytes.Equal([]byte(status), []byte("Critical")) {
		return float64(3), true
	}
	return float64(0), false
}

func parseCommonStatusState(status redfishcommon.State) (float64, bool) {

	if bytes.Equal([]byte(status), []byte("")) {
		return float64(0), false
	} else if bytes.Equal([]byte(status), []byte("Enabled")) {
		return float64(1), true
	} else if bytes.Equal([]byte(status), []byte("Disabled")) {
		return float64(2), true
	} else if bytes.Equal([]byte(status), []byte("StandbyOffinline")) {
		return float64(3), true
	} else if bytes.Equal([]byte(status), []byte("StandbySpare")) {
		return float64(4), true
	} else if bytes.Equal([]byte(status), []byte("InTest")) {
		return float64(5), true
	} else if bytes.Equal([]byte(status), []byte("Starting")) {
		return float64(6), true
	} else if bytes.Equal([]byte(status), []byte("Absent")) {
		return float64(7), true
	} else if bytes.Equal([]byte(status), []byte("UnavailableOffline")) {
		return float64(8), true
	} else if bytes.Equal([]byte(status), []byte("Deferring")) {
		return float64(9), true
	} else if bytes.Equal([]byte(status), []byte("Quiesced")) {
		return float64(10), true
	} else if bytes.Equal([]byte(status), []byte("Updating")) {
		return float64(11), true
	}
	return float64(0), false
}

func parseCommonPowerState(status redfishapi.PowerState) (float64, bool) {
	if bytes.Equal([]byte(status), []byte("On")) {
		return float64(1), true
	} else if bytes.Equal([]byte(status), []byte("Off")) {
		return float64(2), true
	} else if bytes.Equal([]byte(status), []byte("PoweringOn")) {
		return float64(3), true
	} else if bytes.Equal([]byte(status), []byte("PoweringOff")) {
		return float64(4), true
	}
	return float64(0), false
}

func boolToFloat64(data bool) float64 {

	if data {
		return float64(1)
	}
	return float64(0)

}
