package repo

import (
	"context"
	"github.com/spxzx/project-project/pkg/model"
)

type MenuRepo interface {
	FindAll(ctx context.Context) ([]*model.ProjectMenu, error)
}
