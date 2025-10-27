// Package models 定义通用响应结构
package models

// Response 通用API响应结构
type Response struct {
	Code    int         `json:"code"`           // 状态码：0=成功，其他=失败
	Message string      `json:"message"`        // 提示信息
	Data    interface{} `json:"data,omitempty"` // 响应数据（可选）
}

// 状态码常量
const (
	CodeSuccess         = 0    // 成功
	CodeInvalidParams   = 1001 // 参数错误
	CodeServerError     = 1002 // 服务器错误
	CodeUnauthorized    = 1003 // 未授权
	CodeNotFound        = 1004 // 资源不存在
	CodeAlreadyExists   = 1005 // 资源已存在
	CodeTooManyRequests = 1006 // 请求过于频繁
)

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(data interface{}) Response {
	return Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}
