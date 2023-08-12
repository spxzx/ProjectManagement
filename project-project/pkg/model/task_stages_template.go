package model

type TaskStagesTemplate struct {
	Id                  int
	Name                string
	ProjectTemplateCode int
	CreateTime          int64
	Sort                int
}

type TaskStagesTmplNames struct {
	Name string
}

func ToNamesMap(task []TaskStagesTemplate) map[int][]*TaskStagesTmplNames {
	tm := make(map[int][]*TaskStagesTmplNames)
	for _, v := range task {
		t := &TaskStagesTmplNames{}
		t.Name = v.Name
		tm[v.ProjectTemplateCode] = append(tm[v.ProjectTemplateCode], t)
	}
	return tm
}
