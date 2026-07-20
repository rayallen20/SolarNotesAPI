# 天体简介接口

## 接口地址

```http
GET /api/v1/planet/list
```

## 请求方式

GET

## 请求 Header

无

## 请求示例

无

## 请求参数

无

## 响应示例

```json
{
    "code": 200,
    "message": "success",
    "data": {
        "planets": [
            {
                "id": 1,
                "intro": "操作系统相关介绍",
                "name": "Sun",
                "title": "操作系统"
            },
            {
                "id": 2,
                "intro": "数据库相关介绍",
                "name": "Mercury",
                "title": "数据库"
            },
            {
                "id": 3,
                "intro": "数据库相关介绍",
                "name": "Venus",
                "title": "数据库"
            },
            {
                "id": 4,
                "intro": "暂无",
                "name": "Earth",
                "title": "暂无"
            },
            {
                "id": 5,
                "intro": "暂无",
                "name": "Mars",
                "title": "暂无"
            },
            {
                "id": 6,
                "intro": "暂无",
                "name": "Jupiter",
                "title": "暂无"
            },
            {
                "id": 7,
                "intro": "暂无",
                "name": "Saturn",
                "title": "暂无"
            },
            {
                "id": 8,
                "intro": "暂无",
                "name": "Uranus",
                "title": "暂无"
            },
            {
                "id": 9,
                "intro": "暂无",
                "name": "Neptune",
                "title": "暂无"
            }
        ]
    }
}
```

## 响应参数

|  参数名   |   类型   |   说明   |
|:---------:|:--------:|:--------:|
|  `code`   |  `int`   |  状态码  |
| `message` | `string` | 提示信息 |
|  `data`   | `object` | 返回数据 |

## `data`部分参数说明

|  参数名   |      类型       |       说明       |
|:---------:|:---------------:|:----------------:|
| `planets` | `Array<Object>` | 天体简介对象列表 |

### `planets`部分参数说明

| 参数名  |   类型   |                         说明                         |
|:-------:|:--------:|:----------------------------------------------------:|
|  `id`   |  `int`   |                        天体id                        |
| `name`  | `string` | 天体名称(前端据此字段值与对应的天体配置对象组合字段) |
| `title` | `string` |                   天体label的标题                    |
| `intro` | `string` |                   天体label的介绍                    |

## 备注

无