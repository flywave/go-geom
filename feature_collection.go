package geom

import (
	"encoding/json"
)

type FeatureCollection struct {
	Type        string                 `json:"type"`
	BoundingBox *BoundingBox           `json:"bbox,omitempty"`
	Features    []*Feature             `json:"features"`
	CRS         map[string]interface{} `json:"crs,omitempty"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
}

func NewFeatureCollection() *FeatureCollection {
	return &FeatureCollection{
		Type:       "FeatureCollection",
		Features:   make([]*Feature, 0),
		Properties: make(map[string]interface{}),
	}
}

func (fc *FeatureCollection) AddFeature(feature *Feature) *FeatureCollection {
	fc.Features = append(fc.Features, feature)
	return fc
}

func (fc FeatureCollection) MarshalJSON() ([]byte, error) {
	type featureCollection FeatureCollection

	fcol := &featureCollection{
		Type: "FeatureCollection",
	}

	if fc.BoundingBox != nil && len(fc.BoundingBox) != 0 {
		fcol.BoundingBox = fc.BoundingBox
	}

	fcol.Features = fc.Features
	if fcol.Features == nil {
		fcol.Features = make([]*Feature, 0)
	}

	if fc.CRS != nil && len(fc.CRS) != 0 {
		fcol.CRS = fc.CRS
	}

	if fc.Properties != nil && len(fc.Properties) != 0 {
		fcol.Properties = fc.Properties
	}

	return json.Marshal(fcol)
}
