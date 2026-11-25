package sys_error

const DatabaseMessage = "数据库操作失败,请稍后重试"

// DatabaseError 保留查询错误的原始原因，Error 返回对外消息。
type DatabaseError struct {
	Cause error
}

func (e *DatabaseError) Error() string { return DatabaseMessage }

func (e *DatabaseError) Unwrap() error { return e.Cause }
