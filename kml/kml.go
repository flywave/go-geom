package kml

import (
	"fmt"

	"github.com/flywave/go-geom"
	"github.com/twpayne/go-kml/v3"
)

func Encode(g geom.Geometry) (kml.Element, error) {
	switch g := g.(type) {
	case geom.Point3:
		return EncodePoint3(g), nil
	case geom.Point:
		return EncodePoint(g), nil
	case geom.LineString3:
		return EncodeLineString3(g), nil
	case geom.LineString:
		return EncodeLineString(g), nil
	case geom.MultiLine3:
		return EncodeMultiLineString3(g), nil
	case geom.MultiLine:
		return EncodeMultiLineString(g), nil
	case geom.MultiPoint3:
		return EncodeMultiPoint3(g), nil
	case geom.MultiPoint:
		return EncodeMultiPoint(g), nil
	case geom.MultiPolygon3:
		return EncodeMultiPolygon3(g), nil
	case geom.MultiPolygon:
		return EncodeMultiPolygon(g), nil
	case geom.Polygon3:
		return EncodePolygon3(g), nil
	case geom.Polygon:
		return EncodePolygon(g), nil
	case geom.Collection:
		return EncodeGeometryCollection(g)
	default:
		return nil, fmt.Errorf("unsupport geom")
	}
}

func FlatCoords(pts [][]float64, dim int) []float64 {
	if dim == 2 {
		ret := make([]float64, len(pts)*2)
		for i := range pts {
			ret[i*2] = pts[i][0]
			ret[i*2+1] = pts[i][1]
		}
		return ret
	} else if dim == 3 {
		ret := make([]float64, len(pts)*3)
		for i := range pts {
			ret[i*3] = pts[i][0]
			ret[i*3+1] = pts[i][1]
			ret[i*3+2] = pts[i][2]
		}
		return ret
	}
	return nil
}

func EncodeLineString(ls geom.LineString) kml.Element {
	flatCoords := FlatCoords(ls.Data(), 2)
	return kml.LineString(kml.CoordinatesFlat(flatCoords, 0, len(flatCoords), 2, 2))
}

func EncodeLineString3(ls geom.LineString3) kml.Element {
	flatCoords := FlatCoords(ls.Data(), 3)
	return kml.LineString(kml.CoordinatesFlat(flatCoords, 0, len(flatCoords), 3, 3))
}

func EncodeMultiLineString(mls geom.MultiLine) kml.Element {
	num := len(mls.Lines())
	lineStrings := make([]kml.Element, num)
	for i, ls := range mls.Lines() {
		flatCoords := FlatCoords(ls.Data(), 2)
		lineStrings[i] = kml.LineString(kml.CoordinatesFlat(flatCoords, 0, len(flatCoords), 2, 2))
	}
	return kml.MultiGeometry(lineStrings...)
}

func EncodeMultiLineString3(mls geom.MultiLine3) kml.Element {
	num := len(mls.Lines())
	lineStrings := make([]kml.Element, num)
	for i, ls := range mls.Lines() {
		flatCoords := FlatCoords(ls.Data(), 3)
		lineStrings[i] = kml.LineString(kml.CoordinatesFlat(flatCoords, 0, len(flatCoords), 3, 3))
	}
	return kml.MultiGeometry(lineStrings...)
}

func EncodeMultiPoint(mp geom.MultiPoint) kml.Element {
	num := len(mp.Points())
	points := make([]kml.Element, num)
	for i, pt := range mp.Points() {
		points[i] = kml.Point(kml.CoordinatesFlat(pt.Data(), 0, 2, 2, 2))
	}
	return kml.MultiGeometry(points...)
}

func EncodeMultiPoint3(mp geom.MultiPoint3) kml.Element {
	num := len(mp.Points())
	points := make([]kml.Element, num)
	for i, pt := range mp.Points() {
		points[i] = kml.Point(kml.CoordinatesFlat(pt.Data(), 0, 3, 3, 3))
	}
	return kml.MultiGeometry(points...)
}

func EncodeMultiPolygon(mp geom.MultiPolygon) kml.Element {
	num := len(mp.Polygons())
	polygons := make([]kml.Element, num)
	for i, pls := range mp.Polygons() {
		numb := len(pls.Sublines())
		boundaries := make([]kml.Element, numb)
		for j, ls := range pls.Sublines() {
			flatCoords := FlatCoords(ls.Data(), 2)
			linearRing := kml.LinearRing(kml.CoordinatesFlat(flatCoords, 0, len(flatCoords), 2, 2))
			if j == 0 {
				boundaries[j] = kml.OuterBoundaryIs(linearRing)
			} else {
				boundaries[j] = kml.InnerBoundaryIs(linearRing)
			}
		}
		polygons[i] = kml.Polygon(boundaries...)
	}
	return kml.MultiGeometry(polygons...)
}

func EncodeMultiPolygon3(mp geom.MultiPolygon3) kml.Element {
	num := len(mp.Polygons())
	polygons := make([]kml.Element, num)
	for i, pls := range mp.Polygons() {
		numb := len(pls.Sublines())
		boundaries := make([]kml.Element, numb)
		for j, ls := range pls.Sublines() {
			flatCoords := FlatCoords(ls.Data(), 3)
			linearRing := kml.LinearRing(kml.CoordinatesFlat(flatCoords, 0, len(flatCoords), 3, 3))
			if j == 0 {
				boundaries[j] = kml.OuterBoundaryIs(linearRing)
			} else {
				boundaries[j] = kml.InnerBoundaryIs(linearRing)
			}
		}
		polygons[i] = kml.Polygon(boundaries...)
	}
	return kml.MultiGeometry(polygons...)
}

func EncodePoint(p geom.Point) kml.Element {
	return kml.Point(kml.CoordinatesFlat(p.Data(), 0, 2, 2, 2))
}

func EncodePoint3(p geom.Point3) kml.Element {
	return kml.Point(kml.CoordinatesFlat(p.Data(), 0, 3, 3, 3))
}

func EncodePolygon(p geom.Polygon) kml.Element {
	num := len(p.Sublines())
	boundaries := make([]kml.Element, num)
	for i, ls := range p.Sublines() {
		flatCoords := FlatCoords(ls.Data(), 2)
		linearRing := kml.LinearRing(kml.CoordinatesFlat(flatCoords, 0, len(flatCoords), 2, 2))
		if i == 0 {
			boundaries[i] = kml.OuterBoundaryIs(linearRing)
		} else {
			boundaries[i] = kml.InnerBoundaryIs(linearRing)
		}
	}
	return kml.Polygon(boundaries...)
}

func EncodePolygon3(p geom.Polygon3) kml.Element {
	num := len(p.Sublines())
	boundaries := make([]kml.Element, num)
	for i, ls := range p.Sublines() {
		flatCoords := FlatCoords(ls.Data(), 3)
		linearRing := kml.LinearRing(kml.CoordinatesFlat(flatCoords, 0, len(flatCoords), 3, 3))
		if i == 0 {
			boundaries[i] = kml.OuterBoundaryIs(linearRing)
		} else {
			boundaries[i] = kml.InnerBoundaryIs(linearRing)
		}
	}
	return kml.Polygon(boundaries...)
}

func EncodeGeometryCollection(geoms geom.Collection) (kml.Element, error) {
	geometries := make([]kml.Element, len(geoms))
	for i, g := range geoms {
		var err error
		geometries[i], err = Encode(g)
		if err != nil {
			return nil, err
		}
	}
	return kml.MultiGeometry(geometries...), nil
}
