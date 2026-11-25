package resp_body

import "SolarNotesAPI/sys_error"

// RespBody 本结构体用于定义响应体结构
type RespBody struct {
	Code    int                    `json:"code"`    // Code 响应体中的状态码部分
	Message string                 `json:"message"` // Message 响应体中的消息部分
	Data    map[string]interface{} `json:"data"`    // Data 响应体中的数据载荷部分
}

// 响应码规则:
// 200: 成功响应
// 10XXX: 系统级别错误
// 101XX: 参数校验错误
const (
	ServiceInternalErrorCode = 500
	DBErrorCode              = 10201
	SUCCESS                  = 200   // SUCCESS 成功响应
	ValidateFailed           = 10101 // ValidateFailed 参数校验失败
)

// Success 本方法用于定义成功响应时的响应体
func (r *RespBody) Success(data map[string]interface{}) {
	r.Code = SUCCESS
	r.Message = "success"
	r.Data = data
}

// ValidateFailed 本方法用于定义参数校验失败时的响应体
func (r *RespBody) ValidateFailed(err error) {
	r.Code = ValidateFailed
	r.Message = err.Error()
	r.Data = map[string]interface{}{}
}

// DBError 和 InternalError 仅输出公开消息。
func (r *RespBody) DBError() {
	r.failure(DBErrorCode, sys_error.DatabaseMessage)
}

func (r *RespBody) InternalError() {
	r.failure(ServiceInternalErrorCode, sys_error.InternalMessage)
}

func (r *RespBody) failure(code int, message string) {
	r.Code = code
	r.Message = message
	r.Data = map[string]any{}
}
