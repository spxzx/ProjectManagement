package dao

import (
	"context"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/pkg/model"
	"gorm.io/gorm"
)

type MemberAccountDao struct {
	conn *db.GORMConn
}

func (m *MemberAccountDao) FindByMemberId(ctx context.Context, memId int64) (ma *model.MemberAccount, err error) {
	err = m.conn.Session(ctx).Where("member_code=?", memId).Take(&ma).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (m *MemberAccountDao) FindList(ctx context.Context, condition string,
	organizationCode int64, departmentCode int64, page int64, pageSize int64,
) (list []*model.MemberAccount, total int64, err error) {
	s := m.conn.Session(ctx)
	offset := (page - 1) * pageSize
	err = s.Model(&model.MemberAccount{}).
		Where("organization_code=?", organizationCode).
		Where(condition).Limit(int(pageSize)).
		Offset(int(offset)).Find(&list).Error
	s.Model(&model.MemberAccount{}).
		Where("organization_code=?", organizationCode).
		Where(condition).Limit(int(pageSize)).
		Offset(int(offset)).Count(&total)
	return
}

func NewMemberAccountDao() *MemberAccountDao {
	return &MemberAccountDao{conn: db.NewGORM()}
}
