package department

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/spxzx/project-common/encrypts"
	"github.com/spxzx/project-grpc/department/department"
	"github.com/spxzx/project-project/internal/dao"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/internal/domain"
	"github.com/spxzx/project-project/internal/repo"
	"github.com/spxzx/project-project/pkg/data"
)

type Service struct {
	cache repo.Cache
	tsc   db.Transaction
	dd    *domain.DepartmentDomain
	department.UnimplementedDepartmentServiceServer
}

func New() *Service {
	return &Service{
		cache: dao.Rc,
		dd:    domain.NewDepartmentDomain(),
		tsc:   dao.NewTransaction(),
	}
}

func (s *Service) GetDepartmentList(ctx context.Context, req *department.DepartmentRpcRequest) (*department.ListDepartment, error) {
	organizationCode := encrypts.DecryptNoErr(req.OrganizationCode, data.AESKey)
	var parentDepartmentCode int64
	if req.ParentDepartmentCode != "" {
		parentDepartmentCode = encrypts.DecryptNoErr(req.ParentDepartmentCode, data.AESKey)
	}
	dps, total, err := s.dd.GetDepartmentList(
		organizationCode, parentDepartmentCode, req.Page, req.PageSize)
	if err != nil {
		return &department.ListDepartment{}, err
	}
	var list []*department.Department
	_ = copier.Copy(&list, dps)
	for _, v := range list {
		v.Code = encrypts.EncryptNoErr(v.Id, data.AESKey)
	}
	return &department.ListDepartment{List: list, Total: total}, nil
}

func (s *Service) SaveDepartment(ctx context.Context, req *department.DepartmentRpcRequest) (*department.Department, error) {
	organizationCode := encrypts.DecryptNoErr(req.OrganizationCode, data.AESKey)
	var departmentCode int64
	if req.DepartmentCode != "" {
		departmentCode = encrypts.DecryptNoErr(req.DepartmentCode, data.AESKey)
	}
	var parentDepartmentCode int64
	if req.ParentDepartmentCode != "" {
		parentDepartmentCode = encrypts.DecryptNoErr(req.ParentDepartmentCode, data.AESKey)
	}
	dp, err := s.dd.Save(ctx, organizationCode, departmentCode, parentDepartmentCode, req.Name)
	if err != nil {
		return &department.Department{}, err
	}
	var res = &department.Department{}
	_ = copier.Copy(res, dp)
	return res, nil
}

func (s *Service) ReadDepartment(ctx context.Context, req *department.DepartmentRpcRequest) (*department.Department, error) {
	departmentCode := encrypts.DecryptNoErr(req.DepartmentCode, data.AESKey)
	dp, err := s.dd.FindDepartmentById(ctx, departmentCode)
	if err != nil {
		return &department.Department{}, err
	}
	var res = &department.Department{}
	_ = copier.Copy(res, dp.ToDisplay())
	return res, nil
}
