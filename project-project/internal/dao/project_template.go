package dao

import (
	"context"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/pkg/model"
)

type ProjectTemplateDao struct {
	conn *db.GORMConn
}

func (p *ProjectTemplateDao) FindPojTmplSystem(ctx context.Context, page, pageSize int64) (
	list []model.ProjectTemplate, total int64, err error) {
	session := p.conn.Session(ctx)
	err = session.Where("is_system=1").
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		Find(&list).Error
	session.Model(&model.ProjectTemplate{}).Where("is_system=1").Count(&total)
	return
}

func (p *ProjectTemplateDao) FindPojTmplCustom(ctx context.Context, memId, orgCode, page, pageSize int64) (
	list []model.ProjectTemplate, total int64, err error) {
	session := p.conn.Session(ctx)
	err = session.Where("is_system=0 and member_code=? and organization_code=?",
		memId, orgCode).
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		Find(&list).Error
	session.Model(&model.ProjectTemplate{}).
		Where("is_system=0 and member_code=? and organization_code=?", memId, orgCode).Count(&total)
	return
}

func (p *ProjectTemplateDao) FindPojTmplAll(ctx context.Context, orgCode, page, pageSize int64) (
	list []model.ProjectTemplate, total int64, err error) {
	session := p.conn.Session(ctx)
	err = session.Where("organization_code=?",
		orgCode).
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		Find(&list).Error
	session.Model(&model.ProjectTemplate{}).
		Where("organization_code=?", orgCode).Count(&total)
	return
}

func NewProjectTemplateDao() *ProjectTemplateDao {
	return &ProjectTemplateDao{conn: db.NewGORM()}
}
