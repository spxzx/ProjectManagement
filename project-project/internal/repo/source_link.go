package repo

import (
	"context"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/pkg/model"
)

type SourceLinkRepo interface {
	Save(ctx context.Context, conn db.Conn, link *model.SourceLink) error
	FindLinkByTaskCode(ctx context.Context, taskCode int64) ([]*model.SourceLink, error)
}
