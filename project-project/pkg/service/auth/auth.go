package auth

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/spxzx/project-common/encrypts"
	"github.com/spxzx/project-grpc/auth/auth"
	"github.com/spxzx/project-project/internal/dao"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/internal/domain"
	"github.com/spxzx/project-project/internal/repo"
	"github.com/spxzx/project-project/pkg/data"
)

type Service struct {
	cache repo.Cache
	tsc   db.Transaction
	pad   *domain.ProjectAuthDomain
	auth.UnimplementedAuthServiceServer
}

func New() *Service {
	return &Service{
		cache: dao.Rc,
		tsc:   dao.NewTransaction(),
		pad:   domain.NewProjectAuthDomain(),
	}
}

func (s *Service) AuthList(ctx context.Context, req *auth.AuthRpcRequest) (*auth.ListAuth, error) {
	organizationCode := encrypts.DecryptNoErr(req.OrganizationCode, data.AESKey)
	listPage, total, err := s.pad.AuthListPage(organizationCode, req.Page, req.PageSize)
	if err != nil {
		return &auth.ListAuth{}, err
	}
	var prList []*auth.ProjectAuth
	_ = copier.Copy(&prList, listPage)
	return &auth.ListAuth{List: prList, Total: total}, nil
}

func (s *Service) Apply(ctx context.Context, req *auth.AuthRpcRequest) (*auth.ApplyResponse, error) {
	if req.Action == "getnode" {
		//获取列表
		list, checkedList, err := s.pad.AllNodeAndAuth(ctx, req.Id)
		if err != nil {
			return &auth.ApplyResponse{}, err
		}
		var prList []*auth.ProjectNode
		_ = copier.Copy(&prList, list)
		return &auth.ApplyResponse{List: prList, CheckedList: checkedList}, nil
	}
	if req.Action == "save" {
		//保存
		//先删在存 加事务
		if err := s.tsc.Action(func(conn db.Conn) error {
			err := s.pad.Save(conn, req.Id, req.Nodes)
			return err
		}); err != nil {
			return &auth.ApplyResponse{}, err
		}
	}
	return &auth.ApplyResponse{}, nil
}

func (s *Service) AuthNodesByMemberId(ctx context.Context, req *auth.AuthRpcRequest) (*auth.AuthNodesResponse, error) {
	list, err := s.pad.AuthNodes(ctx, req.MemberId)
	if err != nil {
		return nil, err
	}
	return &auth.AuthNodesResponse{List: list}, nil
}
