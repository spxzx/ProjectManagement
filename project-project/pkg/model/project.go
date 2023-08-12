package model

type Project struct {
	Id                 int64
	Cover              string
	Name               string
	Description        string
	AccessControlType  int
	WhiteList          string
	Sort               int
	Deleted            int
	TemplateCode       int
	Schedule           float64
	CreateTime         int64
	OrganizationCode   int64
	DeletedTime        string
	Private            int
	Prefix             string
	OpenPrefix         int
	Archive            int
	ArchiveTime        int64
	OpenBeginTime      int
	OpenTaskPrivate    int
	TaskBoardTheme     string
	BeginTime          int64
	EndTime            int64
	AutoUpdateSchedule int
}

func (pmu *Project) GetAccessControlType() string {
	if pmu.AccessControlType == 0 {
		return "open"
	}
	if pmu.AccessControlType == 1 {
		return "private"
	}
	if pmu.AccessControlType == 2 {
		return "custom"
	}
	return ""
}

func ToProjectMap(list []*Project) map[int64]*Project {
	m:=make(map[int64]*Project, len(list))
	for _, v := range list {
		m[v.Id] = v
	}
	return m
}

type ProjectMemberUnion struct {
	Project
	ProjectCode int64
	MemberCode  int64
	JoinTime    int64
	IsOwner     int64
	Authorize   string
	OwnerName   string
	Collected   int
}

func ToMap(pmu []*ProjectMemberUnion) map[int64]*ProjectMemberUnion {
	m := make(map[int64]*ProjectMemberUnion)
	for _, v := range pmu {
		m[v.ProjectCode] = v
	}
	return m
}

func (pmu *ProjectMemberUnion) GetAccessControlType() string {
	if pmu.AccessControlType == 0 {
		return "open"
	}
	if pmu.AccessControlType == 1 {
		return "private"
	}
	if pmu.AccessControlType == 2 {
		return "custom"
	}
	return ""
}
