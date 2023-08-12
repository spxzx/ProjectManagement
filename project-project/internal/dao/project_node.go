package dao

import (
	"context"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/pkg/model"
)

type ProjectNodeDao struct {
	conn *db.GORMConn
}

func NewProjectNodeDao() *ProjectNodeDao {
	return &ProjectNodeDao{conn: db.NewGORM()}
}

func (p *ProjectNodeDao) FindAll(ctx context.Context) (pms []*model.ProjectNode, err error) {
	err = p.conn.Session(ctx).Model(&model.ProjectNode{}).Find(&pms).Error
	return
}
