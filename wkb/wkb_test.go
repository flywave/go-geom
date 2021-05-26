package wkb

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/devork/geom"
)

func TestTile(t *testing.T) {
	datasets := []struct {
		data string
		srid uint32
		dim  geom.Dimension
	}{
		{"002000000700006c340000000200000000014010000000000000401800000000000000000000020000000240100000000000004018000000000000401c0000000000004024000000000000", 27700, geom.XY}, // ewkb_xdr
		{"0107000020346c000002000000010100000000000000000010400000000000001840010200000002000000000000000000104000000000000018400000000000001c400000000000002440", 27700, geom.XY}, // ewkb_hdr
		{"00000000070000000200000000014010000000000000401800000000000000000000020000000240100000000000004018000000000000401c0000000000004024000000000000", 0, geom.XY},             // wkb_hdr
		{"010700000002000000010100000000000000000010400000000000001840010200000002000000000000000000104000000000000018400000000000001c400000000000002440", 0, geom.XY},             // wkb_xdr
	}
	for _, dataset := range datasets {

		data, err := hex.DecodeString(dataset.data)

		r := bytes.NewReader(data)

		geo, sid, err := DecodeWKB(r)

		if err != nil {
			t.Error(err)
		}

		if geo == nil {
			t.Error("err")
		}

		var geobuf bytes.Buffer

		EncodeWKB(geo, &sid, &geobuf)

		str := string(geobuf.Bytes())

		if len(str) == 0 {
			t.Error("err")
		}
	}
}
