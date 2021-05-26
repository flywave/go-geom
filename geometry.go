package geom

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
)

type BoundingBox []float64

type Geometry interface {
	GetType() string
}

type Point interface {
	Geometry
	X() float64
	Y() float64
	Data() []float64
}

type Point3 interface {
	Point
	Z() float64
}

type MultiPoint interface {
	Geometry
	Points() []Point
	Data() [][]float64
}

type MultiPoint3 interface {
	Geometry
	Points() []Point3
	Data() [][]float64
}

type LineString interface {
	Geometry
	Subpoints() []Point
	Data() [][]float64
}

type LineString3 interface {
	Geometry
	Subpoints() []Point3
	Data() [][]float64
}

type MultiLine interface {
	Geometry
	Lines() []LineString
	Data() [][][]float64
}

type MultiLine3 interface {
	Geometry
	Lines() []LineString3
	Data() [][][]float64
}

type Polygon interface {
	Geometry
	Sublines() []LineString
	Data() [][][]float64
}

type Polygon3 interface {
	Geometry
	Sublines() []LineString3
	Data() [][][]float64
}

type MultiPolygon interface {
	Geometry
	Polygons() []Polygon
	Data() [][][][]float64
}

type MultiPolygon3 interface {
	Geometry
	Polygons() []Polygon3
	Data() [][][][]float64
}

type Collection []Geometry

func (c Collection) Geometries() []Geometry {
	return []Geometry(c)
}

func (c Collection) GetType() string {
	return string(GeometryCollection)
}

func GeometryAsString(g Geometry) string {
	switch geo := g.(type) {
	case LineString:
		rstring := "["
		for _, p := range geo.Subpoints() {
			rstring = fmt.Sprintf("%v ( %v %v )", rstring, p.X(), p.Y())
		}
		rstring += "]"
		return rstring

	default:
		return fmt.Sprintf("%v", g)
	}
}

func GeometryAsMap(g Geometry) map[string]interface{} {
	js := make(map[string]interface{})
	var vals []map[string]interface{}
	switch geo := g.(type) {
	case Point:
		js["type"] = g.GetType()
		js["value"] = []float64{geo.X(), geo.Y()}
	case Point3:
		js["type"] = g.GetType()
		js["value"] = []float64{geo.X(), geo.Y(), geo.Z()}
	case MultiPoint:
		js["type"] = g.GetType()
		for _, p := range geo.Points() {
			vals = append(vals, GeometryAsMap(p))
		}
		js["value"] = vals
	case MultiPoint3:
		js["type"] = g.GetType()
		for _, p := range geo.Points() {
			vals = append(vals, GeometryAsMap(p))
		}
		js["value"] = vals
	case LineString:
		js["type"] = g.GetType()
		var fv []float64
		for _, p := range geo.Subpoints() {
			fv = append(fv, p.X(), p.Y())
		}
		js["value"] = fv
	case LineString3:
		js["type"] = g.GetType()
		var fv []float64
		for _, p := range geo.Subpoints() {
			fv = append(fv, p.X(), p.Y())
		}
		js["value"] = fv
	case MultiLine:
		js["type"] = g.GetType()
		for _, l := range geo.Lines() {
			vals = append(vals, GeometryAsMap(l))
		}
		js["value"] = vals
	case MultiLine3:
		js["type"] = g.GetType()
		for _, l := range geo.Lines() {
			vals = append(vals, GeometryAsMap(l))
		}
		js["value"] = vals
	case Polygon:
		js["type"] = g.GetType()
		for _, l := range geo.Sublines() {
			vals = append(vals, GeometryAsMap(l))
		}
		js["value"] = vals
	case Polygon3:
		js["type"] = g.GetType()
		for _, l := range geo.Sublines() {
			vals = append(vals, GeometryAsMap(l))
		}
		js["value"] = vals
	case MultiPolygon:
		js["type"] = g.GetType()
		for _, p := range geo.Polygons() {
			vals = append(vals, GeometryAsMap(p))
		}
		js["value"] = vals
	case MultiPolygon3:
		js["type"] = g.GetType()
		for _, p := range geo.Polygons() {
			vals = append(vals, GeometryAsMap(p))
		}
		js["value"] = vals
	}
	return js
}

