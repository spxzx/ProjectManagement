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

type FileDomain struct {
	file repo.FileRepo
	tsc  db.Transaction
}

func NewFileDomain() *FileDomain {
	return &FileDomain{
		file: dao.NewFileDao(),
		tsc:  dao.NewTransaction(),
	}
}

func (fd *FileDomain) SaveFile(ctx context.Context, file *model.File) error {
	return fd.tsc.Action(func(conn db.Conn) error {
		if err := fd.file.Save(ctx, conn, file); err != nil {
			zap.L().Error("tasks UploadFile db file Save error, cause by: ", zap.Error(err))
			return errs.GrpcError(data.DBError)
		}
		return nil
	})
}

func (fd *FileDomain) FindFileByIds(ctx context.Context, fIdList []int64) ([]*model.File, map[int64]*model.File, error) {
	files, err := fd.file.FindFileByIds(ctx, fIdList)
	if err != nil {
		zap.L().Error("tasks GetTaskSources FindFileByIds error, cause by: ", zap.Error(err))
		return nil, nil, errs.GrpcError(data.DBError)
	}
	fMap := make(map[int64]*model.File)
	for _, v := range files {
		fMap[v.Id] = v
	}
	return files, fMap, nil
}
