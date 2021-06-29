package general

import (
	"math"

	"github.com/flywave/go-geom"
)

type Extenter interface {
	Extent() (extent [4]float64)
}

type MinMaxer interface {
	MinX() float64
	MinY() float64
	MaxX() float64
	MaxY() float64
}

type Extent [4]float64

/* ========================= ATTRIBUTES ========================= */

// Vertices return the vertices of the Bounding Box. The vertices are ordered in the following maner.
// (minx,miny), (maxx,miny), (maxx,maxy), (minx,maxy)
func (e *Extent) Vertices() [][]float64 {
	return [][]float64{
		{e.MinX(), e.MinY()},
		{e.MaxX(), e.MinY()},
		{e.MaxX(), e.MaxY()},
		{e.MinX(), e.MaxY()},
	}
}

// Verticies is the misspelled version of Vertices to match the interface
func (e *Extent) Verticies() [][]float64 { return e.Vertices() }

// ClockwiseFunc returns weather the set of points should be considered clockwise or counterclockwise. The last point is not the same as the first point, and the function should connect these points as needed.
type ClockwiseFunc func(...[]float64) bool

func (e *Extent) Edges(cwfn ClockwiseFunc) [][2][]float64 {
	v := e.Vertices()
	if cwfn != nil && !cwfn(v...) {
		v[0], v[1], v[2], v[3] = v[3], v[2], v[1], v[0]
	}
	return [][2][]float64{
		{v[0], v[1]},
		{v[1], v[2]},
		{v[2], v[3]},
		{v[3], v[0]},
	}
}

func (e *Extent) MaxX() float64 {
	if e == nil {
		return math.MaxFloat64
	}
	return e[2]
}

func (e *Extent) MinX() float64 {
	if e == nil {
		return -math.MaxFloat64
	}
	return e[0]
}

func (e *Extent) MaxY() float64 {
	if e == nil {
		return math.MaxFloat64
	}
	return e[3]
}

func (e *Extent) MinY() float64 {
	if e == nil {
		return -math.MaxFloat64
	}
	return e[1]
}

func (e *Extent) Min() [2]float64 {
	return [2]float64{e[0], e[1]}
}

func (e *Extent) Max() [2]float64 {
	return [2]float64{e[2], e[3]}
}

func (e *Extent) XSpan() float64 {
	if e == nil {
		return math.Inf(1)
	}
	return e[2] - e[0]
}

func (e *Extent) YSpan() float64 {
	if e == nil {
		return math.Inf(1)
	}
	return e[3] - e[1]
}

func (e *Extent) Extent() [4]float64 {
	return [4]float64{e.MinX(), e.MinY(), e.MaxX(), e.MaxY()}
}

func (e *Extent) Add(extent MinMaxer) {
	if e == nil {
		return
	}
	e[0] = math.Min(e[0], extent.MinX())
	e[2] = math.Max(e[2], extent.MaxX())
	e[1] = math.Min(e[1], extent.MinY())
	e[3] = math.Max(e[3], extent.MaxY())
}

func (e *Extent) AddPoints(points ...[]float64) {
	if e == nil {
		return
	}
	if len(points) == 0 {
		return
	}
	for _, pt := range points {
		if len(pt) < 2 {
			continue
		}
		e[0] = math.Min(pt[0], e[0])
		e[1] = math.Min(pt[1], e[1])
		e[2] = math.Max(pt[0], e[2])
		e[3] = math.Max(pt[1], e[3])
	}
}

func (e *Extent) AddPointers(pts ...geom.Point) {
	for i := range pts {
		e.AddPoints(pts[i].Data())
	}
}

func (e *Extent) AddGeometry(g geom.Geometry) error {
	return getExtent(g, e)
}

func (e *Extent) Area() float64 {
	return math.Abs((e.MaxY() - e.MinY()) * (e.MaxX() - e.MinX()))
}

func NewExtent(points ...[]float64) *Extent {
	var xy []float64
	if len(points) == 0 {
		return nil
	}

	extent := Extent{points[0][0], points[0][1], points[0][0], points[0][1]}
	if len(points) == 1 {
		return &extent
	}
	for i := 1; i < len(points); i++ {
		xy = points[i]
		switch {
		case xy[0] < extent[0]:
			extent[0] = xy[0]
		case xy[0] > extent[2]:
			extent[2] = xy[0]
		}
		switch {
		case xy[1] < extent[1]:
			extent[1] = xy[1]
		case xy[1] > extent[3]:
			extent[3] = xy[1]
		}
	}
	return &extent
}

