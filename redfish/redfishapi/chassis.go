package redfishapi

import (
	"encoding/json"
	"reflect"

	"github.com/magicst0ne/rackserver_exporter/redfish/common"
)

// ChassisType is a physical form of the chassis
type ChassisType string

const (
	// BladeChassisType is an enclosed or semi-enclosed, typically vertically-oriented, system
	// chassis which must be plugged into a multi-system chassis to function normally.
	BladeChassisType ChassisType = "Blade"
	// CardChassisType is a loose device or circuit board intended to be installed in a system
	// or other enclosure.
	CardChassisType ChassisType = "Card"
	// CartridgeChassisType is a small self-contained system intended to be plugged into a multi-system
	// chassis.
	CartridgeChassisType ChassisType = "Cartridge"
	// ComponentChassisType is a small chassis, card, or device which contains devices for a particular
	// subsystem or function.
	ComponentChassisType ChassisType = "Component"
	// DrawerChassisType is an enclosed or semi-enclosed, typically horizontally-oriented, system
	// chassis which may be slid into a multi-system chassis.
	DrawerChassisType ChassisType = "Drawer"
	// EnclosureChassisType is a generic term for a chassis that does not fit any other description.
	EnclosureChassisType ChassisType = "Enclosure"
	// ExpansionChassisType is a chassis which expands the capabilities or capacity of another
	// chassis.
	ExpansionChassisType ChassisType = "Expansion"
	// IPBasedDriveChassisType is a chassis in a drive form factor with IP-based network connections.
	IPBasedDriveChassisType ChassisType = "IPBasedDrive"
	// ModuleChassisType is a small, typically removable, chassis or card which contains devices
	// for a particular subsystem or function.
	ModuleChassisType ChassisType = "Module"
	// OtherChassisType is a chassis that does not fit any of these definitions.
	OtherChassisType ChassisType = "Other"
	// PodChassisType is a collection of equipment racks in a large, likely transportable, container.
	PodChassisType ChassisType = "Pod"
	// RackChassisType is an equipment rack, typically a 19-inch wide freestanding unit.
	RackChassisType ChassisType = "Rack"
	// RackGroupChassisType is a group of racks which form a single entity or share infrastructure.
	RackGroupChassisType ChassisType = "RackGroup"
	// RackMountChassisType is a single system chassis designed specifically for mounting in an
	// equipment rack.
	RackMountChassisType ChassisType = "RackMount"
	// RowChassisType is a collection of equipment racks.
	RowChassisType ChassisType = "Row"
	// ShelfChassisType is an enclosed or semi-enclosed, typically horizontally-oriented, system
	// chassis which must be plugged into a multi-system chassis to function normally.
	ShelfChassisType ChassisType = "Shelf"
	// SidecarChassisType is a chassis that mates mechanically with another chassis to expand
	// its capabilities or capacity.
	SidecarChassisType ChassisType = "Sidecar"
	// SledChassisType is an enclosed or semi-enclosed, system chassis which must be plugged into
	// a multi-system chassis to function normally similar to a blade type chassis.
	SledChassisType ChassisType = "Sled"
	// StandAloneChassisType is a single, free-standing system, commonly called a tower or desktop
	// chassis.
	StandAloneChassisType ChassisType = "StandAlone"
	// StorageEnclosureChassisType is a chassis which encloses storage.
	StorageEnclosureChassisType ChassisType = "StorageEnclosure"
	// ZoneChassisType is a logical division or portion of a physical chassis that contains multiple
	// devices or systems that cannot be physically separated.
	ZoneChassisType ChassisType = "Zone"
)

// EnvironmentalClass is
type EnvironmentalClass string

const (

	// A1EnvironmentalClass ASHRAE Environmental Class 'A1'.
	A1EnvironmentalClass EnvironmentalClass = "A1"
	// A2EnvironmentalClass ASHRAE Environmental Class 'A2'.
	A2EnvironmentalClass EnvironmentalClass = "A2"
	// A3EnvironmentalClass ASHRAE Environmental Class 'A3'.
	A3EnvironmentalClass EnvironmentalClass = "A3"
	// A4EnvironmentalClass ASHRAE Environmental Class 'A4'.
	A4EnvironmentalClass EnvironmentalClass = "A4"
)

// IntrusionSensor is
type IntrusionSensor string

const (

	// NormalIntrusionSensor No abnormal physical security condition is
	// detected at this time.
	NormalIntrusionSensor IntrusionSensor = "Normal"
	// HardwareIntrusionIntrusionSensor A door, lock, or other mechanism
	// protecting the internal system hardware from being accessed is
	// detected to be in an insecure state.
	HardwareIntrusionIntrusionSensor IntrusionSensor = "HardwareIntrusion"
	// TamperingDetectedIntrusionSensor Physical tampering of the monitored
	// entity is detected.
	TamperingDetectedIntrusionSensor IntrusionSensor = "TamperingDetected"
)

// IntrusionSensorReArm is
type IntrusionSensorReArm string

const (

	// ManualIntrusionSensorReArm A manual re-arm of this sensor restores it
	// to the normal state.
	ManualIntrusionSensorReArm IntrusionSensorReArm = "Manual"
	// AutomaticIntrusionSensorReArm Because no abnormal physical security
	// condition is detected, this sensor is automatically restored to the
	// normal state.
	AutomaticIntrusionSensorReArm IntrusionSensorReArm = "Automatic"
)

