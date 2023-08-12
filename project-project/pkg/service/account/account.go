package account

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/spxzx/project-common/encrypts"
	"github.com/spxzx/project-grpc/account/account"
	"github.com/spxzx/project-project/internal/dao"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/internal/domain"
	"github.com/spxzx/project-project/internal/repo"
	"github.com/spxzx/project-project/pkg/data"
)

type Service struct {
	cache repo.Cache
	tsc   db.Transaction
	mad   *domain.MemberAccountDomain
	pad   *domain.ProjectAuthDomain
	account.UnimplementedAccountServiceServer
}

func New() *Service {
	return &Service{
		cache: dao.Rc,
		tsc:   dao.NewTransaction(),
		mad:   domain.NewMemberAccountDomain(),
		pad:   domain.NewProjectAuthDomain(),
	}
}

func (s *Service) GetAccountList(ctx context.Context, req *account.AccountRpcRequest) (*account.AccountResponse, error) {
	aList, total, err := s.mad.GetAccountList(ctx,
		req.OrganizationCode, req.MemberId, req.Page,
		req.PageSize, req.DepartmentCode, req.SearchType,
	)
	if err != nil {
		return &account.AccountResponse{}, err
	}
	authList, err := s.pad.GetAuthList(ctx, encrypts.DecryptNoErr(req.OrganizationCode, data.AESKey))
	if err != nil {
		return &account.AccountResponse{}, err
	}
	var maList []*account.MemberAccount
	_ = copier.Copy(&maList, aList)
	var prList []*account.ProjectAuth
	_ = copier.Copy(&prList, authList)
	return &account.AccountResponse{
		AccountList: maList,
		AuthList:    prList,
		Total:       total,
	}, nil
}
