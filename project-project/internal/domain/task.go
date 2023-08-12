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

type TaskDomain struct {
	task    repo.TaskRepo
	tsc     db.Transaction
	userRpc *UserRpcDomain
}

func NewTaskDomain() *TaskDomain {
	return &TaskDomain{
		task:    dao.NewTaskDao(),
		tsc:     dao.NewTransaction(),
		userRpc: NewUserRpcDomain(),
	}
}

func (t *TaskDomain) FindTaskByIds(ctx context.Context, ids []int64) ([]*model.Task, map[int64]*model.Task, error) {
	tasks, err := t.task.FindTaskByIds(ctx, ids)
	if err != nil {
		zap.L().Error("project GetLogBySelfProject FindTaskByIds error, cause by: ", zap.Error(err))
		return nil, nil, errs.GrpcError(data.DBError)
	}
	tMap := make(map[int64]*model.Task)
	for _, v := range tasks {
		tMap[v.Id] = v
	}
	return tasks, tMap, nil
}

func (t *TaskDomain) GetTaskStageDetailList(ctx context.Context, req *task.TaskRpcRequest) ([]*model.TaskDisplay, error) {
	stageCode := encrypts.DecryptNoErr(req.StageCode, data.AESKey)
	taskList, err := t.task.FindTaskByStageCode(ctx, stageCode)
	if err != nil {
		zap.L().Error("task GetTaskStageDetailList db FindTaskByStageCode error, cause by", zap.Error(err))
		return nil, errs.GrpcError(data.DBError)
	}
	var tdList []*model.TaskDisplay
	var mIds []int64
	for _, v := range taskList {
		d := v.ToTaskDisplay()
		if v.Private == 1 {
			tm, err_ := t.task.FindTaskMemberByTid(ctx, v.Id, req.MemberId)
			if err_ != nil {
				zap.L().Error("task GetTaskStageDetailList db FindTaskMemberByTid error, cause by", zap.Error(err_))
				return nil, errs.GrpcError(data.DBError)
			}
			if tm == nil {
				d.CanRead = data.NotCanRead
			}
		}
		tdList = append(tdList, d)
		mIds = append(mIds, v.AssignTo)
	}

	if mIds == nil || len(mIds) <= 0 {
		return nil, nil
	}

	_, memMap, err := t.userRpc.FindMemInfoByIds(ctx, &login.MemRequest{MemIds: mIds})
	if err != nil {
		zap.L().Error("task GetTaskStageDetailList db LoginServiceClient.FindMemInfoByIds error, cause by", zap.Error(err))
		return nil, err
	}
	for _, v := range tdList {
		mem := memMap[encrypts.DecryptNoErr(v.AssignTo, data.AESKey)]
		ex := model.Executor{
			Name:   mem.Name,
			Avatar: mem.Avatar,
		}
		v.Executor = ex
	}
	return tdList, nil
}

func (t *TaskDomain) FindTaskMaxIdNum(ctx context.Context, pCode int64) (*int, error) {
	maxIdNum, err := t.task.FindTaskMaxIdNum(ctx, pCode)
	if err != nil {
		zap.L().Error("tasks SaveTask FindTaskMaxIdNum error, cause by: ", zap.Error(err))
		return nil, errs.GrpcError(data.DBError)
	}
	if maxIdNum == nil {
		a := 1
		maxIdNum = &a
	}
	return maxIdNum, nil
}

func (t *TaskDomain) FindTaskMaxSort(ctx context.Context, pCode, sCode int64) (*int, error) {
	maxSort, err := t.task.FindTaskMaxSort(ctx, pCode, sCode)
	if err != nil {
		zap.L().Error("tasks SaveTask FindTaskMaxSort error, cause by: ", zap.Error(err))
		return nil, errs.GrpcError(data.DBError)
	}
	if maxSort == nil {
		a := 0
		maxSort = &a
	}
	return maxSort, nil
}

func (t *TaskDomain) SaveTask(ctx context.Context, tsk *model.Task, memId int64) error {
	return t.tsc.Action(func(conn db.Conn) error {
		if err := t.task.SaveTask(ctx, conn, tsk); err != nil {
			zap.L().Error("tasks SaveTask SaveTask error, cause by: ", zap.Error(err))
			return err
		}
		tm := &model.TaskMember{
			MemberCode: memId,
			TaskCode:   tsk.Id,
			JoinTime:   time.Now().UnixMilli(),
			IsOwner:    data.Owner,
		}
		if tsk.AssignTo == memId {
			tm.IsExecutor = data.Executor
		}
		if err := t.task.SaveTaskMember(ctx, conn, tm); err != nil {
			zap.L().Error("tasks SaveTask SaveTaskMember error, cause by: ", zap.Error(err))
			return err
		}
		return nil
	})
}

func (t *TaskDomain) FindTaskById(ctx context.Context, id int64) (*model.Task, error) {
	tk, err := t.task.FindTaskById(ctx, id)
	if err != nil {
		zap.L().Error("tasks TaskDomain db FindTaskById error", zap.Error(err))
		return nil, errs.GrpcError(data.DBError)
	}
	return tk, nil
}

func (t *TaskDomain) UpdateTaskSort(ctx context.Context, conn db.Conn, tk *model.Task) error {
	if err := t.task.UpdateTaskSort(ctx, conn, tk); err != nil {
		zap.L().Error("tasks TaskDomain db UpdateTaskSort error", zap.Error(err))
		return errs.GrpcError(data.DBError)
	}
	return nil
}

