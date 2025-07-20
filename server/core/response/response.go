package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RespType 响应类型
type RespType struct {
	code int
	msg  string
	data interface{}
}

// Response 响应格式结构
type Response struct {
	Code    int         `json:"code"`    // 状态码
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 数据
}

var (
	Success = RespType{code: 0, msg: "操作成功"}
	Failed  = RespType{code: 1, msg: "操作失败"}

	// 参数相关错误
	ParamsValidError    = RespType{code: 400, msg: "参数校验错误"}
	ParamsTypeError     = RespType{code: 400, msg: "参数类型错误"}
	RequestMethodError  = RespType{code: 405, msg: "请求方法错误"}
	AssertArgumentError = RespType{code: 400, msg: "断言参数错误"}

	// 认证相关错误
	LoginAccountError = RespType{code: 401, msg: "账号或密码错误"}
	LoginDisableError = RespType{code: 403, msg: "账号已被禁用"}
	TokenEmpty        = RespType{code: 401, msg: "token不能为空"}
	TokenInvalid      = RespType{code: 401, msg: "token无效或已过期"}
	TokenExpired      = RespType{code: 401, msg: "token已过期"}

	// 权限相关错误
	NoPermission    = RespType{code: 403, msg: "无权限访问"}
	Request404Error = RespType{code: 404, msg: "请求资源不存在"}
	Request405Error = RespType{code: 405, msg: "请求方法不允许"}

	// 系统错误
	SystemError = RespType{code: 500, msg: "系统内部错误"}
)

// Error 实现error方法
func (rt RespType) Error() string {
	return rt.msg
}

// Make 以响应类型生成信息
func (rt RespType) Make(msg string) RespType {
	rt.msg = msg
	return rt
}

// MakeData 以响应类型生成数据
func (rt RespType) MakeData(data interface{}) RespType {
	rt.data = data
	return rt
}

// Code 获取code
func (rt RespType) Code() int {
	return rt.code
}

// Msg 获取msg
func (rt RespType) Msg() string {
	return rt.msg
}

// Data 获取data
func (rt RespType) Data() interface{} {
	return rt.data
}

// Result 统一响应
func Result(c *gin.Context, resp RespType, data interface{}) {
	if data == nil {
		data = resp.data
	}
	c.JSON(http.StatusOK, Response{
		Code:    resp.code,
		Message: resp.msg,
		Data:    data,
	})
}

// Ok 成功响应
func Ok(c *gin.Context) {
	Result(c, Success, map[string]bool{"success": true})
}

// OkWithMsg 成功响应附带msg
func OkWithMsg(c *gin.Context, msg string) {
	resp := Success
	resp.msg = msg
	Result(c, resp, map[string]bool{"success": true})
}

// OkWithData 成功响应附带data
func OkWithData(c *gin.Context, data interface{}) {
	Result(c, Success, data)
}

// Fail 失败响应
func Fail(c *gin.Context, resp RespType) {
	Result(c, resp, nil)
}

// FailWithMsg 失败响应附带msg
func FailWithMsg(c *gin.Context, resp RespType, msg string) {
	resp.msg = msg
	Result(c, resp, nil)
}

// FailWithData 失败响应附带data
func FailWithData(c *gin.Context, resp RespType, data interface{}) {
	Result(c, resp, data)
}

// NoAuth 无权限响应
func NoAuth(c *gin.Context, msg string) {
	Result(c, NoPermission.Make(msg), nil)
}

// CheckAndResp 判断是否出现错误，并返回对应响应
func CheckAndResp(c *gin.Context, err error) {
	if err != nil {
		switch v := err.(type) {
		case RespType:
			data := v.Data()
			if data == nil {
				data = map[string]bool{"success": false}
			}
			FailWithData(c, v, data)
		default:
			Fail(c, SystemError)
		}
		return
	}
	Ok(c)
}

// CheckAndRespWithData 判断是否出现错误，并返回对应响应（带data数据）
func CheckAndRespWithData(c *gin.Context, data interface{}, err error) {
	if err != nil {
		switch v := err.(type) {
		case RespType:
			respData := v.Data()
			if respData == nil {
				respData = map[string]bool{"success": false}
			}
			FailWithData(c, v, respData)
		default:
			Fail(c, SystemError)
		}
		return
	}
	OkWithData(c, data)
}