// PhysicalSecurity shall describe the sensor state of the physical
// security.
type PhysicalSecurity struct {

	// IntrusionSensor is This property shall represent the state of this
	// physical security sensor.  Hardware intrusion indicates the internal
	// hardware is detected as being accessed in an insecure state.
	// Tampering detected indicates the physical tampering of the monitored
	// entity is detected.
	IntrusionSensor IntrusionSensor
	// IntrusionSensorNumber is This property shall contain a numerical
	// identifier for this physical security sensor that is unique within
	// this resource.
	IntrusionSensorNumber int
	// IntrusionSensorReArm is This property shall represent the method that
	// restores this physical security sensor to the normal state.  Manual
	// indicates manual re-arm is needed.  Automatic indicates the state is
	// restored automatically because no abnormal physical security
	// conditions are detected.
	IntrusionSensorReArm IntrusionSensorReArm
}

// Chassis represents the physical components of a system. This
// resource represents the sheet-metal confined spaces and logical zones such
// as racks, enclosures, chassis and all other containers. Subsystems (like sensors)
// that operate outside of a system's data plane (meaning the resources are not
// accessible to software running on the system) are linked either directly or
// indirectly through this resource.
type Chassis struct {
	common.Entity

	// ODataContext is the odata context.
	ODataContext string `json:"@odata.context"`
	// ODataType is the odata type.
	ODataType string `json:"@odata.type"`
	// AssetTag shall contain an identifying string that
	// tracks the chassis for inventory purposes.
	AssetTag string
	// ChassisType shall indicate the physical form factor
	// for the type of chassis.
	ChassisType ChassisType
	// Description provides a description of this resource.
	Description string
	// Location shall contain location information of the
	// associated chassis.
	Location common.Location
	// Manufacturer shall contain the name of the
	// organization responsible for producing the chassis. This organization
	// might be the entity from whom the chassis is purchased, but this is
	// not necessarily true.
	Manufacturer string
	// Model shall contain the name by which the
	// manufacturer generally refers to the chassis.
	Model string
	// PartNumber shall contain a part number assigned by
	// the organization that is responsible for producing or manufacturing
	// the chassis.
	PartNumber string
	// PowerState shall contain the power state of the
	// chassis.
	PowerState string
	// SerialNumber shall contain a manufacturer-allocated
	// number that identifies the chassis.
	SerialNumber string
	// Status shall contain any status or health properties
	// of the resource.
	Status common.Status
	// UUID shall contain the universal unique identifier
	// number for this chassis.
	UUID string
	thermal         string
	power           string
	rawData []byte
}

// UnmarshalJSON unmarshals a Chassis object from the raw JSON.
func (chassis *Chassis) UnmarshalJSON(b []byte) error {
	type temp Chassis
	type linkReference struct {
	}

	var t struct {
		temp
		Thermal         common.Link
		Power           common.Link
		Links           linkReference
	}

	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}

	*chassis = Chassis(t.temp)

	// Extract the links to other entities for later

	chassis.thermal = string(t.Thermal)
	chassis.power = string(t.Power)

	// This is a read/write object, so we need to save the raw object data for later
	chassis.rawData = b

	return nil
}

// Update commits updates to this object's properties to the running system.
func (chassis *Chassis) Update() error {
	// Get a representation of the object's original state so we can find what
	// to update.
	original := new(Chassis)
	err := original.UnmarshalJSON(chassis.rawData)
	if err != nil {
		return err
	}

	readWriteFields := []string{
		"AssetTag",
		"IndicatorLED",
	}

	originalElement := reflect.ValueOf(original).Elem()
	currentElement := reflect.ValueOf(chassis).Elem()

	return chassis.Entity.Update(originalElement, currentElement, readWriteFields)
}

// GetChassis will get a Chassis instance from the Redfish service.
func GetChassis(c common.Client, uri string) (*Chassis, error) {
	resp, err := c.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var chassis Chassis
	err = json.NewDecoder(resp.Body).Decode(&chassis)
	if err != nil {
		return nil, err
	}

	chassis.SetClient(c)
	return &chassis, nil
}

// ListReferencedChassis gets the collection of Chassis from a provided reference.
func ListReferencedChassis(c common.Client, link string) ([]*Chassis, error) {
	var result []*Chassis
	links, err := common.GetCollection(c, link)
	if err != nil {
		return result, err
	}

	collectionError := common.NewCollectionError()
	for _, chassisLink := range links.ItemLinks {
		chassis, err := GetChassis(c, chassisLink)
		if err != nil {
			collectionError.Failures[chassisLink] = err
		} else {
			result = append(result, chassis)
		}
	}

	if collectionError.Empty() {
		return result, nil
	}

	return result, collectionError
}

// Thermal gets the thermal temperature and cooling information for the chassis
func (chassis *Chassis) Thermal() (*Thermal, error) {
	if chassis.thermal == "" {
		return nil, nil
	}

	thermal, err := GetThermal(chassis.Client, chassis.thermal)
	if err != nil {
		return nil, err
	}

	return thermal, nil
}

// Power gets the power information for the chassis
func (chassis *Chassis) Power() (*Power, error) {
	if chassis.power == "" {
		return nil, nil
	}

	power, err := GetPower(chassis.Client, chassis.power)
	if err != nil {
		return nil, err
	}

	return power, nil
}
