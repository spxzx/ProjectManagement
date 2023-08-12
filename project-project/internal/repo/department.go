package repo

import (
	"context"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/pkg/model"
)

type DepartmentRepo interface {
	FindDepartmentById(ctx context.Context, id int64) (*model.Department, error)
	FindDepartment(ctx context.Context, organizationCode int64, parentDepartmentCode int64, name string) (*model.Department, error)
	Save(ctx context.Context, conn db.Conn, dpm *model.Department) error
	ListDepartment(organizationCode int64, parentDepartmentCode int64, page int64, size int64) (list []*model.Department, total int64, err error)
}
