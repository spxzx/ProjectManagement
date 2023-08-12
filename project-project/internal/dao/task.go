package dao

import (
	"context"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/pkg/model"
	"gorm.io/gorm"
)

type TaskDao struct {
	conn *db.GORMConn
}

func (t *TaskDao) FindTaskByIds(ctx context.Context, ids []int64) (list []*model.Task, err error) {
	err = t.conn.Session(ctx).Model(&model.Task{}).Where("id in (?)", ids).Find(&list).Error
	return
}

func (t *TaskDao) FindTaskMemberList(ctx context.Context, code int64, page int64, size int64) (list []*model.TaskMember, total int64, err error) {
	s := t.conn.Session(ctx)
	offset := (page - 1) * size
	err = s.Model(&model.TaskMember{}).
		Where("task_code=?", code).
		Limit(int(size)).
		Offset(int(offset)).Find(&list).
		Error
	s.Model(&model.TaskMember{}).Where("task_code=?", code).Count(&total)
	return
}

func (t *TaskDao) FindTaskByAssignTo(ctx context.Context, id int64, type_ int32, page, size int64) (list []*model.Task, total int64, err error) {
	s := t.conn.Session(ctx)
	err = s.Model(&model.Task{}).
		Where("create_by=? and deleted=0 and done=?", id, type_).
		Limit(int(size)).Offset(int((page - 1) * size)).
		Find(&list).Error
	s.Model(&model.Task{}).Where("create_by=? and deleted=0 and done=?", id, type_).Count(&total)
	return
}

func (t *TaskDao) FindTaskByMemCode(ctx context.Context, id int64, type_ int32, page, size int64) (list []*model.Task, total int64, err error) {
	s := t.conn.Session(ctx)
	sql := "select a.* from pm_task a,pm_task_member b where a.id=b.task_code and member_code=? and a.deleted=0 and a.done=? limit ?,?"
	err = s.Model(&model.Task{}).Raw(sql, id, type_, (page-1)*size, size).Scan(&list).Error
	if err != nil {
		return nil, 0, err
	}
	sqlCount := "select count(*) from pm_task a,pm_task_member b where a.id=b.task_code and member_code=? and a.deleted=0 and a.done=?"
	s.Model(&model.Task{}).Raw(sqlCount, id, type_).Scan(&total)
	return
}

func (t *TaskDao) FindTaskByCreateBy(ctx context.Context, id int64, type_ int32, page, size int64) (list []*model.Task, total int64, err error) {
	s := t.conn.Session(ctx)
	err = s.Model(&model.Task{}).Where("assign_to=? and deleted=0 and done=?", id, type_).
		Limit(int(size)).Offset(int((page - 1) * size)).
		Find(&list).Error
	s.Model(&model.Task{}).Where("assign_to=? and deleted=0 and done=?", id, type_).Count(&total)
	return
}

func (t *TaskDao) FindTaskLessThenSortByStageCode(ctx context.Context, code, sort int) (tk *model.Task, err error) {
	err = t.conn.Session(ctx).Model(&model.Task{}).
		Where("stage_code=? and sort < ?", code, sort).
		Order("sort desc").First(&tk).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (t *TaskDao) UpdateTaskSort(ctx context.Context, conn db.Conn, next *model.Task) error {
	t.conn = conn.(*db.GORMConn)
	return t.conn.Tx(ctx).Select("sort", "stage_code").Updates(&next).Error
}

func (t *TaskDao) FindTaskById(ctx context.Context, id int64) (tk *model.Task, err error) {
	err = t.conn.Session(ctx).Model(&model.Task{}).Where("id=?", id).First(&tk).Error
	return
}

func (t *TaskDao) SaveTaskMember(ctx context.Context, conn db.Conn, tm *model.TaskMember) error {
	t.conn = conn.(*db.GORMConn)
	return t.conn.Tx(ctx).Save(&tm).Error
}

func (t *TaskDao) SaveTask(ctx context.Context, conn db.Conn, ts *model.Task) error {
	t.conn = conn.(*db.GORMConn)
	return t.conn.Tx(ctx).Save(&ts).Error
}

func (t *TaskDao) FindTaskMaxSort(ctx context.Context, pCode int64, TSCode int64) (maxSort *int, err error) {
	err = t.conn.Session(ctx).
		Model(&model.Task{}).
		Where("project_code=? and stage_code=?", pCode, TSCode).
		Select("max(sort) as sort").
		Scan(&maxSort).Error
	return
}

func (t *TaskDao) FindTaskMaxIdNum(ctx context.Context, code int64) (maxIdNum *int, err error) {
	err = t.conn.Session(ctx).
		Model(&model.Task{}).
		Where("project_code=?", code).
		Select("max(id_num) as maxIdNum").
		Scan(&maxIdNum).Error
	return
}

func (t *TaskDao) FindTaskMemberByTid(ctx context.Context, tid, mid int64) (mem *model.TaskMember, err error) {
	err = t.conn.Session(ctx).Where("task_code=? and member_code=?", tid, mid).First(&mem).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (t *TaskDao) FindTaskByStageCode(ctx context.Context, code int64) (list []*model.Task, err error) {
	err = t.conn.Session(ctx).Model(&model.Task{}).
		Where("stage_code=? and deleted=0", code).
		Order("sort asc").Find(&list).Error
	return
}

func NewTaskDao() *TaskDao {
	return &TaskDao{conn: db.NewGORM()}
}
