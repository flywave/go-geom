package geom

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试NewGeometryData函数
func TestNewGeometryData(t *testing.T) {
	// 测试Point类型
	pointData := NewPointGeometryData([]float64{1.0, 2.0})
	assert.Equal(t, GeometryPoint, pointData.Type)
	assert.Equal(t, []float64{1.0, 2.0}, pointData.Point)

	// 测试MultiPoint类型
	multiPointData := NewMultiPointGeometryData([]float64{1.0, 2.0}, []float64{3.0, 4.0})
	assert.Equal(t, GeometryMultiPoint, multiPointData.Type)
	assert.Equal(t, [][]float64{{1.0, 2.0}, {3.0, 4.0}}, multiPointData.MultiPoint)
}

// 测试GeometryData的MarshalJSON和UnmarshalJSON方法
func TestGeometryDataMarshalUnmarshal(t *testing.T) {
	// 创建一个Point类型的GeometryData
	pointData := &GeometryData{
		Type:  GeometryPoint,
		Point: []float64{1.0, 2.0},
	}

	// 测试MarshalJSON
	jsonData, err := json.Marshal(pointData)
	assert.NoError(t, err)

	// 测试UnmarshalJSON
	unmarshaledData := &GeometryData{}
	err = json.Unmarshal(jsonData, unmarshaledData)
	assert.NoError(t, err)

	// 验证结果
	assert.Equal(t, pointData.Type, unmarshaledData.Type)
	assert.Equal(t, pointData.Point, unmarshaledData.Point)

	// 测试MultiPoint类型
	multiPointData := &GeometryData{
		Type:       GeometryMultiPoint,
		MultiPoint: [][]float64{{1.0, 2.0}, {3.0, 4.0}},
	}

	jsonData, err = json.Marshal(multiPointData)
	assert.NoError(t, err)

	unmarshaledData = &GeometryData{}
	err = json.Unmarshal(jsonData, unmarshaledData)
	assert.NoError(t, err)

	assert.Equal(t, multiPointData.Type, unmarshaledData.Type)
	assert.Equal(t, multiPointData.MultiPoint, unmarshaledData.MultiPoint)
}

// 测试各种GeometryData构造函数
func TestGeometryDataConstructors(t *testing.T) {
	// 测试NewPointGeometryData
	pointData := NewPointGeometryData([]float64{1.0, 2.0})
	assert.Equal(t, GeometryPoint, pointData.Type)
	assert.Equal(t, []float64{1.0, 2.0}, pointData.Point)

	// 测试NewMultiPointGeometryData
	multiPointData := NewMultiPointGeometryData([][]float64{{1.0, 2.0}, {3.0, 4.0}}...)
	assert.Equal(t, GeometryMultiPoint, multiPointData.Type)
	assert.Equal(t, [][]float64{{1.0, 2.0}, {3.0, 4.0}}, multiPointData.MultiPoint)

	// 测试NewLineStringGeometryData
	lineStringData := NewLineStringGeometryData([][]float64{{1.0, 2.0}, {3.0, 4.0}, {5.0, 6.0}})
	assert.Equal(t, GeometryLineString, lineStringData.Type)
	assert.Equal(t, [][]float64{{1.0, 2.0}, {3.0, 4.0}, {5.0, 6.0}}, lineStringData.LineString)

	// 测试NewMultiLineStringGeometryData
	multiLineStringData := NewMultiLineStringGeometryData([][]float64{{1.0, 2.0}, {3.0, 4.0}})
	assert.Equal(t, GeometryMultiLineString, multiLineStringData.Type)
	assert.Equal(t, [][][]float64{{{1.0, 2.0}, {3.0, 4.0}}}, multiLineStringData.MultiLineString)
}

// 测试DecodeGeometry函数
func TestDecodeGeometry(t *testing.T) {
	// 测试Point类型
	pointObj := map[string]interface{}{
		"type":        "Point",
		"coordinates": []interface{}{1.0, 2.0},
	}
	pointData := &GeometryData{}
	err := DecodeGeometry(pointData, pointObj)
	assert.NoError(t, err)
	assert.Equal(t, GeometryPoint, pointData.Type)
	assert.Equal(t, []float64{1.0, 2.0}, pointData.Point)

	// 测试MultiPoint类型
	multiPointObj := map[string]interface{}{
		"type": "MultiPoint",
		"coordinates": []interface{}{
			[]interface{}{1.0, 2.0},
			[]interface{}{3.0, 4.0},
		},
	}
	multiPointData := &GeometryData{}
	err = DecodeGeometry(multiPointData, multiPointObj)
	assert.NoError(t, err)
	assert.Equal(t, GeometryMultiPoint, multiPointData.Type)
	assert.Equal(t, [][]float64{{1.0, 2.0}, {3.0, 4.0}}, multiPointData.MultiPoint)
}

