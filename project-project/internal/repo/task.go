package repo

import (
	"context"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/pkg/model"
)

type TaskRepo interface {
	FindTaskByStageCode(ctx context.Context, code int64) ([]*model.Task, error)
	FindTaskMemberByTid(ctx context.Context, tid, mid int64) (*model.TaskMember, error)
	FindTaskMaxIdNum(ctx context.Context, code int64) (*int, error)
	FindTaskMaxSort(ctx context.Context, pCode int64, TSCode int64) (*int, error)
	SaveTask(ctx context.Context, conn db.Conn, ts *model.Task) error
	SaveTaskMember(ctx context.Context, conn db.Conn, tm *model.TaskMember) error
	FindTaskById(ctx context.Context, id int64) (*model.Task, error)
	UpdateTaskSort(ctx context.Context, conn db.Conn, next *model.Task) error
	FindTaskLessThenSortByStageCode(ctx context.Context, code, sort int) (*model.Task, error)
	FindTaskByAssignTo(ctx context.Context, id int64, t int32, page, size int64) ([]*model.Task, int64, error)
	FindTaskByMemCode(ctx context.Context, id int64, t int32, page, size int64) ([]*model.Task, int64, error)
	FindTaskByCreateBy(ctx context.Context, id int64, t int32, page, size int64) ([]*model.Task, int64, error)
	FindTaskMemberList(ctx context.Context, code int64, page int64, size int64) ([]*model.TaskMember, int64, error)
	FindTaskByIds(ctx context.Context, list []int64) ([]*model.Task, error)
}
