package dao

import (
	"context"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/pkg/model"
)

type TaskStagesTmplDao struct {
	conn *db.GORMConn
}

func (t *TaskStagesTmplDao) FindTaskStagesTplByPid(ctx context.Context, pid int64) (list []*model.TaskStagesTemplate, err error) {
	err = t.conn.Session(ctx).
		Model(&model.TaskStagesTemplate{}).
		Where("project_template_code=?", pid).
		Order("sort desc, id asc").
		Find(&list).Error
	return
}

func (t *TaskStagesTmplDao) FindInPojTmplIds(ctx context.Context, ids []int) (list []model.TaskStagesTemplate, err error) {
	err = t.conn.Session(ctx).Where("project_template_code in ?", ids).Find(&list).Error
	return
}

func NewTaskStagesTmplDao() *TaskStagesTmplDao {
	return &TaskStagesTmplDao{conn: db.NewGORM()}
}
