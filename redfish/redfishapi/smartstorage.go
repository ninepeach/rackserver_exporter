package redfishapi

import (
	"encoding/json"
	"fmt"

	"github.com/magicst0ne/rackserver_exporter/redfish/common"
)


// SmartStorage is used to represent a storage controller and its
// directly-attached devices.
type SmartStorage struct {
	common.Entity

	// ODataContext is the odata context.
	ODataContext string `json:"@odata.context"`
	// ODataType is the odata type.
	ODataType string `json:"@odata.type"`
	Id string
	Name string
	// Description provides a description of this resource.
	Description string
	// Devices shall contain a list of storage devices
	// associated with this resource.
	drives string
	Model string
	SerialNumber string
	Location string
	CurrentOperatingMode string
	Status common.Status
}

// UnmarshalJSON unmarshals a SmartStorage object from the raw JSON.
func (smartstorage *SmartStorage) UnmarshalJSON(b []byte) error {
	type temp SmartStorage
	var t struct {
		temp
		Links struct {
			LogicalDrives common.Link
			PhysicalDrives common.Link
		}
	}

	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}

	// Extract the links to other entities for later
	*smartstorage = SmartStorage(t.temp)
	smartstorage.drives = string(t.Links.PhysicalDrives)

	return nil
}


func (smartstorage *SmartStorage) Drives() ([]*Drive, error) {
	return ListReferencedDrives(smartstorage.Client, smartstorage.drives)
}


// GetSmartStorage will get a SmartStorage instance from the service.
func GetSmartStorage(c common.Client, uri string) (*SmartStorage, error) {
	resp, err := c.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var smartstorage SmartStorage
	err = json.NewDecoder(resp.Body).Decode(&smartstorage)
	if err != nil {
		return nil, err
	}

	smartstorage.SetClient(c)
	return &smartstorage, nil
}

// ListReferencedSmartStorages gets the collection of SmartStorage from
// a provided reference.
func ListReferencedSmartStorages(c common.Client, link string) ([]*SmartStorage, error) { //nolint:dupl
	var result []*SmartStorage
	if link == "" {
		return result, nil
	}

	link = fmt.Sprintf("%sArrayControllers/", link)

	links, err := common.GetCollection(c, link)
	if err != nil {
		return result, err
	}

	collectionError := common.NewCollectionError()
	for _, smartstorageLink := range links.ItemLinks {
		smartstorage, err := GetSmartStorage(c, smartstorageLink)
		if err != nil {
			collectionError.Failures[smartstorageLink] = err
		} else {
			result = append(result, smartstorage)
		}
	}

	if collectionError.Empty() {
		return result, nil
	}

	return result, collectionError
}