func GeometryAsJSON(g Geometry, w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(GeometryAsMap(g))
}

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func RoundPoint(point Point) []float64 {
	x, y := point.X(), point.Y()
	return []float64{Round(x, .5, 7), Round(y, .5, 7)}
}

func RoundPoint3(point Point3) []float64 {
	x, y, z := point.X(), point.Y(), point.Z()
	return []float64{Round(x, .5, 7), Round(y, .5, 7), Round(z, .5, 7)}
}

func IsPointEqual(p1, p2 Point) bool {
	if p1 == nil || p2 == nil {
		return p1 == p2
	}
	xdim := math.Abs(p1.X()-p2.X()) > math.Pow10(-6)
	ydim := math.Abs(p1.Y()-p2.Y()) > math.Pow10(-6)

	return !xdim && !ydim
}

func IsPoint3Equal(p1, p2 Point3) bool {
	if p1 == nil || p2 == nil {
		return p1 == p2
	}
	xdim := math.Abs(p1.X()-p2.X()) > math.Pow10(-6)
	ydim := math.Abs(p1.Y()-p2.Y()) > math.Pow10(-6)
	zdim := math.Abs(p1.Z()-p2.Z()) > math.Pow10(-6)

	return !xdim && !ydim && !zdim
}

func IsMultiPointEqual(mp1, mp2 MultiPoint) bool {
	pts1, pts2 := mp1.Points(), mp2.Points()
	if len(pts1) != len(pts2) {
		return false
	}
	for i, pt := range pts1 {
		if !IsPointEqual(pt, pts2[i]) {
			return false
		}
	}
	return true
}

func IsMultiPoint3Equal(mp1, mp2 MultiPoint3) bool {
	pts1, pts2 := mp1.Points(), mp2.Points()
	if len(pts1) != len(pts2) {
		return false
	}
	for i, pt := range pts1 {
		if !IsPointEqual(pt, pts2[i]) {
			return false
		}
	}
	return true
}

func IsLineStringEqual(l1, l2 LineString) bool {
	pts1, pts2 := l1.Subpoints(), l2.Subpoints()
	if len(pts1) != len(pts2) {
		return false
	}
	for i, pt := range pts1 {
		if !IsPointEqual(pt, pts2[i]) {
			return false
		}
	}
	return true
}

func IsLineString3Equal(l1, l2 LineString3) bool {
	pts1, pts2 := l1.Subpoints(), l2.Subpoints()
	if len(pts1) != len(pts2) {
		return false
	}
	for i, pt := range pts1 {
		if !IsPointEqual(pt, pts2[i]) {
			return false
		}
	}
	return true
}

func IsMultiLineEqual(ml1, ml2 MultiLine) bool {
	lns1, lns2 := ml1.Lines(), ml2.Lines()
	if len(lns1) != len(lns2) {
		return false
	}
	for i, ln := range lns1 {
		if !IsLineStringEqual(ln, lns2[i]) {
			return false
		}
	}
	return true
}

func IsMultiLine3Equal(ml1, ml2 MultiLine3) bool {
	lns1, lns2 := ml1.Lines(), ml2.Lines()
	if len(lns1) != len(lns2) {
		return false
	}
	for i, ln := range lns1 {
		if !IsLineString3Equal(ln, lns2[i]) {
			return false
		}
	}
	return true
}

func IsPolygonEqual(p1, p2 Polygon) bool {
	lns1, lns2 := p1.Sublines(), p2.Sublines()
	if len(lns1) != len(lns2) {
		return false
	}
	for i, ln := range lns1 {
		if !IsLineStringEqual(ln, lns2[i]) {
			return false
		}
	}
	return true
}

func IsPolygon3Equal(p1, p2 Polygon3) bool {
	lns1, lns2 := p1.Sublines(), p2.Sublines()
	if len(lns1) != len(lns2) {
		return false
	}
	for i, ln := range lns1 {
		if !IsLineString3Equal(ln, lns2[i]) {
			return false
		}
	}
	return true
}

