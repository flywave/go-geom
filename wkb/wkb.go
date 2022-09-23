package wkb

import (
	"errors"
	"io"

	"github.com/devork/geom"
	"github.com/devork/geom/ewkb"

	geom_ "github.com/flywave/go-geom"
)

const (
	DEFAULT_SRSID = 4326
)

func ConvertToGeom(g *geom_.GeometryData, srsid *uint32) (geom.Geometry, geom.Dimension) {
	var sid uint32
	if srsid != nil {
		sid = *srsid
	} else {
		sid = DEFAULT_SRSID
	}
	var geo geom.Geometry
	var dim geom.Dimension
	switch g.Type {
	case "Point":
		if len(g.Point) == 3 {
			dim = geom.XYZ
			geo = &geom.Point{Hdr: geom.Hdr{Dim: geom.XYZ, Srid: sid}, Coordinate: g.Point}
		} else {
			dim = geom.XY
			geo = &geom.Point{Hdr: geom.Hdr{Dim: geom.XY, Srid: sid}, Coordinate: g.Point}
		}
	case "MultiPoint":
		var mp geom.MultiPoint
		if len(g.MultiPoint) > 0 && len(g.MultiPoint[0]) == 3 {
			dim = geom.XYZ
			mp = geom.MultiPoint{Hdr: geom.Hdr{Dim: geom.XYZ, Srid: sid}}
		} else {
			dim = geom.XY
			mp = geom.MultiPoint{Hdr: geom.Hdr{Dim: geom.XY, Srid: sid}}
		}
		for i := range g.MultiPoint {
			mp.Points = append(mp.Points, geom.Point{Hdr: geom.Hdr{Dim: dim, Srid: sid}, Coordinate: g.MultiPoint[i]})
		}
		geo = &mp
	case "LineString":
		var mp geom.LineString
		if len(g.LineString) > 0 && len(g.LineString[0]) == 3 {
			dim = geom.XYZ
			mp = geom.LineString{Hdr: geom.Hdr{Dim: geom.XYZ, Srid: sid}}
		} else {
			dim = geom.XY
			mp = geom.LineString{Hdr: geom.Hdr{Dim: geom.XY, Srid: sid}}
		}
		for i := range g.LineString {
			mp.Coordinates = append(mp.Coordinates, g.LineString[i])
		}
		geo = &mp
	case "MultiLineString":
		var mp geom.MultiLineString
		if len(g.MultiLineString) > 0 && len(g.MultiLineString[0]) > 0 && len(g.MultiLineString[0][0]) == 3 {
			dim = geom.XYZ
			mp = geom.MultiLineString{Hdr: geom.Hdr{Dim: dim, Srid: sid}}
		} else {
			dim = geom.XY
			mp = geom.MultiLineString{Hdr: geom.Hdr{Dim: dim, Srid: sid}}
		}
		for i := range g.MultiLineString {
			l := geom.LineString{Hdr: geom.Hdr{Dim: dim, Srid: 0}}
			for j := range g.MultiLineString[i] {
				l.Coordinates = append(l.Coordinates, g.MultiLineString[i][j])
			}
			mp.LineStrings = append(mp.LineStrings, l)
		}
		geo = &mp
	case "Polygon":
		var mp geom.Polygon
		if len(g.Polygon) > 0 && len(g.Polygon[0]) > 0 && len(g.Polygon[0][0]) == 3 {
			dim = geom.XYZ
			mp = geom.Polygon{Hdr: geom.Hdr{Dim: geom.XYZ, Srid: sid}}
		} else {
			dim = geom.XY
			mp = geom.Polygon{Hdr: geom.Hdr{Dim: geom.XY, Srid: sid}}
		}
		for i := range g.Polygon {
			l := geom.LinearRing{}
			for j := range g.Polygon[i] {
				l.Coordinates = append(l.Coordinates, g.Polygon[i][j])
			}
			mp.Rings = append(mp.Rings, l)
		}
		geo = &mp
	case "MultiPolygon":
		var mp geom.MultiPolygon
		if len(g.MultiPolygon) > 0 && len(g.MultiPolygon[0]) > 0 && len(g.MultiPolygon[0][0]) > 0 && len(g.MultiPolygon[0][0][0]) == 3 {
			dim = geom.XYZ
			mp = geom.MultiPolygon{Hdr: geom.Hdr{Dim: dim, Srid: sid}}
		} else {
			dim = geom.XY
			mp = geom.MultiPolygon{Hdr: geom.Hdr{Dim: dim, Srid: sid}}
		}
		for i := range g.MultiPolygon {
			pol := geom.Polygon{Hdr: geom.Hdr{Dim: dim, Srid: 0}}
			for j := range g.MultiPolygon[i] {

				l := geom.LinearRing{}
				for k := range g.MultiPolygon[i][j] {
					l.Coordinates = append(l.Coordinates, g.MultiPolygon[i][j][k])
				}
				pol.Rings = append(pol.Rings, l)
			}
			mp.Polygons = append(mp.Polygons, pol)
		}
		geo = &mp
	case "GeometryCollection":
		mp := geom.GeometryCollection{Hdr: geom.Hdr{Dim: geom.XYZ, Srid: sid}}
		for i := range g.Geometries {
			ge, d := ConvertToGeom(g.Geometries[i], srsid)
			mp.Geometries = append(mp.Geometries, ge)
			dim = d
		}
		mp.Dim = dim
		geo = &mp
	}
	return geo, dim
}

