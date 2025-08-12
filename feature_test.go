package geom

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试NewFeature函数
func TestNewFeature(t *testing.T) {
	// 创建一个测试用的Geometry
	point := MockPoint{Point: []float64{1.0, 2.0}}

	// 测试NewFeature
	feature := NewFeature(point)

	// 验证基本字段
	assert.Equal(t, "Feature", feature.Type)
	assert.Equal(t, point, feature.Geometry)
	assert.NotNil(t, feature.GeometryData)
	assert.Equal(t, point.Data(), feature.GeometryData.Point)
	assert.NotNil(t, feature.BoundingBox)
	assert.NotNil(t, feature.Properties)
	assert.NotNil(t, feature.ExtData)
}

// 测试NewPointFeature函数
func TestNewPointFeature(t *testing.T) {
	coordinate := []float64{1.0, 2.0}
	feature := NewPointFeature(coordinate)

	assert.Equal(t, "Feature", feature.Type)
	assert.Equal(t, GeometryPoint, feature.GeometryData.Type)
	assert.Equal(t, coordinate, feature.GeometryData.Point)
	assert.NotNil(t, feature.BoundingBox)
	assert.Equal(t, coordinate[0], feature.BoundingBox[0][0])
	assert.Equal(t, coordinate[1], feature.BoundingBox[0][1])
	assert.Equal(t, coordinate[0], feature.BoundingBox[1][0])
	assert.Equal(t, coordinate[1], feature.BoundingBox[1][1])
}

// 测试NewMultiPointFeature函数
func TestNewMultiPointFeature(t *testing.T) {
	coordinates := [][]float64{{1.0, 2.0}, {3.0, 4.0}}
	feature := NewMultiPointFeature(coordinates...)

	assert.Equal(t, "Feature", feature.Type)
	assert.Equal(t, GeometryMultiPoint, feature.GeometryData.Type)
	assert.Equal(t, coordinates, feature.GeometryData.MultiPoint)
	assert.NotNil(t, feature.BoundingBox)
	assert.Equal(t, 1.0, feature.BoundingBox[0][0])
	assert.Equal(t, 2.0, feature.BoundingBox[0][1])
	assert.Equal(t, 3.0, feature.BoundingBox[1][0])
	assert.Equal(t, 4.0, feature.BoundingBox[1][1])
}

// 测试NewLineStringFeature函数
func TestNewLineStringFeature(t *testing.T) {
	coordinates := [][]float64{{1.0, 2.0}, {3.0, 4.0}, {5.0, 6.0}}
	feature := NewLineStringFeature(coordinates)

	assert.Equal(t, "Feature", feature.Type)
	assert.Equal(t, GeometryLineString, feature.GeometryData.Type)
	assert.Equal(t, coordinates, feature.GeometryData.LineString)
	assert.NotNil(t, feature.BoundingBox)
	assert.Equal(t, 1.0, feature.BoundingBox[0][0])
	assert.Equal(t, 2.0, feature.BoundingBox[0][1])
	assert.Equal(t, 5.0, feature.BoundingBox[1][0])
	assert.Equal(t, 6.0, feature.BoundingBox[1][1])
}

