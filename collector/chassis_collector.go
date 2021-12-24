package collector

import (
	"sync"
	"strings"

	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/magicst0ne/rackserver_exporter/redfish"
	"github.com/magicst0ne/rackserver_exporter/redfish/redfishapi"
)

// ChassisSubsystem is the chassis subsystem
var (
	ChassisSubsystem                  = "chassis"
	ChassisLabelNames                 = []string{"sn", "mfr","resource", "chassis_id"}
	ChassisTemperatureLabelNames      = []string{"sn", "mfr","resource", "chassis_id", "sensor", "sensor_id"}
	ChassisFanLabelNames              = []string{"sn", "mfr","resource", "chassis_id", "fan", "fan_id"}
	ChassisPowerSupplyLabelNames      = []string{"sn", "mfr","resource", "chassis_id", "power_supply", "power_supply_id"}

	chassisMetrics = map[string]chassisMetric{
		"chassis_health": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, ChassisSubsystem, "health"),
				"health of chassis, 1(OK),2(Warning),3(Critical)",
				ChassisLabelNames,
				nil,
			),
		},
		"chassis_state": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, ChassisSubsystem, "state"),
				"state of chassis,1(Enabled),2(Disabled),3(StandbyOffinline),4(StandbySpare),5(InTest),6(Starting),7(Absent),8(UnavailableOffline),9(Deferring),10(Quiesced),11(Updating)",
				ChassisLabelNames,
				nil,
			),
		},
		"chassis_temperature_sensor_state": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, ChassisSubsystem, "temperature_sensor_state"),
				"status state of temperature on this chassis component,1(Enabled),2(Disabled),3(StandbyOffinline),4(StandbySpare),5(InTest),6(Starting),7(Absent),8(UnavailableOffline),9(Deferring),10(Quiesced),11(Updating)",
				ChassisTemperatureLabelNames,
				nil,
			),
		},
		"chassis_temperature_celsius": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, ChassisSubsystem, "temperature_celsius"),
				"celsius of temperature on this chassis component",
				ChassisTemperatureLabelNames,
				nil,
			),
		},
		"chassis_fan_health": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, ChassisSubsystem, "fan_health"),
				"fan health on this chassis component,1(OK),2(Warning),3(Critical)",
				ChassisFanLabelNames,
				nil,
			),
		},
		"chassis_fan_state": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, ChassisSubsystem, "fan_state"),
				"fan state on this chassis component,1(Enabled),2(Disabled),3(StandbyOffinline),4(StandbySpare),5(InTest),6(Starting),7(Absent),8(UnavailableOffline),9(Deferring),10(Quiesced),11(Updating)",
				ChassisFanLabelNames,
				nil,
			),
		},
		"chassis_fan_rpm": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, ChassisSubsystem, "fan_rpm_percentage"),
				"fan rpm percentage on this chassis component",
				ChassisFanLabelNames,
				nil,
			),
		},
		"chassis_power_powersupply_state": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, ChassisSubsystem, "power_powersupply_state"),
				"powersupply state of chassis component,1(Enabled),2(Disabled),3(StandbyOffinline),4(StandbySpare),5(InTest),6(Starting),7(Absent),8(UnavailableOffline),9(Deferring),10(Quiesced),11(Updating)",
				ChassisPowerSupplyLabelNames,
				nil,
			),
		},
		"chassis_power_powersupply_health_status": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, ChassisSubsystem, "power_powersupply_health_status"),
				"powersupply health of chassis component,1(OK),2(Warning),3(Critical)",
				ChassisPowerSupplyLabelNames,
				nil,
			),
		},
		"chassis_power_powersupply_last_power_output_watts": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, ChassisSubsystem, "power_powersupply_last_power_output_watts"),
				"last_power_output_watts of powersupply on this chassis",
				ChassisPowerSupplyLabelNames,
				nil,
			),
		},
		"chassis_power_powersupply_power_capacity_watts": {
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, ChassisSubsystem, "power_powersupply_power_capacity_watts"),
				"power_capacity_watts of powersupply on this chassis",
				ChassisPowerSupplyLabelNames,
				nil,
			),
		},
	}
)

//ChassisCollector implements the prometheus.Collector.
type ChassisCollector struct {
	redfishClient         *redfish.APIClient
	metrics               map[string]chassisMetric
	collectorScrapeStatus *prometheus.GaugeVec
	Log                   *log.Entry
}

type chassisMetric struct {
	desc *prometheus.Desc
}

