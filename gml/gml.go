package gml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html"
	"regexp"
	"strconv"

	"github.com/flywave/go-geom"
)

const (
	NAMESPACE = "http://www.opengis.net/gml/3.2"
	CRS84     = "urn:ogc:def:crs:OGC::CRS84"
)

type Position []float64

func (p Position) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	buf := bytes.Buffer{}
	for i := 0; i < len(p)-1; i++ {
		buf.WriteString(fmt.Sprintf("%f ", p[i]))
	}
	if len(p) > 0 {
		buf.WriteString(fmt.Sprintf("%f", p[len(p)-1]))
	}
	return e.EncodeElement(buf.String(), start)
}

func (p *Position) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var position Position
	for {
		token, err := d.Token()
		if err != nil {
			return err
		}
		switch el := token.(type) {
		case xml.CharData:
			position = getPositionFromString(string([]byte(el)))
		case xml.EndElement:
			if el == start.End() {
				*p = position
				return nil
			}
		}
	}
}

func getPositionFromString(position string) []float64 {
	regex := regexp.MustCompile(` `)
	result := regex.Split(position, -1)
	var ps []float64

	for _, fs := range result {
		if fs == "" {
			break
		}
		f, err := strconv.ParseFloat(fs, 64)
		if err != nil {
			return nil
		}
		ps = append(ps, f)
	}
	return ps
}

type Point struct {
	XMLName     xml.Name  `xml:"gml:Point"`
	ID          string    `xml:"gml:id,attr,omitempty"`
	GML         string    `xml:"xmlns:gml,attr,omitempty"`
	SrsName     string    `xml:"srsName,attr,omitempty"`
	Pos         *Position `xml:"gml:pos,omitempty"`
	Coordinates *Position `xml:"gml:coordinates,omitempty"`
}

func (g Point) Marshal() (string, error) {
	si, err := xml.MarshalIndent(g, "", "")

	if err != nil {
		return "", err
	}

	xmls := unescapeXML(string(si))

	return xmls, nil
}

type LineString struct {
	XMLName     xml.Name   `xml:"gml:LineString"`
	ID          string     `xml:"gml:id,attr,omitempty"`
	GML         string     `xml:"xmlns:gml,attr,omitempty"`
	SrsName     string     `xml:"srsName,attr,omitempty"`
	Pos         []Position `xml:"gml:pos,omitempty"`
	Coordinates *Position  `xml:"gml:coordinates,omitempty"`
	PosList     *Position  `xml:"gml:posList,omitempty"`
}

func (g LineString) Marshal() (string, error) {
	si, err := xml.MarshalIndent(g, "", "")

	if err != nil {
		return "", err
	}

	xmls := unescapeXML(string(si))

	return xmls, nil
}

type MultiPoint struct {
	XMLName xml.Name       `xml:"gml:MultiPoint"`
	ID      string         `xml:"gml:id,attr,omitempty"`
	GML     string         `xml:"xmlns:gml,attr,omitempty"`
	SrsName string         `xml:"srsName,attr,omitempty"`
	Members []PointMembers `xml:"gml:pointMembers,omitempty"`
}

func (g MultiPoint) Marshal() (string, error) {
	si, err := xml.MarshalIndent(g, "", "")

	if err != nil {
		return "", err
	}

	xmls := unescapeXML(string(si))

	return xmls, nil
}

type PointMembers struct {
	Points []Point
}

type LinearRing struct {
	Pos         []Position `xml:"gml:pos,omitempty"`
	Coordinates *Position  `xml:"gml:coordinates,omitempty"`
	PosList     *Position  `xml:"gml:posList,omitempty"`
}

type MultiCurve struct {
	XMLName xml.Name      `xml:"gml:MultiCurve"`
	ID      string        `xml:"gml:id,attr,omitempty"`
	GML     string        `xml:"xmlns:gml,attr,omitempty"`
	SrsName string        `xml:"srsName,attr,omitempty"`
	Members []CurveMember `xml:"gml:curveMember,omitempty"`
}

func (g MultiCurve) Marshal() (string, error) {
	si, err := xml.MarshalIndent(g, "", "")

	if err != nil {
		return "", err
	}

	xmls := unescapeXML(string(si))

	return xmls, nil
}

type CurveMember struct {
	Lines []LineString
}

type Polygon struct {
	XMLName  xml.Name  `xml:"gml:Polygon"`
	ID       string    `xml:"gml:id,attr,omitempty"`
	GML      string    `xml:"xmlns:gml,attr,omitempty"`
	SrsName  string    `xml:"srsName,attr,omitempty"`
	Exterior *Exterior `xml:"gml:exterior,omitempty"`
	Interior *Interior `xml:"gml:interior,omitempty"`
}

func (g Polygon) Marshal() (string, error) {
	si, err := xml.MarshalIndent(g, "", "")

	if err != nil {
		return "", err
	}

	xmls := unescapeXML(string(si))

	return xmls, nil
}

