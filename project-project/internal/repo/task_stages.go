package repo

import (
	"context"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/pkg/model"
)

type TaskStagesRepo interface {
	SaveTaskStages(ctx context.Context, conn db.Conn, ts *model.TaskStages) error
	FindTaskStagesByPCode(ctx context.Context, code, page, pageSize int64) ([]*model.TaskStages, int64, error)
	FindTaskStagesById(ctx context.Context, code int64) (*model.TaskStages, error)
}
