package collector

import (
	"sync"
	"strings"
	"fmt"

	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/magicst0ne/rackserver_exporter/redfish"
	"github.com/magicst0ne/rackserver_exporter/redfish/redfishapi"
)

// A SystemCollector implements the prometheus.Collector.

type systemMetric struct {
	desc *prometheus.Desc
}

// SystemSubsystem is the system subsystem
var (
	SystemSubsystem                   = "system"
	SystemLabelNames                  = []string{"sn","mfr", "resource", "system_id", "hw_model"}
	SystemMemoryLabelNames            = []string{"sn","mfr", "resource", "memory", "memory_id"}
	SystemProcessorLabelNames         = []string{"sn", "resource", "processor_id", "processor_model"}
	SystemDriveLabelNames             = []string{"sn", "resource", "drive_name", "drive_model"}

	systemMetrics                     = map[string]systemMetric{
		"system_state": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, SystemSubsystem, "state"),
				"system state,1(Enabled),2(Disabled),3(StandbyOffinline),4(StandbySpare),5(InTest),6(Starting),7(Absent),8(UnavailableOffline),9(Deferring),10(Quiesced),11(Updating)",
				SystemLabelNames,
				nil,
			),
		},
		"system_health_status": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, SystemSubsystem, "health_status"),
				"system health,1(OK),2(Warning),3(Critical)",
				SystemLabelNames,
				nil,
			),
		},
		"system_power_state": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, SystemSubsystem, "power_state"),
				"system power state",
				SystemLabelNames,
				nil,
			),
		},
		"system_processor_summary_state": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, SystemSubsystem, "processor_summary_state"),
				"system overall processor state,1(Enabled),2(Disabled),3(StandbyOffinline),4(StandbySpare),5(InTest),6(Starting),7(Absent),8(UnavailableOffline),9(Deferring),10(Quiesced),11(Updating)",
				SystemLabelNames,
				nil,
			),
		},
		"system_processor_summary_health_status": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, SystemSubsystem, "processor_summary_health_status"),
				"system overall processor health,1(OK),2(Warning),3(Critical)",
				SystemLabelNames,
				nil,
			),
		},
		"system_processor_summary_count": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, SystemSubsystem, "processor_summary_count"),
				"system total processor count",
				SystemLabelNames,
				nil,
			),
		},
		"system_memory_summary_state": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, SystemSubsystem, "memory_summary_state"),
				"system memory state,1(Enabled),2(Disabled),3(StandbyOffinline),4(StandbySpare),5(InTest),6(Starting),7(Absent),8(UnavailableOffline),9(Deferring),10(Quiesced),11(Updating)",
				SystemMemoryLabelNames,
				nil,
			),
		},
		"system_memory_summary_health_status": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, SystemSubsystem, "memory_summary_health_status"),
				"system overall memory health,1(OK),2(Warning),3(Critical)",
				SystemLabelNames,
				nil,
			),
		},
		"system_memory_summary_size": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, SystemSubsystem, "memory_summary_size"),
				"system total memory size, GiB",
				SystemLabelNames,
				nil,
			),
		},

		"system_processor_state": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, SystemSubsystem, "processor_state"),
				"system processor state,1(Enabled),2(Disabled),3(StandbyOffinline),4(StandbySpare),5(InTest),6(Starting),7(Absent),8(UnavailableOffline),9(Deferring),10(Quiesced),11(Updating)",
				SystemProcessorLabelNames,
				nil,
			),
		},
		"system_processor_health_status": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, SystemSubsystem, "processor_health_status"),
				"system processor health state,1(OK),2(Warning),3(Critical)",
				SystemProcessorLabelNames,
				nil,
			),
		},
		"system_processor_total_threads": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, SystemSubsystem, "processor_total_threads"),
				"system processor total threads",
				SystemProcessorLabelNames,
				nil,
			),
		},
		"system_processor_total_cores": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, SystemSubsystem, "processor_total_cores"),
				"system processor total cores",
				SystemProcessorLabelNames,
				nil,
			),
		},
		"system_storage_drive_state": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, SystemSubsystem, "storage_drive_state"),
				"system storage drive state,1(Enabled),2(Disabled),3(StandbyOffinline),4(StandbySpare),5(InTest),6(Starting),7(Absent),8(UnavailableOffline),9(Deferring),10(Quiesced),11(Updating)",
				SystemDriveLabelNames,
				nil,
			),
		},
		"system_storage_drive_health_state": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, SystemSubsystem, "storage_drive_health_state"),
				"system storage volume health state,1(OK),2(Warning),3(Critical)",
				SystemDriveLabelNames,
				nil,
			),
		},
		"system_storage_drive_capacity": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, SystemSubsystem, "storage_drive_capacity"),
				"system storage drive capacity,Bytes",
				SystemDriveLabelNames,
				nil,
			),
		},
	}
)