type Exterior struct {
	Line LinearRing
}

type Interior struct {
	Lines []LinearRing
}

type MultiSurface struct {
	XMLName xml.Name       `xml:"gml:MultiSurface"`
	ID      string         `xml:"gml:id,attr,omitempty"`
	GML     string         `xml:"xmlns:gml,attr,omitempty"`
	SrsName string         `xml:"srsName,attr,omitempty"`
	Members SurfaceMembers `xml:"gml:surfaceMembers,attr,omitempty"`
}

func (g MultiSurface) Marshal() (string, error) {
	si, err := xml.MarshalIndent(g, "", "")

	if err != nil {
		return "", err
	}

	xmls := unescapeXML(string(si))

	return xmls, nil
}

type SurfaceMembers struct {
	Polygons []Polygon
}

type MultiGeometry struct {
	XMLName xml.Name          `xml:"gml:MultiGeometry"`
	ID      string            `xml:"gml:id,attr,omitempty"`
	GML     string            `xml:"xmlns:gml,attr,omitempty"`
	SrsName string            `xml:"srsName,attr,omitempty"`
	Members []GeometryMembers `xml:"gml:geometryMembers,omitempty"`
}

type GeometryMembers []interface{}

func unescapeXML(XML string) string {
	return html.UnescapeString(XML)
}

func (g GeometryMembers) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	buf := bytes.Buffer{}
	for i := range g {
		switch t := g[i].(type) {
		case *Point, Point:
			si, _ := xml.Marshal(t)
			buf.Write(si)
		case *LineString, LineString:
			si, _ := xml.Marshal(t)
			buf.Write(si)
		case *MultiPoint, MultiPoint:
			si, _ := xml.Marshal(t)
			buf.Write(si)
		case *Polygon, Polygon:
			si, _ := xml.Marshal(t)
			buf.Write(si)
		case *MultiSurface, MultiSurface:
			si, _ := xml.Marshal(t)
			buf.Write(si)
		}
	}
	return e.EncodeElement(buf.String(), start)
}

func (g GeometryMembers) Marshal() (string, error) {
	si, err := xml.MarshalIndent(g, "", "")

	if err != nil {
		return "", err
	}

	xmls := unescapeXML(string(si))

	return xmls, nil
}

func Encode(g geom.Geometry) (interface{}, error) {
	switch g := g.(type) {
	case geom.Point3:
		return EncodePoint3(g)
	case geom.Point:
		return EncodePoint(g)
	case geom.LineString3:
		return EncodeLineString3(g)
	case geom.LineString:
		return EncodeLineString(g)
	case geom.MultiLine3:
		return EncodeMultiLineString3(g)
	case geom.MultiLine:
		return EncodeMultiLineString(g)
	case geom.MultiPoint3:
		return EncodeMultiPoint3(g)
	case geom.MultiPoint:
		return EncodeMultiPoint(g)
	case geom.MultiPolygon3:
		return EncodeMultiPolygon3(g)
	case geom.MultiPolygon:
		return EncodeMultiPolygon(g)
	case geom.Polygon3:
		return EncodePolygon3(g)
	case geom.Polygon:
		return EncodePolygon(g)
	case geom.Collection:
		return EncodeGeometryCollection(g)
	default:
		return nil, fmt.Errorf("unsupport geom")
	}
}

func FlatCoords(pts [][]float64, dim int) Position {
	if dim == 2 {
		ret := make(Position, len(pts)*2)
		for i := range pts {
			ret[i*2] = pts[i][0]
			ret[i*2+1] = pts[i][1]
		}
		return ret
	} else if dim == 3 {
		ret := make(Position, len(pts)*2)
		for i := range pts {
			ret[i*3] = pts[i][0]
			ret[i*3+1] = pts[i][1]
		}
		return ret
	}
	return nil
}

func EncodeLineString(ls geom.LineString) (*LineString, error) {
	flatCoords := FlatCoords(ls.Data(), 2)
	pp := &LineString{Coordinates: &flatCoords}
	return pp, nil
}

func EncodeLineString3(ls geom.LineString3) (*LineString, error) {
	flatCoords := FlatCoords(ls.Data(), 3)
	pp := &LineString{Coordinates: &flatCoords}
	return pp, nil
}

func EncodeMultiLineString(mls geom.MultiLine) (*MultiCurve, error) {
	pp := &MultiCurve{Members: make([]CurveMember, len(mls.Lines()))}
	for i, ls := range mls.Lines() {
		flatCoords := FlatCoords(ls.Data(), 2)
		pp.Members[i].Lines = []LineString{{Coordinates: &flatCoords}}
	}
	return pp, nil
}

func EncodeMultiLineString3(mls geom.MultiLine3) (*MultiCurve, error) {
	pp := &MultiCurve{Members: make([]CurveMember, len(mls.Lines()))}
	for i, ls := range mls.Lines() {
		flatCoords := FlatCoords(ls.Data(), 3)
		pp.Members[i].Lines = []LineString{{Coordinates: &flatCoords}}
	}
	return pp, nil
}