func NewExtentFromPoints(points ...Point) *Extent {
	if len(points) == 0 {
		return nil
	}

	extent := Extent{points[0].X(), points[0].Y(), points[0].X(), points[0].Y()}
	if len(points) == 1 {
		return &extent
	}
	for _, pt := range points[1:] {
		switch {
		case pt.X() < extent[0]:
			extent[0] = pt.X()
		case pt.X() > extent[2]:
			extent[2] = pt.X()
		}
		switch {
		case pt.Y() < extent[1]:
			extent[1] = pt.Y()
		case pt.Y() > extent[3]:
			extent[3] = pt.Y()
		}
	}
	return &extent
}

func NewExtentFromGeometry(g geom.Geometry) (*Extent, error) {
	var pts []geom.Point
	if err := getCoordinates(g, &pts); err != nil {
		return nil, err
	}
	if len(pts) == 0 {
		return nil, nil
	}
	e := Extent{pts[0].X(), pts[0].Y(), pts[0].X(), pts[0].Y()}
	for _, pt := range pts {
		e.AddPoints(pt.Data())
	}

	return &e, nil
}

func (e *Extent) Contains(ne MinMaxer) bool {
	if e == nil {
		return true
	}
	if ne == nil {
		return false
	}
	return e.MinX() <= ne.MinX() &&
		e.MaxX() >= ne.MaxX() &&
		e.MinY() <= ne.MinY() &&
		e.MaxY() >= ne.MaxY()
}

func cmpFloat64(f1, f2, tolerance float64) bool {
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
	diff := math.Abs(f1 - f2)
	return diff <= tolerance
}

func floatLessOrEqual(pt1, pt2 float64) bool {
	if cmpFloat64(pt1, pt2, 0.001) {
		return true
	}
	return pt1 < pt2
}

func (e *Extent) ContainsPoint(pt []float64) bool {
	if e == nil {
		return true
	}

	return floatLessOrEqual(e.MinX(), pt[0]) && floatLessOrEqual(pt[0], e.MaxX()) &&
		floatLessOrEqual(e.MinY(), pt[1]) && floatLessOrEqual(pt[1], e.MaxY())

}

func (e *Extent) ContainsLine(l [2][]float64) bool {
	if e == nil {
		return true
	}
	return e.ContainsPoint(l[0]) && e.ContainsPoint(l[1])
}

func (e *Extent) ContainsGeom(g geom.Geometry) (bool, error) {
	if e.IsUniverse() {
		return true, nil
	}
	if extenter, ok := g.(MinMaxer); ok {
		return e.Contains(extenter), nil
	}
	var ne = new(Extent)
	if err := ne.AddGeometry(g); err != nil {
		return false, err
	}
	return e.Contains(ne), nil
}

func (e *Extent) ScaleBy(s float64) *Extent {
	if e == nil {
		return nil
	}
	return NewExtent(
		[]float64{e[0] * s, e[1] * s},
		[]float64{e[2] * s, e[3] * s},
	)
}

func (e *Extent) ExpandBy(s float64) *Extent {
	if e == nil {
		return nil
	}
	return NewExtent(
		[]float64{e[0] - s, e[1] - s},
		[]float64{e[2] + s, e[3] + s},
	)
}

func (e *Extent) Clone() *Extent {
	if e == nil {
		return nil
	}
	return &Extent{e[0], e[1], e[2], e[3]}
}

func (e *Extent) Intersect(ne *Extent) (*Extent, bool) {
	if e == nil {
		return ne.Clone(), true
	}
	if ne == nil {
		return e.Clone(), true
	}

	minx := e.MinX()
	if minx < ne.MinX() {
		minx = ne.MinX()
	}
	maxx := e.MaxX()
	if maxx > ne.MaxX() {
		maxx = ne.MaxX()
	}
	if minx >= maxx {
		return nil, false
	}
	miny := e.MinY()
	if miny < ne.MinY() {
		miny = ne.MinY()
	}
	maxy := e.MaxY()
	if maxy > ne.MaxY() {
		maxy = ne.MaxY()
	}

	if miny >= maxy {
		return nil, false
	}
	return &Extent{minx, miny, maxx, maxy}, true
}

func (e *Extent) IsUniverse() bool {
	return e == nil || (e.MinX() == -math.MaxFloat64 && e.MaxX() == math.MaxFloat64 &&
		e.MinY() == -math.MaxFloat64 && e.MaxY() == math.MaxFloat64)
}
