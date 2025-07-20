package response

import (
	"net/http"
	"strconv"

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
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var (
	Success = RespType{code: 200, msg: "成功"}
	Failed  = RespType{code: 300, msg: "失败"}

	ParamsValidError    = RespType{code: 310, msg: "参数校验错误"}
	ParamsTypeError     = RespType{code: 311, msg: "参数类型错误"}
	RequestMethodError  = RespType{code: 312, msg: "请求方法错误"}
	AssertArgumentError = RespType{code: 313, msg: "断言参数错误"}

	LoginAccountError = RespType{code: 330, msg: "登录账号或密码错误"}
	LoginDisableError = RespType{code: 331, msg: "登录账号已被禁用了"}
	TokenEmpty        = RespType{code: 332, msg: "token参数为空"}
	TokenInvalid      = RespType{code: 333, msg: "token参数无效"}
	TokenExpired      = RespType{code: 334, msg: "token已过期"}

	NoPermission    = RespType{code: 403, msg: "无相关权限"}
	Request404Error = RespType{code: 404, msg: "请求接口不存在"}
	Request405Error = RespType{code: 405, msg: "请求方法不允许"}

	SystemError = RespType{code: 500, msg: "系统错误"}
)

// Error 实现error方法
func (rt RespType) Error() string {
	return strconv.Itoa(rt.code) + ":" + rt.msg
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
		Code: resp.code,
		Msg:  resp.msg,
		Data: data,
	})
}

// Ok 正常响应
func Ok(c *gin.Context) {
	Result(c, Success, []string{})
}

// OkWithMsg 正常响应附带msg
func OkWithMsg(c *gin.Context, msg string) {
	resp := Success
	resp.msg = msg
	Result(c, resp, []string{})
}

// OkWithData 正常响应附带data
func OkWithData(c *gin.Context, data interface{}) {
	Result(c, Success, data)
}

// Fail 错误响应
func Fail(c *gin.Context, resp RespType) {
	Result(c, resp, []string{})
}

// FailWithMsg 错误响应附带msg
func FailWithMsg(c *gin.Context, resp RespType, msg string) {
	resp.msg = msg
	Result(c, resp, []string{})
}

// FailWithData 错误响应附带data
func FailWithData(c *gin.Context, resp RespType, data interface{}) {
	Result(c, resp, data)
}

// NoAuth 无权限响应
func NoAuth(c *gin.Context, msg string) {
	Result(c, NoPermission.Make(msg), []string{})
}

// CheckAndResp 判断是否出现错误，并返回对应响应
func CheckAndResp(c *gin.Context, err error) {
	if err != nil {
		switch v := err.(type) {
		case RespType:
			data := v.Data()
			if data == nil {
				data = []string{}
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
				respData = []string{}
			}
			FailWithData(c, v, respData)
		default:
			Fail(c, SystemError)
		}
		return
	}
	OkWithData(c, data)
}
