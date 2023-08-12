package member

import (
	"errors"
	"github.com/spxzx/project-api/pkg/model/organization"
	common "github.com/spxzx/project-common"
)

type RegisterReq struct {
	Email     string `json:"email" form:"email"`
	Name      string `json:"name" form:"name"`
	Password  string `json:"password" form:"password"`
	Password2 string `json:"password2" form:"password2"`
	Mobile    string `json:"mobile" form:"mobile"`
	Captcha   string `json:"captcha" form:"captcha"`
}

func (r *RegisterReq) verifyPassword() bool {
	return r.Password == r.Password2
}

func (r *RegisterReq) Verify() error {
	if !common.VerifyEmail(r.Email) {
		return errors.New("邮箱格式不正确")
	}
	if !common.VerifyMobile(r.Mobile) {
		return errors.New("手机号格式不正确")
	}
	if !r.verifyPassword() {
		return errors.New("两次输入密码不一致")
	}
	return nil
}

type Member struct {
	Id               int64  `json:"id"`
	Code             string `json:"code"`
	Name             string `json:"name"`
	Mobile           string `json:"mobile"`
	Status           int    `json:"status"`
	Email            string `json:"email"`
	CreateTime       string `json:"create_time"`
	LastLoginTime    string `json:"last_login_time"`
	OrganizationCode string `json:"organization_code"`
}

type Token struct {
	AccessToken    string `json:"accessToken"`
	RefreshToken   string `json:"refreshToken"`
	TokenType      string `json:"tokenType"`
	AccessTokenExp int64  `json:"accessTokenExp"`
}

type LoginReq struct {
	Account  string `json:"account" form:"account"`
	Password string `json:"password" form:"password"`
}

type LoginResp struct {
	Member           Member                      `json:"member"`
	TokenList        Token                       `json:"tokenList"`
	OrganizationList []organization.Organization `json:"organizationList"`
}
