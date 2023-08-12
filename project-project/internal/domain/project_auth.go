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
	"strconv"
	"time"
)

type ProjectAuthDomain struct {
	pa  repo.ProjectAuthRepo
	pan repo.ProjectAuthNodeRepo
	pnd *ProjectNodeDomain
	mad *MemberAccountDomain
}

func NewProjectAuthDomain() *ProjectAuthDomain {
	return &ProjectAuthDomain{
		pa:  dao.NewProjectAuthDao(),
		pan: dao.NewProjectAuthNodeDao(),
		pnd: NewProjectNodeDomain(),
		mad: NewMemberAccountDomain(),
	}
}

func (d *ProjectAuthDomain) GetAuthList(ctx context.Context, orgCode int64) ([]*model.ProjectAuthDisplay, error) {
	list, err := d.pa.FindAuthList(ctx, orgCode)
	if err != nil {
		zap.L().Error("project AuthList projectAuthRepo.FindAuthList error", zap.Error(err))
		return nil, errs.GrpcError(data.DBError)
	}
	var pdList []*model.ProjectAuthDisplay
	for _, v := range list {
		display := v.ToDisplay()
		pdList = append(pdList, display)
	}
	return pdList, nil
}

func (d *ProjectAuthDomain) AuthListPage(orgCode int64, page int64, pageSize int64) ([]*model.ProjectAuthDisplay, int64, error) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	list, total, err := d.pa.FindAuthListPage(c, orgCode, page, pageSize)
	if err != nil {
		zap.L().Error("project AuthList projectAuthRepo.FindAuthList error", zap.Error(err))
		return nil, 0, errs.GrpcError(data.DBError)
	}
	var pdList []*model.ProjectAuthDisplay
	for _, v := range list {
		display := v.ToDisplay()
		pdList = append(pdList, display)
	}
	return pdList, total, nil
}

func (d *ProjectAuthDomain) AllNodeAndAuth(ctx context.Context, authId int64) ([]*model.ProjectNodeAuthTree, []string, error) {
	treeList, err := d.pnd.AllNodeList()
	if err != nil {
		return nil, nil, err
	}

	authNodeList, err := d.pan.FindNodeStringList(ctx, authId)
	if err != nil {
		return nil, nil, err
	}
	list := model.ToAuthNodeTreeList(treeList, authNodeList)
	return list, authNodeList, nil
}

func (d *ProjectAuthDomain) Save(conn db.Conn, authId int64, nodes []string) error {
	err := d.pan.DeleteByAuthId(context.Background(), conn, authId)
	if err != nil {
		return errs.GrpcError(data.DBError)
	}
	err = d.pan.Save(context.Background(), conn, authId, nodes)
	if err != nil {
		return errs.GrpcError(data.DBError)
	}
	return nil
}

func (d *ProjectAuthDomain) AuthNodes(ctx context.Context, memberId int64) ([]string, error) {
	account, err := d.mad.FindAccount(memberId)
	if err != nil {
		return nil, err
	}
	authorize := account.Authorize
	authId, _ := strconv.ParseInt(authorize, 10, 64)
	authNodeList, dbErr := d.pan.FindNodeStringList(ctx, authId)
	if dbErr != nil {
		return nil, dbErr
	}
	return authNodeList, nil
}
