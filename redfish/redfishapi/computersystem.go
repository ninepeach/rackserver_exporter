package redfishapi

import (
	"encoding/json"
	"strings"
	"reflect"
	"github.com/magicst0ne/rackserver_exporter/redfish/common"
)

// PowerState is the power state of the system.
type PowerState string

const (

	// OnPowerState the system is powered on.
	OnPowerState PowerState = "On"
	// OffPowerState the system is powered off, although some components may
	// continue to have AUX power such as management controller.
	OffPowerState PowerState = "Off"
	// PoweringOnPowerState A temporary state between Off and On. This
	// temporary state can be very short.
	PoweringOnPowerState PowerState = "PoweringOn"
	// PoweringOffPowerState A temporary state between On and Off. The power
	// off action can take time while the OS is in the shutdown process.
	PoweringOffPowerState PowerState = "PoweringOff"
)

// ComputerSystem is used to represent resources that represent a
// computing system in the Redfish specification.
type ComputerSystem struct {
	common.Entity

	// ODataContext is the @odata.context
	ODataContext string `json:"@odata.context"`
	// ODataType is the @odata.type
	ODataType string `json:"@odata.type"`

	// Description is the resource description.
	Description string
	// EthernetInterfaces shall be a link to a
	// collection of type EthernetInterfaceCollection.
	ethernetInterfaces string

	// Manufacturer shall contain a value that represents the manufacturer of the system.
	Manufacturer string
	// Memory shall be a link to a collection of type MemoryCollection.
	memory string

	MemorySummary MemorySummary
	// Model shall contain the information
	// about how the manufacturer references this system. This is typically
	// the product name, without the manufacturer name.
	Model string
	// Name is the resource name.
	Name string
	// PowerState shall contain the power state of the system.
	PowerState PowerState
	// ProcessorSummary shall contain properties which
	// describe the central processors for the current resource.
	ProcessorSummary ProcessorSummary
	// Processors shall be a link to a collection of type ProcessorCollection.
	processors string
	// Redundancy references a redundancy
	// entity that specifies a kind and level of redundancy and a collection
	// (RedundancySet) of other ComputerSystems that provide the specified
	// redundancy to this ComputerSystem.
	Redundancy string
	// RedundancyCount is the number of Redundancy objects.
	RedundancyCount string `json:"Redundancy@odata.count"`
	SerialNumber string
	// SimpleStorage shall be a link to a collection of type SimpleStorageCollection.
	simpleStorage string
	smartStorage string
	// Status shall contain any status or health properties
	// of the resource.
	Status common.Status
	BatteryHealth string
	Oem json.RawMessage
	// rawData holds the original serialized JSON so we can compare updates.
	rawData []byte
}

// UnmarshalJSON unmarshals a ComputerSystem object from the raw JSON.
func (computersystem *ComputerSystem) UnmarshalJSON(b []byte) error {

	type temp ComputerSystem

	type t_links struct {
		Processors common.HpLink
	}

	var t struct {
		temp
		Processors         common.Link
		Memory             common.Link
		SimpleStorage      common.Link
		Links              t_links `json:"links"`
	}

	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}

	*computersystem = ComputerSystem(t.temp)

	// Extract the links to other entities for later
	computersystem.processors = string(t.Processors)
	computersystem.memory = string(t.Memory)
	computersystem.simpleStorage = string(t.SimpleStorage)


    if computersystem.Manufacturer != "" {
        tmpStr := strings.Split(computersystem.Manufacturer, " ")
        computersystem.Manufacturer = tmpStr[0]
    } else {
    	computersystem.Manufacturer = "Unknown"
    }

    if t.Oem!=nil {
        jsonMap := make(map[string]interface{})
        json.Unmarshal(t.Oem, &jsonMap)

        if hpMap, ok := jsonMap["Hp"]; ok {
            if links, ok1 := hpMap.(map[string]interface{})["links"]; ok1 {
                for k, v := range links.(map[string]interface{}) {
                        if link, ok2 := v.(map[string]interface{})["href"]; ok2 {
                            if k=="SmartStorage" {
                                computersystem.smartStorage = link.(string)
                            } else if k=="Memory" {
                                computersystem.memory = link.(string)
                            }
                        }
                }
            }

            if battery, ok1 := hpMap.(map[string]interface{})["Battery"]; ok1 {
                for _, item  := range battery.([]interface{}) {
                    if battery_health, ok2 := item.(map[string]interface{})["Condition"]; ok2 {
                        computersystem.BatteryHealth = battery_health.(string)
                    }
                }
            }
        }
    }

    if computersystem.processors=="" {
    	computersystem.processors = string(t.Links.Processors)
    }

	// This is a read/write object, so we need to save the raw object data for later
	computersystem.rawData = b

	return nil
}

