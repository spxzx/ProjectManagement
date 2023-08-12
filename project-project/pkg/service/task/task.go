package task

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/spxzx/project-common/encrypts"
	"github.com/spxzx/project-common/errs"
	"github.com/spxzx/project-common/tms"
	"github.com/spxzx/project-grpc/task/task"
	"github.com/spxzx/project-grpc/user/login"
	"github.com/spxzx/project-project/internal/dao"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/internal/domain"
	"github.com/spxzx/project-project/internal/repo"
	"github.com/spxzx/project-project/pkg/data"
	"github.com/spxzx/project-project/pkg/model"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type Service struct {
	cache   repo.Cache
	tsc     db.Transaction
	tsd     *domain.TaskStagesDomain
	pd      *domain.ProjectDomain
	td      *domain.TaskDomain
	pld     *domain.ProjectLogDomain
	twtd    *domain.TaskWorkTimeDomain
	fd      *domain.FileDomain
	sld     *domain.SourceLinkDomain
	userRpc *domain.UserRpcDomain
	task.UnimplementedTaskServiceServer
}

func New() *Service {
	return &Service{
		cache:   dao.Rc,
		tsc:     dao.NewTransaction(),
		tsd:     domain.NewTaskStagesDomain(),
		pd:      domain.NewProjectDomain(),
		td:      domain.NewTaskDomain(),
		pld:     domain.NewProjectLogDomain(),
		twtd:    domain.NewTaskWorkTimeDomain(),
		fd:      domain.NewFileDomain(),
		sld:     domain.NewSourceLinkDomain(),
		userRpc: domain.NewUserRpcDomain(),
	}
}

func (s *Service) GetTaskStages(ctx context.Context, req *task.TaskRpcRequest) (*task.TaskStagesResponse, error) {
	pojCode := encrypts.DecryptNoErr(req.ProjectCode, data.AESKey)
	list, tssMap, total, err := s.tsd.FindTaskStagesByPCode(ctx, pojCode, req.Page, req.PageSize)
	if err != nil {
		return &task.TaskStagesResponse{}, err
	}
	var resp []*task.TaskStages
	_ = copier.Copy(&resp, list)
	for _, v := range resp {
		tss := tssMap[int(v.Id)]
		v.Code = encrypts.EncryptNoErr(int64(v.Id), data.AESKey)
		v.CreateTime = tms.FormatByMill(tss.CreateTime)
		v.ProjectCode = req.ProjectCode
	}
	return &task.TaskStagesResponse{
		List:  resp,
		Total: total,
	}, nil
}

func (s *Service) GetProjectMember(ctx context.Context, req *task.TaskRpcRequest) (*task.ProjectMemberResponse, error) {
	_, total, ids, pmMap, err := s.
		pd.FindProjectMemberByPid(ctx, encrypts.DecryptNoErr(req.ProjectCode, data.AESKey))
	if err != nil {
		return &task.ProjectMemberResponse{}, err
	}
	infoResp, _, err := s.userRpc.FindMemInfoByIds(ctx, &login.MemRequest{MemIds: ids})
	if err != nil {
		return &task.ProjectMemberResponse{}, err
	}

	var resp []*task.ProjectMember
	for _, v := range infoResp.List {
		pm := &task.ProjectMember{
			Name:       v.Name,
			Avatar:     v.Avatar,
			MemberCode: v.Id,
			Code:       encrypts.EncryptNoErr(v.Id, data.AESKey),
			Email:      v.Email,
		}
		if int(v.Id) == pmMap[v.Id].IsOwner {
			pm.IsOwner = data.Owner
		}
		resp = append(resp, pm)
	}

	return &task.ProjectMemberResponse{
		List:  resp,
		Total: total,
	}, nil
}

func (s *Service) GetTaskStageDetailList(ctx context.Context, req *task.TaskRpcRequest) (*task.TaskStageDetailListResponse, error) {
	tdList, err := s.td.GetTaskStageDetailList(ctx, req)
	if err != nil {
		return &task.TaskStageDetailListResponse{}, err
	}
	var resp []*task.Task
	_ = copier.Copy(&resp, tdList)
	return &task.TaskStageDetailListResponse{
		List: resp,
	}, nil
}

