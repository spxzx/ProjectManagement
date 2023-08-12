package model

import (
	"github.com/jinzhu/copier"
	"github.com/spxzx/project-common/encrypts"
	"github.com/spxzx/project-common/tms"
	"github.com/spxzx/project-project/pkg/data"
)

type Department struct {
	Id               int64
	OrganizationCode int64
	Name             string
	Sort             int
	Pcode            int64
	icon             string
	CreateTime       int64
	Path             string
}

func (d *Department) ToDisplay() *DepartmentDisplay {
	dp := &DepartmentDisplay{}
	_ = copier.Copy(dp, d)
	dp.CreateTime = tms.FormatByMill(d.CreateTime)
	dp.OrganizationCode = encrypts.EncryptNoErr(d.OrganizationCode, data.AESKey)
	if d.Pcode > 0 {
		dp.Pcode = encrypts.EncryptNoErr(d.Pcode, data.AESKey)
	} else {
		dp.Pcode = ""
	}
	return dp
}

type DepartmentDisplay struct {
	Id               int64
	OrganizationCode string
	Name             string
	Sort             int
	Pcode            string
	icon             string
	CreateTime       string
	Path             string
}
