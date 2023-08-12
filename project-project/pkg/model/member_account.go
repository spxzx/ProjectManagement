package model

import (
	"github.com/jinzhu/copier"
	"github.com/spxzx/project-common/encrypts"
	"github.com/spxzx/project-common/tms"
	"github.com/spxzx/project-project/pkg/data"
)

type MemberAccount struct {
	Id               int64
	OrganizationCode int64
	DepartmentCode   int64
	MemberCode       int64
	Authorize        string
	IsOwner          int
	Name             string
	Mobile           string
	Email            string
	CreateTime       int64
	LastLoginTime    int64
	Status           int
	Description      string
	Avatar           string
	Position         string
	Department       string
}

func (a *MemberAccount) ToDisplay() *MemberAccountDisplay {
	md := &MemberAccountDisplay{}
	_ = copier.Copy(md, a)
	md.Code = encrypts.EncryptNoErr(a.Id, data.AESKey)
	md.MemberAccountCode = md.Code
	md.MemberCode = encrypts.EncryptNoErr(a.MemberCode, data.AESKey)
	md.OrganizationCode = encrypts.EncryptNoErr(a.OrganizationCode, data.AESKey)
	md.DepartmentCode = encrypts.EncryptNoErr(a.DepartmentCode, data.AESKey)
	md.CreateTime = tms.FormatByMill(a.CreateTime)
	md.LastLoginTime = tms.FormatByMill(a.LastLoginTime)
	md.StatusText = a.StatusText()
	md.AuthorizeArr = []string{a.Authorize}
	return md
}

func (a *MemberAccount) StatusText() string {
	if a.Status == 0 {
		return "禁用"
	}
	if a.Status == 1 {
		return "使用中"
	}
	return ""
}

type MemberAccountDisplay struct {
	Id                int64  `json:"id"`
	Code              string `json:"code"`
	MemberCode        string
	OrganizationCode  string   `json:"organization_code"`
	DepartmentCode    string   `json:"department_code"`
	Authorize         string   `json:"authorize"`
	IsOwner           int      `json:"is_owner"`
	Name              string   `json:"name"`
	Mobile            string   `json:"mobile"`
	Email             string   `json:"email"`
	CreateTime        string   `json:"create_time"`
	LastLoginTime     string   `json:"last_login_time"`
	Status            int      `json:"status"`
	Description       string   `json:"description"`
	Avatar            string   `json:"avatar"`
	Position          string   `json:"position"`
	Department        string   `json:"department"`
	MemberAccountCode string   `json:"member_account_code"`
	Departments       string   `json:"departments"`
	StatusText        string   `json:"statusText"`
	AuthorizeArr      []string `json:"authorizeArr"`
}
