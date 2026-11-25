package sys_error

import "errors"

var (
	BookNotFoundErr      = errors.New("书籍不存在或未公开")
	CatalogueNotFoundErr = errors.New("这本书暂无公开目录")
)