// Update commits updates to this object's properties to the running system.
func (computersystem *ComputerSystem) Update() error {
	// Get a representation of the object's original state so we can find what
	// to update.
	cs := new(ComputerSystem)
	err := cs.UnmarshalJSON(computersystem.rawData)
	if err != nil {
		return err
	}

	readWriteFields := []string{
	}

	originalElement := reflect.ValueOf(cs).Elem()
	currentElement := reflect.ValueOf(computersystem).Elem()

	return computersystem.Entity.Update(originalElement, currentElement, readWriteFields)
}

// GetComputerSystem will get a ComputerSystem instance from the service.
func GetComputerSystem(c common.Client, uri string) (*ComputerSystem, error) {
	resp, err := c.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var computersystem ComputerSystem
	err = json.NewDecoder(resp.Body).Decode(&computersystem)
	if err != nil {
		return nil, err
	}

	computersystem.SetClient(c)
	return &computersystem, nil
}

// ListReferencedComputerSystems gets the collection of ComputerSystem from
// a provided reference.
func ListReferencedComputerSystems(c common.Client, link string) ([]*ComputerSystem, error) {
	var result []*ComputerSystem
	links, err := common.GetCollection(c, link)
	if err != nil {
		return result, err
	}

	collectionError := common.NewCollectionError()
	for _, computersystemLink := range links.ItemLinks {
		computersystem, err := GetComputerSystem(c, computersystemLink)
		if err != nil {
			collectionError.Failures[computersystemLink] = err
		} else {
			result = append(result, computersystem)
		}
	}

	if collectionError.Empty() {
		return result, nil
	}

	return result, collectionError
}

// MemorySummary contains properties which describe the central memory for a system.
type MemorySummary struct {
	// Status is the status or health properties of the resource.
	Status common.Status
	// TotalSystemMemoryGiB is the amount of configured system general purpose
	// volatile (RAM) memory as measured in gibibytes.
	TotalSystemMemoryGiB float32
	// TotalSystemPersistentMemoryGiB is the total amount of configured
	// persistent memory available to the system as measured in gibibytes.
	TotalSystemPersistentMemoryGiB float32
}

// ProcessorSummary is This type shall contain properties which describe
// the central processors for a system.
type ProcessorSummary struct {
	// Count is the number of physical central processors in the system.
	Count int
	// LogicalProcessorCount is the number of logical central processors in the system.
	LogicalProcessorCount int
	// Model is the processor model for the central processors in the system,
	// per the description in the Processor Information - Processor Family
	// section of the SMBIOS Specification DSP0134 2.8 or later.
	Model string
	// Status is any status or health properties of the resource.
	Status common.Status
}

// Processors returns a collection of processors from this system
func (computersystem *ComputerSystem) Processors() ([]*Processor, error) {
        return ListReferencedProcessors(computersystem.Client, computersystem.processors)
}

// SimpleStorages gets all simple storage services of this system.
func (computersystem *ComputerSystem) SimpleStorages() ([]*SimpleStorage, error) {
        return ListReferencedSimpleStorages(computersystem.Client, computersystem.simpleStorage)
}

// Storage gets the storage associated with this system.
func (computersystem *ComputerSystem) SmartStorages() ([]*SmartStorage, error) {
        return ListReferencedSmartStorages(computersystem.Client, computersystem.smartStorage)
}