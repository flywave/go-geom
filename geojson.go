package geom

import (
	"encoding/json"
	"errors"
	"fmt"
)

type GeometryType string

const (
	GeometryPoint           GeometryType = "Point"
	GeometryMultiPoint      GeometryType = "MultiPoint"
	GeometryLineString      GeometryType = "LineString"
	GeometryMultiLineString GeometryType = "MultiLineString"
	GeometryPolygon         GeometryType = "Polygon"
	GeometryMultiPolygon    GeometryType = "MultiPolygon"
	GeometryCollection      GeometryType = "GeometryCollection"
)

type GeometryData struct {
	Type            GeometryType `json:"type"`
	BoundingBox     BoundingBox  `json:"bbox,omitempty"`
	Point           []float64
	MultiPoint      [][]float64
	LineString      [][]float64
	MultiLineString [][][]float64
	Polygon         [][][]float64
	MultiPolygon    [][][][]float64
	Geometries      []*GeometryData
	EPSG            int `json:"epsg,omitempty"`
}

func NewGeometryData(geometry Geometry) *GeometryData {
	switch geo := geometry.(type) {
	default:
		return nil
	case Point:
	case Point3:
		var ret GeometryData
		ret.Type = GeometryPoint
		ret.Point = geo.Data()
		return &ret
	case MultiPoint:
	case MultiPoint3:
		var ret GeometryData
		ret.Type = GeometryMultiPoint
		ret.MultiPoint = geo.Data()
		return &ret
	case LineString:
	case LineString3:
		var ret GeometryData
		ret.Type = GeometryLineString
		ret.LineString = geo.Data()
		return &ret
	case MultiLine:
	case MultiLine3:
		var ret GeometryData
		ret.Type = GeometryMultiLineString
		ret.MultiLineString = geo.Data()
		return &ret
	case Polygon:
	case Polygon3:
		var ret GeometryData
		ret.Type = GeometryPolygon
		ret.Polygon = geo.Data()
		return &ret
	case MultiPolygon:
	case MultiPolygon3:
		var ret GeometryData
		ret.Type = GeometryMultiPolygon
		ret.MultiPolygon = geo.Data()
		return &ret
	case Collection:
		var ret GeometryData
		ret.Type = GeometryCollection
		ret.Geometries = make([]*GeometryData, len(geo.Geometries()))
		for i := range geo.Geometries() {
			ret.Geometries[i] = NewGeometryData(geo.Geometries()[i])
		}
		return &ret
	}
	return nil
}

func NewPointGeometryData(coordinate []float64) *GeometryData {
	return &GeometryData{
		Type:  GeometryPoint,
		Point: coordinate,
	}
}

func NewMultiPointGeometryData(coordinates ...[]float64) *GeometryData {
	return &GeometryData{
		Type:       GeometryMultiPoint,
		MultiPoint: coordinates,
	}
}

func NewLineStringGeometryData(coordinates [][]float64) *GeometryData {
	return &GeometryData{
		Type:       GeometryLineString,
		LineString: coordinates,
	}
}

func NewMultiLineStringGeometryData(lines ...[][]float64) *GeometryData {
	return &GeometryData{
		Type:            GeometryMultiLineString,
		MultiLineString: lines,
	}
}

func NewPolygonGeometryData(polygon [][][]float64) *GeometryData {
	return &GeometryData{
		Type:    GeometryPolygon,
		Polygon: polygon,
	}
}

func NewMultiPolygonGeometryData(polygons ...[][][]float64) *GeometryData {
	return &GeometryData{
		Type:         GeometryMultiPolygon,
		MultiPolygon: polygons,
	}
}

func NewCollectionGeometryData(geometries ...*GeometryData) *GeometryData {
	return &GeometryData{
		Type:       GeometryCollection,
		Geometries: geometries,
	}
}

func (g GeometryData) MarshalJSON() ([]byte, error) {
	type geometry struct {
		Type        GeometryType           `json:"type"`
		BoundingBox []float64              `json:"bbox,omitempty"`
		Coordinates interface{}            `json:"coordinates,omitempty"`
		Geometries  interface{}            `json:"geometries,omitempty"`
		CRS         map[string]interface{} `json:"crs,omitempty"`
	}

	geo := &geometry{
		Type: g.Type,
	}

	if g.BoundingBox != nil && len(g.BoundingBox) != 0 {
		geo.BoundingBox = g.BoundingBox
	}

	switch g.Type {
	case GeometryPoint:
		geo.Coordinates = g.Point
	case GeometryMultiPoint:
		geo.Coordinates = g.MultiPoint
	case GeometryLineString:
		geo.Coordinates = g.LineString
	case GeometryMultiLineString:
		geo.Coordinates = g.MultiLineString
	case GeometryPolygon:
		geo.Coordinates = g.Polygon
	case GeometryMultiPolygon:
		geo.Coordinates = g.MultiPolygon
	case GeometryCollection:
		geo.Geometries = g.Geometries
	}

	return json.Marshal(geo)
}

