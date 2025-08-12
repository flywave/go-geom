package geom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试SridToUrn函数
func TestSridToUrn(t *testing.T) {
	// 测试SRID为4326的特殊情况
	assert.Equal(t, "urn:ogc:def:crs:OGC:1.3:CRS84", SridToUrn(4326))

	// 测试其他SRID值
	assert.Equal(t, "urn:ogc:def:crs:EPSG::3857", SridToUrn(3857))
	assert.Equal(t, "urn:ogc:def:crs:EPSG::4490", SridToUrn(4490))
	assert.Equal(t, "urn:ogc:def:crs:EPSG::1234", SridToUrn(1234))
}

// 测试UrnToSrid函数
func TestUrnToSrid(t *testing.T) {
	// 测试CRS84的特殊情况
	assert.Equal(t, 4326, UrnToSrid("urn:ogc:def:crs:OGC:1.3:CRS84"))

	// 测试EPSG格式的URN
	assert.Equal(t, 3857, UrnToSrid("urn:ogc:def:crs:EPSG::3857"))
	assert.Equal(t, 4490, UrnToSrid("urn:ogc:def:crs:EPSG::4490"))
	assert.Equal(t, 1234, UrnToSrid("urn:ogc:def:crs:EPSG::1234"))

	// 测试无效格式的URN
	assert.Equal(t, -1, UrnToSrid("invalid_urn"))
	assert.Equal(t, -1, UrnToSrid("urn:ogc:def:crs:EPSG:1234")) // 格式错误
	assert.Equal(t, -1, UrnToSrid("urn:ogc:def:crs:EPSG::abc")) // 非数字SRID
}
