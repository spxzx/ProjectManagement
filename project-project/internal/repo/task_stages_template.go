package repo

import (
	"context"
	"github.com/spxzx/project-project/pkg/model"
)

type TaskStagesTemplateRepo interface {
	FindInPojTmplIds(ctx context.Context, ids []int) ([]model.TaskStagesTemplate, error)
	FindTaskStagesTplByPid(ctx context.Context, code int64) ([]*model.TaskStagesTemplate, error)
}