// NewChassisCollector returns a collector that collecting chassis statistics
func NewChassisCollector(namespace string, redfishClient *redfish.APIClient, logger *log.Entry) *ChassisCollector {
	// get service from redfish client

	return &ChassisCollector{
		redfishClient: redfishClient,
		metrics:       chassisMetrics,
		Log: logger.WithFields(log.Fields{
			"collector": "ChassisCollector",
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

// Describe implemented prometheus.Collector
func (c *ChassisCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric.desc
	}
	c.collectorScrapeStatus.Describe(ch)

}

// Collect implemented prometheus.Collector
func (c *ChassisCollector) Collect(ch chan<- prometheus.Metric) {
	collectorLogContext := c.Log
	service := c.redfishClient.Service

	// get a list of chassis from service
	if chassises, err := service.Chassis(); err != nil {
		collectorLogContext.WithField("operation", "service.Chassis()").WithError(err).Error("error getting chassis from service")
	} else {
		// process the chassises
		for _, chassis := range chassises {
			chassisLogContext := collectorLogContext.WithField("Chassis", chassis.ID)
			chassisLogContext.Info("collector scrape started")

			SerialNumber := chassis.SerialNumber
			Sku := chassis.SKU
			if Sku!="" {
				SerialNumber = Sku
			}

			systemManufacturer := "Unknown"
			if chassis.Manufacturer != "" {
				tmpStr := strings.Split(chassis.Manufacturer, " ")
				systemManufacturer = tmpStr[0]
			}

			chassisID := chassis.ID
			chassisStatus := chassis.Status
			chassisStatusState := chassisStatus.State
			chassisStatusHealth := chassisStatus.Health
			ChassisLabelValues := []string{SerialNumber, systemManufacturer, "chassis", chassisID}

			if chassisStatusHealthValue, ok := parseCommonStatusHealth(chassisStatusHealth); ok {
				ch <- prometheus.MustNewConstMetric(c.metrics["chassis_health"].desc, prometheus.GaugeValue, chassisStatusHealthValue, ChassisLabelValues...)
			}
			if chassisStatusStateValue, ok := parseCommonStatusState(chassisStatusState); ok {
				ch <- prometheus.MustNewConstMetric(c.metrics["chassis_state"].desc, prometheus.GaugeValue, chassisStatusStateValue, ChassisLabelValues...)
			}

			chassisThermal, err := chassis.Thermal()
			if err != nil {
				chassisLogContext.WithField("operation", "chassis.Thermal()").WithError(err).Error("error getting thermal data from chassis")
			} else if chassisThermal == nil {
				chassisLogContext.WithField("operation", "chassis.Thermal()").Info("no thermal data found")
			} else {
				// process temperature
				chassisTemperatures := chassisThermal.Temperatures
				wg := &sync.WaitGroup{}
				wg.Add(len(chassisTemperatures))

				for _, chassisTemperature := range chassisTemperatures {
					go parseChassisTemperature(ch, SerialNumber, systemManufacturer, chassisID, chassisTemperature, wg)
				}

				// process fans
				chassisFans := chassisThermal.Fans
				wg2 := &sync.WaitGroup{}
				wg2.Add(len(chassisFans))
				for _, chassisFan := range chassisFans {
					go parseChassisFan(ch, SerialNumber, systemManufacturer, chassisID, chassisFan, wg2)
				}
			}


			chassisPowerInfo, err := chassis.Power()
			if err != nil {
				chassisLogContext.WithField("operation", "chassis.Power()").WithError(err).Error("error getting power data from chassis")
			} else if chassisPowerInfo == nil {
				chassisLogContext.WithField("operation", "chassis.Power()").Info("no power data found")
			} else {
				// powerSupply
				chassisPowerInfoPowerSupplies := chassisPowerInfo.PowerSupplies
				wg5 := &sync.WaitGroup{}
				wg5.Add(len(chassisPowerInfoPowerSupplies))
				for _, chassisPowerInfoPowerSupply := range chassisPowerInfoPowerSupplies {
					go parseChassisPowerInfoPowerSupply(ch, SerialNumber, systemManufacturer, chassisID, chassisPowerInfoPowerSupply, wg5)
				}
			}
			chassisLogContext.Info("collector scrape completed")

		}
	}

	c.collectorScrapeStatus.WithLabelValues("chassis").Set(float64(1))
}

func parseChassisTemperature(ch chan<- prometheus.Metric, SerialNumber, systemManufacturer, chassisID string, chassisTemperature redfishapi.Temperature, wg *sync.WaitGroup) {
	defer wg.Done()
	chassisTemperatureSensorName := chassisTemperature.Name
	chassisTemperatureSensorID := chassisTemperature.MemberID
	chassisTemperatureStatus := chassisTemperature.Status
	//			chassisTemperatureStatusHealth :=chassisTemperatureStatus.Health
	chassisTemperatureStatusState := chassisTemperatureStatus.State
	//			chassisTemperatureStatusLabelNames :=[]string{BaseLabelNames,"temperature_sensor_name","temperature_sensor_member_id")
	chassisTemperatureLabelvalues := []string{SerialNumber, systemManufacturer, "temperature", chassisID, chassisTemperatureSensorName, chassisTemperatureSensorID}

	//		ch <- prometheus.MustNewConstMetric(chassisMetrics["chassis_temperature_status_health"].desc, prometheus.GaugeValue, parseCommonStatusHealth(chassisTemperatureStatusHealth), chassisTemperatureLabelvalues...)
	if chassisTemperatureStatusStateValue, ok := parseCommonStatusState(chassisTemperatureStatusState); ok {
		ch <- prometheus.MustNewConstMetric(chassisMetrics["chassis_temperature_sensor_state"].desc, prometheus.GaugeValue, chassisTemperatureStatusStateValue, chassisTemperatureLabelvalues...)
	}

	chassisTemperatureReadingCelsius := chassisTemperature.ReadingCelsius
	ch <- prometheus.MustNewConstMetric(chassisMetrics["chassis_temperature_celsius"].desc, prometheus.GaugeValue, float64(chassisTemperatureReadingCelsius), chassisTemperatureLabelvalues...)
}

func parseChassisFan(ch chan<- prometheus.Metric, SerialNumber, systemManufacturer, chassisID string, chassisFan redfishapi.Fan, wg *sync.WaitGroup) {
	defer wg.Done()
	chassisFanID := chassisFan.MemberID
	chassisFanName := chassisFan.Name
	chassisFanStaus := chassisFan.Status
	chassisFanStausHealth := chassisFanStaus.Health
	chassisFanStausState := chassisFanStaus.State
	chassisFanRPM := chassisFan.Reading

	chassisFanLabelvalues := []string{SerialNumber, systemManufacturer, "fan", chassisID, chassisFanName, chassisFanID}

	if chassisFanStausHealthValue, ok := parseCommonStatusHealth(chassisFanStausHealth); ok {
		ch <- prometheus.MustNewConstMetric(chassisMetrics["chassis_fan_health"].desc, prometheus.GaugeValue, chassisFanStausHealthValue, chassisFanLabelvalues...)
	}

	if chassisFanStausStateValue, ok := parseCommonStatusState(chassisFanStausState); ok {
		ch <- prometheus.MustNewConstMetric(chassisMetrics["chassis_fan_state"].desc, prometheus.GaugeValue, chassisFanStausStateValue, chassisFanLabelvalues...)
	}
	ch <- prometheus.MustNewConstMetric(chassisMetrics["chassis_fan_rpm"].desc, prometheus.GaugeValue, float64(chassisFanRPM), chassisFanLabelvalues...)

}

func parseChassisPowerInfoPowerSupply(ch chan<- prometheus.Metric, SerialNumber, systemManufacturer, chassisID string, chassisPowerInfoPowerSupply redfishapi.PowerSupply, wg *sync.WaitGroup) {

	defer wg.Done()
	chassisPowerInfoPowerSupplyName := chassisPowerInfoPowerSupply.Name
	chassisPowerInfoPowerSupplyID := chassisPowerInfoPowerSupply.MemberID

	if chassisPowerInfoPowerSupplyID=="" {
		chassisPowerInfoPowerSupplyID = chassisPowerInfoPowerSupply.SerialNumber
	}

	chassisPowerInfoPowerSupplyPowerCapacityWatts := chassisPowerInfoPowerSupply.PowerCapacityWatts
	chassisPowerInfoPowerSupplyLastPowerOutputWatts := chassisPowerInfoPowerSupply.LastPowerOutputWatts
	chassisPowerInfoPowerSupplyState := chassisPowerInfoPowerSupply.Status.State
	chassisPowerInfoPowerSupplyHealthStatus := chassisPowerInfoPowerSupply.Status.Health
	chassisPowerSupplyLabelvalues := []string{SerialNumber, systemManufacturer, "power_supply", chassisID, chassisPowerInfoPowerSupplyName, chassisPowerInfoPowerSupplyID}
	if chassisPowerInfoPowerSupplyStateValue, ok := parseCommonStatusState(chassisPowerInfoPowerSupplyState); ok {
		ch <- prometheus.MustNewConstMetric(chassisMetrics["chassis_power_powersupply_state"].desc, prometheus.GaugeValue, chassisPowerInfoPowerSupplyStateValue, chassisPowerSupplyLabelvalues...)
	}
	if chassisPowerInfoPowerSupplyHealthStatusValue, ok := parseCommonStatusHealth(chassisPowerInfoPowerSupplyHealthStatus); ok {
		ch <- prometheus.MustNewConstMetric(chassisMetrics["chassis_power_powersupply_health_status"].desc, prometheus.GaugeValue, chassisPowerInfoPowerSupplyHealthStatusValue, chassisPowerSupplyLabelvalues...)
	}
	ch <- prometheus.MustNewConstMetric(chassisMetrics["chassis_power_powersupply_last_power_output_watts"].desc, prometheus.GaugeValue, float64(chassisPowerInfoPowerSupplyLastPowerOutputWatts), chassisPowerSupplyLabelvalues...)
	ch <- prometheus.MustNewConstMetric(chassisMetrics["chassis_power_powersupply_power_capacity_watts"].desc, prometheus.GaugeValue, float64(chassisPowerInfoPowerSupplyPowerCapacityWatts), chassisPowerSupplyLabelvalues...)
}