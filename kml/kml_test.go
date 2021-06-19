package kml

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/flywave/go-geom/general"
)

func TestKmlPoint(t *testing.T) {
	pt := general.NewPoint([]float64{0, 0})
	want := "<Point><coordinates>0,0</coordinates></Point>"
	sb := &strings.Builder{}
	e := xml.NewEncoder(sb)
	element, _ := Encode(pt)
	e.Encode(element)
	if want != sb.String() {
		t.FailNow()
	}
}

func TestKmlPoint3(t *testing.T) {
	pt := general.NewPoint3([]float64{0, 0, 1})
	want := "<Point><coordinates>0,0,1</coordinates></Point>"
	sb := &strings.Builder{}
	e := xml.NewEncoder(sb)
	element, _ := Encode(pt)
	e.Encode(element)
	if want != sb.String() {
		t.FailNow()
	}
}

func TestKmlMultiPoint(t *testing.T) {
	pt := general.NewMultiPoint([][]float64{{1, 2}, {3, 4}, {5, 6}})
	want := "<MultiGeometry>" +
		"<Point>" +
		"<coordinates>1,2</coordinates>" +
		"</Point>" +
		"<Point>" +
		"<coordinates>3,4</coordinates>" +
		"</Point>" +
		"<Point>" +
		"<coordinates>5,6</coordinates>" +
		"</Point>" +
		"</MultiGeometry>"
	sb := &strings.Builder{}
	e := xml.NewEncoder(sb)
	element, _ := Encode(pt)
	e.Encode(element)
	if want != sb.String() {
		t.FailNow()
	}
}

func TestKmlMultiPoint3(t *testing.T) {
	pt := general.NewMultiPoint3([][]float64{{1, 2, 1}, {3, 4, 1}, {5, 6, 1}})
	want := "<MultiGeometry>" +
		"<Point>" +
		"<coordinates>1,2,1</coordinates>" +
		"</Point>" +
		"<Point>" +
		"<coordinates>3,4,1</coordinates>" +
		"</Point>" +
		"<Point>" +
		"<coordinates>5,6,1</coordinates>" +
		"</Point>" +
		"</MultiGeometry>"
	sb := &strings.Builder{}
	e := xml.NewEncoder(sb)
	element, _ := Encode(pt)
	e.Encode(element)
	str := sb.String()
	if want != str {
		t.FailNow()
	}
}

func TestKmlLineString(t *testing.T) {
	pt := general.NewLineString([][]float64{{0, 0}, {1, 1}})
	want := "<LineString><coordinates>0,0 1,1</coordinates></LineString>"
	sb := &strings.Builder{}
	e := xml.NewEncoder(sb)
	element, _ := Encode(pt)
	e.Encode(element)
	str := sb.String()
	if want != str {
		t.FailNow()
	}
}

func TestKmlLineString3(t *testing.T) {
	pt := general.NewLineString3([][]float64{{0, 0, 1}, {1, 1, 1}})
	want := "<LineString><coordinates>0,0,1 1,1,1</coordinates></LineString>"
	sb := &strings.Builder{}
	e := xml.NewEncoder(sb)
	element, _ := Encode(pt)
	e.Encode(element)
	if want != sb.String() {
		t.FailNow()
	}
}

func TestKmlMultiLineString(t *testing.T) {
	pt := general.NewMultiLineString([][][]float64{{{1, 2}, {3, 4}, {5, 6}, {7, 8}}})
	want := "<MultiGeometry>" +
		"<LineString>" +
		"<coordinates>1,2 3,4 5,6 7,8</coordinates>" +
		"</LineString>" +
		"</MultiGeometry>"
	sb := &strings.Builder{}
	e := xml.NewEncoder(sb)
	element, _ := Encode(pt)
	e.Encode(element)
	str := sb.String()
	if want != str {
		t.FailNow()
	}
}

func TestKmlMultiLineString3(t *testing.T) {
	pt := general.NewMultiLineString3([][][]float64{{{1, 2, 1}, {3, 4, 1}, {5, 6, 1}, {7, 8, 1}}})
	want := "<MultiGeometry>" +
		"<LineString>" +
		"<coordinates>1,2,1 3,4,1 5,6,1 7,8,1</coordinates>" +
		"</LineString>" +
		"</MultiGeometry>"
	sb := &strings.Builder{}
	e := xml.NewEncoder(sb)
	element, _ := Encode(pt)
	e.Encode(element)
	str := sb.String()
	if want != str {
		t.FailNow()
	}
}

