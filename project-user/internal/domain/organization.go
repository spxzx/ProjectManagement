package domain

import (
	"context"
	"github.com/spxzx/project-common/errs"
	"github.com/spxzx/project-grpc/user/login"
	"github.com/spxzx/project-user/internal/dao"
	"github.com/spxzx/project-user/internal/db"
	"github.com/spxzx/project-user/internal/repo"
	"github.com/spxzx/project-user/pkg/data"
	"github.com/spxzx/project-user/pkg/model"
	"go.uber.org/zap"
	"time"
)

type OrganizationDomain struct {
	org repo.OrganizationRepo
	tsc db.Transaction
}

func NewOrganizationDomain() *OrganizationDomain {
	return &OrganizationDomain{
		org: dao.NewOrganizationDao(),
		tsc: dao.NewTransaction(),
	}
}

func (o *OrganizationDomain) SaveOrganization(req *login.RegisterRequest, mem *model.Member) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return o.tsc.Action(func(conn db.Conn) error {
		if err := o.org.SaveOrganization(ctx, conn, &model.Organization{
			Name:       req.Name + "个人项目",
			Avatar:     "",
			MemberId:   mem.Id,
			CreateTime: time.Now().UnixMilli(),
			Personal:   data.Personal,
		}); err != nil {
			zap.L().Error("register db SaveOrganization error, cause by: ", zap.Error(err))
			return errs.GrpcError(data.DBError)
		}
		return nil
	})
}

func (o *OrganizationDomain) FindOrganizationByMemberId(id int64) ([]model.Organization, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	org, err := o.org.FindOrganizationByMemberId(ctx, id)
	if err != nil {
		zap.L().Error("login db FindOrganizationByMemberId error, cause by: ", zap.Error(err))
		return org, err
	}
	return org, nil
}
