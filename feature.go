package geom

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type Feature struct {
	ID           interface{}            `json:"id,omitempty"`
	Type         string                 `json:"type"`
	BoundingBox  BoundingBox            `json:"bbox,omitempty"`
	Geometry     Geometry               `json:"_"`
	Properties   map[string]interface{} `json:"properties"`
	CRS          map[string]interface{} `json:"crs,omitempty"`
	GeometryData GeometryData           `json:"geometry"`
}

func NewFeature(geometry Geometry) *Feature {
	return &Feature{
		Type:        "Feature",
		Geometry:    geometry,
		BoundingBox: BoundingBoxFromGeometry(geometry),
		Properties:  make(map[string]interface{}),
	}
}

func NewFeatureFromGeometryData(geometry *GeometryData) *Feature {
	return &Feature{
		Type:         "Feature",
		GeometryData: *geometry,
		BoundingBox:  BoundingBoxFromGeometryData(geometry),
		Properties:   make(map[string]interface{}),
	}
}

func NewPointFeature(coordinate []float64) *Feature {
	return NewFeatureFromGeometryData(NewPointGeometryData(coordinate))
}

func NewMultiPointFeature(coordinates ...[]float64) *Feature {
	return NewFeatureFromGeometryData(NewMultiPointGeometryData(coordinates...))
}

func NewLineStringFeature(coordinates [][]float64) *Feature {
	return NewFeatureFromGeometryData(NewLineStringGeometryData(coordinates))
}

func NewMultiLineStringFeature(lines ...[][]float64) *Feature {
	return NewFeatureFromGeometryData(NewMultiLineStringGeometryData(lines...))
}

func NewPolygonFeature(polygon [][][]float64) *Feature {
	return NewFeatureFromGeometryData(NewPolygonGeometryData(polygon))
}

func NewMultiPolygonFeature(polygons ...[][][]float64) *Feature {
	return NewFeatureFromGeometryData(NewMultiPolygonGeometryData(polygons...))
}

func NewCollectionFeature(geometries ...*GeometryData) *Feature {
	return NewFeatureFromGeometryData(NewCollectionGeometryData(geometries...))
}

func (f Feature) MarshalJSON() ([]byte, error) {
	type feature Feature

	var data GeometryData
	if f.Geometry != nil {
		data = *NewGeometryData(f.Geometry)
	} else {
		data = f.GeometryData
	}

	fea := &feature{
		ID:           f.ID,
		Type:         "Feature",
		GeometryData: data,
	}

	if f.BoundingBox != nil && len(f.BoundingBox) != 0 {
		fea.BoundingBox = f.BoundingBox
	}
	if f.Properties != nil && len(f.Properties) != 0 {
		fea.Properties = f.Properties
	}
	if f.CRS != nil && len(f.CRS) != 0 {
		fea.CRS = f.CRS
	}

	return json.Marshal(fea)
}

