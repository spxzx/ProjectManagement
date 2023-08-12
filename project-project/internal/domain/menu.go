package domain

import (
	"context"
	"github.com/spxzx/project-common/errs"
	"github.com/spxzx/project-project/internal/dao"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/internal/repo"
	"github.com/spxzx/project-project/pkg/data"
	"github.com/spxzx/project-project/pkg/model"
	"go.uber.org/zap"
)

type MenuDomain struct {
	menu repo.MenuRepo
	tsc  db.Transaction
}

func NewMenuDomain() *MenuDomain {
	return &MenuDomain{
		menu: dao.NewMenuDao(),
		tsc:  dao.NewTransaction(),
	}
}

func (m *MenuDomain) FindAll(ctx context.Context) ([]*model.ProjectMenu, error) {
	menuList, err := m.menu.FindAll(ctx)
	if err != nil {
		zap.L().Error("index db FindAll error, cause by: ", zap.Error(err))
		return nil, errs.GrpcError(data.DBError)
	}
	return menuList, nil
}

func (m *MenuDomain) MenuTreeList(ctx context.Context) ([]*model.ProjectMenuChild, error) {
	menus, err := m.menu.FindAll(ctx)
	if err != nil {
		zap.L().Error("MenuList error", zap.Error(err))
		return nil, errs.GrpcError(data.DBError)
	}
	menuChildren := model.ConvertChild(menus)
	return menuChildren, nil
}
