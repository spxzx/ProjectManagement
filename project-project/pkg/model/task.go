package model

import (
	"github.com/jinzhu/copier"
	"github.com/spxzx/project-common/encrypts"
	"github.com/spxzx/project-common/tms"
	"github.com/spxzx/project-project/pkg/data"
)

const (
	Wait = iota
	Doing
	Done
	Pause
	Cancel
	Closed
)
const (
	NoStarted = iota
	Started
)
const (
	Normal = iota
	Urgent
	VeryUrgent
)

type Task struct {
	Id            int64
	ProjectCode   int64
	Name          string
	Pri           int
	ExecuteStatus int
	Description   string
	CreateBy      int64
	DoneBy        int64
	DoneTime      int64
	CreateTime    int64
	AssignTo      int64
	Deleted       int
	StageCode     int
	TaskTag       string
	Done          int
	BeginTime     int64
	EndTime       int64
	RemindTime    int64
	Pcode         int64
	Sort          int
	Like          int
	Star          int
	DeletedTime   int64
	Private       int
	IdNum         int
	Path          string
	Schedule      int
	VersionCode   int64
	FeaturesCode  int64
	WorkTime      int
	Status        int
}

func (t *Task) GetExecuteStatusStr() string {
	switch t.ExecuteStatus {
	case Wait:
		return "wait"
	case Doing:
		return "doing"
	case Done:
		return "done"
	case Pause:
		return "pause"
	case Cancel:
		return "cancel"
	case Closed:
		return "closed"
	default:
		return ""
	}
}

func (t *Task) GetStatusStr() string {
	status := t.Status
	if status == NoStarted {
		return "未开始"
	}
	if status == Started {
		return "开始"
	}
	return ""
}

func (t *Task) GetPriStr() string {
	status := t.Pri
	if status == Normal {
		return "普通"
	}
	if status == Urgent {
		return "紧急"
	}
	if status == VeryUrgent {
		return "非常紧急"
	}
	return ""
}

type Executor struct {
	Name   string
	Avatar string
	Code   string
}

type TaskDisplay struct {
	Id            int64
	ProjectCode   string
	Name          string
	Pri           int
	ExecuteStatus string
	Description   string
	CreateBy      string
	DoneBy        string
	DoneTime      string
	CreateTime    string
	AssignTo      string
	Deleted       int
	StageCode     string
	TaskTag       string
	Done          int
	BeginTime     string
	EndTime       string
	RemindTime    string
	Pcode         string
	Sort          int
	Like          int
	Star          int
	DeletedTime   string
	Private       int
	IdNum         int
	Path          string
	Schedule      int
	VersionCode   string
	FeaturesCode  string
	WorkTime      int
	Status        int
	Code          string
	CanRead       int
	Executor      Executor
	ProjectName   string
	StageName     string
	PriText       string
	StatusText    string
}

func (t *Task) ToTaskDisplay() *TaskDisplay {
	td := &TaskDisplay{}
	_ = copier.Copy(td, t)
	td.CreateTime = tms.FormatByMill(t.CreateTime)
	td.DoneTime = tms.FormatByMill(t.DoneTime)
	td.BeginTime = tms.FormatByMill(t.BeginTime)
	td.EndTime = tms.FormatByMill(t.EndTime)
	td.RemindTime = tms.FormatByMill(t.RemindTime)
	td.DeletedTime = tms.FormatByMill(t.DeletedTime)
	td.CreateBy = encrypts.EncryptNoErr(t.CreateBy, data.AESKey)
	td.ProjectCode = encrypts.EncryptNoErr(t.ProjectCode, data.AESKey)
	td.DoneBy = encrypts.EncryptNoErr(t.DoneBy, data.AESKey)
	td.AssignTo = encrypts.EncryptNoErr(t.AssignTo, data.AESKey)
	td.StageCode = encrypts.EncryptNoErr(int64(t.StageCode), data.AESKey)
	td.Pcode = encrypts.EncryptNoErr(t.Pcode, data.AESKey)
	td.VersionCode = encrypts.EncryptNoErr(t.VersionCode, data.AESKey)
	td.FeaturesCode = encrypts.EncryptNoErr(t.FeaturesCode, data.AESKey)
	td.ExecuteStatus = t.GetExecuteStatusStr()
	td.Code = encrypts.EncryptNoErr(t.Id, data.AESKey)
	td.CanRead = 1
	td.StatusText = t.GetStatusStr()
	td.PriText = t.GetPriStr()
	return td
}

