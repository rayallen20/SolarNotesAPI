package main

import (
	"SolarNotesAPI/controller/v1/article"
	"SolarNotesAPI/controller/v1/book"
	"SolarNotesAPI/controller/v1/catalogue"
	"SolarNotesAPI/controller/v1/planet"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 使用跨域中间件
	corsConfig := cors.Config{
		AllowAllOrigins: true,
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
	}
	r.Use(cors.New(corsConfig))

	// 路由分组
	v1 := r.Group("/api/v1")

	{
		// 天体简介API
		v1.GET("/planet/list", planet.List)

		// 书籍列表API
		v1.POST("/book/list", book.List)

		// 获取目录API
		v1.POST("/catalogue/show", catalogue.Show)

		// 获取文章API
		v1.POST("/article/show", article.Show)
	}

	r.Run("0.0.0.0:4061")
}
