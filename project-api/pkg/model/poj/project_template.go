package poj

type ProjectTemplate struct {
	Id               int                   `json:"id"`
	Code             string                `json:"code"`
	Name             string                `json:"name"`
	Description      string                `json:"description"`
	Sort             int                   `json:"sort"`
	CreateTime       string                `json:"create_time"`
	OrganizationCode string                `json:"organization_code"`
	Cover            string                `json:"cover"`
	MemberCode       string                `json:"member_code"`
	IsSystem         int                   `json:"is_system"`
	TaskStages       []*TaskStagesOnlyName `json:"task_stages"`
}

type TaskStagesOnlyName struct {
	Name string `json:"name"`
}
