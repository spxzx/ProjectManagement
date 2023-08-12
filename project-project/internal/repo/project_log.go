package repo

import (
	"context"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/pkg/model"
)

type ProjectLogRepo interface {
	SaveProjectLog(ctx context.Context, conn db.Conn, pl *model.ProjectLog) error
	FindLogByTaskCode(ctx context.Context, taskCode int64, comment int) ([]*model.ProjectLog, int64, error)
	FindLogByTaskCodePage(ctx context.Context, taskCode int64, comment, page, size int) ([]*model.ProjectLog, int64, error)
	FindLogByMemberCode(background context.Context, memberId int64, page int64, size int64) ([]*model.ProjectLog, int64, error)
}
