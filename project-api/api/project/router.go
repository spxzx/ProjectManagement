package project

import (
	"github.com/gin-gonic/gin"
	"github.com/spxzx/project-api/api/mid"
	"github.com/spxzx/project-api/api/rpc"
	"github.com/spxzx/project-api/router"
)

func init() {
	router.Register(
		&RouterProject{},
		&RouterTask{},
		&RouterAccount{},
		&RouterDepartment{},
		&RouterAuth{},
		&RouterMenu{},
	)
}

type RouterProject struct{}

func (*RouterProject) Route(r *gin.Engine) {
	rpc.InitProjectRpcClient()
	h := NewProject()
	p := r.Group("/project", mid.TokenVerify(), Auth())
	p.POST("/index", h.index)
	p.POST("/project/selfList", h.getProjectList)
	p.POST("/project", h.getProjectList)
	p.POST("/project_template", h.getProjectTemplates)
	p.POST("/project/save", h.saveProject)
	p.POST("/project/read", h.readProject)
	p.POST("/project/recycle", h.recycleProject)
	p.POST("/project/recovery", h.recoveryProject)
	p.POST("/project_collect/collect", h.collectProject)
	p.POST("/project/edit", h.editProject)
	p.POST("/project/getLogBySelfProject", h.getLogBySelfProject)
	p.POST("/node", h.getNodeList)
}

type RouterTask struct{}

func (*RouterTask) Route(r *gin.Engine) {
	h := NewTask()
	p := r.Group("/project", mid.TokenVerify(), Auth())
	p.POST("/task_stages", h.getTaskStages)
	p.POST("/project_member/index", h.getProjectMember)
	p.POST("/task_stages/tasks", h.getTaskStageDetailList)
	p.POST("/task/save", h.saveTask)
	p.POST("/task/sort", h.moveTask)
	p.POST("/task/selfList", h.getSelfTaskList)
	p.POST("/task/read", h.readTask)
	p.POST("/task_member", h.getTaskMember)
	p.POST("/task/taskLog", h.getTaskLog)
	p.POST("/task/_taskWorkTimeList", h.getTaskWorkTimeList)
	p.POST("/task/saveTaskWorkTime", h.saveTaskWorkTime)
	p.POST("/file/uploadFiles", h.uploadFile)
	p.POST("/task/taskSources", h.getTaskSourceLink)
	p.POST("/task/createComment", h.createComment)
}

type RouterAccount struct{}

func (*RouterAccount) Route(r *gin.Engine) {
	h := NewAccount()
	p := r.Group("/project", mid.TokenVerify(), Auth())
	p.POST("/account", h.getAccountList)
}

type RouterDepartment struct{}

func (*RouterDepartment) Route(r *gin.Engine) {
	h := NewDepartment()
	p := r.Group("/project", mid.TokenVerify(), Auth())
	p.POST("/department", h.getDepartmentList)
	p.POST("/department/save", h.saveDepartment)
	p.POST("/department/read", h.readDepartment)
}

type RouterAuth struct{}

func (*RouterAuth) Route(r *gin.Engine) {
	h := NewAuth()
	p := r.Group("/project", mid.TokenVerify(), Auth())
	p.POST("/auth", h.getAuthList)
	p.POST("/auth/apply", h.apply)
}

type RouterMenu struct{}

func (o *RouterMenu) Route(r *gin.Engine) {
	h := NewMenu()
	p := r.Group("/project", mid.TokenVerify(), Auth())
	p.POST("/menu/menu", h.getMenuList)
}
