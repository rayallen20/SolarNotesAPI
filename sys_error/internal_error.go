package sys_error

const InternalMessage = "服务暂时不可用,请稍后重试"

type InternalError struct {
	Cause error
}

func (e *InternalError) Error() string { return InternalMessage }

func (e *InternalError) Unwrap() error { return e.Cause }
