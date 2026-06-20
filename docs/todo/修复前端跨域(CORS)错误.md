# 修复前端跨域(CORS)错误

## 一、问题现象

前端在浏览器中请求本工程的 API 时,控制台报错 `CORS Error`(跨域错误),导致请求失败、拿不到接口数据。

## 二、问题原因

CORS(Cross-Origin Resource Sharing,跨域资源共享)是浏览器的一种安全机制:当前端页面的「源」(协议 + 域名 + 端口)与所请求接口的「源」不一致时,浏览器会拦截响应,除非服务端在响应头中明确返回 `Access-Control-Allow-Origin` 等允许跨域的字段。

本工程当前的情况:

- `go.mod` 中**已经引入**了 `github.com/gin-contrib/cors v1.7.7` 依赖;
- 但 `main.go` 中**并没有注册**该跨域中间件。

因此服务端从未向浏览器返回任何 CORS 响应头,浏览器便拦截了响应,前端表现为 `CORS Error`。

> 说明:这并不是前端代码的问题,而是后端服务缺少跨域响应头。即使接口本身能正常返回数据,浏览器也会在收到响应后将其拦截。

## 三、修复方案

在 `main.go` 中注册 `gin-contrib/cors` 中间件,让服务端为所有响应附带跨域响应头。该中间件需要在**注册路由之前**通过 `r.Use(...)` 挂载到全局。

### 修改文件:`main.go`

#### 1. 增加 import

```go
import (
	"SolarNotesAPI/controller/article"
	"SolarNotesAPI/controller/v1/book"
	"SolarNotesAPI/controller/v1/catalogue"
	"SolarNotesAPI/controller/v1/planet"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)
```

#### 2. 在创建路由引擎后、注册路由前挂载中间件

```go
func main() {
	r := gin.Default()

	// 跨域中间件
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
	}))

	// 路由分组
	v1 := r.Group("/api/v1")

	{
		// 天体简介API
		v1.GET("/planet/list", planet.List)

		// 书籍列表API
		v1.GET("/book/list", book.List)

		// 获取目录API
		v1.GET("/catalogue/show", catalogue.Show)

		// 获取文章API
		v1.GET("/article/show", article.Show)
	}

	r.Run("0.0.0.0:4061")
}
```

## 四、配置项说明

| 配置项 | 含义 |
| --- | --- |
| `AllowAllOrigins` | 是否允许任意来源访问。设为 `true` 表示允许所有源,适合开发或公开只读接口。 |
| `AllowMethods` | 允许的 HTTP 请求方法。需要包含 `OPTIONS`,以便正确响应浏览器的预检(preflight)请求。 |
| `AllowHeaders` | 允许前端携带的请求头字段。如有自定义请求头(如 token),需在此处补充。 |

## 五、关于生产环境的建议

`AllowAllOrigins: true` 允许任意来源访问,使用方便,适合开发阶段或公开的只读接口。若上线后希望收紧权限,可改为只允许指定的前端域名:

```go
r.Use(cors.New(cors.Config{
	AllowOrigins:  []string{"https://你的前端域名.com"},
	AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
}))
```

> 注意:`AllowAllOrigins: true` 与 `AllowOrigins: []string{...}` 二者互斥,不能同时设置。

## 六、验证步骤

1. 应用上述修改后,重新编译并启动服务:
   ```bash
   go run main.go
   ```
2. 在浏览器中重新发起前端请求,确认 `CORS Error` 不再出现、接口数据正常返回。
3. (可选)用 `curl` 检查响应头中是否包含跨域字段:
   ```bash
   curl -I -H "Origin: http://localhost:3000" http://localhost:4061/api/v1/planet/list
   ```
   响应头中应出现 `Access-Control-Allow-Origin` 字段。