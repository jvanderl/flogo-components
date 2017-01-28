package trigger

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// Metadata is the metadata for a Trigger
type Metadata struct {
	ID       string
	Endpoint EndpointMetadata
	Settings map[string]*data.Attribute
	Outputs  map[string]*data.Attribute
}

// EndpointMetadata is the metadata for a Trigger Endpoint
type EndpointMetadata struct {
	Settings []*data.Attribute `json:"settings"`
}

// NewMetadata creates a Metadata object from the json representation
func NewMetadata(jsonMetadata string) *Metadata {
	md := &Metadata{}
	err := json.Unmarshal([]byte(jsonMetadata), md)
	if err != nil {
		panic("Unable to parse trigger metadata: " + err.Error())
	}

	return md
}

// MarshalJSON overrides the default MarshalJSON for TaskEnv
func (md *Metadata) MarshalJSON() ([]byte, error) {

	settings := make([]*data.Attribute, 0, len(md.Settings))

	for _, value := range md.Settings {
		settings = append(settings, value)
	}

	outputs := make([]*data.Attribute, 0, len(md.Outputs))

	for _, value := range md.Outputs {
		outputs = append(outputs, value)
	}

	return json.Marshal(&struct {
		Name     string            `json:"name"`
		Endpoint EndpointMetadata  `json:"endpoint"`
		Settings []*data.Attribute `json:"settings"`
		Outputs  []*data.Attribute `json:"outputs"`
	}{
		Name:     md.ID,
		Endpoint: md.Endpoint,
		Settings: settings,
		Outputs:  outputs,
	})
}

// UnmarshalJSON overrides the default UnmarshalJSON for TaskEnv
func (md *Metadata) UnmarshalJSON(b []byte) error {

	ser := &struct {
		Name     string            `json:"name"`
		Endpoint EndpointMetadata  `json:"endpoint"`
		Settings []*data.Attribute `json:"settings"`
		Outputs  []*data.Attribute `json:"outputs"`
	}{}

	if err := json.Unmarshal(b, ser); err != nil {
		return err
	}

	md.ID = ser.Name
	md.Settings = make(map[string]*data.Attribute, len(ser.Settings))
	md.Outputs = make(map[string]*data.Attribute, len(ser.Outputs))

	for _, attr := range ser.Settings {
		md.Settings[attr.Name] = attr
	}

	for _, attr := range ser.Outputs {
		md.Outputs[attr.Name] = attr
	}

	return nil
}

// OutputsToAttrs converts the supplied output data to attributes
func (md *Metadata) OutputsToAttrs(outputData map[string]interface{}, coerce bool) ([]*data.Attribute, error) {

	attrs := make([]*data.Attribute, 0, len(md.Outputs))

	for k, a := range md.Outputs {
		v, _ := outputData[k]

		if coerce {
			var err error
			v, err = data.CoerceToValue(v, a.Type)

			if err != nil {
				return nil, err
			}
		}

		attrs = append(attrs, data.NewAttribute(a.Name, a.Type, v))
	}

	return attrs, nil
}
