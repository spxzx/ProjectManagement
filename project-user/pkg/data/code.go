package data

import "github.com/spxzx/project-common/errs"

var (
	RedisError             = errs.NewError(1001, "redis错误")
	DBError                = errs.NewError(1002, "数据库错误")
	NoLegalMobile          = errs.NewError(2001, "手机号不合法")
	CaptchaNotExist        = errs.NewError(2002, "验证码不存在，请尝试重新发送验证码")
	CaptchaError           = errs.NewError(2003, "验证码错误")
	EmailExist             = errs.NewError(2004, "该邮箱已被注册")
	AccountExist           = errs.NewError(2005, "该账号已被注册")
	MobileExist            = errs.NewError(2006, "该手机号已被注册")
	AccountOrPasswordError = errs.NewError(2007, "账号或密码错误")
	NoLogin                = errs.NewError(2008, "未登录")
)
