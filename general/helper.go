package general

type Lesser interface {
	Len() int
	Less(i, j int) bool
}

func FindMinIdx(ln Lesser) (min int) {
	for i := 1; i < ln.Len(); i++ {
		if ln.Less(i, min) {
			min = i
		}
	}
	return min
}

func XYLessPoint(pt1, pt2 []float64) bool {
	if pt1[0] != pt2[0] {
		return pt1[0] < pt2[0]
	}
	return pt1[1] < pt2[1]
}

func FindMinPointIdx(ln [][]float64) (min int) {
	if len(ln) < 2 {
		return 0
	}
	for i := range ln[1:] {
		if XYLessPoint(ln[i+1], ln[min]) {
			min = i + 1
		}
	}
	return min
}

func RotateToIdx(idx int, ln [][]float64) {
	if len(ln) == 0 {
		return
	}
	tmp := make([][]float64, len(ln))
	copy(tmp, ln[idx:])
	copy(tmp[len(ln[idx:]):], ln)
	copy(ln, tmp)
}

func RotateToLeftMostPoint(ln [][]float64) {
	idx := FindMinPointIdx(ln)
	RotateToIdx(idx, ln)
}
