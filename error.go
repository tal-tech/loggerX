package logger

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/cast"
	"github.com/tal-tech/loggerX/logutils"
	"github.com/tal-tech/loggerX/stackerr"
)

type XesError struct {
	Code int
	Msg  string
}

var (
	//参数校验
	PARAM_MISSING     XesError = XesError{10001, "参数校验缺失"}
	PARAM_ERROR       XesError = XesError{10002, "参数校验错误"}
	PARAM_USER_MISSIG XesError = XesError{10101, "用户名缺失"}
	PARAM_USER_ERROR  XesError = XesError{10102, "用户名错误"}
	PARAM_MOBILEPHONE XesError = XesError{10200, "手机号错误"}
	PARAM_TELEPHONE   XesError = XesError{10300, "电话错误"}
	PARAM_EMAIL       XesError = XesError{10400, "邮箱错误"}

	//登录校验
	LOGIN_NOTLOGGEDIN         XesError = XesError{20000, "未登录"}
	LOGIN_SESSIONTIMEOUT      XesError = XesError{20100, "会话超时"}
	LOGIN_KICKED              XesError = XesError{20200, "已被踢"}
	LOGIN_PASSWORDMODIFIED    XesError = XesError{20300, "密码被修改"}
	LOGIN_NAMEMODIFIED        XesError = XesError{20400, "登录名被修改"}
	LOGIN_MOBILEPHONEMODIFIED XesError = XesError{20500, "手机号被修改"}

	//版本检测
	VERSION_NOTSUPPORT_CLOSE   XesError = XesError{30100, "版本不支持"}
	VERSION_NOTSUPPORT_RETURN  XesError = XesError{30200, "版本不支持"}
	VERSION_NOTSUPPORT_UPGRADE XesError = XesError{30300, "版本不支持"}

	//权限控制
	PERMISSION_VIEW   XesError = XesError{40100, "无权限查看"}
	PERMISSION_MODIFY XesError = XesError{40200, "无权限修改"}
	PERMISSION_ADD    XesError = XesError{40300, "无权限增加"}
	PERMISSION_DELETE XesError = XesError{40400, "无权限删除"}

	//系统异常
	SYSTEM_DEFAULT       XesError = XesError{50000, "系统异常"}
	SYSTEM_NOTSUPPORT    XesError = XesError{50100, "系统未支持"}
	SYSTEM_CONNECT_API   XesError = XesError{50201, "系统连接异常"}
	SYSTEM_CONNECT_MYSQL XesError = XesError{50202, "系统连接异常"}
	SYSTEM_CONNECT_REDIS XesError = XesError{50203, "系统连接异常"}
	SYSTEM_TIMEOUT_API   XesError = XesError{50401, "系统连接超时"}
	SYSTEM_TIMEOUT_MYSQL XesError = XesError{50402, "系统连接超时"}
	SYSTEM_TIMEOUT_REDIS XesError = XesError{50403, "系统连接超时"}
)

/*
* NewError 构造错误
* err 如果err的类型是err或string,将错误信息写入ErrorMessage
* 	   如果err是StackErr,直接返回
* ext ext[0]:错误XesError
* ext ext[0]:错误code  ext[1]:返回给调用端的错误信息
 */
func NewError(err interface{}, ext ...XesError) *stackerr.StackErr {
	return newError(err, ext...)
}

func NewErrorWithLevel(err interface{}, lvl int, ext ...XesError) *stackerr.StackErr {
	e := newError(err, ext...)
	e.Level = lvl
	return e
}

func newError(err interface{}, ext ...XesError) *stackerr.StackErr {

	var errInfo string
	switch t := err.(type) {
	case *stackerr.StackErr:
		return t
	case string:
		errInfo = logutils.Filter(t)
	case error:
		errInfo = logutils.Filter(t.Error())
		if xesCode, ok := transfer(t); ok {
			ext = make([]XesError, 1)
			ext[0] = xesCode
		}
	default:
		errInfo = logutils.Filter(fmt.Sprintf("%v", t))
	}
	stackErr := &stackerr.StackErr{}

	stackErr.Info = errInfo
	_, file, line, ok := runtime.Caller(2)
	if ok {
		stackErr.Line = line
		components := strings.Split(file, "/")
		stackErr.Filename = components[(len(components) - 1)]
		stackErr.Position = filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	const size = 1 << 12
	buf := make([]byte, size)
	n := runtime.Stack(buf, false)
	stackErr.StackTrace = logutils.Filter(string(buf[:n]), " ")

	if len(ext) >= 1 {
		c := ext[0]
		stackErr.Code = c.Code
		stackErr.Message = c.Msg
	} else {
		stackErr.Code = SYSTEM_DEFAULT.Code
		stackErr.Message = errInfo
	}

	return stackErr
}

func transfer(err error) (xCode XesError, ok bool) {
	segments := strings.SplitN(err.Error(), "|", 2)
	code := cast.ToInt(segments[0])
	if len(segments) > 1 && code > 0 {
		xCode.Code = code
		xCode.Msg = segments[1]
		ok = true
	}
	return
}
