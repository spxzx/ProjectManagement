package repo

import (
	"context"
	"github.com/spxzx/project-project/pkg/model"
)

type ProjectNodeRepo interface {
	FindAll(ctx context.Context) (pms []*model.ProjectNode, err error)
}