type SelfTaskDisplay struct {
	Id                 int64
	ProjectCode        string
	Name               string
	Pri                int
	ExecuteStatus      string
	Description        string
	CreateBy           string
	DoneBy             string
	DoneTime           string
	CreateTime         string
	AssignTo           string
	Deleted            int
	StageCode          string
	TaskTag            string
	Done               int
	BeginTime          string
	EndTime            string
	RemindTime         string
	Pcode              string
	Sort               int
	Like               int
	Star               int
	DeletedTime        string
	Private            int
	IdNum              int
	Path               string
	Schedule           int
	VersionCode        string
	FeaturesCode       string
	WorkTime           int
	Status             int
	Code               string
	Cover              string `json:"cover"`
	AccessControlType  string `json:"access_control_type"`
	WhiteList          string `json:"white_list"`
	Order              int    `json:"order"`
	TemplateCode       string `json:"template_code"`
	OrganizationCode   string `json:"organization_code"`
	Prefix             string `json:"prefix"`
	OpenPrefix         int    `json:"open_prefix"`
	Archive            int    `json:"archive"`
	ArchiveTime        string `json:"archive_time"`
	OpenBeginTime      int    `json:"open_begin_time"`
	OpenTaskPrivate    int    `json:"open_task_private"`
	TaskBoardTheme     string `json:"task_board_theme"`
	AutoUpdateSchedule int    `json:"auto_update_schedule"`
	HasUnDone          int    `json:"hasUnDone"`
	ParentDone         int    `json:"parentDone"`
	PriText            string `json:"priText"`
	ProjectName        string
	Executor           *Executor
}

func (t *Task) ToSelfTaskDisplay(p *Project, name string, avatar string) *SelfTaskDisplay {
	td := &SelfTaskDisplay{}
	copier.Copy(td, p)
	copier.Copy(td, t)
	td.Executor = &Executor{
		Name:   name,
		Avatar: avatar,
	}
	td.ProjectName = p.Name
	td.CreateTime = tms.FormatByMill(t.CreateTime)
	td.DoneTime = tms.FormatByMill(t.DoneTime)
	td.BeginTime = tms.FormatByMill(t.BeginTime)
	td.EndTime = tms.FormatByMill(t.EndTime)
	td.RemindTime = tms.FormatByMill(t.RemindTime)
	td.DeletedTime = tms.FormatByMill(t.DeletedTime)
	td.CreateBy = encrypts.EncryptNoErr(t.CreateBy, data.AESKey)
	td.ProjectCode = encrypts.EncryptNoErr(t.ProjectCode, data.AESKey)
	td.DoneBy = encrypts.EncryptNoErr(t.DoneBy, data.AESKey)
	td.AssignTo = encrypts.EncryptNoErr(t.AssignTo, data.AESKey)
	td.StageCode = encrypts.EncryptNoErr(int64(t.StageCode), data.AESKey)
	td.Pcode = encrypts.EncryptNoErr(t.Pcode, data.AESKey)
	td.VersionCode = encrypts.EncryptNoErr(t.VersionCode, data.AESKey)
	td.FeaturesCode = encrypts.EncryptNoErr(t.FeaturesCode, data.AESKey)
	td.ExecuteStatus = t.GetExecuteStatusStr()
	td.Code = encrypts.EncryptNoErr(t.Id, data.AESKey)
	td.AccessControlType = p.GetAccessControlType()
	td.ArchiveTime = tms.FormatByMill(p.ArchiveTime)
	td.TemplateCode = encrypts.EncryptNoErr(int64(p.TemplateCode), data.AESKey)
	td.OrganizationCode = encrypts.EncryptNoErr(p.OrganizationCode, data.AESKey)
	return td
}
