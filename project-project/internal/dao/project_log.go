package dao

import (
	"context"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/pkg/model"
)

type ProjectLogDao struct {
	conn *db.GORMConn
}

func (p *ProjectLogDao) FindLogByMemberCode(ctx context.Context, memberId int64,
	page int64, size int64) (list []*model.ProjectLog, total int64, err error) {
	s := p.conn.Session(ctx)
	offset := (page - 1) * size
	err = s.Model(&model.ProjectLog{}).
		Where("member_code=?", memberId).
		Limit(int(size)).
		Offset(int(offset)).
		Order("create_time desc").
		Find(&list).Error
	s.Model(&model.ProjectLog{}).Where("member_code=?", memberId).Count(&total)
	return
}

func (p *ProjectLogDao) FindLogByTaskCode(ctx context.Context, taskCode int64, comment int,
) (list []*model.ProjectLog, total int64, err error) {
	s := p.conn.Session(ctx)
	m := s.Model(&model.ProjectLog{})
	if comment == 1 {
		err = m.Where("source_code=? and is_comment=?", taskCode, comment).Find(&list).Error
		m.Where("source_code=? and is_comment=?", taskCode, comment).Count(&total)
	} else {
		err = m.Where("source_code=?", taskCode).Find(&list).Error
		m.Where("source_code=?", taskCode).Count(&total)
	}
	return
}

func (p *ProjectLogDao) FindLogByTaskCodePage(ctx context.Context, taskCode int64, comment, page, size int,
) (list []*model.ProjectLog, total int64, err error) {
	s := p.conn.Session(ctx)
	m := s.Model(&model.ProjectLog{})
	offset := (page - 1) * size
	if comment == 1 {
		err = m.Where("source_code=? and is_comment=?", taskCode, comment).Limit(size).Offset(offset).Find(&list).Error
		m.Where("source_code=? and is_comment=?", taskCode, comment).Count(&total)
	} else {
		err = m.Where("source_code=?", taskCode).Limit(size).Offset(offset).Find(&list).Error
		m.Where("source_code=?", taskCode).Count(&total)
	}
	return
}

func (p *ProjectLogDao) SaveProjectLog(ctx context.Context, conn db.Conn, pl *model.ProjectLog) (err error) {
	p.conn = conn.(*db.GORMConn)
	err = p.conn.Tx(ctx).Save(&pl).Error
	return
}

func NewProjectLogDao() *ProjectLogDao {
	return &ProjectLogDao{conn: db.NewGORM()}
}
