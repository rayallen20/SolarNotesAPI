package main

import (
	"log"

	"SolarNotesAPI/conf"
	"SolarNotesAPI/controller/v1/article"
	"SolarNotesAPI/controller/v1/book"
	"SolarNotesAPI/controller/v1/catalogue"
	"SolarNotesAPI/controller/v1/planet"
	"SolarNotesAPI/db"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := conf.Load()
	if err != nil {
		return err
	}
	if err := db.Connect(cfg.Database); err != nil {
		return err
	}
	pool, err := db.DB.DB()
	if err != nil {
		return err
	}
	defer pool.Close()

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

	return r.Run("0.0.0.0:" + cfg.Service.Port)
}
