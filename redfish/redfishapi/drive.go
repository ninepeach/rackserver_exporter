package redfishapi

import (
	"encoding/json"

	"github.com/magicst0ne/rackserver_exporter/redfish/common"
)

// Drive is used to represent a disk drive or other physical storage
// medium for a Redfish implementation.
type Drive struct {
	common.Entity

	// ODataContext is the odata context.
	ODataContext string `json:"@odata.context"`
	// ODataType is the odata type.
	ODataType string `json:"@odata.type"`

	Id string
	Name string
	Model string
	Description string
	SerialNumber string
	Location string
	InterfaceType string
	CapacityGB int
	Status common.Status
}

// UnmarshalJSON unmarshals a Drive object from the raw JSON.
func (drive *Drive) UnmarshalJSON(b []byte) error {
	type temp Drive
	var t struct {
		temp
	}

	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}

	// Extract the links to other entities for later
	*drive = Drive(t.temp)

	return nil
}

// GetDrive will get a Drive instance from the service.
func GetDrive(c common.Client, uri string) (*Drive, error) {
	resp, err := c.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var drive Drive
	err = json.NewDecoder(resp.Body).Decode(&drive)
	if err != nil {
		return nil, err
	}

	drive.SetClient(c)
	return &drive, nil
}

// ListReferencedDrives gets the collection of Drives from a provided reference.
func ListReferencedDrives(c common.Client, link string) ([]*Drive, error) { //nolint:dupl
	var result []*Drive
	if link == "" {
		return result, nil
	}

	links, err := common.GetCollection(c, link)
	if err != nil {
		return result, err
	}

	collectionError := common.NewCollectionError()
	for _, driveLink := range links.ItemLinks {
		drive, err := GetDrive(c, driveLink)
		if err != nil {
			collectionError.Failures[driveLink] = err
		} else {
			result = append(result, drive)
		}
	}

	if collectionError.Empty() {
		return result, nil
	}

	return result, collectionError
}