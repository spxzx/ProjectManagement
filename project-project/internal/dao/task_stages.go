package dao

import (
	"context"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/pkg/model"
)

type TaskStagesDao struct {
	conn *db.GORMConn
}

func (t *TaskStagesDao) FindTaskStagesById(ctx context.Context, code int64) (ts *model.TaskStages, err error) {
	err = t.conn.Session(ctx).Model(&model.TaskStages{}).Where("id=?", code).First(&ts).Error
	return
}

func (t *TaskStagesDao) FindTaskStagesByPCode(ctx context.Context, code, page, pageSize int64) (list []*model.TaskStages, total int64, err error) {
	s := t.conn.Session(ctx)
	err = s.Model(&model.TaskStages{}).
		Where("project_code=? and deleted=0", code).
		Order("sort asc").
		Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).
		Find(&list).Error
	s.Model(&model.TaskStages{}).Where("project_code=? and deleted=0", code).Count(&total)
	return
}

func (t *TaskStagesDao) SaveTaskStages(ctx context.Context, conn db.Conn, ts *model.TaskStages) error {
	t.conn = conn.(*db.GORMConn)
	return t.conn.Tx(ctx).Model(&model.TaskStages{}).Save(&ts).Error
}

func NewTaskStagesDao() *TaskStagesDao {
	return &TaskStagesDao{conn: db.NewGORM()}
}
