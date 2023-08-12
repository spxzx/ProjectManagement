package domain

import (
	"context"
	"github.com/spxzx/project-common/encrypts"
	"github.com/spxzx/project-common/errs"
	"github.com/spxzx/project-grpc/project/project"
	"github.com/spxzx/project-project/internal/dao"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/internal/repo"
	"github.com/spxzx/project-project/pkg/data"
	"github.com/spxzx/project-project/pkg/model"
	"go.uber.org/zap"
	"strconv"
)

type ProjectTemplateDomain struct {
	tpl repo.ProjectTemplateRepo
	tsc db.Transaction
}

func NewProjectTemplateDomain() *ProjectTemplateDomain {
	return &ProjectTemplateDomain{
		tpl: dao.NewProjectTemplateDao(),
		tsc: dao.NewTransaction(),
	}
}
func (p *ProjectTemplateDomain) FindProjectTemplates(ctx context.Context, req *project.ProjectRpcRequest) ([]model.ProjectTemplate, int64, error) {
	orgCodeStr, _ := encrypts.Decrypt(req.OrganizationCode, data.AESKey)
	orgCode, _ := strconv.ParseInt(orgCodeStr, 10, 64)
	var ptList []model.ProjectTemplate
	var total int64
	var err error
	switch req.ViewType {
	case -1:
		ptList, total, err = p.tpl.FindPojTmplAll(ctx, orgCode, req.Page, req.PageSize)
	case 0:
		ptList, total, err = p.tpl.FindPojTmplCustom(ctx, req.MemberId, orgCode, req.Page, req.PageSize)
	case 1:
		ptList, total, err = p.tpl.FindPojTmplSystem(ctx, req.Page, req.PageSize)
	}
	if err != nil {
		zap.L().Error("getProjectTemplates db find with view_type error, cause by: ", zap.Error(err))
		return []model.ProjectTemplate{}, 0, errs.GrpcError(data.DBError)
	}
	return ptList, total, nil
}
