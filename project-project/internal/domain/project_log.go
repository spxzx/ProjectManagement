package domain

import (
	"context"
	"github.com/spxzx/project-common/encrypts"
	"github.com/spxzx/project-common/errs"
	"github.com/spxzx/project-grpc/task/task"
	"github.com/spxzx/project-grpc/user/login"
	"github.com/spxzx/project-project/internal/dao"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/internal/repo"
	"github.com/spxzx/project-project/pkg/data"
	"github.com/spxzx/project-project/pkg/model"
	"go.uber.org/zap"
	"time"
)

type ProjectLogDomain struct {
	log     repo.ProjectLogRepo
	tsc     db.Transaction
	userRpc *UserRpcDomain
}

func NewProjectLogDomain() *ProjectLogDomain {
	return &ProjectLogDomain{
		log:     dao.NewProjectLogDao(),
		tsc:     dao.NewTransaction(),
		userRpc: NewUserRpcDomain(),
	}
}

func (p *ProjectLogDomain) FindLogByMemberCode(ctx context.Context, mid, page,
	size int64) ([]*model.ProjectLog, int64, []int64, []int64, []int64, error) {
	projectLogs, total, err := p.log.FindLogByMemberCode(ctx, mid, page, size)
	if err != nil {
		zap.L().Error("project GetLogBySelfProject FindLogByMemberCode error, cause by: ", zap.Error(err))
		return nil, 0, nil, nil, nil, errs.GrpcError(data.DBError)
	}
	pIdList := make([]int64, len(projectLogs))
	mIdList := make([]int64, len(projectLogs))
	taskIdList := make([]int64, len(projectLogs))
	for _, v := range projectLogs {
		pIdList = append(pIdList, v.ProjectCode)
		mIdList = append(mIdList, v.MemberCode)
		taskIdList = append(taskIdList, v.SourceCode)
	}
	return projectLogs, total, pIdList, mIdList, taskIdList, nil
}

func (p *ProjectLogDomain) CreateProjectLog(ctx context.Context, pCode int64,
	taskCode int64, taskName string, assignTo int64, logType string, actionType string) error {
	return p.tsc.Action(func(conn db.Conn) error {
		remark := ""
		if logType == "create" {
			remark = "创建了任务"
		}
		if err := p.log.SaveProjectLog(ctx, conn, &model.ProjectLog{
			MemberCode:  assignTo,
			Content:     taskName,
			Remark:      remark,
			Type:        logType,
			CreateTime:  time.Now().UnixMilli(),
			SourceCode:  taskCode,
			ActionType:  actionType,
			IsComment:   0,
			ProjectCode: pCode,
			Icon:        "plus",
			IsRobot:     0,
		}); err != nil {
			zap.L().Error("tasks SaveTask createProjectLog error, cause by: ", zap.Error(err))
			return err
		}
		return nil
	})
}

func (p *ProjectLogDomain) GetTaskLog(ctx context.Context, req *task.TaskRpcRequest) ([]*model.ProjectLogDisplay, int64, error) {
	taskCode := encrypts.DecryptNoErr(req.TaskCode, data.AESKey)
	var list []*model.ProjectLog
	var total int64
	var err error
	switch req.All {
	case 0:
		list, total, err = p.log.FindLogByTaskCodePage(ctx, taskCode, int(req.Comment), int(req.Page), int(req.PageSize))
	case 1:
		list, total, err = p.log.FindLogByTaskCode(ctx, taskCode, int(req.Comment))
	}
	if err != nil {
		zap.L().Error("tasks GetTaskLog FindLogByTaskCode error, cause by: ", zap.Error(err))
		return nil, 0, errs.GrpcError(data.DBError)
	}
	if total == 0 {
		return nil, 0, nil
	}
	var displayList []*model.ProjectLogDisplay
	var mIdList []int64
	for _, v := range list {
		mIdList = append(mIdList, v.MemberCode)
	}
	MemList, _, err := p.userRpc.FindMemInfoByIds(ctx, &login.MemRequest{MemIds: mIdList})
	if err != nil {
		return nil, 0, err
	}
	mMap := make(map[int64]*login.Member)
	for _, v := range MemList.List {
		mMap[v.Id] = v
	}
	for _, v := range list {
		display := v.ToDisplay()
		message := mMap[v.MemberCode]
		m := model.Member{}
		m.Name = message.Name
		m.Id = message.Id
		m.Avatar = message.Avatar
		m.Code = message.Code
		display.Member = m
		displayList = append(displayList, display)
	}
	return displayList, total, nil
}

func (p *ProjectLogDomain) CreateComment(ctx context.Context, pl *model.ProjectLog) error {
	return p.tsc.Action(func(conn db.Conn) error {
		if err := p.log.SaveProjectLog(ctx, conn, pl); err != nil {
			zap.L().Error("tasks CreateComment db SaveProjectLog error, cause by: ", zap.Error(err))
			return errs.GrpcError(data.DBError)
		}
		return nil
	})
}
