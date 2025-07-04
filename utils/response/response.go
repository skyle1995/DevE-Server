package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构体
type Response struct {
	Code      int         `json:"code"`      // 状态码：200-成功，非200-失败
	Message   string      `json:"message"`   // 状态描述
	Data      interface{} `json:"data"`      // 响应数据
	Timestamp int64       `json:"timestamp"` // 响应时间戳
	Sign      string      `json:"sign"`      // 签名（当应用开启签名验证时返回）
}

// Result 返回统一响应结构体
func Result(code int, message string, data interface{}, c *gin.Context) {
	// 获取当前时间戳
	timestamp := time.Now().Unix()

	// 构造响应结构体
	response := Response{
		Code:      code,
		Message:   message,
		Data:      data,
		Timestamp: timestamp,
		Sign:      "", // 暂不实现签名
	}

	// 返回JSON响应
	c.JSON(http.StatusOK, response)
}

// Ok 返回成功响应
func Ok(c *gin.Context) {
	Result(200, "success", nil, c)
}

// OkWithMessage 返回带消息的成功响应
func OkWithMessage(message string, c *gin.Context) {
	Result(200, message, nil, c)
}

// OkWithData 返回带数据的成功响应
func OkWithData(data interface{}, c *gin.Context) {
	Result(200, "success", data, c)
}

// OkWithDetailed 返回带详细信息的成功响应
func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(200, message, data, c)
}

// Fail 返回失败响应
func Fail(c *gin.Context) {
	Result(400, "failed", nil, c)
}

// FailWithMessage 返回带消息的失败响应
func FailWithMessage(message string, c *gin.Context) {
	Result(400, message, nil, c)
}

// FailWithDetailed 返回带详细信息的失败响应
func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(400, message, data, c)
}

// FailWithCode 返回带状态码和消息的失败响应
func FailWithCode(code int, message string, c *gin.Context) {
	Result(code, message, nil, c)
}

// PaginationResponse 分页响应结构体
type PaginationResponse struct {
	List     interface{} `json:"list"`      // 数据列表
	Page     int         `json:"page"`      // 当前页码
	PageSize int         `json:"page_size"` // 每页条数
	Total    int64       `json:"total"`     // 总条数
}

// SuccessWithPagination 返回带分页的成功响应
func SuccessWithPagination(c *gin.Context, message string, list interface{}, page, pageSize int, total int64) {
	paginationData := PaginationResponse{
		List:     list,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}
	Result(200, message, paginationData, c)
}

// Success 返回成功响应（带数据和消息）
func Success(c *gin.Context, message string, data interface{}) {
	Result(200, message, data, c)
}

// Unauthorized 返回未授权响应
func Unauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, Response{
		Code:      401,
		Message:   message,
		Data:      nil,
		Timestamp: time.Now().Unix(),
		Sign:      "",
	})
}

// Forbidden 返回禁止访问响应
func Forbidden(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusForbidden, Response{
		Code:      403,
		Message:   message,
		Data:      nil,
		Timestamp: time.Now().Unix(),
		Sign:      "",
	})
}

// InternalServerError 返回服务器内部错误响应
func InternalServerError(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
		Code:      500,
		Message:   message,
		Data:      nil,
		Timestamp: time.Now().Unix(),
		Sign:      "",
	})
}

// BadRequest 返回请求参数错误响应
func BadRequest(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, Response{
		Code:      400,
		Message:   message,
		Data:      nil,
		Timestamp: time.Now().Unix(),
		Sign:      "",
	})
}

// NotFound 返回资源未找到响应
func NotFound(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusNotFound, Response{
		Code:      404,
		Message:   message,
		Data:      nil,
		Timestamp: time.Now().Unix(),
		Sign:      "",
	})
}