// 测试属性设置和获取方法
func TestFeatureProperties(t *testing.T) {
	feature := NewPointFeature([]float64{1.0, 2.0})

	// 测试SetProperty
	feature.SetProperty("name", "test")
	feature.SetProperty("value", 123)
	feature.SetProperty("flag", true)
	feature.SetProperty("score", 45.67)

	// 测试PropertyString
	name, err := feature.PropertyString("name")
	assert.NoError(t, err)
	assert.Equal(t, "test", name)

	// 测试PropertyInt
	value, err := feature.PropertyInt("value")
	assert.NoError(t, err)
	assert.Equal(t, 123, value)

	// 测试PropertyBool
	flag, err := feature.PropertyBool("flag")
	assert.NoError(t, err)
	assert.Equal(t, true, flag)

	// 测试PropertyFloat64
	score, err := feature.PropertyFloat64("score")
	assert.NoError(t, err)
	assert.Equal(t, 45.67, score)

	// 测试不存在的属性
	_, err = feature.PropertyString("notexist")
	assert.Error(t, err)

	// 测试PropertyMust系列方法
	assert.Equal(t, "test", feature.PropertyMustString("name"))
	assert.Equal(t, 123, feature.PropertyMustInt("value"))
	assert.Equal(t, true, feature.PropertyMustBool("flag"))
	assert.Equal(t, 45.67, feature.PropertyMustFloat64("score"))

	// 测试默认值
	assert.Equal(t, "default", feature.PropertyMustString("notexist", "default"))
	assert.Equal(t, 456, feature.PropertyMustInt("notexist", 456))
	assert.Equal(t, false, feature.PropertyMustBool("notexist", false))
	assert.Equal(t, 78.9, feature.PropertyMustFloat64("notexist", 78.9))
}

// 测试MarshalJSON方法
func TestFeatureMarshalJSON(t *testing.T) {
	feature := NewPointFeature([]float64{1.0, 2.0})
	feature.SetProperty("name", "test")
	feature.ID = 123

	jsonData, err := json.Marshal(feature)
	assert.NoError(t, err)

	var unmarshaled map[string]interface{}
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)

	assert.Equal(t, "Feature", unmarshaled["type"])
	assert.Equal(t, 123.0, unmarshaled["id"])

	geometry, ok := unmarshaled["geometry"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Point", geometry["type"])

	properties, ok := unmarshaled["properties"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "test", properties["name"])
}

// 测试BoundingBoxFromPoints函数
func TestBoundingBoxFromPoints(t *testing.T) {
	pts := [][]float64{{1.0, 2.0}, {3.0, 4.0}, {5.0, 1.0}}
	bbox := BoundingBoxFromPoints(pts)

	assert.Equal(t, 1.0, bbox[0][0]) // west
	assert.Equal(t, 1.0, bbox[0][1]) // south
	assert.Equal(t, 0.0, bbox[0][2]) // buttom
	assert.Equal(t, 5.0, bbox[1][0]) // east
	assert.Equal(t, 4.0, bbox[1][1]) // north
	assert.Equal(t, 0.0, bbox[1][2]) // top

	// 测试带z坐标的点
	pts3D := [][]float64{{1.0, 2.0, 3.0}, {4.0, 5.0, 6.0}}
	bbox3D := BoundingBoxFromPoints(pts3D)

	assert.Equal(t, 1.0, bbox3D[0][0])
	assert.Equal(t, 2.0, bbox3D[0][1])
	assert.Equal(t, 3.0, bbox3D[0][2])
	assert.Equal(t, 4.0, bbox3D[1][0])
	assert.Equal(t, 5.0, bbox3D[1][1])
	assert.Equal(t, 6.0, bbox3D[1][2])
}

// 测试BoundingBoxsFromTwoBBox函数
func TestBoundingBoxsFromTwoBBox(t *testing.T) {
	bb1 := &BoundingBox{{1.0, 2.0, 3.0}, {4.0, 5.0, 6.0}}
	bb2 := &BoundingBox{{2.0, 1.0, 0.0}, {5.0, 6.0, 7.0}}

	combined := BoundingBoxsFromTwoBBox(bb1, bb2)

	assert.Equal(t, 1.0, combined[0][0]) // west
	assert.Equal(t, 1.0, combined[0][1]) // south
	assert.Equal(t, 0.0, combined[0][2]) // buttom
	assert.Equal(t, 5.0, combined[1][0]) // east
	assert.Equal(t, 6.0, combined[1][1]) // north
	assert.Equal(t, 7.0, combined[1][2]) // top
}

