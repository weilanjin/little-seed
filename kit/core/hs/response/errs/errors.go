package errs

import (
	"kit/core/hs/response/codes"
	"kit/core/hs/response/status"
)

var (
	BadRequest    = status.New(codes.BadRequest, "请求参数错误")
	InternalError = status.New(codes.InternalServerError, "服务器繁忙")
)
