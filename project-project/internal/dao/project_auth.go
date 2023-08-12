package dao

import (
	"context"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/pkg/model"
)

type ProjectAuthDao struct {
	conn *db.GORMConn
}

func (p *ProjectAuthDao) FindAuthListPage(ctx context.Context, code int64, page int64, size int64) (list []*model.ProjectAuth, total int64, err error) {
	session := p.conn.Session(ctx)
	err = session.Model(&model.ProjectAuth{}).
		Where("organization_code=?", code).
		Limit(int(size)).
		Offset(int((page - 1) * size)).
		Find(&list).Error
	session.Model(&model.ProjectAuth{}).
		Where("organization_code=?", code).
		Count(&total)
	return
}

func (p *ProjectAuthDao) FindAuthList(ctx context.Context, orgCode int64) (list []*model.ProjectAuth, err error) {
	err = p.conn.Session(ctx).Model(&model.ProjectAuth{}).
		Where("organization_code=? and status=1", orgCode).Find(&list).Error
	return
}

func NewProjectAuthDao() *ProjectAuthDao {
	return &ProjectAuthDao{conn: db.NewGORM()}
}
