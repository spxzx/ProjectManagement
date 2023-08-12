package repo

import (
	"context"
	"github.com/spxzx/project-project/pkg/model"
)

type ProjectTemplateRepo interface {
	FindPojTmplSystem(ctx context.Context, page, pageSize int64) ([]model.ProjectTemplate, int64, error)
	FindPojTmplCustom(ctx context.Context, memId, orgCode, page, pageSize int64) ([]model.ProjectTemplate, int64, error)
	FindPojTmplAll(ctx context.Context, orgCode, page, pageSize int64) ([]model.ProjectTemplate, int64, error)
}
