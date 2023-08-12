package dao

import (
	"context"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/pkg/model"
)

type ProjectAuthNodeDao struct {
	conn *db.GORMConn
}

func (p *ProjectAuthNodeDao) DeleteByAuthId(ctx context.Context, conn db.Conn, authId int64) error {
	p.conn = conn.(*db.GORMConn)
	tx := p.conn.Tx(ctx)
	err := tx.Where("auth=?", authId).Delete(&model.ProjectAuthNode{}).Error
	return err
}

func (p *ProjectAuthNodeDao) Save(ctx context.Context, conn db.Conn, authId int64, nodes []string) error {
	p.conn = conn.(*db.GORMConn)
	tx := p.conn.Tx(ctx)
	var list []*model.ProjectAuthNode
	for _, v := range nodes {
		pn := &model.ProjectAuthNode{
			Auth: authId,
			Node: v,
		}
		list = append(list, pn)
	}
	err := tx.Create(list).Error
	return err
}

func NewProjectAuthNodeDao() *ProjectAuthNodeDao {
	return &ProjectAuthNodeDao{
		conn: db.NewGORM(),
	}
}

func (p *ProjectAuthNodeDao) FindNodeStringList(ctx context.Context, authId int64) (list []string, err error) {
	err = p.conn.Session(ctx).Model(&model.ProjectAuthNode{}).Where("auth=?", authId).Select("node").Find(&list).Error
	return
}
