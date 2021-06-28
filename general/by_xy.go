package general

type ByXY [][]float64

func (xy ByXY) Less(i, j int) bool { return XYLessPoint(xy[i], xy[j]) }
func (xy ByXY) Swap(i, j int)      { xy[i], xy[j] = xy[j], xy[i] }
func (xy ByXY) Len() int           { return len(xy) }

type PointByXY []Point

func (xy PointByXY) Less(i, j int) bool {
	return XYLessPoint(xy[i].Data(), xy[j].Data())
}
func (xy PointByXY) Swap(i, j int) { xy[i], xy[j] = xy[j], xy[i] }
func (xy PointByXY) Len() int      { return len(xy) }

type bySubRingSizeXY [][][]float64

func (xy bySubRingSizeXY) Less(i, j int) bool {
	switch {
	case i == 0:
		return true
	case j == 0:
		return false
	case len(xy[i]) != len(xy[j]):
		return len(xy[i]) < len(xy[j])
	default:
		mi, mj := FindMinPointIdx(xy[i]), FindMinPointIdx(xy[j])
		return XYLessPoint(xy[i][mi], xy[j][mj])
	}
}

func (xy bySubRingSizeXY) Len() int      { return len(xy) }
func (xy bySubRingSizeXY) Swap(i, j int) { xy[i], xy[j] = xy[j], xy[i] }

type byPolygonMainSizeXY [][][][]float64

func (xy byPolygonMainSizeXY) Less(i, j int) bool {
	if len(xy[i]) == 0 {
		return true
	}
	if len(xy[j]) == 0 {
		return false
	}
	if len(xy[i][0]) != len(xy[j][0]) {
		return len(xy[i][0]) < len(xy[j][0])
	}
	mi, mj := FindMinPointIdx(xy[i][0]), FindMinPointIdx(xy[j][0])
	return XYLessPoint(xy[i][0][mi], xy[j][0][mj])
}
func (xy byPolygonMainSizeXY) Len() int      { return len(xy) }
func (xy byPolygonMainSizeXY) Swap(i, j int) { xy[i], xy[j] = xy[j], xy[i] }
