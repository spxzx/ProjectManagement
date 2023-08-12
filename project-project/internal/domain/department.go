package domain

import (
	"context"
	"github.com/spxzx/project-common/errs"
	"github.com/spxzx/project-project/internal/dao"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/internal/repo"
	"github.com/spxzx/project-project/pkg/data"
	"github.com/spxzx/project-project/pkg/model"
	"time"
)

type DepartmentDomain struct {
	dt  repo.DepartmentRepo
	tsc db.Transaction
}

func NewDepartmentDomain() *DepartmentDomain {
	return &DepartmentDomain{
		dt:  dao.NewDepartmentDao(),
		tsc: dao.NewTransaction(),
	}
}

func (d *DepartmentDomain) FindDepartmentById(ctx context.Context, id int64) (*model.Department, error) {
	return d.dt.FindDepartmentById(ctx, id)
}

func (d *DepartmentDomain) GetDepartmentList(organizationCode int64, parentDepartmentCode int64, page int64, size int64) ([]*model.DepartmentDisplay, int64, error) {
	list, total, err := d.dt.ListDepartment(organizationCode, parentDepartmentCode, page, size)
	if err != nil {
		return nil, 0, errs.GrpcError(data.DBError)
	}
	var dList []*model.DepartmentDisplay
	for _, v := range list {
		dList = append(dList, v.ToDisplay())
	}
	return dList, total, nil
}

func (d *DepartmentDomain) Save(ctx context.Context, organizationCode int64, departmentCode int64,
	parentDepartmentCode int64, name string) (*model.DepartmentDisplay, error) {
	dpm, err := d.dt.FindDepartment(ctx, organizationCode, parentDepartmentCode, name)
	if err != nil {
		return nil, errs.GrpcError(data.DBError)
	}
	if dpm == nil {
		dpm = &model.Department{
			Name:             name,
			OrganizationCode: organizationCode,
			CreateTime:       time.Now().UnixMilli(),
		}
		if parentDepartmentCode > 0 {
			dpm.Pcode = parentDepartmentCode
		}
		err := d.tsc.Action(func(conn db.Conn) error {
			return d.dt.Save(ctx, conn, dpm)
		})
		if err != nil {
			return nil, errs.GrpcError(data.DBError)
		}
		return dpm.ToDisplay(), nil
	}
	return dpm.ToDisplay(), nil
}
