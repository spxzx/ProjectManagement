package model

import (
	"github.com/jinzhu/copier"
	"github.com/spxzx/project-common/encrypts"
	"github.com/spxzx/project-common/tms"
	"github.com/spxzx/project-project/pkg/data"
)

type ProjectAuth struct {
	Id               int64  `json:"id"`
	OrganizationCode int64  `json:"organization_code"`
	Title            string `json:"title"`
	CreateAt         int64  `json:"create_at"`
	Sort             int    `json:"sort"`
	Status           int    `json:"status"`
	Desc             string `json:"desc"`
	CreateBy         int64  `json:"create_by"`
	IsDefault        int    `json:"is_default"`
	Type             string `json:"type"`
}

func (a *ProjectAuth) ToDisplay() *ProjectAuthDisplay {
	p := &ProjectAuthDisplay{}
	_ = copier.Copy(p, a)
	p.OrganizationCode = encrypts.EncryptNoErr(a.OrganizationCode, data.AESKey)
	p.CreateAt = tms.FormatByMill(a.CreateAt)
	if a.Type == "admin" || a.Type == "member" {
		p.CanDelete = 0
	} else {
		p.CanDelete = 1
	}
	return p
}

type ProjectAuthDisplay struct {
	Id               int64  `json:"id"`
	OrganizationCode string `json:"organization_code"`
	Title            string `json:"title"`
	CreateAt         string `json:"create_at"`
	Sort             int    `json:"sort"`
	Status           int    `json:"status"`
	Desc             string `json:"desc"`
	CreateBy         int64  `json:"create_by"`
	IsDefault        int    `json:"is_default"`
	Type             string `json:"type"`
	CanDelete        int    `json:"canDelete"`
}
