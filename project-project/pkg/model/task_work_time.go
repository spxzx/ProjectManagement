package model

import (
	"github.com/jinzhu/copier"
	"github.com/spxzx/project-common/encrypts"
	"github.com/spxzx/project-common/tms"
	"github.com/spxzx/project-project/pkg/data"
)

type TaskWorkTime struct {
	Id         int64
	TaskCode   int64
	MemberCode int64
	CreateTime int64
	Content    string
	BeginTime  int64
	Num        int
}

func (t *TaskWorkTime) ToDisplay() *TaskWorkTimeDisplay {
	td := &TaskWorkTimeDisplay{}
	_ = copier.Copy(td, t)
	td.MemberCode = encrypts.EncryptNoErr(t.MemberCode, data.AESKey)
	td.TaskCode = encrypts.EncryptNoErr(t.TaskCode, data.AESKey)
	td.CreateTime = tms.FormatByMill(t.CreateTime)
	td.BeginTime = tms.FormatByMill(t.BeginTime)
	return td
}

type TaskWorkTimeDisplay struct {
	Id         int64
	TaskCode   string
	MemberCode string
	CreateTime string
	Content    string
	BeginTime  string
	Num        int
	Member     Member
}
