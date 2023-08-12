package poj

type Project struct {
	//Id                 int64   `json:"id"`
	Code               string  `json:"code"`
	Cover              string  `json:"cover"`
	Name               string  `json:"name"`
	Description        string  `json:"description"`
	AccessControlType  string  `json:"access_control_type"`
	WhiteList          string  `json:"white_list"`
	Order              int     `json:"order"`
	Deleted            int     `json:"deleted"`
	TemplateCode       string  `json:"template_code"`
	Schedule           float64 `json:"schedule"`
	CreateTime         string  `json:"create_time"`
	OrganizationCode   string  `json:"organization_code"`
	DeletedTime        string  `json:"deleted_time"`
	Private            int     `json:"private"`
	Prefix             string  `json:"prefix"`
	OpenPrefix         int     `json:"open_prefix"`
	Archive            int     `json:"archive"`
	ArchiveTime        int64   `json:"archive_time"`
	OpenBeginTime      int     `json:"open_begin_time"`
	OpenTaskPrivate    int     `json:"open_task_private"`
	TaskBoardTheme     string  `json:"task_board_theme"`
	BeginTime          int64   `json:"begin_time"`
	EndTime            int64   `json:"end_time"`
	AutoUpdateSchedule int     `json:"auto_update_schedule"`
}

type ProjectMemberUnion struct {
	Project
	ProjectCode int64  `json:"project_code"`
	MemberCode  int64  `json:"member_code"`
	JoinTime    string `json:"join_time"`
	IsOwner     int64  `json:"is_owner"`
	Authorize   string `json:"authorize"`
	OwnerName   string `json:"owner_name"`
	Collected   int    `json:"collected"`
}

type SaveProjectReq struct {
	Name         string `json:"name" form:"name"`
	TemplateCode string `json:"templateCode" form:"templateCode"`
	Description  string `json:"description" form:"description"`
	Id           int64  `json:"id" form:"id"`
}

type SaveProjectResp struct {
	Id               int64  `json:"id"`
	Cover            string `json:"cover"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Code             string `json:"code"`
	CreateTime       string `json:"create_time"`
	TaskBoardTheme   string `json:"task_board_theme"`
	OrganizationCode string `json:"organization_code"`
}

type ReadProjectResp struct {
	Project
	OwnerName   string `json:"owner_name"`
	Collected   int    `json:"collected"`
	OwnerAvatar string `json:"owner_avatar"`
}

type EditProjectReq struct {
	ProjectCode         string `json:"projectCode" form:"projectCode"`
	Name                string `json:"name" form:"name"`
	Description         string `json:"description" form:"description"`
	Cover               string `json:"cover" form:"cover"`
	Private             int    `json:"private" form:"private"`
	Prefix              string `json:"prefix" form:"prefix"`
	TaskBoardTheme      string `json:"task_board_theme" form:"task_board_theme"`
	OpenPrefix          int    `json:"open_prefix" form:"open_prefix"`
	OpenBeginTime       int64  `json:"open_begin_time" form:"open_begin_time"`
	OpenTaskPrivate     int    `json:"open_task_private" form:"open_task_private"`
	Schedule            int    `json:"schedule" form:"schedule"`
	AutoUpdateScheduler int    `json:"auto_update_scheduler" form:"auto_update_scheduler"`
}
