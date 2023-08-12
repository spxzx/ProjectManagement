package project

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/spxzx/project-common/encrypts"
	"github.com/spxzx/project-common/tms"
	"github.com/spxzx/project-grpc/project/project"
	"github.com/spxzx/project-grpc/user/login"
	"github.com/spxzx/project-project/internal/dao"
	"github.com/spxzx/project-project/internal/domain"
	"github.com/spxzx/project-project/internal/repo"
	"github.com/spxzx/project-project/pkg/data"
	"github.com/spxzx/project-project/pkg/model"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
	"time"
)

type Service struct {
	cache   repo.Cache
	md      *domain.MenuDomain
	pd      *domain.ProjectDomain
	tpld    *domain.ProjectTemplateDomain
	tstd    *domain.TaskStagesTemplateDomain
	tsd     *domain.TaskStagesDomain
	pld     *domain.ProjectLogDomain
	td      *domain.TaskDomain
	pnd     *domain.ProjectNodeDomain
	userRpc *domain.UserRpcDomain
	project.UnimplementedProjectServiceServer
}

func New() *Service {
	return &Service{
		cache:   dao.Rc,
		md:      domain.NewMenuDomain(),
		pd:      domain.NewProjectDomain(),
		tpld:    domain.NewProjectTemplateDomain(),
		tstd:    domain.NewTaskStagesTemplateDomain(),
		tsd:     domain.NewTaskStagesDomain(),
		pld:     domain.NewProjectLogDomain(),
		td:      domain.NewTaskDomain(),
		pnd:     domain.NewProjectNodeDomain(),
		userRpc: domain.NewUserRpcDomain(),
	}
}

func (s *Service) Index(ctx context.Context, _ *emptypb.Empty) (*project.IndexResponse, error) {
	menuList, err := s.md.FindAll(ctx)
	if err != nil {
		return &project.IndexResponse{}, err
	}
	var menus []*project.Menu
	_ = copier.Copy(&menus, model.ConvertChild(menuList))
	return &project.IndexResponse{
		Menus: menus,
	}, nil
}

func (s *Service) GetProjectList(ctx context.Context, req *project.ProjectRpcRequest) (*project.ProjectRpcResponse, error) {
	list, total, err := s.pd.GetProjectList(ctx, req)
	if err != nil {
		return &project.ProjectRpcResponse{}, err
	}

	var pmList []*project.Project
	_ = copier.Copy(&pmList, list)
	m := model.ToMap(list)
	for _, v := range pmList {
		v.Code, _ = encrypts.EncryptInt64(v.ProjectCode, data.AESKey)
		lm := m[v.ProjectCode]
		v.AccessControlType = lm.GetAccessControlType()
		v.OrganizationCode, _ = encrypts.EncryptInt64(lm.OrganizationCode, data.AESKey)
		v.OwnerName = req.MemberName
		v.JoinTime = tms.FormatByMill(lm.JoinTime)
		v.Order = int32(lm.Sort)
		v.CreateTime = tms.FormatByMill(lm.CreateTime)
		v.Collected = int32(lm.Collected)
	}
	return &project.ProjectRpcResponse{Pm: pmList, Total: total}, nil
}

func (s *Service) GetProjectTemplates(ctx context.Context, req *project.ProjectRpcRequest) (*project.ProjectTemplateResponse, error) {
	// 1.根据viewType查询项目模板表 得到list
	ptList, total, err := s.tpld.FindProjectTemplates(ctx, req)
	if err != nil {
		return &project.ProjectTemplateResponse{}, nil
	}
	// 2.模型转换 拿到模板id列表去任务步骤模板表查询
	taskList, err := s.tstd.FindInPojTmplIds(ctx, model.ToPojTmplIds(ptList))
	if err != nil {
		return &project.ProjectTemplateResponse{}, err
	}
	// 3.数据整合
	var ptDetail []*model.PojTmplDetail
	for _, v := range ptList {
		ptDetail = append(ptDetail, v.Combine(model.ToNamesMap(taskList)[v.Id]))
	}
	// 4.返回数据
	var resp []*project.ProjectTemplate
	_ = copier.Copy(&resp, ptDetail)
	return &project.ProjectTemplateResponse{Pt: resp, Total: total}, nil
}

func (s *Service) SaveProject(ctx context.Context, req *project.ProjectRpcRequest) (*project.SaveProjectResponse, error) {
	orgCodeStr, _ := encrypts.Decrypt(req.OrganizationCode, data.AESKey)
	orgCode, _ := strconv.ParseInt(orgCodeStr, 10, 64)
	tplCodeStr, _ := encrypts.Decrypt(req.TemplateCode, data.AESKey)
	tplCode, _ := strconv.ParseInt(tplCodeStr, 10, 64)
	// 获取模板信息
	tsTmpl, err := s.tstd.FindTaskStagesTplByPid(ctx, tplCode)
	if err != nil {
		return &project.SaveProjectResponse{}, err
	}
	// 保存项目表
	pj := &model.Project{
		Name:              req.Name,
		Description:       req.Description,
		TemplateCode:      int(tplCode),
		CreateTime:        time.Now().UnixMilli(),
		Cover:             "https://img2.baidu.com/it/u=792555388,2449797505&fm=253&fmt=auto&app=138&f=JPEG?w=667&h=500",
		Deleted:           data.NotDeleted,
		Archive:           data.NotArchive,
		OrganizationCode:  orgCode,
		AccessControlType: data.Open,
		TaskBoardTheme:    data.Simple,
	}
	if err = s.pd.SaveProject(ctx, req, pj); err != nil {
		return &project.SaveProjectResponse{}, err
	}
	for i, v := range tsTmpl {
		if err := s.tsd.SaveTaskStages(ctx, pj, i, v); err != nil {
			return &project.SaveProjectResponse{}, err
		}
	}
	code, _ := encrypts.EncryptInt64(pj.Id, data.AESKey)
	return &project.SaveProjectResponse{
		Id:               pj.Id,
		Cover:            pj.Cover,
		Name:             pj.Name,
		Description:      pj.Description,
		Code:             code,
		CreateTime:       tms.FormatByMill(pj.CreateTime),
		TaskBoardTheme:   pj.TaskBoardTheme,
		OrganizationCode: req.OrganizationCode,
	}, nil
}

