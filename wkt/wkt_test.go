package wkt

import (
	"bytes"
	"testing"
)

func TestWKT(t *testing.T) {
	//data := "SRID=4312;GEOMETRYCOLLECTION(POINT(4 6),LINESTRING(4 6,7 10))"
	//data := "SRID=4312;POLYGONZ((0 0 0,4 0 0,4 4 0,0 4 0,0 0 0),(1 1 0,2 1 0,2 2 0,1 2 0,1 1 0))"
	data := "SRID=4312;MULTIPOLYGONZ(((0 0 0,4 0 0,4 4 0,0 4 0,0 0 0),(1 1 0,2 1 0,2 2 0,1 2 0,1 1 0)),((-1 -1 0,-1 -2 0,-2 -2 0,-2 -1 0,-1 -1 0)))"
	geo, srid, err := DecodeWKT([]byte(data))
	if err != nil {
		t.Error(err)
	}

	if srid != 4312 {
		t.Error("err")
	}

	if geo == nil {
		t.Error("err")
	}

	var geobuf bytes.Buffer

	EncodeWKT(geo, &srid, &geobuf)

	str := string(geobuf.Bytes())

	if len(str) == 0 {
		t.Error("err")
	}
}
