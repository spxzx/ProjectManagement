package domain

import (
	"context"
	"github.com/spxzx/project-common/errs"
	"github.com/spxzx/project-project/internal/dao"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/internal/repo"
	"github.com/spxzx/project-project/pkg/data"
	"github.com/spxzx/project-project/pkg/model"
	"go.uber.org/zap"
	"time"
)

type TaskStagesDomain struct {
	tss repo.TaskStagesRepo
	tsc db.Transaction
}

func NewTaskStagesDomain() *TaskStagesDomain {
	return &TaskStagesDomain{
		tss: dao.NewTaskStagesDao(),
		tsc: dao.NewTransaction(),
	}
}

func (t *TaskStagesDomain) SaveTaskStages(ctx context.Context, pj *model.Project, idx int, tst *model.TaskStagesTemplate) error {
	return t.tsc.Action(func(conn db.Conn) error {
		if err := t.tss.SaveTaskStages(ctx, conn, &model.TaskStages{
			Name:        tst.Name,
			ProjectCode: pj.Id,
			Sort:        idx + 1,
			CreateTime:  time.Now().UnixMilli(),
			Deleted:     data.NotDeleted,
		}); err != nil {
			zap.L().Error("saveProject db SaveTaskStages error, cause by: ", zap.Error(err))
			return data.DBError
		}
		return nil
	})
}

func (t *TaskStagesDomain) FindTaskStagesByPCode(ctx context.Context, pCode, page, size int64,
) ([]*model.TaskStages, map[int]*model.TaskStages, int64, error) {
	list, total, err := t.tss.FindTaskStagesByPCode(ctx, pCode, page, size)
	if err != nil {
		zap.L().Error("TaskStagesDomain db FindTaskStagesByPCode error, cause by: ", zap.Error(err))
		return nil, nil, 0, errs.GrpcError(data.DBError)
	}
	tssMap := model.ToTSSMap(list)
	return list, tssMap, total, nil
}

func (t *TaskStagesDomain) FindTaskStagesById(ctx context.Context, id int64) (*model.TaskStages, error) {
	taskStage, err := t.tss.FindTaskStagesById(ctx, id)
	if err != nil {
		zap.L().Error("tasks SaveTask FindTaskStagesById error, cause by: ", zap.Error(err))
		return nil, errs.GrpcError(data.DBError)
	}
	if taskStage == nil {
		return nil, errs.GrpcError(data.TaskStagesNotNull)
	}
	return taskStage, nil
}
