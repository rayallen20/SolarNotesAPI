package book

import (
	"SolarNotesAPI/param/req/v1/book"
	"SolarNotesAPI/resp_body"
	"net/http"

	"github.com/gin-gonic/gin"
)

// List 书籍列表API
func List(c *gin.Context) {
	req := &book.ListReq{}
	respBody := &resp_body.RespBody{}

	err := c.ShouldBindJSON(req)
	if err != nil {
		respBody.ValidateFailed(err)
		c.JSON(http.StatusOK, respBody)
		return
	}

	data := map[string]interface{}{
		"bookList": []map[string]interface{}{
			{
				"id":    1,
				"title": "详解操作系统",
			},
			{
				"id":    2,
				"title": "操作系统解密",
			},
			{
				"id":    3,
				"title": "数据库教程",
			},
			{
				"id":    4,
				"title": "算法基础",
			},
			{
				"id":    5,
				"title": "计算机网络",
			},
			{
				"id":    6,
				"title": "设计模式",
			},
			{
				"id":    7,
				"title": "编码规范",
			},
			{
				"id":    8,
				"title": "微服务设计模式",
			},
		},
	}

	respBody.Success(data)
	c.JSON(http.StatusOK, respBody)
	return
}
