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
            "id": "1",
            "type": "folder",
            "name": "书籍名称",
            "intro": "书籍简介",
            "createdAt": "2025-12-31",
            "children": [
                {
                    "id": "2",
                    "type": "folder",
                    "name": "第1章: XXXX",
                    "intro": "第1章的简介",
                    "createdAt": "2026-01-01",
                    "children": [
                        {
                            "id": "9017199254740996",
                            "type": "file",
                            "name": "第1章第1篇文章: XYYY"
                        },
                        {
                            "id": "4",
                            "type": "file",
                            "name": "第1章第2篇文章: XYXY"
                        }
                    ]
                },
                {
                    "id": "5",
                    "type": "folder",
                    "name": "第2章: YYYY",
                    "intro": "第2章的简介",
                    "createdAt": "2026-01-02",
                    "children": [
                        {
                            "id": "6",
                            "type": "folder",
                            "name": "第2章第1节",
                            "intro": "第2章第1节的简介",
                            "createdAt": "2026-01-03",
                            "children": [
                                {
                                    "id": "7",
                                    "type": "file",
                                    "name": "第2章第1节的第1篇文章"
                                },
                                {
                                    "id": "8",
                                    "type": "file",
                                    "name": "第2章第1节的第2篇文章"
                                }
                            ]
                        },
                        {
                            "id": "9",
                            "type": "folder",
                            "name": "第2章第2节",
                            "intro": "第2章第2节的简介",
                            "createdAt": "2026-01-04",
                            "children": [
                                {
                                    "id": "10",
                                    "type": "file",
                                    "name": "第2章第2节的第1篇文章"
                                },
                                {
                                    "id": "11",
                                    "type": "file",
                                    "name": "第2章第2节的第2篇文章"
                                }
                            ]
                        }
                    ]
                },
                {
                    "id": "12",
                    "type": "folder",
                    "name": "第3章: ZZZZ",
                    "intro": "第3章的简介",
                    "createdAt": "2026-01-11",
                    "children": [
                        {
                            "id": "13",
                            "type": "folder",
                            "name": "第3章第1节",
                            "intro": "第3章第1节的简介",
                            "createdAt": "2026-01-12",
                            "children": [
                                {
                                    "id": "14",
                                    "type": "file",
                                    "name": "第3章第1节的第1篇文章"
                                },
                                {
                                    "id": "15",
                                    "type": "file",
                                    "name": "第3章第1节的第2篇文章"
                                }
                            ]
                        },
                        {
                            "id": "16",
                            "type": "folder",
                            "name": "第3章第2节",
                            "intro": "第3章第2节的简介",
                            "createdAt": "2026-01-15",
                            "children": [
                                {
                                    "id": "17",
                                    "type": "file",
                                    "name": "第3章第2节的第1篇文章"
                                },
                                {
                                    "id": "18",
                                    "type": "file",
                                    "name": "第3章第2节的第2篇文章"
                                }
                            ]
                        }
                    ]
                },
                {
                    "id": "19",
                    "type": "folder",
                    "name": "第4章",
                    "intro": "第4章的简介",
                    "createdAt": "2026-01-11",
                    "children": [
                        {
                            "id": "20",
                            "type": "folder",
                            "name": "第4章第1节",
                            "intro": "第4章第1节的简介",
                            "createdAt": "2026-01-12",
                            "children": [
                                {
                                    "id": "21",
                                    "type": "file",
                                    "name": "第4章第1节的第1篇文章"
                                },
                                {
                                    "id": "22",
                                    "type": "file",
                                    "name": "第4章第1节的第2篇文章"
                                }
                            ]
                        },
                        {
                            "id": "23",
                            "type": "folder",
                            "name": "第4章第2节",
                            "intro": "第4章第2节的简介",
                            "createdAt": "2026-01-15",
                            "children": [
                                {
                                    "id": "24",
                                    "type": "file",
                                    "name": "第4章第2节的第1篇文章"
                                },
                                {
                                    "id": "25",
                                    "type": "file",
                                    "name": "第4章第2节的第2篇文章"
                                }
                            ]
                        }
                    ]
                },
                {
                    "id": "26",
                    "type": "folder",
                    "name": "第5章",
                    "intro": "第5章的简介",
                    "createdAt": "2026-01-11",
                    "children": [
                        {
                            "id": "27",
                            "type": "folder",
                            "name": "第5章第1节",
                            "intro": "第5章第1节的简介",
                            "createdAt": "2026-01-12",
                            "children": [
                                {
                                    "id": "28",
                                    "type": "file",
                                    "name": "第5章第1节的第1篇文章"
                                },
                                {
                                    "id": "29",
                                    "type": "file",
                                    "name": "第5章第1节的第2篇文章"
                                }
                            ]
                        },
                        {
                            "id": "30",
                            "type": "folder",
                            "name": "第5章第2节",
                            "intro": "第5章第2节的简介",
                            "createdAt": "2026-01-15",
                            "children": [
                                {
                                    "id": "31",
                                    "type": "file",
                                    "name": "第5章第2节的第1篇文章"
                                },
                                {
                                    "id": "32",
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
| `catalogue` | `Object` | 所请求书籍的单个公开根文件夹 |

### `catalogue`部分参数说明

| 参数名 | 类型 | 说明 |
|:------:|:----:|:-----|
| `id` | `string` | 目录节点 ID，正 `int64` 的十进制字符串。读取文章时原样回传所选文件的 ID。 |
| `type` | `string` | `folder` 为文件夹，`file` 为文件。 |
| `name` | `string` | 节点的 `catalogue.name`。根也使用目录行名称，书单改名后两处名称可能不同。 |
| `intro` | `string` | 仅文件夹返回，取自 `catalogue.intro`。根的简介介绍整本书，子文件夹的简介介绍相应章节。未填写时为 `""`。 |
| `createdAt` | `string` | 仅文件夹返回，取该目录行的 `created_at`，保留数据库日历日期并按 `YYYY-MM-DD` 输出。 |
| `children` | `Array<catalogue>` | 仅文件夹返回，包含按阅读顺序排列的直接公开子节点。没有公开子节点时为 `[]`。 |

## 目录显隐和阅读顺序

接口只读取所请求的公开书籍。节点及其所有祖先均为 `visible` 时才可展示。名称恰为小写 `img` 的文件夹连同整支子树也会被排除。名为 `img` 的文件和名为 `IMG` 的文件夹按普通节点处理。

每个文件夹的公开子项按同级 `sort` 数值升序排列，文件和文件夹共同排序。隐藏节点占用的序号可以造成跳号。若只有根可展示，接口仍返回根的完整字段及 `"children": []`。

## 响应码

以下各项的 HTTP 状态均为 200。调用方通过响应体的 `code` 判断结果，失败时 `data` 为 `{}`。

| 业务码 | 含义 |
|:-----:|:-----|
| `200` | 返回所请求书籍的单个公开根文件夹。 |
| `10101` | 缺少 `book` 或 `id`，或 ID 不是正 `int64` 的十进制字符串，例如裸数字、零、负数、小数或越界值。 |
| `20202` | 书籍不存在或未公开。 |
| `20301` | 书籍已公开但暂无公开目录，包括尚未导入目录、根被隐藏或根为小写 `img` 文件夹。 |
| `500` | 非空目录的结构损坏，例如多根、孤儿或环，隐藏分支也参与校验。 |
| `10201` | 数据库连接、查询或事务失败。 |
