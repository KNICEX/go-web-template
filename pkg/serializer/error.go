package serializer

import (
	"github.com/gin-gonic/gin"
	"myBulebell/pkg/logger"
)

type Error struct {
	Code     ResCode
	Msg      string
	RawError error
}

func (e Error) Error() string {
	return e.Msg
}

func (e Error) WithError(err error) Error {
	e.RawError = err
	return e
}

type ResCode int

const (
	CodeSuccess      ResCode = 0
	CodeParamErr             = 400
	CodeUnauthorized         = 401
	CodeForbidden            = 403
	CodeServerErr            = 500
)

const (
	CodeEmailExist ResCode = iota + 1000
	CodeEncryptError
	CodeInvalidSign
	CodeSignTimeout
	CodeNoPermission
	CodeEmailSendErr
	CodeCredentialInvalid
	CodeServerBusy
	CodePasswordErr
	CodeUserBaned
	CodeUserNotActivated
)

var codeMsgMap = map[ResCode]string{
	CodeParamErr:     "请求参数错误",
	CodeServerErr:    "服务器内部错误",
	CodeSuccess:      "成功",
	CodeUnauthorized: "未登录",
	CodeForbidden:    "无权限",

	CodeEmailExist:        "邮箱已被注册",
	CodeEncryptError:      "加密错误",
	CodeInvalidSign:       "无效签名",
	CodeSignTimeout:       "签名超时",
	CodeNoPermission:      "无权限",
	CodeEmailSendErr:      "邮件发送失败",
	CodeCredentialInvalid: "凭证无效",
	CodeServerBusy:        "服务器繁忙",
	CodePasswordErr:       "用户名或密码错误",
	CodeUserBaned:         "用户被封禁",
	CodeUserNotActivated:  "用户未激活",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		return codeMsgMap[CodeServerBusy]
	}
	return msg
}

func NewError(code ResCode, msg string, err error) Error {
	return Error{
		Code:     code,
		Msg:      msg,
		RawError: err,
	}
}

func ErrResponse(code ResCode, err error) Response {
	return Err(code, code.Msg(), err)
}

func ErrWithMsg(code ResCode, msg string, err error) Response {
	return Err(code, msg, err)
}

func Err(code ResCode, msg string, err error) Response {
	res := Response{
		Code: code,
		Msg:  msg,
	}

	if err != nil && gin.Mode() != gin.ReleaseMode {
		res.Error = err.Error()
		logger.L().Error(err)
	}
	return res
}
