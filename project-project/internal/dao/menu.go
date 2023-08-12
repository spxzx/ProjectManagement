package dao

import (
	"context"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/pkg/model"
)

type MenuDao struct {
	conn *db.GORMConn
}

func (m *MenuDao) FindAll(ctx context.Context) (list []*model.ProjectMenu, err error) {
	err = m.conn.Session(ctx).Order("pid,sort asc, id asc").Find(&list).Error
	return
}

func NewMenuDao() *MenuDao {
	return &MenuDao{conn: db.NewGORM()}
}
