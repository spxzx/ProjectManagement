package model

import (
	"github.com/spxzx/project-common/encrypts"
	"github.com/spxzx/project-common/tms"
	"github.com/spxzx/project-project/pkg/data"
)

type ProjectTemplate struct {
	Id               int
	Name             string
	Description      string
	Sort             int
	CreateTime       int64
	OrganizationCode int64
	Cover            string
	MemberCode       int64
	IsSystem         int
}

type PojTmplDetail struct {
	Id               int
	Code             string
	Name             string
	Description      string
	Sort             int
	CreateTime       string
	OrganizationCode string
	Cover            string
	MemberCode       string
	IsSystem         int
	TaskStages       []*TaskStagesTmplNames
}

func (p *ProjectTemplate) Combine(task []*TaskStagesTmplNames) *PojTmplDetail {
	orgCode, _ := encrypts.EncryptInt64(p.OrganizationCode, data.AESKey)
	memCode, _ := encrypts.EncryptInt64(p.MemberCode, data.AESKey)
	code, _ := encrypts.EncryptInt64(int64(p.Id), data.AESKey)
	return &PojTmplDetail{
		Id:               p.Id,
		Code:             code,
		Name:             p.Name,
		Description:      p.Description,
		Sort:             p.Sort,
		CreateTime:       tms.FormatByMill(p.CreateTime),
		OrganizationCode: orgCode,
		Cover:            p.Cover,
		MemberCode:       memCode,
		IsSystem:         p.IsSystem,
		TaskStages:       task,
	}
}

func ToPojTmplIds(pt []ProjectTemplate) (ids []int) {
	for _, v := range pt {
		ids = append(ids, v.Id)
	}
	return
}