func EncodeMultiPoint(mp geom.MultiPoint) (*MultiPoint, error) {
	pp := &MultiPoint{Members: make([]PointMembers, len(mp.Points()))}
	for i, ps := range mp.Points() {
		pos := Position(ps.Data())
		pp.Members[i].Points = []Point{{Coordinates: &pos}}
	}
	return pp, nil
}

func EncodeMultiPoint3(mp geom.MultiPoint3) (*MultiPoint, error) {
	pp := &MultiPoint{Members: make([]PointMembers, len(mp.Points()))}
	for i, ps := range mp.Points() {
		pos := Position(ps.Data()[:2])
		pp.Members[i].Points = []Point{{Coordinates: &pos}}
	}
	return pp, nil
}

func EncodeMultiPolygon(mp geom.MultiPolygon) (*MultiSurface, error) {
	pp := &MultiSurface{Members: SurfaceMembers{Polygons: make([]Polygon, len(mp.Polygons()))}}
	for i, p := range mp.Polygons() {
		ppp := &Polygon{Exterior: &Exterior{}}

		flatCoords := FlatCoords(p.Sublines()[0].Data(), 2)
		ppp.Exterior.Line = LinearRing{Coordinates: &flatCoords}

		if len(p.Sublines()) > 1 {
			ppp.Interior = &Interior{Lines: make([]LinearRing, len(p.Sublines())-1)}
			for i, ls := range p.Sublines()[1:] {
				flatCoords := FlatCoords(ls.Data(), 2)
				ppp.Interior.Lines[i] = LinearRing{Coordinates: &flatCoords}
			}
		}
		pp.Members.Polygons[i] = *ppp
	}
	return pp, nil
}

func EncodeMultiPolygon3(mp geom.MultiPolygon3) (*MultiSurface, error) {
	pp := &MultiSurface{Members: SurfaceMembers{Polygons: make([]Polygon, len(mp.Polygons()))}}
	for i, p := range mp.Polygons() {
		ppp := &Polygon{Exterior: &Exterior{}}

		flatCoords := FlatCoords(p.Sublines()[0].Data(), 3)
		ppp.Exterior.Line = LinearRing{Coordinates: &flatCoords}

		if len(p.Sublines()) > 1 {
			ppp.Interior = &Interior{Lines: make([]LinearRing, len(p.Sublines())-1)}
			for i, ls := range p.Sublines()[1:] {
				flatCoords := FlatCoords(ls.Data(), 3)
				ppp.Interior.Lines[i] = LinearRing{Coordinates: &flatCoords}
			}
		}
		pp.Members.Polygons[i] = *ppp
	}
	return pp, nil
}

func EncodePoint(p geom.Point) (*Point, error) {
	pos := Position(p.Data())
	pp := &Point{Coordinates: &pos}
	return pp, nil
}

func EncodePoint3(p geom.Point3) (*Point, error) {
	pos := Position(p.Data()[:2])
	pp := &Point{Coordinates: &pos}
	return pp, nil
}

func EncodePolygon(p geom.Polygon) (*Polygon, error) {
	pp := &Polygon{Exterior: &Exterior{}}

	flatCoords := FlatCoords(p.Sublines()[0].Data(), 2)
	pp.Exterior.Line = LinearRing{Coordinates: &flatCoords}

	if len(p.Sublines()) > 1 {
		pp.Interior = &Interior{Lines: make([]LinearRing, len(p.Sublines())-1)}
		for i, ls := range p.Sublines()[1:] {
			flatCoords := FlatCoords(ls.Data(), 2)
			pp.Interior.Lines[i] = LinearRing{Coordinates: &flatCoords}
		}
	}

	return pp, nil
}

func EncodePolygon3(p geom.Polygon3) (*Polygon, error) {
	pp := &Polygon{Exterior: &Exterior{}}

	flatCoords := FlatCoords(p.Sublines()[0].Data(), 3)
	pp.Exterior.Line = LinearRing{Coordinates: &flatCoords}

	if len(p.Sublines()) > 1 {
		pp.Interior = &Interior{Lines: make([]LinearRing, len(p.Sublines())-1)}
		for i, ls := range p.Sublines()[1:] {
			flatCoords := FlatCoords(ls.Data(), 3)
			pp.Interior.Lines[i] = LinearRing{Coordinates: &flatCoords}
		}
	}

	return pp, nil
}

func EncodeGeometryCollection(geoms geom.Collection) (*MultiGeometry, error) {
	pp := &MultiGeometry{Members: []GeometryMembers{make(GeometryMembers, len(geoms))}}
	for i, g := range geoms {
		var err error
		pp.Members[0][i], err = Encode(g)
		if err != nil {
			return nil, err
		}
	}
	return pp, nil
}
