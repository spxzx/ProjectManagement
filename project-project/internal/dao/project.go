package dao

import (
	"context"
	"fmt"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/pkg/model"
)

type ProjectDao struct {
	conn *db.GORMConn
}

func (p *ProjectDao) FindProjectByIds(ctx context.Context, pids []int64) (list []*model.Project, err error) {
	err = p.conn.Session(ctx).Model(&model.Project{}).Where("id in (?)", pids).Find(&list).Error
	return
}

func (p *ProjectDao) FindProjectById(ctx context.Context, code int64) (pj *model.Project, err error) {
	err = p.conn.Session(ctx).Model(&model.Project{}).Where("id=?", code).First(&pj).Error
	return
}

func (p *ProjectDao) FindProjectMemberByPid(ctx context.Context, code int64) (
	list []*model.ProjectMember, total int64, err error) {
	s := p.conn.Session(ctx)
	err = s.Model(&model.ProjectMember{}).Where("project_code=?", code).Find(&list).Error
	s.Model(&model.ProjectMember{}).Where("project_code=?", code).Count(&total)
	return
}

func (p *ProjectDao) UpdateProject(ctx context.Context, pj *model.Project) error {
	return p.conn.Session(ctx).Model(&model.Project{}).Where("id=?", pj.Id).Updates(&pj).Error
}

func (p *ProjectDao) SaveProjectCollected(ctx context.Context, pc *model.ProjectCollection) error {
	return p.conn.Session(ctx).Save(&pc).Error
}

func (p *ProjectDao) DeleteProjectCollected(ctx context.Context, mid, pid int64) error {
	return p.conn.Session(ctx).Where("member_code=? and project_code=?", mid, pid).Delete(&model.ProjectCollection{}).Error
}

func (p *ProjectDao) UpdateProjectDeleted(ctx context.Context, pid int64, deleted bool) (err error) {
	s := p.conn.Session(ctx)
	if deleted {
		err = s.Model(&model.Project{}).Where("id=?", pid).Update("deleted", 1).Error
	} else {
		err = s.Model(&model.Project{}).Where("id=?", pid).Update("deleted", 0).Error
	}
	return
}

func (p *ProjectDao) FindCollectByPidAndMemId(ctx context.Context, pid, mid int64) (bool, error) {
	var cnt int64
	err := p.conn.Session(ctx).Model(&model.ProjectCollection{}).Where("project_code=? and member_code=?", pid, mid).Count(&cnt).Error
	return cnt > 0, err
}

func (p *ProjectDao) FindProjectByPidAndMemId(ctx context.Context, pid, mid int64) (pam *model.ProjectMemberUnion, err error) {
	err = p.conn.Session(ctx).Raw(
		"select * from pm_project a, pm_project_member b where a.id=b.project_code and a.id=? and b.member_code=? limit 1",
		pid, mid,
	).Scan(&pam).Error
	return
}

func (p *ProjectDao) SaveProjectMember(ctx context.Context, conn db.Conn, pm *model.ProjectMember) error {
	p.conn = conn.(*db.GORMConn)
	return p.conn.Tx(ctx).Create(pm).Error
}

func (p *ProjectDao) SaveProject(ctx context.Context, conn db.Conn, pj *model.Project) error {
	p.conn = conn.(*db.GORMConn)
	return p.conn.Tx(ctx).Create(pj).Error
}

func (p *ProjectDao) FindCollectProjectListByMemId(ctx context.Context, id int64, page int64, pageSize int64) (
	list []*model.ProjectMemberUnion, total int64, err error) {
	session := p.conn.Session(ctx)
	idx := (page - 1) * pageSize
	sql := "select * from pm_project a, pm_project_member b where a.id=b.project_code and a.id in (select project_code from pm_project_collection where member_code=?) and deleted=0 order by sort limit ?,?"
	d := session.Raw(sql, id, idx, pageSize)
	err = d.Scan(&list).Error
	cnt := session.Raw("select count(*) from pm_project where id in (select project_code from pm_project_collection where member_code=?) and deleted=0", id)
	cnt.Scan(&total)
	return
}

func (p *ProjectDao) FindMyProjectListByMemId(ctx context.Context, id int64, condition string, page int64, pageSize int64) (
	list []*model.ProjectMemberUnion, total int64, err error) {
	session := p.conn.Session(ctx)
	idx := (page - 1) * pageSize
	switch condition {
	case "":
		condition = "and deleted=0"
	case "archive":
		condition = "and archive=1 and deleted=0"
	default:
		condition = "and deleted=1"
	}
	sql := fmt.Sprintf(
		"select * from pm_project a, pm_project_member b where a.id=b.project_code and b.member_code=? %s order by sort limit ?,?",
		condition)
	d := session.Raw(sql, id, idx, pageSize)
	err = d.Scan(&list).Error
	cSql := fmt.Sprintf("select count(*) from pm_project a, pm_project_member b where a.id=b.project_code and b.member_code=? %s", condition)
	session.Raw(cSql, id).Scan(&total)
	return
}

func NewProjectDao() *ProjectDao {
	return &ProjectDao{conn: db.NewGORM()}
}