func (s *Service) SaveTask(ctx context.Context, req *task.TaskRpcRequest) (*task.Task, error) {
	if req.Name == "" {
		return &task.Task{}, errs.GrpcError(data.TaskNameNotNull)
	}
	stageCode := encrypts.DecryptNoErr(req.StageCode, data.AESKey)
	if _, err := s.tsd.FindTaskStagesById(ctx, stageCode); err != nil {
		return &task.Task{}, err
	}

	pojCode := encrypts.DecryptNoErr(req.ProjectCode, data.AESKey)
	pj, err := s.pd.FindProjectById(ctx, pojCode)
	if err != nil {
		return &task.Task{}, err
	}

	maxIdNum, err := s.td.FindTaskMaxIdNum(ctx, pojCode)
	maxSort, err := s.td.FindTaskMaxSort(ctx, pojCode, stageCode)

	assignTo := encrypts.DecryptNoErr(req.AssignTo, data.AESKey)
	ts := &model.Task{
		Name:        req.Name,
		CreateTime:  time.Now().UnixMilli(),
		CreateBy:    req.MemberId,
		AssignTo:    assignTo,
		ProjectCode: pojCode,
		StageCode:   int(stageCode),
		IdNum:       *maxIdNum + 1,
		Private:     pj.OpenTaskPrivate,
		Sort:        *maxSort + 65536,
		BeginTime:   time.Now().UnixMilli(),
		EndTime:     time.Now().Add(2 * 24 * time.Hour).UnixMilli(),
	}

	if err = s.td.SaveTask(ctx, ts, req.MemberId); err != nil {
		return &task.Task{}, errs.GrpcError(data.DBError)
	}

	if err = s.pld.CreateProjectLog(ctx, ts.ProjectCode, ts.Id, ts.Name,
		ts.AssignTo, "create", "task"); err != nil {
		return &task.Task{}, err
	}

	member, err := s.userRpc.FindMemInfoById(ctx, assignTo)
	if err != nil {
		return &task.Task{}, err
	}

	td := ts.ToTaskDisplay()
	td.Executor = model.Executor{
		Name:   member.Name,
		Avatar: member.Avatar,
		Code:   member.Code,
	}

	resp := &task.Task{}
	_ = copier.Copy(resp, td)
	return resp, nil
}

func (s *Service) MoveTask(ctx context.Context, req *task.TaskRpcRequest) (*emptypb.Empty, error) {
	if req.PreTaskCode == req.NextTaskCode {
		return &emptypb.Empty{}, nil
	}
	preTaskCode := encrypts.DecryptNoErr(req.PreTaskCode, data.AESKey)
	toStageCode := encrypts.DecryptNoErr(req.ToStageCode, data.AESKey)
	if err := s.td.MoveAndSortTask(ctx, preTaskCode, req.NextTaskCode, toStageCode); err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) GetSelfTaskList(ctx context.Context, req *task.TaskRpcRequest) (*task.SelfTaskListResponse, error) {
	tkList, total, pids, mids, err := s.td.FindTaskByCondition(ctx, req)
	if err != nil {
		return &task.SelfTaskListResponse{}, err
	}

	pListChan := make(chan []*model.Project)
	defer close(pListChan)
	mListChan := make(chan *login.MemberInfoResponse)
	defer close(mListChan)

	go func() {
		pList, _, _ := s.pd.FindProjectByIds(ctx, pids)
		pListChan <- pList
	}()

	go func() {
		mList, _, _ := s.userRpc.FindMemInfoByIds(ctx, &login.MemRequest{MemIds: mids})
		mListChan <- mList
	}()

	pList := <-pListChan
	pMap := model.ToProjectMap(pList)
	mList := <-mListChan

	mMap := make(map[int64]*login.Member)
	for _, v := range mList.List {
		mMap[v.Id] = v
	}
	var stdList []*model.SelfTaskDisplay
	for _, v := range tkList {
		memberMessage := mMap[v.AssignTo]
		name := memberMessage.Name
		avatar := memberMessage.Avatar
		mtd := v.ToSelfTaskDisplay(pMap[v.ProjectCode], name, avatar)
		stdList = append(stdList, mtd)
	}
	var resp []*task.SelfTask
	_ = copier.Copy(&resp, stdList)
	return &task.SelfTaskListResponse{List: resp, Total: total}, nil
}

func (s *Service) ReadTask(ctx context.Context, req *task.TaskRpcRequest) (*task.Task, error) {
	taskInfo, td, err := s.td.ReadTask(ctx, req)
	if err != nil {
		return &task.Task{}, err
	}
	if taskInfo == nil {
		return &task.Task{}, nil
	}

	pj, err := s.pd.FindProjectById(ctx, taskInfo.ProjectCode)
	if err != nil {
		return &task.Task{}, err
	}
	td.ProjectName = pj.Name
	taskStages, err := s.tsd.FindTaskStagesById(ctx, int64(taskInfo.StageCode))
	if err != nil {
		return &task.Task{}, err
	}
	td.StageName = taskStages.Name
	mem, err := s.userRpc.FindMemInfoById(ctx, taskInfo.AssignTo)
	if err != nil {
		return &task.Task{}, err
	}
	e := model.Executor{
		Name:   mem.Name,
		Avatar: mem.Avatar,
	}
	td.Executor = e
	taskResp := &task.Task{}
	_ = copier.Copy(taskResp, td)
	//_ = copier.Copy(taskResp.Executor, td.Executor)
	return taskResp, nil
}

