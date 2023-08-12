package tasks

type TaskStagesResp struct {
	Name         string `json:"name"`
	ProjectCode  string `json:"project_code"`
	Sort         int    `json:"sort"`
	Description  string `json:"description"`
	CreateTime   string `json:"create_time"`
	Code         string `json:"code"`
	Deleted      int    `json:"deleted"`
	TasksLoading bool   `json:"tasksLoading"`
	FixedCreator bool   `json:"fixedCreator"`
	ShowTaskCard bool   `json:"showTaskCard"`
	Tasks        []int  `json:"tasks"`
	DoneTasks    []int  `json:"doneTasks"`
	UnDoneTasks  []int  `json:"unDoneTasks"`
}

type ProjectMemberResp struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Avatar  string `json:"avatar"`
	Code    string `json:"code"`
	IsOwner int    `json:"isOwner"`
}

type SaveTaskReq struct {
	Name        string `form:"name"`
	StageCode   string `form:"stage_code"`
	ProjectCode string `form:"project_code"`
	AssignTo    string `form:"assign_to"`
}

type MoveTaskReq struct {
	PreTaskCode  string `form:"preTaskCode"`
	NextTaskCode string `form:"nextTaskCode"`
	ToStageCode  string `form:"toStageCode"`
}

type SelfTaskReq struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
	TaskType int `json:"taskType" form:"taskType"`
	Type     int `json:"type" form:"type"`
}

type TaskMember struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	Avatar            string `json:"avatar"`
	Code              string `json:"code"`
	MemberAccountCode string `json:"member_account_code"`
	IsExecutor        int    `json:"is_executor"`
	IsOwner           int    `json:"is_owner"`
}
