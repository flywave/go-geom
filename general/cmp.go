package general

import (
	"math"
	"sort"

	"github.com/flywave/go-geom"
)

const TOLERANCE = 0.000001

var (
	NilPoint      = (geom.Point)(nil)
	NilMultiPoint = (geom.MultiPoint)(nil)
	NilLineString = (geom.LineString)(nil)
	NilMultiLine  = (geom.MultiLine)(nil)
	NilPoly       = (geom.Polygon)(nil)
	NilMultiPoly  = (geom.MultiPolygon)(nil)
	NilCollection = (geom.Collection)(nil)
)

func FloatSlice(f1, f2 []float64) bool { return Float64Slice(f1, f2, TOLERANCE) }

func Float64Slice(f1, f2 []float64, tolerance float64) bool {
	if len(f1) != len(f2) {
		return false
	}
	if len(f1) == 0 {
		return true
	}
	f1s := make([]float64, len(f1))
	f2s := make([]float64, len(f2))
	copy(f1s, f1)
	copy(f2s, f2)
	sort.Float64s(f1s)
	sort.Float64s(f2s)
	for i := range f1s {
		if !Float64(f1s[i], f2s[i], tolerance) {
			return false
		}
	}
	return true
}

func Float64(f1, f2, tolerance float64) bool {
	if math.IsInf(f1, 1) {
		return math.IsInf(f2, 1)
	}
	if math.IsInf(f2, 1) {
		return math.IsInf(f1, 1)
	}
	if math.IsInf(f1, -1) {
		return math.IsInf(f2, -1)
	}
	if math.IsInf(f2, -1) {
		return math.IsInf(f1, -1)
	}
	return math.Abs(f1-f2) < tolerance
}

func Float(f1, f2 float64) bool { return Float64(f1, f2, TOLERANCE) }

func Extented(extent1, extent2 [4]float64) bool {
	return Float(extent1[0], extent2[0]) && Float(extent1[1], extent2[1]) &&
		Float(extent1[2], extent2[2]) && Float(extent1[3], extent2[3])
}

func GeomExtent(extent1, extent2 Extenter) bool {
	return Extented(extent1.Extent(), extent2.Extent())
}

func PointLess(p1, p2 []float64) bool {
	if p1[0] != p2[0] {
		return p1[0] < p2[0]
	}
	return p1[1] < p2[1]
}

func PointEqual(p1, p2 []float64) bool {
	return Float(p1[0], p2[0]) && Float(p1[1], p2[1])
}

func GeomPointEqual(p1, p2 Point) bool {
	return Float(p1.X(), p2.X()) && Float(p1.Y(), p2.Y())
}

func MultiPointEqual(p1, p2 [][]float64) bool {
	if len(p1) != len(p2) {
		return false
	}
	cv1 := make([][]float64, len(p1))
	copy(cv1, p1)
	cv2 := make([][]float64, len(p2))
	copy(cv2, p2)
	sort.Sort(ByXY(cv1))
	sort.Sort(ByXY(cv2))
	for i := range cv1 {
		if !PointEqual(cv1[i], cv2[i]) {
			return false
		}
	}
	return true
}

func LineStringEqual(v1, v2 [][]float64) bool {
	if len(v1) != len(v2) {
		return false
	}
	cv1 := make([][]float64, len(v1))
	copy(cv1, v1)
	cv2 := make([][]float64, len(v2))
	copy(cv2, v2)
	RotateToLeftMostPoint(cv1)
	RotateToLeftMostPoint(cv2)
	for i := range cv1 {
		if !PointEqual(cv1[i], cv2[i]) {
			return false
		}
	}
	return true
}

func MultiLineEqual(ml1, ml2 [][][]float64) bool {
	if len(ml1) != len(ml2) {
		return false
	}
LOOP:
	for i := range ml1 {
		for j := range ml2 {
			if LineStringEqual(ml1[i], ml2[j]) {
				continue LOOP
			}
		}
		return false
	}
	return true
}

func PolygonEqual(ply1, ply2 [][][]float64) bool {
	if len(ply1) != len(ply2) {
		return false
	}

	var points1, points2 [][]float64
	for i := range ply1 {
		points1 = append(points1, ply1[i]...)
	}
	extent1 := NewExtent(points1...)
	for i := range ply2 {
		points2 = append(points2, ply2[i]...)
	}
	extent2 := NewExtent(points2...)
	if !GeomExtent(extent1, extent2) {
		return false
	}

	sort.Sort(bySubRingSizeXY(ply1))
	sort.Sort(bySubRingSizeXY(ply2))
	for i := range ply1 {
		if !LineStringEqual(ply1[i], ply2[i]) {
			return false
		}
	}
	return true
}

