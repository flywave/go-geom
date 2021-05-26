package wkt

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/flywave/go-geom"
)

type scanner struct {
	raw  []byte
	i    int
	opt  Opt
	srid uint32
}

func (s *scanner) peek() (byte, error) {
	if s.i >= len(s.raw) {
		return '\x00', io.ErrUnexpectedEOF
	}
	return s.raw[s.i], nil
}

func (s *scanner) skipWs() {
	if s.i >= len(s.raw) {
		return
	}
	for i, b := range s.raw[s.i:] {
		if b == ' ' || b == '\n' || b == '\t' || b == '\r' {
			continue
		}
		s.i += i
		return
	}
}

func (s *scanner) scanSrid() error {
	s.skipWs()
	if s.i >= len(s.raw) {
		return io.ErrUnexpectedEOF
	}
	var (
		id []byte
		b  byte
		i  int
	)
	if s.i+5 < len(s.raw) {
		ins := s.raw[s.i : s.i+5]
		if (ins[0] == 'S' || ins[0] == 's') &&
			(ins[1] == 'R' || ins[1] == 'r') &&
			(ins[2] == 'I' || ins[2] == 'i') &&
			(ins[3] == 'D' || ins[3] == 'd') &&
			ins[4] == '=' {
			for i, b = range s.raw[s.i+5:] {
				if b >= '0' && b <= '9' {
					id = append(id, b)
					continue
				}
				if b == ';' {
					break
				}
			}
			if len(id) == 0 {
				return fmt.Errorf("no srid byte %q", b)
			}
			s.i += 5 + i + 1
			sid, err := strconv.ParseUint(string(id), 10, 32)
			if err != nil {
				return err
			}
			s.srid = uint32(sid)
		}
	}

	return nil
}

func (s *scanner) scanStart() error {
	s.skipWs()
	c, err := s.peek()
	if err != nil {
		return err
	}
	if c != '(' {
		return fmt.Errorf("expect '(' got %q", c)
	}
	s.i++
	return nil
}

func (s *scanner) scanContinue() (bool, error) {
	s.skipWs()
	c, err := s.peek()
	if err != nil {
		return false, err
	}
	comma := c == ','
	if !comma && c != ')' {
		return false, fmt.Errorf("expect ',' or ')' got %q", c)
	}
	s.i++
	return comma, nil
}

func (s *scanner) scanIdent() (string, error) {
	s.skipWs()
	if s.i >= len(s.raw) {
		return "", io.ErrUnexpectedEOF
	}
	var (
		ident []byte
		b     byte
		i     int
	)
	for i, b = range s.raw[s.i:] {
		lower := b >= 'a' && b <= 'z'
		if lower || b >= 'A' && b <= 'Z' {
			if lower {
				b = b - 'a' + 'A'
			}
			ident = append(ident, b)
			continue
		}
		break
	}
	if len(ident) == 0 {
		return "", fmt.Errorf("no ident byte %q", b)
	}
	s.i += i
	return string(ident), nil
}

func (s *scanner) scanCoord() (c Coord, comma bool, err error) {
	s.skipWs()
	if s.i >= len(s.raw) {
		return c, false, io.ErrUnexpectedEOF
	}
	r := bytes.NewReader(s.raw[s.i:])
	var fs []*float64
	if s.opt.Is3d() {
		fs = []*float64{&c[0], &c[1], &c[2]}
	} else if s.opt.IsMeasured() {
		fs = []*float64{&c[0], &c[1], &c[3]}
	} else if s.opt.Is3dMeasured() {
		fs = []*float64{&c[0], &c[1], &c[2], &c[3]}
	} else {
		fs = []*float64{&c[0], &c[1]}
	}
	for _, f := range fs {
		_, err := fmt.Fscan(r, f)
		if err != nil {
			return c, false, err
		}
		s.i = len(s.raw) - r.Len()
		s.skipWs()
		b, err := s.peek()
		if err != nil {
			return c, false, io.ErrUnexpectedEOF
		}
		if comma = b == ','; comma || b == ')' {
			s.i++
			break
		}
	}
	return
}

func (s *scanner) scanCoords(multi bool) ([]Coord, error) {
	err := s.scanStart()
	if err != nil {
		return nil, err
	}
	var cs []Coord
	var c Coord
	var comma bool
	if multi {
		err = s.scanStart()
		multi = err == nil
	}
	for {
		c, comma, err = s.scanCoord()
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
		if comma {
			if multi {
				return nil, fmt.Errorf("expect ')' got ','")
			}
			continue
		}
		if multi {
			comma, err = s.scanContinue()
			if err != nil {
				return nil, err
			}
			if comma {
				err = s.scanStart()
				if err != nil {
					return nil, err
				}
				continue
			}
		}
		return cs, nil
	}
}

func (s *scanner) scanMultiLinedata() ([][]Coord, error) {
	err := s.scanStart()
	if err != nil {
		return nil, err
	}
	var poly [][]Coord
	var cs []Coord
	var comma bool
	for {
		cs, err = s.scanCoords(false)
		if err != nil {
			return nil, err
		}
		if len(cs) < 4 {
			return nil, fmt.Errorf("a polygon ring must have at least 4 points, got %d", len(cs))
		}
		poly = append(poly, cs)
		comma, err = s.scanContinue()
		if err != nil {
			return nil, err
		}
		if comma {
			continue
		}
		return poly, nil
	}
}

