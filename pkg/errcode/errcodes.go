package errcode

const (
	SUCCESS       = 200
	ERROR         = 500
	InvalidParams = 400

	InvalidEmail = 1001
	HasSendCode  = 1002

	PasswordNotMatch = 2001
	CodeError        = 2002
	HasRegister      = 2003
	NotRegister      = 2004
	NullEmail        = 2005
	PasswordError    = 2006
	UserNotLogin     = 2007
	CodeNeedLogin    = 2008
	CodeInvalidToken = 2009
	NoPerson         = 3001
)

var codeMsg = map[int]string{
	SUCCESS:       "success",
	ERROR:         "服务器出错",
	InvalidParams: "请求参数错误",

	InvalidEmail: "邮箱地址错误",
	HasSendCode:  "请勿重复发送",

	PasswordNotMatch: "两次密码输入不一致",
	CodeError:        "验证码错误",
	HasRegister:      "此邮箱已被注册",
	NotRegister:      "注册失败",
	NullEmail:        "邮箱尚未注册",
	PasswordError:    "密码输入错误，请重试",
	UserNotLogin:     "用户不在线",
	CodeNeedLogin:    "需要登录",
	CodeInvalidToken: "无效的token",
	NoPerson:         "查无此人",
}

func Msg(code int) string {
	msg, ok := codeMsg[code]
	if ok {
		return msg
	}
	return codeMsg[ERROR]
}
