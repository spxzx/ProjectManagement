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
)

type TaskStagesTemplateDomain struct {
	tst repo.TaskStagesTemplateRepo
	tsc db.Transaction
}

func NewTaskStagesTemplateDomain() *TaskStagesTemplateDomain {
	return &TaskStagesTemplateDomain{
		tst: dao.NewTaskStagesTmplDao(),
		tsc: dao.NewTransaction(),
	}
}

func (t *TaskStagesTemplateDomain) FindInPojTmplIds(ctx context.Context, ids []int) ([]model.TaskStagesTemplate, error) {
	taskList, err := t.tst.FindInPojTmplIds(ctx, ids)
	if err != nil {
		zap.L().Error("getProjectTemplates db FindInPojTmplIds error, cause by: ", zap.Error(err))
		return nil, errs.GrpcError(data.DBError)
	}
	return taskList, nil
}

func (t *TaskStagesTemplateDomain) FindTaskStagesTplByPid(ctx context.Context, tplCode int64) ([]*model.TaskStagesTemplate, error) {
	tsTmpl, err := t.tst.FindTaskStagesTplByPid(ctx, tplCode)
	if err != nil {
		zap.L().Error("saveProject db FindTaskStagesTplByPid error, cause by: ", zap.Error(err))
		return nil, err
	}
	return tsTmpl, nil
}
