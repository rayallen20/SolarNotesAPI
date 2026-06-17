package catalogue

import (
	"SolarNotesAPI/param/req/v1/catalogue"
	"SolarNotesAPI/resp_body"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Show 获取目录API
func Show(c *gin.Context) {
	req := &catalogue.ShowReq{}
	respBody := &resp_body.RespBody{}

	err := c.ShouldBindJSON(req)
	if err != nil {
		respBody.ValidateFailed(err)
		c.JSON(http.StatusOK, respBody)
		return
	}

	data := map[string]interface{}{
		"catalogue": map[string]interface{}{
			"id":        1,
			"type":      "folder",
			"name":      "书籍名称",
			"intro":     "书籍简介",
			"createdAt": "2025-12-31",
			"children": []map[string]interface{}{
				{
					"id":        2,
					"type":      "folder",
					"name":      "第1章: XXXX",
					"intro":     "第1章的简介",
					"createdAt": "2026-01-01",
					"children": []map[string]interface{}{
						{
							"id":   3,
							"type": "file",
							"name": "第1章第1篇文章: XYYY",
						},
						{
							"id":   4,
							"type": "file",
							"name": "第1章第2篇文章: XYXY",
						},
					},
				},
				{
					"id":        5,
					"type":      "folder",
					"name":      "第2章: YYYY",
					"intro":     "第2章的简介",
					"createdAt": "2026-01-02",
					"children": []map[string]interface{}{
						{
							"id":        6,
							"type":      "folder",
							"name":      "第2章第1节",
							"intro":     "第2章第1节的简介",
							"createdAt": "2026-01-03",
							"children": []map[string]interface{}{
								{
									"id":   7,
									"type": "file",
									"name": "第2章第1节的第1篇文章",
								},
								{
									"id":   8,
									"type": "file",
									"name": "第2章第1节的第2篇文章",
								},
							},
						},
						{
							"id":        9,
							"type":      "folder",
							"name":      "第2章第2节",
							"intro":     "第2章第2节的简介",
							"createdAt": "2026-01-04",
							"children": []map[string]interface{}{
								{
									"id":   10,
									"type": "file",
									"name": "第2章第2节的第1篇文章",
								},
								{
									"id":   11,
									"type": "file",
									"name": "第2章第2节的第2篇文章",
								},
							},
						},
					},
				},
				{
					"id":        12,
					"type":      "folder",
					"name":      "第3章: ZZZZ",
					"intro":     "第3章的简介",
					"createdAt": "2026-01-11",
					"children": []map[string]interface{}{
						{
							"id":        13,
							"type":      "folder",
							"name":      "第3章第1节",
							"intro":     "第3章第1节的简介",
							"createdAt": "2026-01-12",
							"children": []map[string]interface{}{
								{
									"id":   14,
									"type": "file",
									"name": "第3章第1节的第1篇文章",
								},
								{
									"id":   15,
									"type": "file",
									"name": "第3章第1节的第2篇文章",
								},
							},
						},
						{
							"id":        16,
							"type":      "folder",
							"name":      "第3章第2节",
							"intro":     "第3章第2节的简介",
							"createdAt": "2026-01-15",
							"children": []map[string]interface{}{
								{
									"id":   17,
									"type": "file",
									"name": "第3章第2节的第1篇文章",
								},
								{
									"id":   18,
									"type": "file",
									"name": "第3章第2节的第2篇文章",
								},
							},
						},
					},
				},
				{
					"id":        19,
					"type":      "folder",
					"name":      "第4章",
					"intro":     "第4章的简介",
					"createdAt": "2026-01-11",
					"children": []map[string]interface{}{
						{
							"id":        20,
							"type":      "folder",
							"name":      "第4章第1节",
							"intro":     "第4章第1节的简介",
							"createdAt": "2026-01-12",
							"children": []map[string]interface{}{
								{
									"id":   21,
									"type": "file",
									"name": "第4章第1节的第1篇文章",
								},
								{
									"id":   22,
									"type": "file",
									"name": "第4章第1节的第2篇文章",
								},
							},
						},
						{
							"id":        23,
							"type":      "folder",
							"name":      "第4章第2节",
							"intro":     "第4章第2节的简介",
							"createdAt": "2026-01-15",
							"children": []map[string]interface{}{
								{
									"id":   24,
									"type": "file",
									"name": "第4章第2节的第1篇文章",
								},
								{
									"id":   25,
									"type": "file",
									"name": "第4章第2节的第2篇文章",
								},
							},
						},
					},
				},
				{
					"id":        26,
					"type":      "folder",
					"name":      "第5章",
					"intro":     "第5章的简介",
					"createdAt": "2026-01-11",
					"children": []map[string]interface{}{
						{
							"id":        27,
							"type":      "folder",
							"name":      "第5章第1节",
							"intro":     "第5章第1节的简介",
							"createdAt": "2026-01-12",
							"children": []map[string]interface{}{
								{
									"id":   28,
									"type": "file",
									"name": "第5章第1节的第1篇文章",
								},
								{
									"id":   29,
									"type": "file",
									"name": "第5章第1节的第2篇文章",
								},
							},
						},
						{
							"id":        30,
							"type":      "folder",
							"name":      "第5章第2节",
							"intro":     "第5章第2节的简介",
							"createdAt": "2026-01-15",
							"children": []map[string]interface{}{
								{
									"id":   31,
									"type": "file",
									"name": "第5章第2节的第1篇文章",
								},
								{
									"id":   32,
									"type": "file",
									"name": "第5章第2节的第2篇文章",
								},
							},
						},
					},
				},
			},
		},
	}

	respBody.Success(data)
	c.JSON(http.StatusOK, respBody)
	return
}
