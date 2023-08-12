package domain

import (
	"context"
	"github.com/spxzx/project-common/errs"
	"github.com/spxzx/project-project/internal/dao"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/internal/repo"
	"github.com/spxzx/project-project/pkg/data"
	"github.com/spxzx/project-project/pkg/model"
)

type ProjectNodeDomain struct {
	pn      repo.ProjectNodeRepo
	tsc     db.Transaction
	userRpc *UserRpcDomain
}

func NewProjectNodeDomain() *ProjectNodeDomain {
	return &ProjectNodeDomain{
		pn:      dao.NewProjectNodeDao(),
		tsc:     dao.NewTransaction(),
		userRpc: NewUserRpcDomain(),
	}
}

func (p *ProjectNodeDomain) TreeList(ctx context.Context) ([]*model.ProjectNodeTree, error) {
	nodes, err := p.pn.FindAll(ctx)
	if err != nil {
		return nil, errs.GrpcError(data.DBError)
	}
	treeList := model.ToNodeTreeList(nodes)
	return treeList, nil
}

func (d *ProjectNodeDomain) AllNodeList() ([]*model.ProjectNode, error) {
	nodes, err := d.pn.FindAll(context.Background())
	if err != nil {
		return nil, errs.GrpcError(data.DBError)
	}
	return nodes, nil
}