func EncodeWKB(g *geom_.GeometryData, srsid *uint32, w io.Writer) error {
	geo, _ := ConvertToGeom(g, srsid)
	return ewkb.Encode(geo, w)
}

func ConvertFromGeom(g geom.Geometry) (*geom_.GeometryData, error) {
	var ret geom_.GeometryData
	switch geo := g.(type) {
	case *geom.Point:
		ret.Type = "Point"
		ret.Point = geo.Coordinate
	case *geom.LineString:
		ret.Type = "LineString"
		for i := range geo.Coordinates {
			ret.LineString = append(ret.LineString, geo.Coordinates[i])
		}
	case *geom.Polygon:
		ret.Type = "Polygon"
		for i := range geo.Rings {
			ring := geo.Rings[i]
			var lr [][]float64
			for j := range ring.Coordinates {
				lr = append(lr, ring.Coordinates[j])
			}
			ret.Polygon = append(ret.Polygon, lr)
		}
	case *geom.MultiPoint:
		ret.Type = "MultiPoint"
		for i := range geo.Points {
			ret.MultiPoint = append(ret.MultiPoint, geo.Points[i].Coordinate)
		}
	case *geom.MultiLineString:
		ret.Type = "MultiLineString"
		for i := range geo.LineStrings {
			ring := geo.LineStrings[i]
			var lr [][]float64
			for j := range ring.Coordinates {
				lr = append(lr, ring.Coordinates[j])
			}
			ret.MultiLineString = append(ret.MultiLineString, lr)
		}
	case *geom.MultiPolygon:
		ret.Type = "MultiPolygon"
		for i := range geo.Polygons {
			var l [][][]float64
			for j := range geo.Polygons[i].Rings {
				ring := geo.Polygons[i].Rings[j]
				var lr [][]float64
				for k := range ring.Coordinates {
					lr = append(lr, ring.Coordinates[k])
				}
				l = append(l, lr)
			}
			ret.MultiPolygon = append(ret.MultiPolygon, l)
		}
	case *geom.GeometryCollection:
		ret.Type = "GeometryCollection"
		for i := range geo.Geometries {
			geoc, err := ConvertFromGeom(geo.Geometries[i])
			if err != nil {
				return nil, err
			}
			ret.Geometries = append(ret.Geometries, geoc)
		}
	default:
		return nil, errors.New("error not support")
	}
	return &ret, nil
}

func DecodeWKB(r io.Reader) (*geom_.GeometryData, uint32, error) {
	geom, err := ewkb.Decode(r)
	if err != nil {
		return nil, 0, err
	}
	g, err := ConvertFromGeom(geom)
	return g, geom.SRID(), err
}
