package repo

import (
	"context"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/pkg/model"
)

type FileRepo interface {
	Save(ctx context.Context, conn db.Conn, file *model.File) error
	FindFileByIds(background context.Context, ids []int64) ([]*model.File, error)
}
