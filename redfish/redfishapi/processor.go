package redfishapi

import (
	"encoding/json"
	"strconv"

	"github.com/magicst0ne/rackserver_exporter/redfish/common"
)


// Processor is used to represent a single processor contained within a
// system.
type Processor struct {
	common.Entity

	// ODataContext is the odata context.
	ODataContext string `json:"@odata.context"`
	// ODataType is the odata type.
	ODataType string `json:"@odata.type"`
	// Description provides a description of this resource.
	Description string
	// Manufacturer shall contain a string which identifies
	// the manufacturer of the processor.
	Manufacturer string
	// MaxSpeedMHz shall indicate the maximum rated clock
	// speed of the processor in MHz.
	MaxSpeedMHz float32
	// MaxTDPWatts shall be the maximum Thermal
	// Design Power (TDP) in watts.
	MaxTDPWatts int
	// Metrics shall be a reference to the Metrics
	// associated with this Processor.
	metrics string
	// Model shall indicate the model information as
	// provided by the manufacturer of this processor.
	Model string
	// ProcessorID shall contain identification information for this processor.
	ProcessorID ProcessorID `json:"ProcessorId"`
	// Socket shall contain the string which identifies the
	// physical location or socket of the processor.
	//czw remove
	//Socket string
	// Status shall contain any status or health properties
	// of the resource.
	Status common.Status
	// TotalCores shall indicate the total count of
	// independent processor cores contained within this processor.
	TotalCores int
	// TotalEnabledCores shall indicate the total count of
	// enabled independent processor cores contained within this processor.
	TotalEnabledCores int
	// TotalThreads shall indicate the total count of
	// independent execution threads supported by this processor.
	TotalThreads int
}

// UnmarshalJSON unmarshals a Processor object from the raw JSON.
func (processor *Processor) UnmarshalJSON(b []byte) error {
	type temp Processor
	type t1 struct {
		temp
	}
	var t t1

	err := json.Unmarshal(b, &t)
	if err != nil {
		// Handle invalid data type returned for MaxSpeedMHz
		var t2 struct {
			t1
			MaxSpeedMHz string
		}
		err2 := json.Unmarshal(b, &t2)

		if err2 != nil {
			// Return the original error
			return err
		}

		// Extract the real Processor struct and replace its MaxSpeedMHz with
		// the parsed string version
		t = t2.t1
		if t2.MaxSpeedMHz != "" {
			bitSize := 32
			mhz, err := strconv.ParseFloat(t2.MaxSpeedMHz, bitSize)
			if err != nil {
				t.MaxSpeedMHz = float32(mhz)
			}
		}
	}

	*processor = Processor(t.temp)

	return nil
}

// GetProcessor will get a Processor instance from the system
func GetProcessor(c common.Client, uri string) (*Processor, error) {
	resp, err := c.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var processor Processor
	err = json.NewDecoder(resp.Body).Decode(&processor)
	if err != nil {
		return nil, err
	}

	processor.SetClient(c)
	return &processor, nil
}

// ListReferencedProcessors gets the collection of Processor from a provided reference.
func ListReferencedProcessors(c common.Client, link string) ([]*Processor, error) {
	var result []*Processor
	links, err := common.GetCollection(c, link)
	if err != nil {
		return result, err
	}

	collectionError := common.NewCollectionError()
	for _, processorLink := range links.ItemLinks {
		processor, err := GetProcessor(c, processorLink)
		if err != nil {
			collectionError.Failures[processorLink] = err
		} else {
			result = append(result, processor)
		}
	}

	if collectionError.Empty() {
		return result, nil
	}

	return result, collectionError
}

// ProcessorID shall contain identification information for a processor.
type ProcessorID struct {
	// EffectiveFamily shall indicate the effective Family
	// information as provided by the manufacturer of this processor.
	EffectiveFamily string
	// EffectiveModel shall indicate the effective Model
	// information as provided by the manufacturer of this processor.
	EffectiveModel string
	// IdentificationRegisters shall include the raw CPUID
	// instruction output as provided by the manufacturer of this processor.
	IdentificationRegisters string
	// MicrocodeInfo shall indicate the Microcode
	// Information as provided by the manufacturer of this processor.
	MicrocodeInfo string
	// Step shall indicate the Step or revision string
	// information as provided by the manufacturer of this processor.
	Step string
	// VendorID shall indicate the Vendor Identification
	// string information as provided by the manufacturer of this processor.
	VendorID string `json:"VendorId"`
}
