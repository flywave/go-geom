package geom

import (
	"encoding/json"
)

type FeatureCollection struct {
	Type        string                 `json:"type"`
	BoundingBox BoundingBox            `json:"bbox,omitempty"`
	Features    []*Feature             `json:"features"`
	CRS         map[string]interface{} `json:"crs,omitempty"`
}

func NewFeatureCollection() *FeatureCollection {
	return &FeatureCollection{
		Type:     "FeatureCollection",
		Features: make([]*Feature, 0),
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

	return json.Marshal(fcol)
}

func UnmarshalFeatureCollection(data []byte) (*FeatureCollection, error) {
	fc := &FeatureCollection{}
	err := json.Unmarshal(data, fc)
	if err != nil {
		return nil, err
	}

	return fc, nil
}
