package geom

import (
	"fmt"
	"strconv"
	"strings"
)

func SridToUrn(srid int) string {
	if int(srid) == 4326 {
		return "urn:ogc:def:crs:OGC:1.3:CRS84"
	}

	return fmt.Sprintf("urn:ogc:def:crs:EPSG::%d", srid)
}

func UrnToSrid(urn string) int {
	if urn == "urn:ogc:def:crs:OGC:1.3:CRS84" {
		return 4326
	}

	if strings.HasPrefix(urn, "urn:ogc:def:crs:EPSG::") {
		estr := strings.TrimPrefix(urn, "urn:ogc:def:crs:EPSG::")
		epsg, err := strconv.Atoi(estr)

		if err == nil {
			return epsg
		}
	}

	return -1
}