func UnmarshalGeometry(data []byte) (*GeometryData, error) {
	g := &GeometryData{}
	err := json.Unmarshal(data, g)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *GeometryData) UnmarshalJSON(data []byte) error {
	var object map[string]interface{}
	err := json.Unmarshal(data, &object)
	if err != nil {
		return err
	}

	return DecodeGeometry(g, object)
}

func (g *GeometryData) GetType() string {
	return string(g.Type)
}

func (g *GeometryData) Scan(value interface{}) error {
	var data []byte

	switch value.(type) {
	case string:
		data = []byte(value.(string))
	case []byte:
		data = value.([]byte)
	default:
		return errors.New("unable to parse this type into geojson")
	}

	return g.UnmarshalJSON(data)
}

func DecodeGeometry(g *GeometryData, object map[string]interface{}) error {
	t, ok := object["type"]
	if !ok {
		return errors.New("type property not defined")
	}

	if s, ok := t.(string); ok {
		g.Type = GeometryType(s)
	} else {
		return errors.New("type property not string")
	}

	if s, ok := object["epsg"]; ok {
		t, ok1 := s.(float64)
		if ok1 {
			g.EPSG = int(t)
		}
	}

	var err error
	switch g.Type {
	case GeometryPoint:
		g.Point, err = decodePosition(object["coordinates"])
	case GeometryMultiPoint:
		g.MultiPoint, err = decodePositionSet(object["coordinates"])
	case GeometryLineString:
		g.LineString, err = decodePositionSet(object["coordinates"])
	case GeometryMultiLineString:
		g.MultiLineString, err = decodePathSet(object["coordinates"])
	case GeometryPolygon:
		g.Polygon, err = decodePathSet(object["coordinates"])
	case GeometryMultiPolygon:
		g.MultiPolygon, err = decodePolygonSet(object["coordinates"])
	case GeometryCollection:
		g.Geometries, err = decodeGeometries(object["geometries"])
	}

	return err
}

func decodePosition(data interface{}) ([]float64, error) {
	coords, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a valid position, got %v", data)
	}

	result := make([]float64, 0, len(coords))
	for _, coord := range coords {
		if f, ok := coord.(float64); ok {
			result = append(result, f)
		} else {
			return nil, fmt.Errorf("not a valid coordinate, got %v", coord)
		}
	}

	return result, nil
}

func decodePositionSet(data interface{}) ([][]float64, error) {
	points, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a valid set of positions, got %v", data)
	}

	result := make([][]float64, 0, len(points))
	for _, point := range points {
		if p, err := decodePosition(point); err == nil {
			result = append(result, p)
		} else {
			return nil, err
		}
	}

	return result, nil
}

func decodePathSet(data interface{}) ([][][]float64, error) {
	sets, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a valid path, got %v", data)
	}

	result := make([][][]float64, 0, len(sets))

	for _, set := range sets {
		if s, err := decodePositionSet(set); err == nil {
			result = append(result, s)
		} else {
			return nil, err
		}
	}

	return result, nil
}

func decodePolygonSet(data interface{}) ([][][][]float64, error) {
	polygons, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a valid polygon, got %v", data)
	}

	result := make([][][][]float64, 0, len(polygons))
	for _, polygon := range polygons {
		if p, err := decodePathSet(polygon); err == nil {
			result = append(result, p)
		} else {
			return nil, err
		}
	}

	return result, nil
}

func decodeGeometries(data interface{}) ([]*GeometryData, error) {
	if vs, ok := data.([]interface{}); ok {
		geometries := make([]*GeometryData, 0, len(vs))
		for _, v := range vs {
			g := &GeometryData{}

			vmap, ok := v.(map[string]interface{})
			if !ok {
				break
			}

			err := DecodeGeometry(g, vmap)
			if err != nil {
				return nil, err
			}

			geometries = append(geometries, g)
		}

		if len(geometries) == len(vs) {
			return geometries, nil
		}
	}

	return nil, fmt.Errorf("not a valid set of geometries, got %v", data)
}

func (g *GeometryData) IsEmpty() bool {
	return g.Type == ""
}

func (g *GeometryData) IsPoint() bool {
	return g.Type == GeometryPoint
}

func (g *GeometryData) IsMultiPoint() bool {
	return g.Type == GeometryMultiPoint
}

func (g *GeometryData) IsLineString() bool {
	return g.Type == GeometryLineString
}

func (g *GeometryData) IsMultiLineString() bool {
	return g.Type == GeometryMultiLineString
}

func (g *GeometryData) IsPolygon() bool {
	return g.Type == GeometryPolygon
}

func (g *GeometryData) IsMultiPolygon() bool {
	return g.Type == GeometryMultiPolygon
}

func (g *GeometryData) IsCollection() bool {
	return g.Type == GeometryCollection
}
