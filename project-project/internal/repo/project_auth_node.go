package repo

import (
	"context"
	"github.com/spxzx/project-project/internal/db"
)

type ProjectAuthNodeRepo interface {
	FindNodeStringList(ctx context.Context, authId int64) ([]string, error)
	DeleteByAuthId(background context.Context, conn db.Conn, authId int64) error
	Save(background context.Context, conn db.Conn, authId int64, nodes []string) error
}
