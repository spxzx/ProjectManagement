package model

import (
	"github.com/jinzhu/copier"
	"github.com/spxzx/project-common/encrypts"
	"github.com/spxzx/project-common/tms"
	"github.com/spxzx/project-project/pkg/data"
)

type SourceLink struct {
	Id               int64
	SourceType       string
	SourceCode       int64
	LinkType         string
	LinkCode         int64
	OrganizationCode int64
	CreateBy         int64
	CreateTime       int64
	Sort             int
}

func (s *SourceLink) ToDisplay(f *File) *SourceLinkDisplay {
	sl := &SourceLinkDisplay{}
	_ = copier.Copy(sl, s)
	sl.SourceDetail = SourceDetail{}
	_ = copier.Copy(&sl.SourceDetail, f)
	sl.LinkCode = encrypts.EncryptNoErr(s.LinkCode, data.AESKey)
	sl.OrganizationCode = encrypts.EncryptNoErr(s.OrganizationCode, data.AESKey)
	sl.CreateTime = tms.FormatByMill(s.CreateTime)
	sl.CreateBy = encrypts.EncryptNoErr(s.CreateBy, data.AESKey)
	sl.SourceCode = encrypts.EncryptNoErr(s.SourceCode, data.AESKey)
	sl.SourceDetail.OrganizationCode = encrypts.EncryptNoErr(f.OrganizationCode, data.AESKey)
	sl.SourceDetail.CreateBy = encrypts.EncryptNoErr(f.CreateBy, data.AESKey)
	sl.SourceDetail.CreateTime = tms.FormatByMill(f.CreateTime)
	sl.SourceDetail.DeletedTime = tms.FormatByMill(f.DeletedTime)
	sl.SourceDetail.TaskCode = encrypts.EncryptNoErr(f.TaskCode, data.AESKey)
	sl.SourceDetail.ProjectCode = encrypts.EncryptNoErr(f.ProjectCode, data.AESKey)
	sl.SourceDetail.FullName = f.Title
	sl.Title = f.Title
	return sl
}

type SourceLinkDisplay struct {
	Id               int64        `json:"id"`
	Code             string       `json:"code"`
	SourceType       string       `json:"source_type"`
	SourceCode       string       `json:"source_code"`
	LinkType         string       `json:"link_type"`
	LinkCode         string       `json:"link_code"`
	OrganizationCode string       `json:"organization_code"`
	CreateBy         string       `json:"create_by"`
	CreateTime       string       `json:"create_time"`
	Sort             int          `json:"sort"`
	Title            string       `json:"title"`
	SourceDetail     SourceDetail `json:"sourceDetail"`
}

type SourceDetail struct {
	Id               int64  `json:"id"`
	Code             string `json:"code"`
	PathName         string `json:"path_name"`
	Title            string `json:"title"`
	Extension        string `json:"extension"`
	Size             int    `json:"size"`
	ObjectType       string `json:"object_type"`
	OrganizationCode string `json:"organization_code"`
	TaskCode         string `json:"task_code"`
	ProjectCode      string `json:"project_code"`
	CreateBy         string `json:"create_by"`
	CreateTime       string `json:"create_time"`
	Downloads        int    `json:"downloads"`
	Extra            string `json:"extra"`
	Deleted          int    `json:"deleted"`
	FileUrl          string `json:"file_url"`
	FileType         string `json:"file_type"`
	DeletedTime      string `json:"deleted_time"`
	ProjectName      string `json:"projectName"`
	FullName         string `json:"fullName"`
}
