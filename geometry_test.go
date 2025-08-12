package geom

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试BoundingBox的UnmarshalJSON方法
func TestBoundingBoxUnmarshalJSON(t *testing.T) {
	// 测试一维数组格式 (4元素)
	bbox1 := &BoundingBox{}
	data1 := []byte(`[1, 2, 3, 4]`)
	err := json.Unmarshal(data1, bbox1)
	assert.NoError(t, err)
	assert.Equal(t, [3]float64{1, 2, 0}, bbox1[0])
	assert.Equal(t, [3]float64{3, 4, 0}, bbox1[1])

	// 测试一维数组格式 (6元素)
	bbox2 := &BoundingBox{}
	data2 := []byte(`[1, 2, 3, 4, 5, 6]`)
	err = json.Unmarshal(data2, bbox2)
	assert.NoError(t, err)
	assert.Equal(t, [3]float64{1, 2, 3}, bbox2[0])
	assert.Equal(t, [3]float64{4, 5, 6}, bbox2[1])

	// 测试二维数组格式
	bbox3 := &BoundingBox{}
	data3 := []byte(`[[1, 2, 3], [4, 5, 6]]`)
	err = json.Unmarshal(data3, bbox3)
	assert.NoError(t, err)
	assert.Equal(t, [3]float64{1, 2, 3}, bbox3[0])
	assert.Equal(t, [3]float64{4, 5, 6}, bbox3[1])

	// 测试错误格式
	bbox4 := &BoundingBox{}
	data4 := []byte(`[1, 2]`)
	err = json.Unmarshal(data4, bbox4)
	assert.Error(t, err)

	// 测试空数组格式
	bbox5 := &BoundingBox{}
	data5 := []byte(`[]`)
	err = json.Unmarshal(data5, bbox5)
	assert.Error(t, err)
}

// 测试GeometryAsString函数
func TestGeometryAsString(t *testing.T) {
	// 创建一个模拟的LineString
	mockLineString := &MockLineString{
		Points: [][]float64{
			{1, 2},
			{3, 4},
		},
	}

	result := GeometryAsString(mockLineString)
	expected := "[ ( 1 2 ) ( 3 4 )]"
	assert.Contains(t, result, expected)

	// 测试其他类型
	mockPoint := &MockPoint{Point: []float64{5, 6}}
	result = GeometryAsString(mockPoint)
	assert.Contains(t, result, "Point")
}

// 测试GeometryAsMap函数
func TestGeometryAsMap(t *testing.T) {
	// 测试Point类型
	mockPoint := &MockPoint{Point: []float64{1, 2}}
	result := GeometryAsMap(mockPoint)
	assert.Equal(t, "Point", result["type"])
	assert.Equal(t, []float64{1, 2}, result["value"])

	// 测试Point3类型
	mockPoint3 := &MockPoint3{Point: []float64{1, 2, 3}}
	result = GeometryAsMap(mockPoint3)
	assert.Equal(t, "Point3", result["type"])
	assert.Equal(t, []float64{1, 2, 3}, result["value"])
}

// 测试Round函数
func TestRound(t *testing.T) {
	// 测试基本四舍五入
	assert.Equal(t, 2.0, Round(1.5, 0.5, 0))

	// 测试保留小数位
	assert.Equal(t, 1.55, Round(1.549, 0.5, 2))
	assert.Equal(t, 1.54, Round(1.544, 0.5, 2))

	// 测试负数
	assert.Equal(t, -2.0, Round(-1.5, 0.5, 0))
}

// 测试IsPointEqual函数
func TestIsPointEqual(t *testing.T) {
	// 测试相等的点
	p1 := &MockPoint{Point: []float64{1, 2}}
	p2 := &MockPoint{Point: []float64{1, 2}}
	assert.True(t, IsPointEqual(p1, p2))

	// 测试不相等的点
	p3 := &MockPoint{Point: []float64{1, 3}}
	assert.False(t, IsPointEqual(p1, p3))

	// 测试容差范围内的点
	p4 := &MockPoint{Point: []float64{1.0000001, 2.0000001}}
	assert.True(t, IsPointEqual(p1, p4))

	// 测试nil值
	assert.True(t, IsPointEqual(nil, nil))
	assert.False(t, IsPointEqual(p1, nil))
}

// 测试IsGeometryEmpty函数
func TestIsGeometryEmpty(t *testing.T) {
	// 测试空几何对象
	emptyPoint := &MockPoint{Point: []float64{}}
	assert.True(t, IsGeometryEmpty(emptyPoint))

	// 测试非空几何对象
	nonEmptyPoint := &MockPoint{Point: []float64{1, 2}}
	assert.False(t, IsGeometryEmpty(nonEmptyPoint))
}

// Mock实现

// MockPoint 是一个模拟Point接口的结构体，用于测试
type MockPoint struct {
	Point []float64
}

// Data 返回点的坐标数据，实现Point接口
func (m MockPoint) Data() []float64 {
	return m.Point
}

// GetType 返回几何类型，实现Geometry接口
func (m MockPoint) GetType() string {
	return "Point"
}

// X 返回点的X坐标，实现Point接口
func (m MockPoint) X() float64 {
	if len(m.Point) > 0 {
		return m.Point[0]
	}
	return 0
}

// Y 返回点的Y坐标，实现Point接口
func (m MockPoint) Y() float64 {
	if len(m.Point) > 1 {
		return m.Point[1]
	}
	return 0
}

// GeometryData 返回几何数据，实现Geometry接口
func (m MockPoint) GeometryData() *GeometryData {
	return NewPointGeometryData(m.Point)
}

// MockPoint3 模拟Point3接口
type MockPoint3 struct {
	Point []float64
}

func (m *MockPoint3) GetType() string {
	return "Point3"
}

func (m *MockPoint3) X() float64 {
	if len(m.Point) > 0 {
		return m.Point[0]
	}
	return 0
}

func (m *MockPoint3) Y() float64 {
	if len(m.Point) > 1 {
		return m.Point[1]
	}
	return 0
}

func (m *MockPoint3) Z() float64 {
	if len(m.Point) > 2 {
		return m.Point[2]
	}
	return 0
}

func (m *MockPoint3) Data() []float64 {
	return m.Point
}

// MockLineString 模拟LineString接口
type MockLineString struct {
	Points [][]float64
}

func (m *MockLineString) GetType() string {
	return "LineString"
}

func (m *MockLineString) Subpoints() []Point {
	points := make([]Point, len(m.Points))
	for i, p := range m.Points {
		points[i] = &MockPoint{Point: p}
	}
	return points
}

func (m *MockLineString) Data() [][]float64 {
	return m.Points
}