// 测试ExpandBoundingBoxs函数
func TestExpandBoundingBoxs(t *testing.T) {
	bboxes := []*BoundingBox{
		{{1.0, 2.0, 3.0}, {4.0, 5.0, 6.0}},
		{{2.0, 1.0, 0.0}, {5.0, 6.0, 7.0}},
		{{0.0, 3.0, 2.0}, {6.0, 4.0, 8.0}},
	}

	expanded := ExpandBoundingBoxs(bboxes)

	assert.Equal(t, 0.0, expanded[0][0]) // west
	assert.Equal(t, 1.0, expanded[0][1]) // south
	assert.Equal(t, 0.0, expanded[0][2]) // buttom
	assert.Equal(t, 6.0, expanded[1][0]) // east
	assert.Equal(t, 6.0, expanded[1][1]) // north
	assert.Equal(t, 8.0, expanded[1][2]) // top

	// 测试空切片
	emptyExpanded := ExpandBoundingBoxs([]*BoundingBox{})
	assert.Nil(t, emptyExpanded)

	// 测试只有一个元素的切片
	singleExpanded := ExpandBoundingBoxs([]*BoundingBox{bboxes[0]})
	assert.Equal(t, bboxes[0], singleExpanded)
}

// 测试GetKeyDifs函数
func TestGetKeyDifs(t *testing.T) {
	f1 := map[string]interface{}{"a": 1, "b": 2, "c": 3}
	f2 := map[string]interface{}{"b": 2, "c": 4, "d": 5}

	k1dif, k2dif := GetKeyDifs(f1, f2)

	assert.Equal(t, []string{"a"}, k1dif)
	assert.Equal(t, []string{"d"}, k2dif)
}

// 测试GetErrorsKeyDif函数
func TestGetErrorsKeyDif(t *testing.T) {
	kd1 := []string{"a"}
	kd2 := []string{"d"}

	errors := GetErrorsKeyDif(kd1, kd2)

	assert.Equal(t, []string{
		"Feature1 Contains field a Feature2 does not.",
		"Feature2 Contains field d Feature1 does not.",
	}, errors)
}

// 测试CheckProperties函数
func TestCheckProperties(t *testing.T) {
	p1 := map[string]interface{}{"a": 1, "b": 2, "c": "test"}
	p2 := map[string]interface{}{"a": 1, "b": 2, "c": "test"}
	p3 := map[string]interface{}{"a": 1, "b": 3, "c": "test"}
	p4 := map[string]interface{}{"a": 1, "b": 2}

	assert.True(t, CheckProperties(p1, p2))
	assert.False(t, CheckProperties(p1, p3))
	assert.False(t, CheckProperties(p1, p4))
}

// 测试ConvertFeatureID函数
func TestConvertFeatureID(t *testing.T) {
	// 测试各种类型的ID转换
	idFloat := 123.45
	idInt64 := int64(123)
	idUint64 := uint64(123)
	idUint := uint(123)
	idInt8 := int8(123)
	idUint8 := uint8(123)
	idUint16 := uint16(123)
	idInt32 := int32(123)
	idUint32 := uint32(123)
	idString := "123"

	// 测试成功转换的情况
	convertedFloat, err := ConvertFeatureID(idFloat)
	assert.NoError(t, err)
	assert.Equal(t, uint64(123), convertedFloat)

	convertedInt64, err := ConvertFeatureID(idInt64)
	assert.NoError(t, err)
	assert.Equal(t, uint64(123), convertedInt64)

	convertedUint64, err := ConvertFeatureID(idUint64)
	assert.NoError(t, err)
	assert.Equal(t, uint64(123), convertedUint64)

	convertedUint, err := ConvertFeatureID(idUint)
	assert.NoError(t, err)
	assert.Equal(t, uint64(123), convertedUint)

	convertedInt8, err := ConvertFeatureID(idInt8)
	assert.NoError(t, err)
	assert.Equal(t, uint64(123), convertedInt8)

	convertedUint8, err := ConvertFeatureID(idUint8)
	assert.NoError(t, err)
	assert.Equal(t, uint64(123), convertedUint8)

	convertedUint16, err := ConvertFeatureID(idUint16)
	assert.NoError(t, err)
	assert.Equal(t, uint64(123), convertedUint16)

	convertedInt32, err := ConvertFeatureID(idInt32)
	assert.NoError(t, err)
	assert.Equal(t, uint64(123), convertedInt32)

	convertedUint32, err := ConvertFeatureID(idUint32)
	assert.NoError(t, err)
	assert.Equal(t, uint64(123), convertedUint32)

	convertedString, err := ConvertFeatureID(idString)
	assert.NoError(t, err)
	assert.Equal(t, uint64(123), convertedString)

	// 测试无法转换的情况
	_, err = ConvertFeatureID(true)
	assert.Error(t, err)
}

