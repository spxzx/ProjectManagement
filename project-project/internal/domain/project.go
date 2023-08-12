package domain

import (
	"context"
	"github.com/spxzx/project-common/encrypts"
	"github.com/spxzx/project-common/errs"
	"github.com/spxzx/project-grpc/project/project"
	"github.com/spxzx/project-grpc/user/login"
	"github.com/spxzx/project-project/internal/dao"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/internal/repo"
	"github.com/spxzx/project-project/pkg/data"
	"github.com/spxzx/project-project/pkg/model"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type ProjectDomain struct {
	poj     repo.ProjectRepo
	tsc     db.Transaction
	userRpc *UserRpcDomain
}

func NewProjectDomain() *ProjectDomain {
	return &ProjectDomain{
		poj:     dao.NewProjectDao(),
		tsc:     dao.NewTransaction(),
		userRpc: NewUserRpcDomain(),
	}
}

func (p *ProjectDomain) GetProjectList(ctx context.Context, req *project.ProjectRpcRequest) ([]*model.ProjectMemberUnion, int64, error) {
	//config.SendLog(kfk.Info("Get", "ProjectDomain.GetProjectList", kfk.FieldMap{
	//	"ProjectRpcRequest": req,
	//}))
	var list, clist []*model.ProjectMemberUnion
	var total, cTotal int64
	var err error
	if req.SelectBy == "" || req.SelectBy == "my" {
		list, total, err = p.poj.FindMyProjectListByMemId(ctx, req.MemberId, "", req.Page, req.PageSize)
	}
	if req.SelectBy == "archive" || req.SelectBy == "deleted" {
		list, total, err = p.poj.FindMyProjectListByMemId(ctx, req.MemberId, req.SelectBy, req.Page, req.PageSize)
	}
	clist, cTotal, err = p.poj.FindCollectProjectListByMemId(ctx, req.MemberId, req.Page, req.PageSize)
	if err != nil {
		zap.L().Error("getProjectList db Find MyProjectList/Collect ByMemId error, cause by: ", zap.Error(err))
		return nil, 0, errs.GrpcError(data.DBError)
	}
	if req.SelectBy == "collect" {
		list = clist
		total = cTotal
		for _, v := range list {
			v.Collected = data.Collected
		}
	} else {
		cMap := model.ToMap(clist)
		for _, pm := range list {
			if cMap[pm.ProjectCode] != nil {
				pm.Collected = data.Collected
			}
		}
	}
	return list, total, nil
}

func (p *ProjectDomain) SaveProject(ctx context.Context, req *project.ProjectRpcRequest, pj *model.Project) error {
	return p.tsc.Action(func(conn db.Conn) error {
		if err := p.poj.SaveProject(ctx, conn, pj); err != nil {
			zap.L().Error("saveProject db SaveProject error, cause by: ", zap.Error(err))
			return data.DBError
		}
		pm := &model.ProjectMember{
			ProjectCode: pj.Id,
			MemberCode:  req.MemberId,
			JoinTime:    time.Now().UnixMilli(),
			IsOwner:     int(req.MemberId),
			Authorize:   "",
		}
		if err := p.poj.SaveProjectMember(ctx, conn, pm); err != nil {
			zap.L().Error("saveProject db SaveProjectMember error, cause by: ", zap.Error(err))
			return data.DBError
		}
		return nil
	})
}

func (p *ProjectDomain) GetProjectDetail(ctx context.Context, req *project.ProjectRpcRequest) (*model.ProjectMemberUnion, *login.Member, error) {
	// 1. 查项目和成员的关联表 查到项目的拥有者
	pojCodeStr, _ := encrypts.Decrypt(req.ProjectCode, data.AESKey)
	pojCode, _ := strconv.ParseInt(pojCodeStr, 10, 64)
	pam, err := p.poj.FindProjectByPidAndMemId(ctx, pojCode, req.MemberId)
	if err != nil {
		zap.L().Error("readProject db FindProjectByPidAndMemId error, cause by:", zap.Error(err))
		return nil, nil, errs.GrpcError(data.DBError)
	}
	// 2. 去member表查名字
	mem, err := p.userRpc.FindMemInfoById(ctx, pam.IsOwner)
	if err != nil {
		return nil, nil, err
	}
	// 3. 查收藏表 判断收藏状态
	isCollect, err := p.poj.FindCollectByPidAndMemId(ctx, pojCode, req.MemberId)
	if err != nil {
		zap.L().Error("readProject db FindCollectByPidAndMemId error, cause by:", zap.Error(err))
		return nil, nil, errs.GrpcError(data.DBError)
	}
	if isCollect {
		pam.Collected = data.Collected
	}
	return pam, mem, nil
}

func (p *ProjectDomain) UpdateProjectDeleted(ctx context.Context, pCode int64, deleted bool) error {
	if err := p.poj.UpdateProjectDeleted(ctx, pCode, deleted); err != nil {
		zap.L().Error("updateProjectDeleted db UpdateProjectDeleted DeleteProject error, cause by: ", zap.Error(err))
		return errs.GrpcError(data.DBError)
	}
	return nil
}

func (p *ProjectDomain) UpdateProjectCollected(ctx context.Context, req *project.ProjectRpcRequest) error {
	pojCode := encrypts.DecryptNoErr(req.ProjectCode, data.AESKey)
	var err error
	switch req.CollectType {
	case "collect":
		err = p.poj.SaveProjectCollected(ctx, &model.ProjectCollection{
			ProjectCode: pojCode,
			MemberCode:  req.MemberId,
			CreateTime:  time.Now().UnixMilli(),
		})
	case "cancel":
		err = p.poj.DeleteProjectCollected(ctx, req.MemberId, pojCode)
	}
	if err != nil {
		zap.L().Error("db UpdateProjectCollected error, cause by: ", zap.Error(err))
		return errs.GrpcError(data.DBError)
	}
	return nil
}

func (p *ProjectDomain) UpdateProject(ctx context.Context, req *project.EditProjectRequest) error {
	pojCode := encrypts.DecryptNoErr(req.ProjectCode, data.AESKey)
	err := p.poj.UpdateProject(ctx, &model.Project{
		Id:                 pojCode,
		Name:               req.Name,
		Description:        req.Description,
		Cover:              req.Cover,
		TaskBoardTheme:     req.TaskBoardTheme,
		Prefix:             req.Prefix,
		Private:            int(req.Private),
		OpenPrefix:         int(req.OpenPrefix),
		OpenTaskPrivate:    int(req.OpenTaskPrivate),
		Schedule:           req.Schedule,
		AutoUpdateSchedule: int(req.AutoUpdateSchedule),
	})
	if err != nil {
		zap.L().Error("editProject db EditProject error, cause by:", zap.Error(err))
		return errs.GrpcError(data.DBError)
	}
	return nil
}

func (p *ProjectDomain) FindProjectByIds(ctx context.Context, ids []int64) ([]*model.Project, map[int64]*model.Project, error) {
	projects, err := p.poj.FindProjectByIds(ctx, ids)
	if err != nil {
		zap.L().Error("project GetLogBySelfProject FindProjectByIds error, cause by: ", zap.Error(err))
		return nil, nil, errs.GrpcError(data.DBError)
	}
	pMap := make(map[int64]*model.Project)
	for _, v := range projects {
		pMap[v.Id] = v
	}
	return projects, pMap, nil
}

func (p *ProjectDomain) FindProjectMemberByPid(ctx context.Context, pCode int64,
) ([]*model.ProjectMember, int64, []int64, map[int64]*model.ProjectMember, error) {
	pmList, total, err := p.poj.FindProjectMemberByPid(ctx, pCode)
	if err != nil {
		zap.L().Error("task GetProjectMember db FindProjectMemberByPid error, cause by", zap.Error(err))
		return nil, 0, nil, nil, errs.GrpcError(data.DBError)
	}
	if len(pmList) == 0 {
		return nil, 0, nil, nil, nil
	}
	ids, pmMap := model.ToPMMemIdsAndMap(pmList)
	return pmList, total, ids, pmMap, nil
}

func (p *ProjectDomain) FindProjectById(ctx context.Context, pid int64) (*model.Project, error) {
	pj, err := p.poj.FindProjectById(ctx, pid)
	if err != nil {
		zap.L().Error("tasks SaveTask FindProjectById error, cause by: ", zap.Error(err))
		return nil, errs.GrpcError(data.DBError)
	}
	if pj.Deleted == data.Deleted {
		return nil, errs.GrpcError(data.TaskStagesNotNull)
	}
	return pj, nil
}
