package wkt

import (
	"bytes"
	"io"
	"strconv"

	"github.com/flywave/go-geom"
)

const (
	_ Opt = iota
	Z
	M
	ZM
)

type Opt uint

func (o Opt) Is3d() bool { return o&Z != 0 }

func (o Opt) IsMeasured() bool { return o&M != 0 }

func (o Opt) Is3dMeasured() bool { return o&ZM != 0 }

type Coord [4]float64

func dumpPoint(buffer *bytes.Buffer, coordinate []float64) int {
	var dim int
	if len(coordinate) == 3 {
		buffer.WriteString("POINTZ(")
		dim = 3
	} else {
		buffer.WriteString("POINT(")
		dim = 2
	}
	for i := range coordinate {
		buffer.WriteString(strconv.FormatFloat(coordinate[i], 'f', -1, 32))
		if i < len(coordinate)-1 {
			buffer.WriteString(" ")
		} else {
			buffer.WriteString(")")
		}
	}
	return dim
}

func dumpMultiPoint(buffer *bytes.Buffer, coordinates ...[]float64) int {
	var dim int
	if len(coordinates) > 0 && len(coordinates[0]) == 3 {
		buffer.WriteString("MULTIPOINTZ(")
		dim = 3
	} else {
		buffer.WriteString("MULTIPOINT(")
		dim = 2
	}

	for i := range coordinates {
		for j := range coordinates[i] {
			buffer.WriteString(strconv.FormatFloat(coordinates[i][j], 'f', -1, 32))
			if j < len(coordinates[i])-1 {
				buffer.WriteString(" ")
			}
		}
		if i < len(coordinates)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString(")")
	return dim
}

func dumpLineString(buffer *bytes.Buffer, coordinates [][]float64) int {
	var dim int
	if len(coordinates) > 0 && len(coordinates[0]) == 3 {
		buffer.WriteString("LINESTRINGZ(")
		dim = 3
	} else {
		buffer.WriteString("LINESTRING(")
		dim = 2
	}

	for i := range coordinates {
		for j := range coordinates[i] {
			buffer.WriteString(strconv.FormatFloat(coordinates[i][j], 'f', -1, 32))
			if j < len(coordinates[i])-1 {
				buffer.WriteString(" ")
			}
		}
		if i < len(coordinates)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString(")")
	return dim
}

func dumpMultiLineString(buffer *bytes.Buffer, lines ...[][]float64) int {
	var dim int
	if len(lines) > 0 && len(lines[0]) > 0 && len(lines[0][0]) == 3 {
		buffer.WriteString("MULTILINESTRINGZ(")
		dim = 3
	} else {
		buffer.WriteString("MULTILINESTRING(")
		dim = 2
	}

	for i := range lines {
		buffer.WriteString("(")
		for j := range lines[i] {
			for k := range lines[i][j] {
				buffer.WriteString(strconv.FormatFloat(lines[i][j][k], 'f', -1, 32))
				if k < len(lines[i][j])-1 {
					buffer.WriteString(" ")
				}
			}
			if j < len(lines[i])-1 {
				buffer.WriteString(",")
			}
		}
		buffer.WriteString(")")
		if i < len(lines)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString(")")
	return dim
}

func dumpPolygon(buffer *bytes.Buffer, polygon [][][]float64) int {
	var dim int
	if len(polygon) > 0 && len(polygon[0]) > 0 && len(polygon[0][0]) == 3 {
		buffer.WriteString("POLYGONZ(")
		dim = 3
	} else {
		buffer.WriteString("POLYGON(")
		dim = 2
	}

	for i := range polygon {
		buffer.WriteString("(")
		for j := range polygon[i] {
			for k := range polygon[i][j] {
				buffer.WriteString(strconv.FormatFloat(polygon[i][j][k], 'f', -1, 32))
				if k < len(polygon[i][j])-1 {
					buffer.WriteString(" ")
				}
			}
			if j < len(polygon[i])-1 {
				buffer.WriteString(",")
			}
		}
		buffer.WriteString(")")
		if i < len(polygon)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString(")")
	return dim
}

func dumpMultiPolygon(buffer *bytes.Buffer, polygons ...[][][]float64) int {
	var dim int
	if len(polygons) > 0 && len(polygons[0]) > 0 && len(polygons[0][0]) > 0 && len(polygons[0][0][0]) == 3 {
		buffer.WriteString("MULTIPOLYGONZ(")
		dim = 3
	} else {
		buffer.WriteString("MULTIPOLYGON(")
		dim = 2
	}

	for o := range polygons {
		buffer.WriteString("(")
		for i := range polygons[o] {
			buffer.WriteString("(")
			for j := range polygons[o][i] {
				for k := range polygons[o][i][j] {
					buffer.WriteString(strconv.FormatFloat(polygons[o][i][j][k], 'f', -1, 32))
					if k < len(polygons[o][i][j])-1 {
						buffer.WriteString(" ")
					}
				}
				if j < len(polygons[o][i])-1 {
					buffer.WriteString(",")
				}
			}
			buffer.WriteString(")")
			if i < len(polygons[o])-1 {
				buffer.WriteString(",")
			}
		}
		buffer.WriteString(")")
		if o < len(polygons)-1 {
			buffer.WriteString(",")
		}
	}

	buffer.WriteString(")")
	return dim
}

func dumpCollection(buffer *bytes.Buffer, geometries ...*geom.GeometryData) int {
	var geobuf bytes.Buffer
	var dim int

	for i := range geometries {
		switch geometries[i].Type {
		case "Point":
			dim = dumpPoint(&geobuf, geometries[i].Point)
		case "MultiPoint":
			dim = dumpMultiPoint(&geobuf, geometries[i].MultiPoint...)
		case "LineString":
			dim = dumpLineString(&geobuf, geometries[i].LineString)
		case "MultiLineString":
			dim = dumpMultiLineString(&geobuf, geometries[i].MultiLineString...)
		case "Polygon":
			dim = dumpPolygon(&geobuf, geometries[i].Polygon)
		case "MultiPolygon":
			dim = dumpMultiPolygon(&geobuf, geometries[i].MultiPolygon...)
		case "GeometryCollection":
			dim = dumpCollection(&geobuf, geometries[i].Geometries...)
		}
		if i < len(geometries)-1 {
			geobuf.WriteString(",")
		}
	}

	if dim == 3 {
		buffer.WriteString("GEOMETRYCOLLECTIONZ(")
	} else {
		buffer.WriteString("GEOMETRYCOLLECTION(")
	}

	buffer.Write(geobuf.Bytes())

	buffer.WriteString(")")

	return dim
}

func EncodeWKT(g *geom.GeometryData, srsid *uint32, w io.Writer) error {
	var geobuf bytes.Buffer

	if srsid != nil {
		geobuf.WriteString("SRID=")
		geobuf.WriteString(strconv.FormatUint(uint64(*srsid), 10))
		geobuf.WriteString(";")
	}

	switch g.Type {
	case "Point":
		_ = dumpPoint(&geobuf, g.Point)
	case "MultiPoint":
		_ = dumpMultiPoint(&geobuf, g.MultiPoint...)
	case "LineString":
		_ = dumpLineString(&geobuf, g.LineString)
	case "MultiLineString":
		_ = dumpMultiLineString(&geobuf, g.MultiLineString...)
	case "Polygon":
		_ = dumpPolygon(&geobuf, g.Polygon)
	case "MultiPolygon":
		_ = dumpMultiPolygon(&geobuf, g.MultiPolygon...)
	case "GeometryCollection":
		_ = dumpCollection(&geobuf, g.Geometries...)
	}

	_, err := w.Write(geobuf.Bytes())
	return err
}

func DecodeWKT(data []byte) (*geom.GeometryData, uint32, error) {
	s := &scanner{raw: data}
	g, err := s.scanGeom()
	return g, s.srid, err
}