// SystemCollector implemented prometheus.Collector
type SystemCollector struct {
	redfishClient           *redfish.APIClient
	metrics                 map[string]systemMetric
	collectorScrapeStatus   *prometheus.GaugeVec
	collectorScrapeDuration *prometheus.SummaryVec
	Log                     *log.Entry
}

// NewSystemCollector returns a collector that collecting memory statistics
func NewSystemCollector(namespace string, redfishClient *redfish.APIClient, logger *log.Entry) *SystemCollector {
	return &SystemCollector{
		redfishClient: redfishClient,
		metrics:       systemMetrics,
		Log: logger.WithFields(log.Fields{
			"collector": "SystemCollector",
		}),
		collectorScrapeStatus: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "collector_scrape_status",
				Help:      "collector_scrape_status",
			},
			[]string{"collector"},
		),
	}
}

// Describe implements prometheus.Collector.
func (s *SystemCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range s.metrics {
		ch <- metric.desc
	}
	s.collectorScrapeStatus.Describe(ch)

}

// Collect implements prometheus.Collector.
func (s *SystemCollector) Collect(ch chan<- prometheus.Metric) {
	collectorLogContext := s.Log
	//get service
	service := s.redfishClient.Service

	// get a list of systems from service
	if systems, err := service.Systems(); err != nil {
		collectorLogContext.WithField("operation", "service.Systems()").WithError(err).Error("error getting systems from service")
	} else {
		for _, system := range systems {
			systemLogContext := collectorLogContext.WithField("System", system.ID)
			systemLogContext.Info("collector scrape started")
			// overall system metrics


			// server info
			SystemID := system.ID
			SerialNumber := system.SerialNumber
			Sku := system.SKU
			if Sku!="" {
				SerialNumber = Sku
			}

			systemModel := system.Model

			systemManufacturer := "Unknown"
			if system.Manufacturer != "" {
				tmpStr := strings.Split(system.Manufacturer, " ")
				systemManufacturer = tmpStr[0]
			}
			
			//common status
			systemState := system.Status.State
			systemHealthStatus := system.Status.Health
			systemPowerState := system.PowerState

			//cpu summary
			systemProcessorSummaryState := system.ProcessorSummary.Status.State
			systemProcessorSummaryHealthStatus := system.ProcessorSummary.Status.Health
			systemProcessorSummaryCount := system.ProcessorSummary.Count

			//mem summary
			systemMemorySummaryState := system.MemorySummary.Status.State
			systemMemorySummaryHealthStatus := system.MemorySummary.Status.Health
			systemMemorySummarySize := system.MemorySummary.TotalSystemMemoryGiB

			systemLabelValues := []string{SerialNumber, systemManufacturer, "system", SystemID, systemModel}

			// system state health
			if systemStateValue, ok := parseCommonStatusState(systemState); ok {
				ch <- prometheus.MustNewConstMetric(s.metrics["system_state"].desc, prometheus.GaugeValue, systemStateValue, systemLabelValues...)
			}
			if systemHealthStatusValue, ok := parseCommonStatusHealth(systemHealthStatus); ok {
				ch <- prometheus.MustNewConstMetric(s.metrics["system_health_status"].desc, prometheus.GaugeValue, systemHealthStatusValue, systemLabelValues...)
			}
			if systemPowerStateValue, ok := parseCommonPowerState(systemPowerState); ok {
				ch <- prometheus.MustNewConstMetric(s.metrics["system_power_state"].desc, prometheus.GaugeValue, systemPowerStateValue, systemLabelValues...)
			}

			// cpu summary
			if systemProcessorSummaryStateValue, ok := parseCommonStatusState(systemProcessorSummaryState); ok {
				ch <- prometheus.MustNewConstMetric(s.metrics["system_processor_summary_state"].desc, prometheus.GaugeValue, systemProcessorSummaryStateValue, systemLabelValues...)
			}
			if systemProcessorSummaryHealthStatusValue, ok := parseCommonStatusHealth(systemProcessorSummaryHealthStatus); ok {
				ch <- prometheus.MustNewConstMetric(s.metrics["system_processor_summary_health_status"].desc, prometheus.GaugeValue, systemProcessorSummaryHealthStatusValue, systemLabelValues...)
			}
			ch <- prometheus.MustNewConstMetric(s.metrics["system_processor_summary_count"].desc, prometheus.GaugeValue, float64(systemProcessorSummaryCount), systemLabelValues...)

			// mem summary
			if systemMemorySummaryStateValue, ok := parseCommonStatusState(systemMemorySummaryState); ok {
				ch <- prometheus.MustNewConstMetric(s.metrics["system_memory_summary_state"].desc, prometheus.GaugeValue, systemMemorySummaryStateValue, systemLabelValues...)
			}
			if systemMemorySummaryHealthStatusValue, ok := parseCommonStatusHealth(systemMemorySummaryHealthStatus); ok {
				ch <- prometheus.MustNewConstMetric(s.metrics["system_memory_summary_health_status"].desc, prometheus.GaugeValue, systemMemorySummaryHealthStatusValue, systemLabelValues...)
			}
			ch <- prometheus.MustNewConstMetric(s.metrics["system_memory_summary_size"].desc, prometheus.GaugeValue, float64(systemMemorySummarySize), systemLabelValues...)


			// process processor metrics
			processors, err := system.Processors()
			if err != nil {
				systemLogContext.WithField("operation", "system.Processors()").WithError(err).Error("error getting processor data from system")
			} else if processors == nil {
				systemLogContext.WithField("operation", "system.Processors()").Info("no processor data found")
			} else {
				wg2 := &sync.WaitGroup{}
				wg2.Add(len(processors))

				for _, processor := range processors {
					go parsePorcessor(ch, SerialNumber, systemManufacturer, processor, wg2, systemLogContext)
				}
			}


			if systemManufacturer=="HPE" {
				hpStorages, err := system.SmartStorages()
				if err != nil {
					systemLogContext.WithField("operation", "system.Storage()").WithError(err).Error("error getting storage data from system")
				} else if hpStorages == nil {
					systemLogContext.WithField("operation", "system.Storage()").Info("no storage data found")
				} else {
					for _, storage := range hpStorages {
						drives, err := storage.Drives()
						if err != nil {
							systemLogContext.WithField("operation", "system.Drives()").WithError(err).Error("error getting drive data from system")
						} else if drives == nil {
							systemLogContext.WithFields(log.Fields{"operation": "system.Drives()", "storage": storage.ID}).Info("no drive data found")
						} else {
							wg4 := &sync.WaitGroup{}
							wg4.Add(len(drives))
							for _, drive := range drives {
								go parseHpDrive(ch, SerialNumber, systemManufacturer, drive, wg4, systemLogContext)
							}
						}
					}
				}
			} else if systemManufacturer=="Dell" {
				dellStorages, err := system.SimpleStorages()
				if err != nil {
					systemLogContext.WithField("operation", "system.Storage()").WithError(err).Error("error getting storage data from system")
				} else if dellStorages == nil {
					systemLogContext.WithField("operation", "system.Storage()").Info("no storage data found")
				} else {
					for _, item := range dellStorages {
						devices := item.Devices

						wg4 := &sync.WaitGroup{}
						wg4.Add(len(devices))
						for _, device := range devices {
							go parseDellDrive(ch, SerialNumber, systemManufacturer, device, wg4, systemLogContext)
						}
					}
				}
			}


			systemLogContext.Info("collector scrape completed")
		}
		s.collectorScrapeStatus.WithLabelValues("system").Set(float64(1))
	}
	
}