// 测试类型判断方法
func TestGeometryDataTypeMethods(t *testing.T) {
	// 测试Point类型判断
	pointData := NewPointGeometryData([]float64{1.0, 2.0})
	assert.True(t, pointData.IsPoint())
	assert.False(t, pointData.IsMultiPoint())
	assert.False(t, pointData.IsLineString())
	assert.False(t, pointData.IsMultiLineString())
	assert.False(t, pointData.IsPolygon())
	assert.False(t, pointData.IsMultiPolygon())
	assert.False(t, pointData.IsCollection())

	// 测试MultiPoint类型判断
	multiPointData := NewMultiPointGeometryData([][]float64{{1.0, 2.0}, {3.0, 4.0}}...)
	assert.False(t, multiPointData.IsPoint())
	assert.True(t, multiPointData.IsMultiPoint())
	assert.False(t, multiPointData.IsLineString())
	assert.False(t, multiPointData.IsMultiLineString())
	assert.False(t, multiPointData.IsPolygon())
	assert.False(t, multiPointData.IsMultiPolygon())
	assert.False(t, multiPointData.IsCollection())
}

// MockMultiPoint 是一个模拟MultiPoint接口的结构体，用于测试
type MockMultiPoint struct {
	Points [][]float64
}

// Data 返回点集数据，实现MultiPoint接口
func (m MockMultiPoint) Data() [][]float64 {
	return m.Points
}

// GetType 返回几何类型，实现Geometry接口
func (m MockMultiPoint) GetType() string {
	return "MultiPoint"
}

// GeometryData 返回几何数据，实现Geometry接口
func (m MockMultiPoint) GeometryData() *GeometryData {
	return NewMultiPointGeometryData(m.Points...)
}

// 测试UnmarshalGeometry函数
func TestUnmarshalGeometry(t *testing.T) {
	// 测试Point类型
	pointJSON := `{"type":"Point","coordinates":[1.0,2.0]}`
	geometryData, err := UnmarshalGeometry([]byte(pointJSON))
	assert.NoError(t, err)
	assert.Equal(t, GeometryPoint, geometryData.Type)
	assert.Equal(t, []float64{1.0, 2.0}, geometryData.Point)

	// 测试MultiPoint类型
	multiPointJSON := `{"type":"MultiPoint","coordinates":[[1.0,2.0],[3.0,4.0]]}`
	geometryData, err = UnmarshalGeometry([]byte(multiPointJSON))
	assert.NoError(t, err)
	assert.Equal(t, GeometryMultiPoint, geometryData.Type)
	assert.Equal(t, [][]float64{{1.0, 2.0}, {3.0, 4.0}}, geometryData.MultiPoint)
}

// 测试IsEmpty方法
func TestGeometryDataIsEmpty(t *testing.T) {
	// 测试空GeometryData
	emptyData := &GeometryData{}
	assert.True(t, emptyData.IsEmpty())

	// 测试非空GeometryData
	pointData := NewPointGeometryData([]float64{1.0, 2.0})
	assert.False(t, pointData.IsEmpty())
}

// 测试GeometryData的Scan方法
func TestGeometryDataScan(t *testing.T) {
	// 测试从字符串扫描
	pointData := &GeometryData{}
	err := pointData.Scan(`{"type":"Point","coordinates":[1.0,2.0]}`)
	assert.NoError(t, err)
	assert.Equal(t, GeometryPoint, pointData.Type)
	assert.Equal(t, []float64{1.0, 2.0}, pointData.Point)

	// 测试从字节数组扫描
	multiPointData := &GeometryData{}
	err = multiPointData.Scan([]byte(`{"type":"MultiPoint","coordinates":[[1.0,2.0],[3.0,4.0]]}`))
	assert.NoError(t, err)
	assert.Equal(t, GeometryMultiPoint, multiPointData.Type)
	assert.Equal(t, [][]float64{{1.0, 2.0}, {3.0, 4.0}}, multiPointData.MultiPoint)

	// 测试不支持的类型
	err = pointData.Scan(123)
	assert.Error(t, err)
}

