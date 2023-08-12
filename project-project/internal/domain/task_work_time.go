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

type TaskWorkTimeDomain struct {
	twt     repo.TaskWorkTimeRepo
	tsc     db.Transaction
	userRpc *UserRpcDomain
}

func NewTaskWorkTimeDomain() *TaskWorkTimeDomain {
	return &TaskWorkTimeDomain{
		twt:     dao.NewTaskWorkTimeDao(),
		tsc:     dao.NewTransaction(),
		userRpc: NewUserRpcDomain(),
	}
}

func (t *TaskWorkTimeDomain) GetTaskWorkTimeList(ctx context.Context, req *task.TaskRpcRequest) ([]*model.TaskWorkTimeDisplay, error) {
	taskCode := encrypts.DecryptNoErr(req.TaskCode, data.AESKey)
	var list []*model.TaskWorkTime
	var err error
	list, err = t.twt.FindWorkTimeList(ctx, taskCode)
	if err != nil {
		zap.L().Error("tasks GetTaskWorkTimeList FindWorkTimeList error, cause by: ", zap.Error(err))
		return nil, errs.GrpcError(data.DBError)
	}
	if len(list) <= 0 {
		return nil, nil
	}
	var displayList []*model.TaskWorkTimeDisplay
	var mIds []int64
	for _, v := range list {
		mIds = append(mIds, v.MemberCode)
	}
	memList, _, err := t.userRpc.FindMemInfoByIds(ctx, &login.MemRequest{MemIds: mIds})
	if err != nil {
		return nil, err
	}
	mMap := make(map[int64]*login.Member)
	for _, v := range memList.List {
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
	return displayList, nil
}

func (t *TaskWorkTimeDomain) SaveWorkTime(ctx context.Context, req *task.TaskRpcRequest) error {
	return t.tsc.Action(func(conn db.Conn) error {
		if err_ := t.twt.SaveWorkTime(ctx, conn, &model.TaskWorkTime{
			TaskCode:   encrypts.DecryptNoErr(req.TaskCode, data.AESKey),
			MemberCode: req.MemberId,
			CreateTime: time.Now().UnixMilli(),
			Content:    req.Content,
			BeginTime:  req.BeginTime,
			Num:        int(req.Num),
		}); err_ != nil {
			zap.L().Error("tasks SaveTaskWorkTimeList SaveWorkTime error, cause by: ", zap.Error(err_))
			return errs.GrpcError(data.DBError)
		}
		return nil
	})
}