func parsePorcessor(ch chan<- prometheus.Metric, SerialNumber string, systemManufacturer string, processor *redfishapi.Processor, wg *sync.WaitGroup, systemLogContext *log.Entry) {
	defer func() {
		wg.Done()
        // recover from panic caused by writing to a closed channel
        if r := recover(); r != nil {
            err := fmt.Errorf("%v", r)
            log.Info(fmt.Sprintf("%s write: error writing on channel: %v\n", systemManufacturer, err))
            return
        }
    }()

	processorID := processor.ID
	processorModel := processor.Model
	processorTotalCores := processor.TotalCores
	processorTotalThreads := processor.TotalThreads
	processorState := processor.Status.State
	processorHealthStatus := processor.Status.Health

	systemProcessorLabelValues := []string{SerialNumber, "processor", processorID, processorModel}

	if processorStateValue, ok := parseCommonStatusState(processorState); ok {
		ch <- prometheus.MustNewConstMetric(systemMetrics["system_processor_state"].desc, prometheus.GaugeValue, processorStateValue, systemProcessorLabelValues...)
	}
	if processorHealthStatusValue, ok := parseCommonStatusHealth(processorHealthStatus); ok {
		ch <- prometheus.MustNewConstMetric(systemMetrics["system_processor_health_status"].desc, prometheus.GaugeValue, processorHealthStatusValue, systemProcessorLabelValues...)
	}
	ch <- prometheus.MustNewConstMetric(systemMetrics["system_processor_total_threads"].desc, prometheus.GaugeValue, float64(processorTotalThreads), systemProcessorLabelValues...)
	ch <- prometheus.MustNewConstMetric(systemMetrics["system_processor_total_cores"].desc, prometheus.GaugeValue, float64(processorTotalCores), systemProcessorLabelValues...)
}

