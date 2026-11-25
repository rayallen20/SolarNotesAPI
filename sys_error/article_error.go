package sys_error

import "errors"

var ArticleUnavailableErr = errors.New("文章不存在或暂不可读")
