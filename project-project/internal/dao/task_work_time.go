package dao

import (
	"context"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/pkg/model"
)

type TaskWorkTimeDao struct {
	conn *db.GORMConn
}

func (t *TaskWorkTimeDao) SaveWorkTime(ctx context.Context, conn db.Conn, twt *model.TaskWorkTime) error {
	t.conn = conn.(*db.GORMConn)
	return t.conn.Tx(ctx).Save(&twt).Error
}

func (t *TaskWorkTimeDao) FindWorkTimeList(ctx context.Context, code int64) (list []*model.TaskWorkTime, err error) {
	err = t.conn.Session(ctx).Model(&model.TaskWorkTime{}).Where("task_code=?", code).Find(&list).Error
	return
}

func NewTaskWorkTimeDao() *TaskWorkTimeDao {
	return &TaskWorkTimeDao{conn: db.NewGORM()}
}
