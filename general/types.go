package general

import "github.com/flywave/go-geom"

type Coordinate []float64

type Point struct {
	geom.Point
	point Coordinate
	srid  int
}

func (c *Point) SetSRID(id int) {
	c.srid = id
}

func (c Point) SRID() int {
	return c.srid
}

func (c Point) Empty() bool {
	return len(c.point) == 0
}

func (p *Point) GetType() string {
	return "Point"
}

func (p *Point) X() float64 {
	return p.point[0]
}

func (p *Point) Y() float64 {
	return p.point[1]
}

func (p *Point) Data() []float64 {
	return p.point
}

func NewPoint(pt []float64) geom.Point {
	return &Point{point: pt}
}

type Point3 struct {
	geom.Point3
	point Coordinate
	srid  int
}

func (c *Point3) SetSRID(id int) {
	c.srid = id
}

func (p *Point3) X() float64 {
	return p.point[0]
}

func (p *Point3) Y() float64 {
	return p.point[1]
}

func (p *Point3) Z() float64 {
	return p.point[2]
}

func (p *Point3) Data() []float64 {
	return p.point
}

func (p *Point3) GetType() string {
	return "Point"
}

func (c Point3) SRID() int {
	return c.srid
}

func (c Point3) Empty() bool {
	return len(c.point) == 0
}

func NewPoint3(pt []float64) geom.Point3 {
	return &Point3{point: pt}
}

type Polygon struct {
	rings []geom.LineString
	srid  int
}

func NewPolygon(pts [][][]float64) *Polygon {
	rets := make([]geom.LineString, len(pts))
	for i := range pts {
		rets[i] = NewLineString(pts[i])
	}
	return &Polygon{rings: rets}
}

func (c *Polygon) SetSRID(id int) {
	c.srid = id
}

func (p *Polygon) Sublines() []geom.LineString {
	return p.rings
}

func (p *Polygon) Data() [][][]float64 {
	ret := make([][][]float64, len(p.rings))
	for i, p := range p.rings {
		ret[i] = p.Data()
	}
	return ret
}

func (p *Polygon) GetType() string {
	return "Polygon"
}

func (c Polygon) SRID() int {
	return c.srid
}

func (c Polygon) Empty() bool {
	return len(c.rings) == 0
}

type Polygon3 struct {
	rings []geom.LineString3
	srid  int
}

func NewPolygon3(pts [][][]float64) *Polygon3 {
	rets := make([]geom.LineString3, len(pts))
	for i := range pts {
		rets[i] = NewLineString3(pts[i])
	}
	return &Polygon3{rings: rets}
}

func (c *Polygon3) SetSRID(id int) {
	c.srid = id
}

func (p *Polygon3) Sublines() []geom.LineString3 {
	return p.rings
}

func (p *Polygon3) Data() [][][]float64 {
	ret := make([][][]float64, len(p.rings))
	for i, p := range p.rings {
		ret[i] = p.Data()
	}
	return ret
}

func (p *Polygon3) GetType() string {
	return "Polygon"
}

func (c Polygon3) SRID() int {
	return c.srid
}

func (c Polygon3) Empty() bool {
	return len(c.rings) == 0
}

type MultiPolygon struct {
	polygons []geom.Polygon
	srid     int
}

func NewMultiPolygon(pols [][][][]float64) *MultiPolygon {
	rets := make([]geom.Polygon, len(pols))
	for i := range pols {
		rets[i] = NewPolygon(pols[i])
	}
	return &MultiPolygon{polygons: rets}
}

func (c *MultiPolygon) SetSRID(id int) {
	c.srid = id
}

func (p *MultiPolygon) Polygons() []geom.Polygon {
	return p.polygons
}

func (p *MultiPolygon) Data() [][][][]float64 {
	ret := make([][][][]float64, len(p.polygons))
	for i, p := range p.polygons {
		ret[i] = p.Data()
	}
	return ret
}

func (p *MultiPolygon) GetType() string {
	return "MultiPolygon"
}

func (c MultiPolygon) SRID() int {
	return c.srid
}

func (c MultiPolygon) Empty() bool {
	return len(c.polygons) == 0
}

type MultiPolygon3 struct {
	polygons []geom.Polygon3
	srid     int
}

func NewMultiPolygon3(pols [][][][]float64) *MultiPolygon3 {
	rets := make([]geom.Polygon3, len(pols))
	for i := range pols {
		rets[i] = NewPolygon3(pols[i])
	}
	return &MultiPolygon3{polygons: rets}
}

func (c *MultiPolygon3) SetSRID(id int) {
	c.srid = id
}

func (p *MultiPolygon3) Polygons() []geom.Polygon3 {
	return p.polygons
}

func (p *MultiPolygon3) Data() [][][][]float64 {
	ret := make([][][][]float64, len(p.polygons))
	for i, p := range p.polygons {
		ret[i] = p.Data()
	}
	return ret
}

func (p *MultiPolygon3) GetType() string {
	return "MultiPolygon"
}

func (c MultiPolygon3) SRID() int {
	return c.srid
}

func (c MultiPolygon3) Empty() bool {
	return len(c.polygons) == 0
}

type MultiPoint struct {
	points []geom.Point
	srid   int
}