func TestKmlPolygon(t *testing.T) {
	pt := general.NewPolygon([][][]float64{{{1, 2}, {3, 4}, {5, 6}, {1, 2}}})
	want := "<Polygon>" +
		"<outerBoundaryIs>" +
		"<LinearRing>" +
		"<coordinates>1,2 3,4 5,6 1,2</coordinates>" +
		"</LinearRing>" +
		"</outerBoundaryIs>" +
		"</Polygon>"
	sb := &strings.Builder{}
	e := xml.NewEncoder(sb)
	element, _ := Encode(pt)
	e.Encode(element)
	str := sb.String()
	if want != str {
		t.FailNow()
	}
}

func TestKmlPolygon3(t *testing.T) {
	pt := general.NewPolygon3([][][]float64{{{1, 2, 1}, {3, 4, 1}, {5, 6, 1}, {1, 2, 1}}})
	want := "<Polygon>" +
		"<outerBoundaryIs>" +
		"<LinearRing>" +
		"<coordinates>1,2,1 3,4,1 5,6,1 1,2,1</coordinates>" +
		"</LinearRing>" +
		"</outerBoundaryIs>" +
		"</Polygon>"
	sb := &strings.Builder{}
	e := xml.NewEncoder(sb)
	element, _ := Encode(pt)
	e.Encode(element)
	str := sb.String()
	if want != str {
		t.FailNow()
	}
}

func TestKmlMultiPolygon(t *testing.T) {
	pt := general.NewMultiPolygon([][][][]float64{
		{
			{{1, 2}, {4, 5}, {7, 8}, {1, 2}},
			{{0.4, 0.5}, {0.7, 0.8}, {0.1, 0.2}, {0.4, 0.5}},
		},
	})
	want := "<MultiGeometry>" +
		"<Polygon>" +
		"<outerBoundaryIs>" +
		"<LinearRing>" +
		"<coordinates>1,2 4,5 7,8 1,2</coordinates>" +
		"</LinearRing>" +
		"</outerBoundaryIs>" +
		"<innerBoundaryIs>" +
		"<LinearRing>" +
		"<coordinates>0.4,0.5 0.7,0.8 0.1,0.2 0.4,0.5</coordinates>" +
		"</LinearRing>" +
		"</innerBoundaryIs>" +
		"</Polygon>" +
		"</MultiGeometry>"
	sb := &strings.Builder{}
	e := xml.NewEncoder(sb)
	element, _ := Encode(pt)
	e.Encode(element)
	str := sb.String()
	t.Log(str)
	if want != str {
		t.FailNow()
	}
}

func TestKmlMultiPolygon3(t *testing.T) {
	pt := general.NewMultiPolygon3([][][][]float64{
		{
			{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {1, 2, 3}},
			{{0.4, 0.5, 0.6}, {0.7, 0.8, 0.9}, {0.1, 0.2, 0.3}, {0.4, 0.5, 0.6}},
		},
	})
	want := "<MultiGeometry>" +
		"<Polygon>" +
		"<outerBoundaryIs>" +
		"<LinearRing>" +
		"<coordinates>1,2,3 4,5,6 7,8,9 1,2,3</coordinates>" +
		"</LinearRing>" +
		"</outerBoundaryIs>" +
		"<innerBoundaryIs>" +
		"<LinearRing>" +
		"<coordinates>0.4,0.5,0.6 0.7,0.8,0.9 0.1,0.2,0.3 0.4,0.5,0.6</coordinates>" +
		"</LinearRing>" +
		"</innerBoundaryIs>" +
		"</Polygon>" +
		"</MultiGeometry>"
	sb := &strings.Builder{}
	e := xml.NewEncoder(sb)
	element, _ := Encode(pt)
	e.Encode(element)
	str := sb.String()
	t.Log(str)
	if want != str {
		t.FailNow()
	}
}

func TestKmlGeometryCollection(t *testing.T) {
	geos := general.NewGeometryCollection(
		general.NewLineString([][]float64{
			{-122.4425587930444, 37.80666418607323},
			{-122.4428379594768, 37.80663578323093},
		}),
		general.NewLineString([][]float64{
			{-122.4425509770566, 37.80662588061205},
			{-122.4428340530617, 37.8065999493009},
		}),
	)

	want := "<MultiGeometry>" +
		"<LineString>" +
		"<coordinates>" +
		"-122.4425587930444,37.80666418607323 -122.4428379594768,37.80663578323093" +
		"</coordinates>" +
		"</LineString>" +
		"<LineString>" +
		"<coordinates>" +
		"-122.4425509770566,37.80662588061205 -122.4428340530617,37.8065999493009" +
		"</coordinates>" +
		"</LineString>" +
		"</MultiGeometry>"

	sb := &strings.Builder{}
	e := xml.NewEncoder(sb)
	element, _ := Encode(geos)
	e.Encode(element)
	str := sb.String()
	t.Log(str)
	if want != str {
		t.FailNow()
	}
}