func UnmarshalFeature(data []byte) (*Feature, error) {
	f := &Feature{}
	err := json.Unmarshal(data, f)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (f *Feature) SetProperty(key string, value interface{}) {
	if f.Properties == nil {
		f.Properties = make(map[string]interface{})
	}
	f.Properties[key] = value
}

func (f *Feature) PropertyBool(key string) (bool, error) {
	if b, ok := (f.Properties[key]).(bool); ok {
		return b, nil
	}
	return false, fmt.Errorf("type assertion of `%s` to bool failed", key)
}

func (f *Feature) PropertyInt(key string) (int, error) {
	if i, ok := (f.Properties[key]).(int); ok {
		return i, nil
	}

	if i, ok := (f.Properties[key]).(float64); ok {
		return int(i), nil
	}

	return 0, fmt.Errorf("type assertion of `%s` to int failed", key)
}

func (f *Feature) PropertyFloat64(key string) (float64, error) {
	if i, ok := (f.Properties[key]).(float64); ok {
		return i, nil
	}
	return 0, fmt.Errorf("type assertion of `%s` to float64 failed", key)
}

func (f *Feature) PropertyString(key string) (string, error) {
	if s, ok := (f.Properties[key]).(string); ok {
		return s, nil
	}
	return "", fmt.Errorf("type assertion of `%s` to string failed", key)
}

func (f *Feature) PropertyMustBool(key string, def ...bool) bool {
	var defaul bool

	b, err := f.PropertyBool(key)
	if err == nil {
		return b
	}

	if len(def) > 0 {
		defaul = def[0]
	}

	return defaul
}

func (f *Feature) PropertyMustInt(key string, def ...int) int {
	var defaul int

	b, err := f.PropertyInt(key)
	if err == nil {
		return b
	}

	if len(def) > 0 {
		defaul = def[0]
	}

	return defaul
}

func (f *Feature) PropertyMustFloat64(key string, def ...float64) float64 {
	var defaul float64

	b, err := f.PropertyFloat64(key)
	if err == nil {
		return b
	}

	if len(def) > 0 {
		defaul = def[0]
	}

	return defaul
}

func (f *Feature) PropertyMustString(key string, def ...string) string {
	var defaul string

	b, err := f.PropertyString(key)
	if err == nil {
		return b
	}

	if len(def) > 0 {
		defaul = def[0]
	}

	return defaul
}

func decodeBoundingBox(bb interface{}) ([]float64, error) {
	if bb == nil {
		return nil, nil
	}

	switch f := bb.(type) {
	case []float64:
		return f, nil
	case []interface{}:
		bb := make([]float64, 0, 4)
		for _, v := range f {
			switch c := v.(type) {
			case float64:
				bb = append(bb, c)
			default:
				return nil, fmt.Errorf("bounding box coordinate not usable, got %T", v)
			}

		}
		return bb, nil
	default:
		return nil, fmt.Errorf("bounding box property not usable, got %T", bb)
	}
}

// BoundingBox implementation as per https://tools.ietf.org/html/rfc7946
// BoundingBox syntax: "bbox": [west, south, east, north]
// BoundingBox defaults "bbox": [-180.0, -90.0, 180.0, 90.0]
func BoundingBoxFromPoints(pts [][]float64) BoundingBox {
	west, south, east, north := 180.0, 90.0, -180.0, -90.0

	for _, pt := range pts {
		if pt == nil {
			continue
		}
		x, y := pt[0], pt[1]
		if x < west {
			west = x
		} else if x > east {
			east = x
		}

		if y < south {
			south = y
		} else if y > north {
			north = y
		}
	}
	return BoundingBox{west, south, east, north}
}

func BoundingBoxsFromTwoBBox(bb1 BoundingBox, bb2 BoundingBox) BoundingBox {
	west, south, east, north := 180.0, 90.0, -180.0, -90.0

	west1, south1, east1, north1 := bb1[0], bb1[1], bb1[2], bb1[3]
	west2, south2, east2, north2 := bb2[0], bb2[1], bb2[2], bb2[3]

	if west1 < west2 {
		west = west1
	} else {
		west = west2
	}

	if south1 < south2 {
		south = south1
	} else {
		south = south2
	}

	if east1 > east2 {
		east = east1
	} else {
		east = east2
	}

	if north1 > north2 {
		north = north1
	} else {
		north = north2
	}

	return BoundingBox{west, south, east, north}
}

func ExpandBoundingBoxs(bboxs []BoundingBox) BoundingBox {
	bbox := bboxs[0]
	for _, temp_bbox := range bboxs[1:] {
		bbox = BoundingBoxsFromTwoBBox(bbox, temp_bbox)
	}
	return bbox
}

func BoundingBoxFromPointGeometry(pt []float64) BoundingBox {
	return BoundingBox{pt[0], pt[1], pt[0], pt[1]}
}

func BoundingBoxFromMultiPointGeometry(pts [][]float64) BoundingBox {
	return BoundingBoxFromPoints(pts)
}

func BoundingBoxFromLineStringGeometry(line [][]float64) BoundingBox {
	return BoundingBoxFromPoints(line)
}

func BoundingBoxFromMultiLineStringGeometry(multiline [][][]float64) BoundingBox {
	bboxs := []BoundingBox{}
	for _, line := range multiline {
		bboxs = append(bboxs, BoundingBoxFromPoints(line))
	}
	return ExpandBoundingBoxs(bboxs)
}

func BoundingBoxFromPolygonGeometry(polygon [][][]float64) BoundingBox {
	bboxs := []BoundingBox{}
	for _, cont := range polygon {
		bboxs = append(bboxs, BoundingBoxFromPoints(cont))
	}
	return ExpandBoundingBoxs(bboxs)
}

func BoundingBoxFromMultiPolygonGeometry(multipolygon [][][][]float64) BoundingBox {
	bboxs := []BoundingBox{}
	for _, polygon := range multipolygon {
		for _, cont := range polygon {
			bboxs = append(bboxs, BoundingBoxFromPoints(cont))
		}
	}
	return ExpandBoundingBoxs(bboxs)
}

func BoundingBoxFromGeometryCollection(gs []Geometry) BoundingBox {
	bboxs := []BoundingBox{}
	for _, g := range gs {
		bboxs = append(bboxs, BoundingBoxFromGeometry(g))
	}
	return ExpandBoundingBoxs(bboxs)
}

func BoundingBoxFromGeometry(g Geometry) BoundingBox {
	switch t := (g).(type) {
	case Point:
	case Point3:
		return BoundingBoxFromPointGeometry(t.Data())
	case MultiPoint:
	case MultiPoint3:
		return BoundingBoxFromMultiPointGeometry(t.Data())
	case LineString:
	case LineString3:
		return BoundingBoxFromLineStringGeometry(t.Data())
	case MultiLine:
	case MultiLine3:
		return BoundingBoxFromMultiLineStringGeometry(t.Data())
	case Polygon:
	case Polygon3:
		return BoundingBoxFromPolygonGeometry(t.Data())
	case MultiPolygon:
	case MultiPolygon3:
		return BoundingBoxFromMultiPolygonGeometry(t.Data())
	}
	return []float64{}
}

func BoundingBoxFromGeometryData(g *GeometryData) BoundingBox {
	switch g.Type {
	case "Point":
		return BoundingBoxFromPointGeometry(g.Point)
	case "MultiPoint":
		return BoundingBoxFromMultiPointGeometry(g.MultiPoint)
	case "LineString":
		return BoundingBoxFromLineStringGeometry(g.LineString)
	case "MultiLineString":
		return BoundingBoxFromMultiLineStringGeometry(g.MultiLineString)
	case "Polygon":
		return BoundingBoxFromPolygonGeometry(g.Polygon)
	case "MultiPolygon":
		return BoundingBoxFromMultiPolygonGeometry(g.MultiPolygon)
	}
	return BoundingBox{}
}

func GetKeyDifs(f1, f2 map[string]interface{}) ([]string, []string) {
	keys1 := []string{}
	k1map := map[string]string{}
	for k := range f1 {
		keys1 = append(keys1, k)
		k1map[k] = ""
	}
	keys2 := []string{}
	k2map := map[string]string{}
	for k := range f2 {
		keys2 = append(keys2, k)
		k2map[k] = ""
	}

	k1dif := []string{}
	for k := range k1map {
		_, boolval := k2map[k]
		if !boolval {
			k1dif = append(k1dif, k)
		}
	}
	k2dif := []string{}
	for k := range k2map {
		_, boolval := k1map[k]
		if !boolval {
			k2dif = append(k2dif, k)
		}
	}
	return k1dif, k2dif
}

func GetErrorsKeyDif(kd1, kd2 []string) []string {
	lines := []string{}
	for _, k := range kd1 {
		lines = append(lines, fmt.Sprintf("Feature1 Contains field %s Feature2 does not.", k))
	}
	for _, k := range kd2 {
		lines = append(lines, fmt.Sprintf("Feature2 Contains field %s Feature1 does not.", k))
	}
	return lines
}

func CheckProperties(p1, p2 map[string]interface{}) bool {
	if len(p1) != len(p2) {
		return false
	}
	for k := range p1 {
		val1, boolval1 := p1[k]
		val2, boolval2 := p2[k]
		if boolval1 && boolval2 {
			if val1 != val2 {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func IsFeatureEqual(feat1, feat2 Feature) bool {
	return IsGeometryEqual(feat1.Geometry, feat2.Geometry) && CheckProperties(feat1.Properties, feat2.Properties)
}

func ConvertFeatureID(v interface{}) (uint64, error) {
	switch aval := v.(type) {
	case float64:
		return uint64(aval), nil
	case int64:
		return uint64(aval), nil
	case uint64:
		return aval, nil
	case uint:
		return uint64(aval), nil
	case int8:
		return uint64(aval), nil
	case uint8:
		return uint64(aval), nil
	case uint16:
		return uint64(aval), nil
	case int32:
		return uint64(aval), nil
	case uint32:
		return uint64(aval), nil
	case string:
		return strconv.ParseUint(aval, 10, 64)
	default:
		return 0, errors.New("no convert feature id")
	}
}
