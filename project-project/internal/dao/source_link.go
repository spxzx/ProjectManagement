package dao

import (
	"context"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/pkg/model"
)

type SourceLinkDao struct {
	conn *db.GORMConn
}

func (l *SourceLinkDao) Save(ctx context.Context, conn db.Conn, link *model.SourceLink) error {
	l.conn = conn.(*db.GORMConn)
	return l.conn.Session(ctx).Save(&link).Error
}

func (l *SourceLinkDao) FindLinkByTaskCode(ctx context.Context, taskCode int64) (list []*model.SourceLink, err error) {
	err = l.conn.Session(ctx).Model(&model.SourceLink{}).Where("link_code=?", taskCode).Find(&list).Error
	return
}

func NewSourceLinkDao() *SourceLinkDao {
	return &SourceLinkDao{conn: db.NewGORM()}
}