func PointerEqual(geo1, geo2 geom.Point) bool {
	if geo1Nil, geo2Nil := geo1 == NilPoint, geo2 == NilPoint; geo1Nil || geo2Nil {
		return geo1Nil && geo2Nil
	}
	return PointEqual(geo1.Data(), geo2.Data())
}

func PointerLess(p1, p2 Point) bool { return PointLess(p1.Data(), p2.Data()) }

func MultiPointerEqual(geo1, geo2 geom.MultiPoint) bool {
	if geo1Nil, geo2Nil := geo1 == NilMultiPoint, geo2 == NilMultiPoint; geo1Nil || geo2Nil {
		return geo1Nil && geo2Nil
	}
	return MultiPointEqual(geo1.Data(), geo2.Data())
}

func LineStringerEqual(geo1, geo2 geom.LineString) bool {
	if geo1Nil, geo2Nil := geo1 == NilLineString, geo2 == NilLineString; geo1Nil || geo2Nil {
		return geo1Nil && geo2Nil
	}
	return LineStringEqual(geo1.Data(), geo2.Data())
}

func MultiLineStringerEqual(geo1, geo2 geom.MultiLine) bool {
	if geo1Nil, geo2Nil := geo1 == NilMultiLine, geo2 == NilMultiLine; geo1Nil || geo2Nil {
		return geo1Nil && geo2Nil
	}
	return MultiLineEqual(geo1.Data(), geo2.Data())
}

func PolygonerEqual(geo1, geo2 geom.Polygon) bool {
	if geo1Nil, geo2Nil := geo1 == NilPoly, geo2 == NilPoly; geo1Nil || geo2Nil {
		return geo1Nil && geo2Nil
	}
	return PolygonEqual(geo1.Data(), geo2.Data())
}

func MultiPolygonerEqual(geo1, geo2 geom.MultiPolygon) bool {
	if geo1Nil, geo2Nil := geo1 == NilMultiPoly, geo2 == NilMultiPoly; geo1Nil || geo2Nil {
		return geo1Nil && geo2Nil
	}

	p1, p2 := geo1.Data(), geo2.Data()
	if len(p1) != len(p2) {
		return false
	}
	sort.Sort(byPolygonMainSizeXY(p1))
	sort.Sort(byPolygonMainSizeXY(p2))
	for i := range p1 {
		if !PolygonEqual(p1[i], p2[i]) {
			return false
		}
	}
	return true
}

func CollectionerEqual(col1, col2 geom.Collection) bool {
	if colNil, col2Nil := col1 == nil, col2 == nil; colNil || col2Nil {
		return colNil && col2Nil
	}

	g1, g2 := col1.Geometries(), col2.Geometries()
	if len(g1) != len(g2) {
		return false
	}
	for i := range g1 {
		if !GeometryEqual(g1[i], g2[i]) {
			return false
		}
	}
	return true
}

func GeometryEqual(g1, g2 geom.Geometry) bool {
	switch pg1 := g1.(type) {
	case geom.Point:
		if pg2, ok := g2.(geom.Point); ok {
			return PointerEqual(pg1, pg2)
		}
	case geom.MultiPoint:
		if pg2, ok := g2.(geom.MultiPoint); ok {
			return MultiPointerEqual(pg1, pg2)
		}
	case geom.LineString:
		if pg2, ok := g2.(geom.LineString); ok {
			return LineStringerEqual(pg1, pg2)
		}
	case geom.MultiLine:
		if pg2, ok := g2.(geom.MultiLine); ok {
			return MultiLineStringerEqual(pg1, pg2)
		}
	case geom.Polygon:
		if pg2, ok := g2.(geom.Polygon); ok {
			return PolygonerEqual(pg1, pg2)
		}
	case geom.MultiPolygon:
		if pg2, ok := g2.(geom.MultiPolygon); ok {
			return MultiPolygonerEqual(pg1, pg2)
		}
	case geom.Collection:
		if pg2, ok := g2.(geom.Collection); ok {
			return CollectionerEqual(pg1, pg2)
		}
	}
	return false
}
