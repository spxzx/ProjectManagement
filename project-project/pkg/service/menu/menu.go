package menu

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/spxzx/project-grpc/menu/menu"
	"github.com/spxzx/project-project/internal/dao"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/internal/domain"
	"github.com/spxzx/project-project/internal/repo"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	cache repo.Cache
	tsc   db.Transaction
	md    *domain.MenuDomain
	menu.UnimplementedMenuServiceServer
}

func New() *Service {
	return &Service{
		cache: dao.Rc,
		tsc:   dao.NewTransaction(),
		md:    domain.NewMenuDomain(),
	}
}

func (s *Service) GetMenuList(ctx context.Context, _ *emptypb.Empty) (*menu.MenuResponse, error) {
	treeList, err := s.md.MenuTreeList(ctx)
	if err != nil {
		return &menu.MenuResponse{}, err
	}
	var list []*menu.Menu
	_ = copier.Copy(&list, treeList)
	return &menu.MenuResponse{List: list}, nil
}
