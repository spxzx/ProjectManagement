package repo

import (
	"context"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/pkg/model"
)

type TaskWorkTimeRepo interface {
	FindWorkTimeList(ctx context.Context, code int64) ([]*model.TaskWorkTime, error)
	SaveWorkTime(ctx context.Context, conn db.Conn, twt *model.TaskWorkTime) error
}