func (t *TaskDomain) MoveAndSortTask(ctx context.Context, preTaskCode int64, reqNextTaskCode string, toStageCode int64) error {
	// 1.从小到大排
	// 2.原有的顺序 比如 1 2 3 4 5 4 排到 2 前面去 4 的序号在 1 和 2 之间
	//   如果 4 是最后一个 保证 4 比 所有的序号都打 如果 排到第一位 直接置为 0
	tk, err := t.FindTaskById(ctx, preTaskCode)
	if err != nil {
		return err
	}
	return t.tsc.Action(func(conn db.Conn) error {
		tk.StageCode = int(toStageCode)
		if reqNextTaskCode != "" { // 有下一个任务步骤，自己不是最后的那一个
			nextTaskCode := encrypts.DecryptNoErr(reqNextTaskCode, data.AESKey)
			next, err_ := t.FindTaskById(ctx, nextTaskCode)
			if err_ != nil {
				return errs.GrpcError(data.DBError)
			}
			low, err_ := t.task.FindTaskLessThenSortByStageCode(ctx, next.StageCode, next.Sort)
			if err_ != nil {
				zap.L().Error("tasks MoveTask db FindTaskLessThenSortByStageCode error", zap.Error(err_))
				return errs.GrpcError(data.DBError)
			}
			if low != nil { // 中间
				tk.Sort = (low.Sort + next.Sort) / 2
			} else { // 第一位
				tk.Sort = 0
			}
		} else { // 最后一位
			maxSort, err_ := t.FindTaskMaxSort(ctx, tk.ProjectCode, int64(tk.StageCode))
			if err_ != nil {
				return err
			}
			tk.Sort = *maxSort + 65536
		}
		if tk.Sort < 50 { // 重置排序，防止出错
			list, err_ := t.task.FindTaskByStageCode(ctx, toStageCode)
			if err != nil {
				zap.L().Error("tasks MoveTask db FindTaskByStageCode error", zap.Error(err_))
				return errs.GrpcError(data.DBError)
			}
			iSort := 65536
			for i, v := range list {
				v.Sort = (i + 1) * iSort
				if err_ = t.UpdateTaskSort(ctx, conn, v); err_ != nil {
					return err_
				}
			}
		}
		if err_ := t.UpdateTaskSort(ctx, conn, tk); err_ != nil {
			return err_
		}
		return nil
	})
}

func (t *TaskDomain) FindTaskByCondition(ctx context.Context, req *task.TaskRpcRequest,
) ([]*model.Task, int64, []int64, []int64, error) {
	var tkList []*model.Task
	var err error
	var total int64
	switch req.TaskType {
	case data.AssignTo:
		tkList, total, err = t.task.FindTaskByAssignTo(ctx, req.MemberId, req.Type, req.Page, req.PageSize)
	case data.MemberCode:
		tkList, total, err = t.task.FindTaskByMemCode(ctx, req.MemberId, req.Type, req.Page, req.PageSize)
	case data.CreateBy:
		tkList, total, err = t.task.FindTaskByCreateBy(ctx, req.MemberId, req.Type, req.Page, req.PageSize)
	}
	if err != nil {
		zap.L().Error("task GetSelfTaskList FindTaskByXXX error, cause by: ", zap.Error(err))
		return nil, 0, nil, nil, errs.GrpcError(data.DBError)
	}
	if tkList == nil || len(tkList) == 0 {
		return nil, 0, nil, nil, nil
	}
	var pids []int64
	var mids []int64
	for _, v := range tkList {
		pids = append(pids, v.ProjectCode)
		mids = append(mids, v.AssignTo)
	}
	return tkList, total, pids, mids, nil
}

func (t *TaskDomain) ReadTask(ctx context.Context, req *task.TaskRpcRequest) (*model.Task, *model.TaskDisplay, error) {
	taskInfo, err := t.FindTaskById(ctx, encrypts.DecryptNoErr(req.TaskCode, data.AESKey))
	if err != nil {
		return nil, nil, err
	}
	if taskInfo == nil {
		return nil, nil, nil
	}
	td := taskInfo.ToTaskDisplay()
	if taskInfo.Private == data.Private {
		taskMember, err := t.task.FindTaskMemberByTid(ctx, taskInfo.Id, req.MemberId)
		if err != nil {
			zap.L().Error("tasks ReadTask db FindTaskMemberByTid error, cause by: ", zap.Error(err))
			return nil, nil, errs.GrpcError(data.DBError)
		}
		if taskMember != nil {
			td.CanRead = data.CanRead
		} else {
			td.CanRead = data.NotCanRead
		}
	}
	return taskInfo, td, nil
}

func (t *TaskDomain) GetTaskMemberList(ctx context.Context, tCode, page, size int64) ([]*model.TaskMember, int64, map[int64]*login.Member, error) {
	tmList, total, err := t.task.FindTaskMemberList(ctx, tCode, page, size)
	if err != nil {
		zap.L().Error("tasks GetTaskMemberList db FindTaskMemberList error, cause by: ", zap.Error(err))
		return nil, 0, nil, errs.GrpcError(data.DBError)
	}
	var mids []int64
	for _, v := range tmList {
		mids = append(mids, v.MemberCode)
	}
	memList, _, err := t.userRpc.FindMemInfoByIds(ctx, &login.MemRequest{MemIds: mids})
	if err != nil {
		zap.L().Error("tasks GetTaskMemberList FindMemInfoByIds error, cause by: ", zap.Error(err))
		return nil, 0, nil, err
	}
	mMap := make(map[int64]*login.Member, len(memList.List))
	for _, v := range memList.List {
		mMap[v.Id] = v
	}
	return tmList, total, mMap, nil
}

func (t *TaskDomain) FindProjectIdByTaskId(taskId int64) (int64, bool, error) {
	tk, err := t.task.FindTaskById(context.Background(), taskId)
	if err != nil {
		return 0, false, errs.GrpcError(data.DBError)
	}
	if tk == nil {
		return 0, false, nil
	}
	return tk.ProjectCode, true, nil
}
