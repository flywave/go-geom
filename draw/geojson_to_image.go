package draw

import (
	"errors"
	"math"

	"github.com/flywave/go-geom"
	vec3d "github.com/flywave/go3d/float64/vec3"
	"github.com/fogleman/gg"
)

type DrawOptions struct {
	BackGroundColor [4]byte
	Color           [4]byte
	Format          string
	Scale           int
	Col             *geom.FeatureCollection
}

func DefaultDrawOptions() *DrawOptions {
	return &DrawOptions{
		BackGroundColor: [4]byte{255, 255, 255, 255},
		Color:           [4]byte{0, 0, 0, 255},
		Format:          "png",
	}
}

type GeojsonRender struct {
	opt       *DrawOptions
	localCol  *geom.FeatureCollection
	boundbox  *geom.BoundingBox
	localSize [3]float64
}

func NewGeojsonRender(opt *DrawOptions) (*GeojsonRender, error) {
	r := &GeojsonRender{opt: opt}
	if opt.Col == nil {
		return nil, errors.New("col is nil")
	}
	r.init()
	return r, nil
}

func (g *GeojsonRender) init() {
	if g.opt.Col.BoundingBox == nil {
		box := &geom.BoundingBox{vec3d.MaxVal, vec3d.MinVal}
		for _, c := range g.opt.Col.Features {
			if c.BoundingBox == nil {
				c.BoundingBox = geom.BoundingBoxFromGeometryData(&c.GeometryData)
			}
			box = geom.BoundingBoxsFromTwoBBox(box, c.BoundingBox)
		}
		g.opt.Col.BoundingBox = box
	}
	g.boundbox = g.opt.Col.BoundingBox
	g.localSize = [3]float64{
		g.boundbox[1][0] - g.boundbox[0][0],
		g.boundbox[1][1] - g.boundbox[0][1],
		g.boundbox[1][2] - g.boundbox[0][2],
	}

	fn := func(pt []float64) []float64 {
		return []float64{pt[0] - g.boundbox[0][0], pt[1] - g.boundbox[0][1], 0}
	}

	col2 := geom.NewFeatureCollection()
	for _, f := range g.opt.Col.Features {
		g := geom.ProcessGeometryData(&f.GeometryData, fn)
		f := geom.NewFeatureFromGeometryData(g)
		col2.Features = append(col2.Features, f)
	}
	g.localCol = col2
}

func processLine(g *geom.GeometryData, fn func([][]float64)) {
	switch g.Type {
	case "LineString":
		fn(g.LineString)
	case "MultiLineString":
		for _, l := range g.MultiLineString {
			fn(l)
		}
	case "Polygon":
		for _, l := range g.Polygon {
			fn(l)
		}
	case "MultiPolygon":
		for _, l := range g.MultiPolygon {
			for _, ll := range l {
				fn(ll)
			}
		}
	}
}

func (g *GeojsonRender) Render(resolution float64) {
	pixSize := [2]int{}
	pixSize[0] = int(math.Ceil(g.localSize[0] / resolution))
	pixSize[1] = int(math.Ceil(g.localSize[1] / resolution))

	dc := gg.NewContext(pixSize[0], pixSize[1])
	dc.SetRGBA255(int(g.opt.BackGroundColor[0]), int(g.opt.BackGroundColor[1]), int(g.opt.BackGroundColor[2]), int(g.opt.BackGroundColor[3]))
	dc.Clear()
	dc.InvertY()

	fn := func(l [][]float64) {
		for i := 0; i < len(l)-1; i++ {
			dc.DrawLine(l[i][0]/resolution, l[i][1]/resolution, l[i+1][0]/resolution, l[i+1][1]/resolution)
			dc.SetLineWidth(9)
		}
		dc.Stroke()
	}

	dc.SetRGBA255(int(g.opt.Color[0]), int(g.opt.Color[1]), int(g.opt.Color[2]), int(g.opt.Color[3]))
	for _, f := range g.localCol.Features {
		processLine(&f.GeometryData, fn)

	}
	dc.SavePNG("./out.png")
}