func parseHpDrive(ch chan<- prometheus.Metric, SerialNumber string, systemManufacturer string, drive *redfishapi.Drive, wg *sync.WaitGroup, systemLogContext *log.Entry) {
	defer func() {
		wg.Done()
        // recover from panic caused by writing to a closed channel
        if r := recover(); r != nil {
            err := fmt.Errorf("%v", r)
            log.Info(fmt.Sprintf("%s write: error writing on channel: %v\n", systemManufacturer, err))
            return
        }
    }()

	driveName := drive.Location
	driveModel := drive.Model
	driveCapacityGB := drive.CapacityGB
	driveHealthStatus := drive.Status.Health


	systemdriveLabelValues := []string{SerialNumber, "drive", driveName, driveModel}

	if driveHealthStatusValue, ok := parseCommonStatusHealth(driveHealthStatus); ok {
		ch <- prometheus.MustNewConstMetric(systemMetrics["system_storage_drive_health_state"].desc, prometheus.GaugeValue, driveHealthStatusValue, systemdriveLabelValues...)

	}
	ch <- prometheus.MustNewConstMetric(systemMetrics["system_storage_drive_capacity"].desc, prometheus.GaugeValue, float64(driveCapacityGB), systemdriveLabelValues...)

}


func parseDellDrive(ch chan<- prometheus.Metric, SerialNumber string, systemManufacturer string, device redfishapi.Device, wg *sync.WaitGroup, systemLogContext *log.Entry) {
	defer func() {
		wg.Done()
        // recover from panic caused by writing to a closed channel
        if r := recover(); r != nil {
            err := fmt.Errorf("%v", r)
            log.Info(fmt.Sprintf("%s write: error writing on channel: %v\n", systemManufacturer, err))
            return
        }
    }()

	driveName := device.Name
	driveModel := device.Model
	driveCapacityGB := device.CapacityBytes / 1024 / 1024
	driveHealthStatus := device.Status.Health


	systemdriveLabelValues := []string{SerialNumber, "drive", driveName, driveModel}

	if driveHealthStatusValue, ok := parseCommonStatusHealth(driveHealthStatus); ok {
		ch <- prometheus.MustNewConstMetric(systemMetrics["system_storage_drive_health_state"].desc, prometheus.GaugeValue, driveHealthStatusValue, systemdriveLabelValues...)

	}
	ch <- prometheus.MustNewConstMetric(systemMetrics["system_storage_drive_capacity"].desc, prometheus.GaugeValue, float64(driveCapacityGB), systemdriveLabelValues...)

}
