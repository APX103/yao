package global

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/gou"
	"github.com/yaoapp/kun/any"
	"github.com/yaoapp/xiang/table"
	"github.com/yaoapp/xun/capsule"
)

func TestProcessPing(t *testing.T) {
	process := gou.NewProcess("xiang.global.ping")
	res, ok := processPing(process).(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, res["version"], VERSION)
}

func TestProcessSearch(t *testing.T) {
	args := []interface{}{
		"service",
		gou.QueryParam{
			Wheres: []gou.QueryWhere{
				{Column: "status", Value: "enabled"},
			},
		},
		1,
		2,
		&gin.Context{},
	}
	process := gou.NewProcess("xiang.table.Search", args...)
	response := table.ProcessSearch(process)
	assert.NotNil(t, response)
	res := any.Of(response).Map()
	assert.True(t, res.Has("data"))
	assert.True(t, res.Has("next"))
	assert.True(t, res.Has("page"))
	assert.True(t, res.Has("pagecnt"))
	assert.True(t, res.Has("pagesize"))
	assert.True(t, res.Has("prev"))
	assert.True(t, res.Has("total"))
	assert.Equal(t, 1, res.Get("page"))
	assert.Equal(t, 2, res.Get("pagesize"))
}

func TestProcessFind(t *testing.T) {
	args := []interface{}{
		"service",
		1,
		&gin.Context{},
	}
	process := gou.NewProcess("xiang.table.Find", args...)
	response := table.ProcessFind(process)
	assert.NotNil(t, response)
	res := any.Of(response).Map()
	assert.Equal(t, any.Of(res.Get("id")).CInt(), 1)
}

func TestProcessSave(t *testing.T) {
	args := []interface{}{
		"service",
		map[string]interface{}{
			"name":          "腾讯黑岩云主机",
			"short_name":    "高性能云主机",
			"kind_id":       3,
			"manu_id":       1,
			"price_options": []string{"按月订阅"},
		},
	}
	process := gou.NewProcess("xiang.table.Save", args...)
	response := table.ProcessSave(process)
	assert.NotNil(t, response)
	assert.True(t, any.Of(response).IsInt())

	id := any.Of(response).CInt()

	// 清空数据
	capsule.Query().Table("service").Where("id", id).Delete()
}

func TestProcessDelete(t *testing.T) {
	args := []interface{}{
		"service",
		map[string]interface{}{
			"name":          "腾讯黑岩云主机",
			"short_name":    "高性能云主机",
			"kind_id":       3,
			"manu_id":       1,
			"price_options": []string{"按月订阅"},
		},
	}
	process := gou.NewProcess("xiang.table.Save", args...)
	response := table.ProcessSave(process)
	assert.NotNil(t, response)
	assert.True(t, any.Of(response).IsInt())

	id := any.Of(response).CInt()
	args = []interface{}{
		"service",
		id,
	}
	process = gou.NewProcess("xiang.table.Delete", args...)
	response = table.ProcessDelete(process)
	assert.Nil(t, response)

	// 清空数据
	capsule.Query().Table("service").Where("id", id).Delete()
}

func TestProcessSetting(t *testing.T) {
	args := []interface{}{"service"}
	process := gou.NewProcess("xiang.table.Setting", args...)
	response := table.ProcessSetting(process)
	assert.NotNil(t, response)
	res := any.Of(response).Map()
	assert.Equal(t, res.Get("name"), "云服务库")
	assert.True(t, res.Has("title"))
	assert.True(t, res.Has("decription"))
	assert.True(t, res.Has("columns"))
	assert.True(t, res.Has("filters"))
	assert.True(t, res.Has("list"))
	assert.True(t, res.Has("edit"))
	assert.True(t, res.Has("view"))
	assert.True(t, res.Has("insert"))
}