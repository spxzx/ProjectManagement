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

type SourceLinkDomain struct {
	link repo.SourceLinkRepo
	tsc  db.Transaction
}

func NewSourceLinkDomain() *SourceLinkDomain {
	return &SourceLinkDomain{
		link: dao.NewSourceLinkDao(),
		tsc:  dao.NewTransaction(),
	}
}

func (s *SourceLinkDomain) SaveSourceLink(ctx context.Context, sl *model.SourceLink) error {
	return s.tsc.Action(func(conn db.Conn) error {
		if err := s.link.Save(ctx, conn, sl); err != nil {
			zap.L().Error("tasks UploadFile db link Save error, cause by: ", zap.Error(err))
			return errs.GrpcError(data.DBError)
		}
		return nil
	})
}

func (s *SourceLinkDomain) FindLinkByTaskCode(ctx context.Context, tCode int64) ([]*model.SourceLink, []int64, error) {
	sourceLinks, err := s.link.FindLinkByTaskCode(ctx, tCode)
	if err != nil {
		zap.L().Error("tasks GetTaskSources db FindLinkByTaskCode error, cause by: ", zap.Error(err))
		return nil, nil, errs.GrpcError(data.DBError)
	}
	var fIdList []int64
	for _, v := range sourceLinks {
		fIdList = append(fIdList, v.SourceCode)
	}
	return sourceLinks, fIdList, nil
}