// 测试ProcessGeometryData函数
func TestProcessGeometryData(t *testing.T) {
	// 创建一个Point类型的GeometryData
	pointData := []float64{1.0, 2.0}
	geometryData := &GeometryData{
		Type:  "Point",
		Point: pointData,
	}

	// 定义一个处理函数，将每个坐标加1
	processFunc := func(coord []float64) []float64 {
		return []float64{coord[0] + 1, coord[1] + 1}
	}

	// 处理GeometryData
	processedData := ProcessGeometryData(geometryData, processFunc)

	// 验证处理结果
	assert.Equal(t, GeometryPoint, processedData.Type)
	assert.Equal(t, []float64{2.0, 3.0}, processedData.Point)

	// 测试MultiPoint类型
	multiPointData := [][]float64{{1.0, 2.0}, {3.0, 4.0}}
	multiPointGeometryData := &GeometryData{
		Type:       "MultiPoint",
		MultiPoint: multiPointData,
	}

	processedMultiPointData := ProcessGeometryData(multiPointGeometryData, processFunc)

	assert.Equal(t, GeometryMultiPoint, processedMultiPointData.Type)
	assert.Equal(t, [][]float64{{2.0, 3.0}, {4.0, 5.0}}, processedMultiPointData.MultiPoint)
}

// 测试IsFeatureEqual函数
func TestIsFeatureEqual(t *testing.T) {
	// 创建两个相同的feature
	point1 := MockPoint{Point: []float64{1.0, 2.0}}
	feature1 := NewFeature(point1)
	feature1.SetProperty("name", "test")

	point2 := MockPoint{Point: []float64{1.0, 2.0}}
	feature2 := NewFeature(point2)
	feature2.SetProperty("name", "test")

	// 创建一个不同的feature
	point3 := MockPoint{Point: []float64{3.0, 4.0}}
	feature3 := NewFeature(point3)
	feature3.SetProperty("name", "test")

	// 创建一个属性不同的feature
	point4 := MockPoint{Point: []float64{1.0, 2.0}}
	feature4 := NewFeature(point4)
	feature4.SetProperty("name", "different")

	assert.True(t, IsFeatureEqual(*feature1, *feature2))
	assert.False(t, IsFeatureEqual(*feature1, *feature3))
	assert.False(t, IsFeatureEqual(*feature1, *feature4))
}

// 测试DecodeBoundingBox函数
func TestDecodeBoundingBox(t *testing.T) {
	// 测试[]float64类型
	bbFloat := []float64{1.0, 2.0, 3.0, 4.0}
	decodedFloat, err := DecodeBoundingBox(bbFloat)
	assert.NoError(t, err)
	assert.Equal(t, bbFloat, decodedFloat)

	// 测试[]interface{}类型
	bbInterface := []interface{}{1.0, 2.0, 3.0, 4.0}
	decodedInterface, err := DecodeBoundingBox(bbInterface)
	assert.NoError(t, err)
	assert.Equal(t, bbFloat, decodedInterface)

	// 测试nil
	decodedNil, err := DecodeBoundingBox(nil)
	assert.NoError(t, err)
	assert.Nil(t, decodedNil)

	// 测试其他类型
	_, err = DecodeBoundingBox(123)
	assert.Error(t, err)
}

