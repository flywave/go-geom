package geom

import "fmt"

func SridToUrn(srid int) string {
	if int(srid) == 4326 {
		return "urn:ogc:def:crs:OGC:1.3:CRS84"
	}

	return fmt.Sprintf("urn:ogc:def:crs:EPSG::%d", srid)
}
