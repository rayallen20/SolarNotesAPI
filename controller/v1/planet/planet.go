package planet

import (
	"SolarNotesAPI/resp_body"
	"net/http"

	"github.com/gin-gonic/gin"
)

// List 天体简介API
func List(c *gin.Context) {
	respBody := &resp_body.RespBody{}

	data := map[string]interface{}{
		"planets": []map[string]interface{}{
			{
				"id":    1,
				"name":  "Sun",
				"title": "操作系统",
				"intro": "操作系统相关介绍",
			},
			{
				"id":    2,
				"name":  "Mercury",
				"title": "数据库",
				"intro": "数据库相关介绍",
			},
			{
				"id":    3,
				"name":  "Venus",
				"title": "数据库",
				"intro": "数据库相关介绍",
			},
			{
				"id":    4,
				"name":  "Earth",
				"title": "暂无",
				"intro": "暂无",
			},
			{
				"id":    5,
				"name":  "Mars",
				"title": "暂无",
				"intro": "暂无",
			},
			{
				"id":    6,
				"name":  "Jupiter",
				"title": "暂无",
				"intro": "暂无",
			},
			{
				"id":    7,
				"name":  "Saturn",
				"title": "暂无",
				"intro": "暂无",
			},
			{
				"id":    8,
				"name":  "Uranus",
				"title": "暂无",
				"intro": "暂无",
			},
			{
				"id":    9,
				"name":  "Neptune",
				"title": "暂无",
				"intro": "暂无",
			},
		},
	}

	respBody.Success(data)
	c.JSON(http.StatusOK, respBody)
}
