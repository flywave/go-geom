package gml

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"testing"
)

func TestPoint(t *testing.T) {
	p := &Point{ID: "ID", GML: "http://www.opengis.net/gml/3.2", SrsName: "EPSG:4326", Pos: &Position{2.0, 1.0}}

	si, _ := xml.MarshalIndent(p, "", "")

	xmls := string(si)

	if xmls == "" {
		t.FailNow()
	}

	p2 := &Point{ID: "ID", GML: "http://www.opengis.net/gml/3.2", SrsName: "EPSG:4326", Coordinates: &Position{2.0, 1.0}}

	si2, _ := xml.MarshalIndent(p2, "", "")

	xml2 := string(si2)

	if xml2 == "" {
		t.FailNow()
	}

}

func TestGeom(t *testing.T) {
	p0 := &Point{ID: "ID", GML: "http://www.opengis.net/gml/3.2", SrsName: "EPSG:4326", Pos: &Position{2.0, 1.0}}
	p1 := &Point{ID: "ID", GML: "http://www.opengis.net/gml/3.2", SrsName: "EPSG:4326", Pos: &Position{3.0, 1.0}}
	p2 := &Point{ID: "ID", GML: "http://www.opengis.net/gml/3.2", SrsName: "EPSG:4326", Pos: &Position{4.0, 1.0}}

	g := &MultiGeometry{ID: "ID", GML: "http://www.opengis.net/gml/3.2", SrsName: "EPSG:4326", Members: []GeometryMembers{[]interface{}{p0, p1}, []interface{}{p2}}}

	si, _ := xml.MarshalIndent(g, "", "")

	xmls := unescapeXML(string(si))

	if xmls == "" {
		t.FailNow()
	}

	ioutil.WriteFile("./test.xml", []byte(xmls), os.ModePerm)

}
