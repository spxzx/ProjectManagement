package model

import (
	"github.com/jinzhu/copier"
	"github.com/spxzx/project-common/encrypts"
	"github.com/spxzx/project-common/tms"
	"github.com/spxzx/project-project/pkg/data"
)

type ProjectLog struct {
	Id           int64
	MemberCode   int64
	Content      string
	Remark       string
	Type         string
	CreateTime   int64
	SourceCode   int64
	ActionType   string
	ToMemberCode int64
	IsComment    int
	ProjectCode  int64
	Icon         string
	IsRobot      int
}

func (l *ProjectLog) ToDisplay() *ProjectLogDisplay {
	pld := &ProjectLogDisplay{}
	_ = copier.Copy(pld, l)
	pld.MemberCode = encrypts.EncryptNoErr(l.MemberCode, data.AESKey)
	pld.ToMemberCode = encrypts.EncryptNoErr(l.ToMemberCode, data.AESKey)
	pld.ProjectCode = encrypts.EncryptNoErr(l.ProjectCode, data.AESKey)
	pld.CreateTime = tms.FormatByMill(l.CreateTime)
	pld.SourceCode = encrypts.EncryptNoErr(l.SourceCode, data.AESKey)
	return pld
}

func (l *ProjectLog) ToIndexDisplay() *IndexProjectLogDisplay {
	pd := &IndexProjectLogDisplay{}
	_ = copier.Copy(pd, l)
	pd.ProjectCode = encrypts.EncryptNoErr(l.ProjectCode, data.AESKey)
	pd.CreateTime = tms.FormatByMill(l.CreateTime)
	pd.SourceCode = encrypts.EncryptNoErr(l.SourceCode, data.AESKey)
	return pd
}

type Member struct {
	Id     int64
	Name   string
	Avatar string
	Code   string
}

type ProjectLogDisplay struct {
	Id           int64
	MemberCode   string
	Content      string
	Remark       string
	Type         string
	CreateTime   string
	SourceCode   string
	ActionType   string
	ToMemberCode string
	IsComment    int
	ProjectCode  string
	Icon         string
	IsRobot      int
	Member       Member
}

type IndexProjectLogDisplay struct {
	Content      string
	Remark       string
	CreateTime   string
	SourceCode   string
	IsComment    int
	ProjectCode  string
	MemberAvatar string
	MemberName   string
	ProjectName  string
	TaskName     string
}