func NewMultiPoint(pts [][]float64) *MultiPoint {
	rets := make([]geom.Point, len(pts))
	for i := range pts {
		rets[i] = NewPoint(pts[i])
	}
	return &MultiPoint{points: rets}
}

func (c *MultiPoint) SetSRID(id int) {
	c.srid = id
}

func (mp *MultiPoint) Points() []geom.Point {
	return mp.points
}

func (mp *MultiPoint) Data() [][]float64 {
	ret := make([][]float64, len(mp.points))
	for i, p := range mp.points {
		ret[i] = p.Data()
	}
	return ret
}

func (p *MultiPoint) GetType() string {
	return "MultiPoint"
}

func (c MultiPoint) SRID() int {
	return c.srid
}

func (c MultiPoint) Empty() bool {
	return len(c.points) == 0
}

type MultiPoint3 struct {
	points []geom.Point3
	srid   int
}

func NewMultiPoint3(pts [][]float64) *MultiPoint3 {
	rets := make([]geom.Point3, len(pts))
	for i := range pts {
		rets[i] = NewPoint3(pts[i])
	}
	return &MultiPoint3{points: rets}
}

func (c *MultiPoint3) SetSRID(id int) {
	c.srid = id
}

func (mp *MultiPoint3) Points() []geom.Point3 {
	return mp.points
}

func (mp *MultiPoint3) Data() [][]float64 {
	ret := make([][]float64, len(mp.points))
	for i, p := range mp.points {
		ret[i] = p.Data()
	}
	return ret
}

func (p *MultiPoint3) GetType() string {
	return "MultiPoint3"
}

func (c MultiPoint3) SRID() int {
	return c.srid
}

func (c MultiPoint3) Empty() bool {
	return len(c.points) == 0
}

type MultiLineString struct {
	lines []geom.LineString
	srid  int
}

func NewMultiLineString(pts [][][]float64) *MultiLineString {
	rets := make([]geom.LineString, len(pts))
	for i := range pts {
		rets[i] = NewLineString(pts[i])
	}
	return &MultiLineString{lines: rets}
}

func (c *MultiLineString) SetSRID(id int) {
	c.srid = id
}

func (ml *MultiLineString) Lines() []geom.LineString {
	return ml.lines
}

func (ml *MultiLineString) Data() [][][]float64 {
	ret := make([][][]float64, len(ml.lines))
	for i, p := range ml.lines {
		ret[i] = p.Data()
	}
	return ret
}

func (p *MultiLineString) GetType() string {
	return "MultiLineString"
}

func (c *MultiLineString) SRID() int {
	return c.srid
}

func (c *MultiLineString) Empty() bool {
	return len(c.lines) == 0
}

type MultiLineString3 struct {
	lines []geom.LineString3
	srid  int
}

func NewMultiLineString3(pts [][][]float64) *MultiLineString3 {
	rets := make([]geom.LineString3, len(pts))
	for i := range pts {
		rets[i] = NewLineString3(pts[i])
	}
	return &MultiLineString3{lines: rets}
}

func (c *MultiLineString3) SetSRID(id int) {
	c.srid = id
}

func (ml *MultiLineString3) Lines() []geom.LineString3 {
	return ml.lines
}

func (ml *MultiLineString3) Data() [][][]float64 {
	ret := make([][][]float64, len(ml.lines))
	for i, p := range ml.lines {
		ret[i] = p.Data()
	}
	return ret
}

func (p *MultiLineString3) GetType() string {
	return "MultiLineString"
}

func (c *MultiLineString3) SRID() int {
	return c.srid
}

func (c *MultiLineString3) Empty() bool {
	return len(c.lines) == 0
}

type LineString struct {
	points []geom.Point
	srid   int
}

func NewLineString(ls [][]float64) *LineString {
	rets := make([]geom.Point, len(ls))
	for i := range ls {
		rets[i] = NewPoint(ls[i])
	}
	return &LineString{points: rets}
}

func (c *LineString) SetSRID(id int) {
	c.srid = id
}

func (p *LineString) Subpoints() []geom.Point {
	return p.points
}

func (p *LineString) Data() [][]float64 {
	ret := make([][]float64, len(p.points))
	for i, p := range p.points {
		ret[i] = p.Data()
	}
	return ret
}

func (p *LineString) GetType() string {
	return "LineString"
}

func (c *LineString) SRID() int {
	return c.srid
}

func (c *LineString) Empty() bool {
	return len(c.points) == 0
}

type LineString3 struct {
	points []geom.Point3
	srid   int
}

func NewLineString3(ls [][]float64) *LineString3 {
	rets := make([]geom.Point3, len(ls))
	for i := range ls {
		rets[i] = NewPoint3(ls[i])
	}
	return &LineString3{points: rets}
}

func (c *LineString3) SetSRID(id int) {
	c.srid = id
}

func (p *LineString3) Subpoints() []geom.Point3 {
	return p.points
}

func (p *LineString3) Data() [][]float64 {
	ret := make([][]float64, len(p.points))
	for i, p := range p.points {
		ret[i] = p.Data()
	}
	return ret
}

func (p *LineString3) GetType() string {
	return "LineString"
}

func (c *LineString3) SRID() int {
	return c.srid
}

func (c *LineString3) Empty() bool {
	return len(c.points) == 0
}

func NewGeometryCollection(geoms ...geom.Geometry) geom.Collection {
	return geoms
}
