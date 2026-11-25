# 获取目录接口

## 接口地址

```http
POST /api/v1/catalogue/show
```

## 请求方式

POST

## 请求`Header`

无

## 请求示例

```json
{
    "book": {
        "id": "9007199254740995"
    }
}
```

## 请求参数

| 参数名 |   类型   | 必填 |   说明   |
|:------:|:--------:|:----:|:--------:|
| `book` | `Object` |  是  | 书籍对象 |

### `book`部分请求参数

| 参数名 | 类型  | 必填 |  说明  |
|:------:|:-----:|:----:|:------:|
|  `id`  | `string` | 是 | 书籍 ID，取正 `int64` 的十进制字符串，范围为 1 到 9223372036854775807 |

将书单返回的 `data.bookList[].id` 原样传入此字段。JSON 使用双引号，例如 `"9007199254740995"`。

## 响应示例

```json
{
    "code": 200,
    "message": "success",
    "data": {
        "catalogue": {
            "id": 1,
            "type": "folder",
            "name": "书籍名称",
            "intro": "书籍简介",
            "createdAt": "2025-12-31",
            "children": [
                {
                    "id": 2,
                    "type": "folder",
                    "name": "第1章: XXXX",
                    "intro": "第1章的简介",
                    "createdAt": "2026-01-01",
                    "children": [
                        {
                            "id": 3,
                            "type": "file",
                            "name": "第1章第1篇文章: XYYY"
                        },
                        {
                            "id": 4,
                            "type": "file",
                            "name": "第1章第2篇文章: XYXY"
                        }
                    ]
                },
                {
                    "id": 5,
                    "type": "folder",
                    "name": "第2章: YYYY",
                    "intro": "第2章的简介",
                    "createdAt": "2026-01-02",
                    "children": [
                        {
                            "id": 6,
                            "type": "folder",
                            "name": "第2章第1节",
                            "intro": "第2章第1节的简介",
                            "createdAt": "2026-01-03",
                            "children": [
                                {
                                    "id": 7,
                                    "type": "file",
                                    "name": "第2章第1节的第1篇文章"
                                },
                                {
                                    "id": 8,
                                    "type": "file",
                                    "name": "第2章第1节的第2篇文章"
                                }
                            ]
                        },
                        {
                            "id": 9,
                            "type": "folder",
                            "name": "第2章第2节",
                            "intro": "第2章第2节的简介",
                            "createdAt": "2026-01-04",
                            "children": [
                                {
                                    "id": 10,
                                    "type": "file",
                                    "name": "第2章第2节的第1篇文章"
                                },
                                {
                                    "id": 11,
                                    "type": "file",
                                    "name": "第2章第2节的第2篇文章"
                                }
                            ]
                        }
                    ]
                },
                {
                    "id": 12,
                    "type": "folder",
                    "name": "第3章: ZZZZ",
                    "intro": "第3章的简介",
                    "createdAt": "2026-01-11",
                    "children": [
                        {
                            "id": 13,
                            "type": "folder",
                            "name": "第3章第1节",
                            "intro": "第3章第1节的简介",
                            "createdAt": "2026-01-12",
                            "children": [
                                {
                                    "id": 14,
                                    "type": "file",
                                    "name": "第3章第1节的第1篇文章"
                                },
                                {
                                    "id": 15,
                                    "type": "file",
                                    "name": "第3章第1节的第2篇文章"
                                }
                            ]
                        },
                        {
                            "id": 16,
                            "type": "folder",
                            "name": "第3章第2节",
                            "intro": "第3章第2节的简介",
                            "createdAt": "2026-01-15",
                            "children": [
                                {
                                    "id": 17,
                                    "type": "file",
                                    "name": "第3章第2节的第1篇文章"
                                },
                                {
                                    "id": 18,
                                    "type": "file",
                                    "name": "第3章第2节的第2篇文章"
                                }
                            ]
                        }
                    ]
                },
                {
                    "id": 19,
                    "type": "folder",
                    "name": "第4章",
                    "intro": "第4章的简介",
                    "createdAt": "2026-01-11",
                    "children": [
                        {
                            "id": 20,
                            "type": "folder",
                            "name": "第4章第1节",
                            "intro": "第4章第1节的简介",
                            "createdAt": "2026-01-12",
                            "children": [
                                {
                                    "id": 21,
                                    "type": "file",
                                    "name": "第4章第1节的第1篇文章"
                                },
                                {
                                    "id": 22,
                                    "type": "file",
                                    "name": "第4章第1节的第2篇文章"
                                }
                            ]
                        },
                        {
                            "id": 23,
                            "type": "folder",
                            "name": "第4章第2节",
                            "intro": "第4章第2节的简介",
                            "createdAt": "2026-01-15",
                            "children": [
                                {
                                    "id": 24,
                                    "type": "file",
                                    "name": "第4章第2节的第1篇文章"
                                },
                                {
                                    "id": 25,
                                    "type": "file",
                                    "name": "第4章第2节的第2篇文章"
                                }
                            ]
                        }
                    ]
                },
                {
                    "id": 26,
                    "type": "folder",
                    "name": "第5章",
                    "intro": "第5章的简介",
                    "createdAt": "2026-01-11",
                    "children": [
                        {
                            "id": 27,
                            "type": "folder",
                            "name": "第5章第1节",
                            "intro": "第5章第1节的简介",
                            "createdAt": "2026-01-12",
                            "children": [
                                {
                                    "id": 28,
                                    "type": "file",
                                    "name": "第5章第1节的第1篇文章"
                                },
                                {
                                    "id": 29,
                                    "type": "file",
                                    "name": "第5章第1节的第2篇文章"
                                }
                            ]
                        },
                        {
                            "id": 30,
                            "type": "folder",
                            "name": "第5章第2节",
                            "intro": "第5章第2节的简介",
                            "createdAt": "2026-01-15",
                            "children": [
                                {
                                    "id": 31,
                                    "type": "file",
                                    "name": "第5章第2节的第1篇文章"
                                },
                                {
                                    "id": 32,
                                    "type": "file",
                                    "name": "第5章第2节的第2篇文章"
                                }
                            ]
                        }
                    ]
                }
            ]
        }
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

|   参数名    |   类型   |   说明   |
|:-----------:|:--------:|:--------:|
| `catalogue` | `Object` | 目录对象 |

### `catalogue`部分参数说明

|   参数名    |        类型        |                                                 说明                                                  |
|:-----------:|:------------------:|:-----------------------------------------------------------------------------------------------------:|
|    `id`     |       `int`        |                                                章节id                                                 |
|   `type`    |      `string`      | 可能的枚举值为:`folder`/`file`.`folder`: 表示当前对象类型非叶节点; `file`: 表示当前对象类型为叶子节点 |
|   `name`    |      `string`      |                                               节点名称                                                |
|   `intro`   |      `string`      |                         章节简介(本字段**仅当**`type`字段值为`folder`时存在)                          |
| `createdAt` |      `string`      |                       节点创建时间(本字段**仅当**`type`字段值为`folder`时存在)                        |
| `children`  | `Array<catalogue>` |                                              子节点列表                                               |

## 备注

无