// 测试NewFeatureCollection函数
func TestNewFeatureCollection(t *testing.T) {
	fc := NewFeatureCollection()

	assert.Equal(t, "FeatureCollection", fc.Type)
	assert.Empty(t, fc.Features)
	assert.NotNil(t, fc.Properties)
	assert.Nil(t, fc.BoundingBox)
	assert.Nil(t, fc.CRS)
}

// 测试AddFeature方法
func TestFeatureCollectionAddFeature(t *testing.T) {
	fc := NewFeatureCollection()

	// 创建测试要素
	point := MockPoint{Point: []float64{1.0, 2.0}}
	feature := NewFeature(point)

	// 添加要素
	fc.AddFeature(feature)

	// 验证
	assert.Len(t, fc.Features, 1)
	assert.Equal(t, feature, fc.Features[0])

	// 添加第二个要素
	point2 := MockPoint{Point: []float64{3.0, 4.0}}
	feature2 := NewFeature(point2)
	fc.AddFeature(feature2)

	// 验证
	assert.Len(t, fc.Features, 2)
	assert.Equal(t, feature2, fc.Features[1])
}

// 测试MarshalJSON方法
func TestFeatureCollectionMarshalJSON(t *testing.T) {
	fc := NewFeatureCollection()

	// 添加属性
	fc.Properties["name"] = "test collection"
	fc.Properties["count"] = 2

	// 创建并添加要素
	point1 := MockPoint{Point: []float64{1.0, 2.0}}
	feature1 := NewFeature(point1)
	feature1.SetProperty("name", "feature1")

	point2 := MockPoint{Point: []float64{3.0, 4.0}}
	feature2 := NewFeature(point2)
	feature2.SetProperty("name", "feature2")

	fc.AddFeature(feature1)
	fc.AddFeature(feature2)

	// 设置边界框
	bbox := &BoundingBox{{1.0, 2.0, 0.0}, {3.0, 4.0, 0.0}}
	fc.BoundingBox = bbox

	// 序列化
	jsonData, err := json.Marshal(fc)
	assert.NoError(t, err)

	// 反序列化以验证
	var unmarshaled map[string]interface{}
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)

	// 验证基本字段
	assert.Equal(t, "FeatureCollection", unmarshaled["type"])

	// 验证属性
	properties, ok := unmarshaled["properties"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "test collection", properties["name"])
	assert.Equal(t, 2.0, properties["count"])

	// 验证边界框
	bboxJSON, ok := unmarshaled["bbox"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, bboxJSON, 2)

	// 验证min坐标
	minJSON, ok := bboxJSON[0].([]interface{})
	assert.True(t, ok)
	assert.Len(t, minJSON, 3)
	assert.Equal(t, 1.0, minJSON[0])
	assert.Equal(t, 2.0, minJSON[1])
	assert.Equal(t, 0.0, minJSON[2])

	// 验证max坐标
	maxJSON, ok := bboxJSON[1].([]interface{})
	assert.True(t, ok)
	assert.Len(t, maxJSON, 3)
	assert.Equal(t, 3.0, maxJSON[0])
	assert.Equal(t, 4.0, maxJSON[1])
	assert.Equal(t, 0.0, maxJSON[2])

	// 验证要素
	features, ok := unmarshaled["features"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, features, 2)

	// 验证第一个要素
	feature1JSON, ok := features[0].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Feature", feature1JSON["type"])

	// 验证第一个要素的属性
	feature1Props, ok := feature1JSON["properties"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "feature1", feature1Props["name"])

	// 验证第二个要素
	feature2JSON, ok := features[1].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Feature", feature2JSON["type"])

	// 验证第二个要素的属性
	feature2Props, ok := feature2JSON["properties"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "feature2", feature2Props["name"])
}
