package repo

import (
	"context"
	"github.com/spxzx/project-project/pkg/model"
)

type ProjectAuthRepo interface {
	FindAuthList(ctx context.Context, orgCode int64) ([]*model.ProjectAuth, error)
	FindAuthListPage(ctx context.Context, code int64, page int64, size int64) (list []*model.ProjectAuth, total int64, err error)
}