func (s *Service) ReadProject(ctx context.Context, req *project.ProjectRpcRequest) (*project.ReadProjectResponse, error) {
	pam, mem, err := s.pd.GetProjectDetail(ctx, req)
	if err != nil {
		return &project.ReadProjectResponse{}, err
	}
	resp := &project.ReadProjectResponse{}
	_ = copier.Copy(resp, pam)
	resp.Code, _ = encrypts.EncryptInt64(pam.ProjectCode, data.AESKey)
	resp.AccessControlType = pam.GetAccessControlType()
	resp.OrganizationCode, _ = encrypts.EncryptInt64(pam.OrganizationCode, data.AESKey)
	resp.Order = int32(pam.Sort)
	resp.CreateTime = tms.FormatByMill(pam.CreateTime)
	resp.OwnerName = mem.Name
	resp.OwnerAvatar = mem.Avatar
	return resp, nil
}

func (s *Service) UpdateProjectDeleted(ctx context.Context, req *project.ProjectRpcRequest) (*emptypb.Empty, error) {
	pojCode := encrypts.DecryptNoErr(req.ProjectCode, data.AESKey)
	return &emptypb.Empty{}, s.pd.UpdateProjectDeleted(ctx, pojCode, req.Deleted)
}

func (s *Service) UpdateProjectCollected(ctx context.Context, req *project.ProjectRpcRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.pd.UpdateProjectCollected(ctx, req)
}

func (s *Service) EditProject(ctx context.Context, req *project.EditProjectRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.pd.UpdateProject(ctx, req)
}

func (s *Service) GetLogBySelfProject(ctx context.Context, req *project.ProjectRpcRequest) (*project.ProjectLogResponse, error) {
	projectLogs, total, pIdList, mIdList, taskIdList, err := s.
		pld.FindLogByMemberCode(ctx, req.MemberId, req.Page, req.PageSize)
	if err != nil {
		return &project.ProjectLogResponse{}, err
	}

	_, pMap, err := s.pd.FindProjectByIds(ctx, pIdList)
	if err != nil {
		return &project.ProjectLogResponse{}, err
	}

	_, mMap, err := s.userRpc.FindMemInfoByIds(ctx, &login.MemRequest{MemIds: mIdList})

	_, tMap, err := s.td.FindTaskByIds(ctx, taskIdList)

	var list []*model.IndexProjectLogDisplay
	for _, v := range projectLogs {
		display := v.ToIndexDisplay()
		display.ProjectName = pMap[v.ProjectCode].Name
		display.MemberAvatar = mMap[v.MemberCode].Avatar
		display.MemberName = mMap[v.MemberCode].Name
		display.TaskName = tMap[v.SourceCode].Name
		list = append(list, display)
	}
	var msgList []*project.ProjectLog
	_ = copier.Copy(&msgList, list)
	return &project.ProjectLogResponse{List: msgList, Total: total}, nil
}

func (s *Service) GetNodeList(ctx context.Context, _ *project.ProjectRpcRequest) (*project.ProjectNodeResponse, error) {
	list, err := s.pnd.TreeList(ctx)
	if err != nil {
		return nil, err
	}
	var nodes []*project.ProjectNode
	_ = copier.Copy(&nodes, list)
	return &project.ProjectNodeResponse{Nodes: nodes}, nil
}

func (s *Service) FindProjectByMemberId(ctx context.Context, req *project.ProjectRpcRequest) (*project.FindProjectByMemberIdResponse, error) {
	isProjectCode := false
	var projectId int64
	if req.ProjectCode != "" {
		projectId = encrypts.DecryptNoErr(req.ProjectCode, data.AESKey)
		isProjectCode = true
	}
	isTaskCode := false
	var taskId int64
	if req.TaskCode != "" {
		taskId = encrypts.DecryptNoErr(req.TaskCode, data.AESKey)
		isTaskCode = true
	}
	if !isProjectCode && isTaskCode {
		projectCode, ok, bError := s.td.FindProjectIdByTaskId(taskId)
		if bError != nil {
			return &project.FindProjectByMemberIdResponse{}, bError
		}
		if !ok {
			return &project.FindProjectByMemberIdResponse{
				Project:  nil,
				IsOwner:  false,
				IsMember: false,
			}, nil
		}
		projectId = projectCode
	}
	if isProjectCode {
		// 根据projectid和memberid查询
		req.ProjectCode = encrypts.EncryptNoErr(projectId, data.AESKey)
		pm, _, err := s.pd.GetProjectDetail(ctx, req)
		if err != nil {
			return &project.FindProjectByMemberIdResponse{}, err
		}
		if pm == nil {
			return &project.FindProjectByMemberIdResponse{
				Project:  nil,
				IsOwner:  false,
				IsMember: false,
			}, nil
		}
		projectMessage := &project.Project{}
		_ = copier.Copy(projectMessage, pm)
		isOwner := false
		if pm.IsOwner == 1 {
			isOwner = true
		}
		return &project.FindProjectByMemberIdResponse{
			Project:  projectMessage,
			IsOwner:  isOwner,
			IsMember: true,
		}, nil
	}
	return &project.FindProjectByMemberIdResponse{}, nil
}