func (s *scanner) scanPolydata() ([][]Coord, error) {
	err := s.scanStart()
	if err != nil {
		return nil, err
	}
	var poly [][]Coord
	var cs []Coord
	var comma bool
	for {
		cs, err = s.scanCoords(false)
		if err != nil {
			return nil, err
		}
		if cs[0] != cs[len(cs)-1] {
			return nil, fmt.Errorf("a polygon ring must be closed")
		}
		poly = append(poly, cs)
		comma, err = s.scanContinue()
		if err != nil {
			return nil, err
		}
		if comma {
			continue
		}
		return poly, nil
	}
}

func (s *scanner) scanMultiPolydata() ([][][]Coord, error) {
	err := s.scanStart()
	if err != nil {
		return nil, err
	}
	var multi [][][]Coord
	var poly [][]Coord
	var comma bool
	for {
		poly, err = s.scanPolydata()
		if err != nil {
			return nil, err
		}
		multi = append(multi, poly)
		comma, err = s.scanContinue()
		if err != nil {
			return nil, err
		}
		if comma {
			continue
		}
		return multi, nil
	}
}

func (s *scanner) scanGeom() (*geom.GeometryData, error) {
	err := s.scanSrid()
	if err != nil {
		return nil, err
	}
	ident, err := s.scanIdent()
	if err != nil {
		return nil, err
	}
	if ident[len(ident)-1] == 'Z' {
		s.opt = Z
	} else if ident[len(ident)-1] == 'M' {
		if ident[len(ident)-2] == 'Z' {
			s.opt = ZM
		} else {
			s.opt = M
		}
	} else {
		s.opt = 0
	}
	var g geom.GeometryData
	switch ident {
	case "POINT", "POINTZ", "MULTIPOINT", "MULTIPOINTZ", "LINESTRING", "LINESTRINGZ":
		var cs []Coord
		cs, err = s.scanCoords(ident == "MULTIPOINT" || ident == "MULTIPOINTZ")
		if err != nil {
			break
		}
		switch ident {
		case "POINT", "POINTZ":
			if len(cs) != 1 {
				return nil, fmt.Errorf("expected 1 got %d points", len(cs))
			}
			g.Type = "Point"
			if s.opt.Is3d() || s.opt.Is3dMeasured() {
				g.Point = cs[0][0:3]
			} else {
				g.Point = cs[0][0:2]
			}
		case "MULTIPOINT", "MULTIPOINTZ":
			g.Type = "MultiPoint"
			for i := range cs {
				if s.opt.Is3d() || s.opt.Is3dMeasured() {
					g.MultiPoint = append(g.MultiPoint, cs[i][0:3])
				} else {
					g.MultiPoint = append(g.MultiPoint, cs[i][0:2])
				}
			}
		case "LINESTRING", "LINESTRINGZ":
			g.Type = "LineString"
			for i := range cs {
				if s.opt.Is3d() || s.opt.Is3dMeasured() {
					g.LineString = append(g.LineString, cs[i][0:3])
				} else {
					g.LineString = append(g.LineString, cs[i][0:2])
				}
			}
		}
	case "MULTILINESTRING", "MULTILINESTRINGZ":
		var rings [][]Coord
		rings, err = s.scanMultiLinedata()
		if err != nil {
			break
		}
		g.Type = "MultiLineString"
		for i := range rings {
			cs := rings[i]
			var l [][]float64
			for j := range cs {
				if s.opt.Is3d() || s.opt.Is3dMeasured() {
					l = append(l, cs[j][0:3])
				} else {
					l = append(l, cs[j][0:2])
				}
			}
			g.MultiLineString = append(g.MultiLineString, l)
		}
	case "POLYGON", "POLYGONZ":
		var rings [][]Coord
		rings, err = s.scanPolydata()
		if err != nil {
			break
		}
		g.Type = "Polygon"
		for i := range rings {
			cs := rings[i]
			var l [][]float64
			for j := range cs {
				if s.opt.Is3d() || s.opt.Is3dMeasured() {
					l = append(l, cs[j][0:3])
				} else {
					l = append(l, cs[j][0:2])
				}
			}
			g.Polygon = append(g.Polygon, l)
		}
	case "MULTIPOLYGON", "MULTIPOLYGONZ":
		var multi [][][]Coord
		multi, err = s.scanMultiPolydata()
		if err != nil {
			break
		}
		g.Type = "MultiPolygon"
		for i := range multi {
			var p [][][]float64
			for j := range multi[i] {
				cs := multi[i][j]
				var l [][]float64
				for k := range cs {
					if s.opt.Is3d() || s.opt.Is3dMeasured() {
						l = append(l, cs[k][0:3])
					} else {
						l = append(l, cs[k][0:2])
					}
				}
				p = append(p, l)
			}
			g.MultiPolygon = append(g.MultiPolygon, p)
		}
	case "GEOMETRYCOLLECTION", "GEOMETRYCOLLECTIONZ":
		err := s.scanStart()
		if err != nil {
			break
		}
		g.Type = "GeometryCollection"
		for {
			geo, err := s.scanGeom()
			if err != nil {
				break
			}

			g.Geometries = append(g.Geometries, geo)
			comma, err := s.scanContinue()
			if err != nil {
				return nil, err
			}
			if comma {
				continue
			} else {
				break
			}
		}
	default:
		err = fmt.Errorf("unknown geom '%s'", ident)
	}
	if err != nil {
		return nil, err
	}
	return &g, nil
}
