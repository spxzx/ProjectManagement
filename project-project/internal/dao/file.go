package dao

import (
	"context"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/pkg/model"
)

type FileDao struct {
	conn *db.GORMConn
}

func (f *FileDao) Save(ctx context.Context, conn db.Conn, file *model.File) error {
	f.conn = conn.(*db.GORMConn)
	return f.conn.Tx(ctx).Save(&file).Error
}

func (f *FileDao) FindFileByIds(ctx context.Context, ids []int64) (list []*model.File, err error) {
	err = f.conn.Session(ctx).Model(&model.File{}).Where("id in (?)", ids).Find(&list).Error
	return
}

func NewFileDao() *FileDao {
	return &FileDao{conn: db.NewGORM()}
}