// 测试Polygon和MultiPolygon类型
func TestGeometryDataPolygonTypes(t *testing.T) {
	// 测试Polygon类型
	polygonCoords := [][][]float64{
		{{0.0, 0.0}, {0.0, 10.0}, {10.0, 10.0}, {10.0, 0.0}, {0.0, 0.0}},
	}
	polygonData := NewPolygonGeometryData(polygonCoords)
	assert.Equal(t, GeometryPolygon, polygonData.Type)
	assert.Equal(t, polygonCoords, polygonData.Polygon)
	assert.True(t, polygonData.IsPolygon())

	// 测试MultiPolygon类型
	multiPolygonCoords := [][][][]float64{
		{
			{{0.0, 0.0}, {0.0, 10.0}, {10.0, 10.0}, {10.0, 0.0}, {0.0, 0.0}},
		},
		{
			{{10.0, 10.0}, {10.0, 20.0}, {20.0, 20.0}, {20.0, 10.0}, {10.0, 10.0}},
		},
	}
	multiPolygonData := NewMultiPolygonGeometryData(multiPolygonCoords...)
	assert.Equal(t, GeometryMultiPolygon, multiPolygonData.Type)
	assert.Equal(t, multiPolygonCoords, multiPolygonData.MultiPolygon)
	assert.True(t, multiPolygonData.IsMultiPolygon())
}

// 测试Collection类型
func TestGeometryDataCollection(t *testing.T) {
	// 创建不同类型的GeometryData
	pointData := NewPointGeometryData([]float64{1.0, 2.0})
	multiPointData := NewMultiPointGeometryData([][]float64{{3.0, 4.0}, {5.0, 6.0}}...)

	// 创建Collection
	collectionData := NewCollectionGeometryData(pointData, multiPointData)
	assert.Equal(t, GeometryCollection, collectionData.Type)
	assert.Equal(t, 2, len(collectionData.Geometries))
	assert.Equal(t, pointData, collectionData.Geometries[0])
	assert.Equal(t, multiPointData, collectionData.Geometries[1])
	assert.True(t, collectionData.IsCollection())

	// 测试Collection的JSON序列化和反序列化
	jsonData, err := json.Marshal(collectionData)
	assert.NoError(t, err)

	unmarshaledData := &GeometryData{}
	err = json.Unmarshal(jsonData, unmarshaledData)
	assert.NoError(t, err)

	assert.Equal(t, GeometryCollection, unmarshaledData.Type)
	assert.Equal(t, 2, len(unmarshaledData.Geometries))
	assert.Equal(t, pointData.Type, unmarshaledData.Geometries[0].Type)
	assert.Equal(t, multiPointData.Type, unmarshaledData.Geometries[1].Type)
}

// 测试带有BoundingBox的GeometryData
func TestGeometryDataWithBoundingBox(t *testing.T) {
	// 创建带BoundingBox的Point
	pointData := NewPointGeometryData([]float64{1.0, 2.0})
	bbox := &BoundingBox{
		{1.0, 2.0, 0.0},
		{1.0, 2.0, 0.0},
	}
	pointData.BoundingBox = bbox

	// 测试JSON序列化和反序列化
	jsonData, err := json.Marshal(pointData)
	assert.NoError(t, err)

	unmarshaledData := &GeometryData{}
	err = json.Unmarshal(jsonData, unmarshaledData)
	assert.NoError(t, err)

	assert.Equal(t, pointData.Type, unmarshaledData.Type)
	assert.Equal(t, pointData.Point, unmarshaledData.Point)
	assert.Equal(t, pointData.BoundingBox, unmarshaledData.BoundingBox)
}

// 测试带有EPSG属性的GeometryData
func TestGeometryDataWithEPSG(t *testing.T) {
	// 创建带EPSG的Point
	pointData := NewPointGeometryData([]float64{1.0, 2.0})
	pointData.EPSG = 4326

	// 测试JSON序列化和反序列化
	jsonData, err := json.Marshal(pointData)
	assert.NoError(t, err)

	unmarshaledData := &GeometryData{}
	err = json.Unmarshal(jsonData, unmarshaledData)
	assert.NoError(t, err)

	assert.Equal(t, pointData.Type, unmarshaledData.Type)
	assert.Equal(t, pointData.Point, unmarshaledData.Point)
	assert.Equal(t, pointData.EPSG, unmarshaledData.EPSG)
}

// 测试错误处理
func TestGeometryDataErrorHandling(t *testing.T) {
	// 测试无效的几何类型
	invalidObj := map[string]interface{}{
		"type": "InvalidType",
	}
	invalidData := &GeometryData{}
	err := DecodeGeometry(invalidData, invalidObj)
	assert.NoError(t, err) // DecodeGeometry目前对未知类型不返回错误
	assert.Equal(t, GeometryType("InvalidType"), invalidData.Type)

	// 测试无效的坐标格式
	invalidCoordsObj := map[string]interface{}{
		"type":        "Point",
		"coordinates": "invalid",
	}
	invalidCoordsData := &GeometryData{}
	err = DecodeGeometry(invalidCoordsData, invalidCoordsObj)
	assert.Error(t, err)

	// 测试无效的JSON
	_, err = UnmarshalGeometry([]byte(`invalid json`))
	assert.Error(t, err)
}
