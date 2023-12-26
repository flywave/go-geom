package draw_test

import (
	"math/rand"
	"os"
	"testing"

	"github.com/flywave/go-geom/draw"
	"github.com/flywave/go-geom/general"
	"github.com/fogleman/gg"
)

func TestDraw(t *testing.T) {
	bt, _ := os.ReadFile("./已完成巷道.json")
	col, err := general.UnmarshalFeatureCollection(bt)
	if err != nil {
		return
	}
	opt := draw.DefaultDrawOptions()
	opt.Col = col

	gr, _ := draw.NewGeojsonRender(opt)
	gr.Render(9)
}

func TestLines(t *testing.T) {
	dc := gg.NewContext(1000, 1000)
	dc.SetRGB(0.5, 0.5, 0.5)
	dc.Clear()
	rnd := rand.New(rand.NewSource(99))
	for i := 0; i < 100; i++ {
		x1 := rnd.Float64() * 100
		y1 := rnd.Float64() * 100
		x2 := rnd.Float64() * 100
		y2 := rnd.Float64() * 100
		dc.DrawLine(x1, y1, x2, y2)
		dc.SetLineWidth(rnd.Float64() * 3)
		dc.SetRGB(rnd.Float64(), rnd.Float64(), rnd.Float64())
		dc.Stroke()
	}
	dc.SavePNG("TestLines.png")
}