func (s *Service) GetTaskMemberList(ctx context.Context, req *task.TaskRpcRequest) (*task.TaskMemberList, error) {
	taskCode := encrypts.DecryptNoErr(req.TaskCode, data.AESKey)
	tmList, total, mMap, err := s.td.GetTaskMemberList(ctx, taskCode, req.Page, req.PageSize)
	if err != nil {
		return &task.TaskMemberList{}, err
	}
	var tmResp []*task.TaskMember
	for _, v := range tmList {
		tm := &task.TaskMember{}
		tm.Code = encrypts.EncryptNoErr(v.MemberCode, data.AESKey)
		tm.Id = v.Id
		msg := mMap[v.MemberCode]
		tm.Name = msg.Name
		tm.Avatar = msg.Avatar
		tm.IsExecutor = int32(v.IsExecutor)
		tm.IsOwner = int32(v.IsOwner)
		tmResp = append(tmResp, tm)
	}
	return &task.TaskMemberList{
		List:  tmResp,
		Total: total,
	}, nil
}

func (s *Service) GetTaskLog(ctx context.Context, req *task.TaskRpcRequest) (*task.TaskLogList, error) {
	displayList, total, err := s.pld.GetTaskLog(ctx, req)
	if err != nil {
		return &task.TaskLogList{}, err
	}
	var tl []*task.TaskLog
	_ = copier.Copy(&tl, displayList)
	return &task.TaskLogList{List: tl, Total: total}, nil
}

func (s *Service) GetTaskWorkTimeList(ctx context.Context, req *task.TaskRpcRequest) (*task.TaskWorkTimeResponse, error) {
	displayList, err := s.twtd.GetTaskWorkTimeList(ctx, req)
	if err != nil {
		return &task.TaskWorkTimeResponse{}, err
	}
	var twtl []*task.TaskWorkTime
	_ = copier.Copy(&twtl, displayList)
	return &task.TaskWorkTimeResponse{List: twtl}, nil
}

func (s *Service) SaveTaskWorkTimeList(ctx context.Context, req *task.TaskRpcRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.twtd.SaveWorkTime(ctx, req)
}

func (s *Service) UploadFile(ctx context.Context, req *task.TaskFileRequest) (*emptypb.Empty, error) {
	orgCode := encrypts.DecryptNoErr(req.OrganizationCode, data.AESKey)
	taskCode := encrypts.DecryptNoErr(req.TaskCode, data.AESKey)
	file := model.File{
		PathName:         req.PathName,
		Title:            req.FileName,
		Extension:        req.Extension,
		Size:             int(req.Size),
		ObjectType:       "",
		OrganizationCode: orgCode,
		TaskCode:         taskCode,
		ProjectCode:      encrypts.DecryptNoErr(req.ProjectCode, data.AESKey),
		CreateBy:         req.MemberId,
		CreateTime:       req.MemberId,
		Downloads:        0,
		Extra:            "",
		Deleted:          data.NotDeleted,
		FileUrl:          req.FileUrl,
		FileType:         req.FileType,
		DeletedTime:      0,
	}
	if err := s.fd.SaveFile(ctx, &file); err != nil {
		return &emptypb.Empty{}, err
	}
	if err := s.sld.SaveSourceLink(ctx, &model.SourceLink{
		SourceType:       "file",
		SourceCode:       file.Id,
		LinkType:         "tasks",
		LinkCode:         taskCode,
		OrganizationCode: orgCode,
		CreateBy:         req.MemberId,
		CreateTime:       time.Now().UnixMilli(),
		Sort:             0,
	}); err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) GetTaskSources(ctx context.Context, req *task.TaskRpcRequest) (*task.TaskSourceResponse, error) {
	taskCode := encrypts.DecryptNoErr(req.TaskCode, data.AESKey)
	sourceLinks, fIdList, err := s.sld.FindLinkByTaskCode(ctx, taskCode)
	if len(sourceLinks) <= 0 {
		return &task.TaskSourceResponse{}, nil
	}
	if err != nil {
		return &task.TaskSourceResponse{}, err
	}
	_, fMap, err := s.fd.FindFileByIds(ctx, fIdList)
	if err != nil {
		return &task.TaskSourceResponse{}, nil
	}
	var list []*model.SourceLinkDisplay
	for _, v := range sourceLinks {
		list = append(list, v.ToDisplay(fMap[v.SourceCode]))
	}
	var sl []*task.TaskSource
	_ = copier.Copy(&sl, list)
	return &task.TaskSourceResponse{List: sl}, nil
}

func (s *Service) CreateComment(ctx context.Context, req *task.TaskRpcRequest) (*emptypb.Empty, error) {
	taskCode := encrypts.DecryptNoErr(req.TaskCode, data.AESKey)
	tk, err := s.td.FindTaskById(ctx, taskCode)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	if err = s.pld.CreateComment(ctx, &model.ProjectLog{
		MemberCode:   req.MemberId,
		Content:      req.CommentContent,
		Remark:       req.CommentContent,
		Type:         "createComment",
		CreateTime:   time.Now().UnixMilli(),
		SourceCode:   taskCode,
		ActionType:   "tasks",
		ToMemberCode: 0,
		IsComment:    data.Comment,
		ProjectCode:  tk.ProjectCode,
		Icon:         "plus",
		IsRobot:      0,
	}); err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}
