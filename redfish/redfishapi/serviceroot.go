//
// SPDX-License-Identifier: BSD-3-Clause
//

package redfishapi

import (
	"encoding/json"
	"fmt"

	"github.com/magicst0ne/rackserver_exporter/redfish/common"
)

// Service represents the root Redfish service. All values for resources
// described by this schema shall comply to the requirements as described in the
// Redfish specification.
type Service struct {
	common.Entity

	// ODataContext is the odata context.
	ODataContext string `json:"@odata.context"`
	// ODataID is the odata identifier.
	ODataID string `json:"@odata.id"`
	// ODataType is the odata type.
	ODataType string `json:"@odata.type"`
	// AccountService shall only contain a reference to a resource that complies
	// to the AccountService schema.
	accountService string
	// CertificateService shall be a link to the CertificateService.
	certificateService string
	// Chassis shall only contain a reference to a collection of resources that
	// comply to the Chassis schema.
	chassis string
	// CompositionService shall only contain a reference to a resource that
	// complies to the CompositionService schema.
	compositionService string
	// Description provides a description of this resource.
	Description string
	// EventService shall only contain a reference to a resource that complies
	// to the EventService schema.
	eventService string
	// Fabrics shall contain references to all Fabric instances.
	fabrics string
	// JobService shall only contain a reference to a resource that conforms to
	// the JobService schema.
	jobService string
	// JsonSchemas shall only contain a reference to a collection of resources
	// that comply to the SchemaFile schema where the files are Json-Schema
	// files.
	jsonSchemas string
	// Managers shall only contain a reference to a collection of resources that
	// comply to the Managers schema.
	managers string
	// Oem contains all the vendor specific actions. It is vendor responsibility to parse
	// this field accordingly
	Oem json.RawMessage
	// Product shall include the name of the product represented by this Redfish
	// service.
	Product string
	// RedfishVersion shall represent the version of the Redfish service. The
	// format of this string shall be of the format
	// majorversion.minorversion.errata in compliance with Protocol Version
	// section of the Redfish specification.
	RedfishVersion string
	// Registries shall contain a reference to Message Registry.
	registries string
	// ResourceBlocks shall contain references to all Resource Block instances.
	resourceBlocks string
	// SessionService shall only contain a reference to a resource that complies
	// to the SessionService schema.
	sessionService string
	// StorageServices shall contain references to all StorageService instances.
	storageServices string
	// StorageSystems shall contain computer systems that act as storage
	// servers. The HostingRoles attribute of each such computer system shall
	// have an entry for StorageServer.
	storageSystems string
	// Systems shall only contain a reference to a collection of resources that
	// comply to the Systems schema.
	systems string
	// Tasks shall only contain a reference to a resource that complies to the
	// TaskService schema.
	tasks string
	// TelemetryService shall be a link to the TelemetryService.
	telemetryService string
	// UUID shall be an exact match of the UUID value returned in a 200OK from
	// an SSDP M-SEARCH request during discovery. RFC4122 describes methods that
	// can be used to create a UUID value. The value should be considered to be
	// opaque. Client software should only treat the overall value as a
	// universally unique identifier and should not interpret any sub-fields
	// within the UUID.
	UUID string
	// UpdateService shall only contain a reference to a resource that complies
	// to the UpdateService schema.
	updateService string
	// Vendor shall include the name of the manufacturer or vendor represented
	// by this Redfish service. If this property is supported, the vendor name
	// shall not be included in the value of the Product property.
	Vendor string
	// Sessions shall contain the link to a collection of Sessions.
	sessions string
}

// UnmarshalJSON unmarshals a Service object from the raw JSON.
func (serviceroot *Service) UnmarshalJSON(b []byte) error {
	type temp Service
	var t struct {
		temp
		CertificateService common.Link
		Chassis            common.Link
		Managers           common.Link
		Tasks              common.Link
		StorageServices    common.Link
		StorageSystems     common.Link
		AccountService     common.Link
		EventService       common.Link
		Registries         common.Link
		Systems            common.Link
		CompositionService common.Link
		Fabrics            common.Link
		JobService         common.Link
		JSONSchemas        common.Link `json:"JsonSchemas"`
		ResourceBlocks     common.Link
		SessionService     common.Link
		TelemetryService   common.Link
		UpdateService      common.Link
		Links              struct {
			Sessions common.Link
		}
		Oem json.RawMessage // OEM message will be stored here
	}

	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}

	// Extract the links to other entities for later
	*serviceroot = Service(t.temp)
	serviceroot.certificateService = string(t.CertificateService)
	serviceroot.chassis = string(t.Chassis)
	serviceroot.managers = string(t.Managers)
	serviceroot.tasks = string(t.Tasks)
	serviceroot.sessions = string(t.Links.Sessions)
	serviceroot.storageServices = string(t.StorageServices)
	serviceroot.storageSystems = string(t.StorageSystems)
	serviceroot.accountService = string(t.AccountService)
	serviceroot.eventService = string(t.EventService)
	serviceroot.registries = string(t.Registries)
	serviceroot.systems = string(t.Systems)
	serviceroot.compositionService = string(t.CompositionService)
	serviceroot.fabrics = string(t.Fabrics)
	serviceroot.jobService = string(t.JobService)
	serviceroot.jsonSchemas = string(t.JSONSchemas)
	serviceroot.resourceBlocks = string(t.ResourceBlocks)
	serviceroot.sessionService = string(t.SessionService)
	serviceroot.telemetryService = string(t.TelemetryService)
	serviceroot.updateService = string(t.UpdateService)
	serviceroot.Oem = t.Oem

        if serviceroot.sessions=="" {
            serviceroot.sessions="/redfish/v1/SessionService/Sessions"
        }

	return nil
}

// ServiceRoot will get a Service instance from the service.
func ServiceRoot(c common.Client) (*Service, error) {
	resp, err := c.Get(common.DefaultServiceRoot)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var serviceroot Service
	err = json.NewDecoder(resp.Body).Decode(&serviceroot)
	if err != nil {
		return nil, err
	}

	serviceroot.SetClient(c)
	return &serviceroot, nil
}

// CreateSession creates a new session and returns the token and id
func (serviceroot *Service) CreateSession(username, password string) (*common.AuthToken, error) {
	return common.CreateSession(serviceroot.Client, serviceroot.sessions, username, password)
}

// Sessions gets the system's active sessions
func (serviceroot *Service) Sessions() ([]*common.Session, error) {
	return common.ListReferencedSessions(serviceroot.Client, serviceroot.sessions)
}

// DeleteSession logout the specified session
func (serviceroot *Service) DeleteSession(url string) error {
	return common.DeleteSession(serviceroot.Client, url)
}

// Systems get the system instances from the service
func (serviceroot *Service) Systems() ([]*ComputerSystem, error) {
        return ListReferencedComputerSystems(serviceroot.Client, serviceroot.systems)
}

func (serviceroot *Service) Chassis() ([]*Chassis, error) {
        return ListReferencedChassis(serviceroot.Client, serviceroot.chassis)
}

func DumpObj(obj interface{}) {

    empJSON, err := json.MarshalIndent(obj, "", "  ")
    if err != nil {
    	//czw debug
        fmt.Println(err.Error())
    }
    fmt.Printf("MarshalIndent funnction output\n %s\n", string(empJSON))
}