func IsMultiPolygonEqual(mp1, mp2 MultiPolygon) bool {
	pgs1, pgs2 := mp1.Polygons(), mp2.Polygons()
	if len(pgs1) != len(pgs2) {
		return false
	}
	for i, pg := range pgs1 {
		if !IsPolygonEqual(pg, pgs2[i]) {
			return false
		}
	}
	return true
}

func IsMultiPolygon3Equal(mp1, mp2 MultiPolygon3) bool {
	pgs1, pgs2 := mp1.Polygons(), mp2.Polygons()
	if len(pgs1) != len(pgs2) {
		return false
	}
	for i, pg := range pgs1 {
		if !IsPolygon3Equal(pg, pgs2[i]) {
			return false
		}
	}
	return true
}

func IsGeometryEqual(g1, g2 Geometry) bool {
	switch geo1 := g1.(type) {
	case Point:
		geo2, ok := g2.(Point)
		if !ok {
			return false
		}
		return IsPointEqual(geo1, geo2)
	case Point3:
		geo2, ok := g2.(Point3)
		if !ok {
			return false
		}
		return IsPoint3Equal(geo1, geo2)
	case MultiPoint:
		geo2, ok := g2.(MultiPoint)
		if !ok {
			return false
		}
		return IsMultiPointEqual(geo1, geo2)
	case MultiPoint3:
		geo2, ok := g2.(MultiPoint3)
		if !ok {
			return false
		}
		return IsMultiPoint3Equal(geo1, geo2)
	case LineString:
		geo2, ok := g2.(LineString)
		if !ok {
			return false
		}
		return IsLineStringEqual(geo1, geo2)
	case LineString3:
		geo2, ok := g2.(LineString3)
		if !ok {
			return false
		}
		return IsLineString3Equal(geo1, geo2)
	case MultiLine:
		geo2, ok := g2.(MultiLine)
		if !ok {
			return false
		}
		return IsMultiLineEqual(geo1, geo2)
	case MultiLine3:
		geo2, ok := g2.(MultiLine3)
		if !ok {
			return false
		}
		return IsMultiLine3Equal(geo1, geo2)
	case Polygon:
		geo2, ok := g2.(Polygon)
		if !ok {
			return false
		}
		return IsPolygonEqual(geo1, geo2)
	case Polygon3:
		geo2, ok := g2.(Polygon3)
		if !ok {
			return false
		}
		return IsPolygon3Equal(geo1, geo2)
	case MultiPolygon:
		geo2, ok := g2.(MultiPolygon)
		if !ok {
			return false
		}
		return IsMultiPolygonEqual(geo1, geo2)
	case MultiPolygon3:
		geo2, ok := g2.(MultiPolygon3)
		if !ok {
			return false
		}
		return IsMultiPolygon3Equal(geo1, geo2)
	case Collection:
		geo2, ok := g2.(Collection)
		if !ok {
			return false
		}
		return IsCollectionEqual(geo1, geo2)
	}
	return false
}

func IsCollectionEqual(c1, c2 Collection) bool {
	geos1, geos2 := c1.Geometries(), c2.Geometries()
	if len(geos1) != len(geos2) {
		return false
	}
	for i, geo := range geos1 {
		if !IsGeometryEqual(geo, geos2[i]) {
			return false
		}
	}
	return true
}

func IsGeometryEmpty(geom Geometry) bool {
	switch t := geom.(type) {
	case Point:
	case Point3:
		return len(t.Data()) == 0
	case MultiPoint:
	case MultiPoint3:
		return len(t.Data()) == 0
	case LineString:
	case LineString3:
		return len(t.Data()) == 0
	case MultiLine:
	case MultiLine3:
		return len(t.Data()) == 0
	case Polygon:
	case Polygon3:
		return len(t.Data()) == 0
	case MultiPolygon:
	case MultiPolygon3:
		return len(t.Data()) == 0
	}
	return false
}
