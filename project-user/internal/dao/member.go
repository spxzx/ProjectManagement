package dao

import (
	"context"
	"github.com/spxzx/project-user/internal/db"
	"github.com/spxzx/project-user/pkg/model"
)

type MemberDao struct {
	conn *db.GORMConn
}

func (m *MemberDao) FindMemberByIds(ctx context.Context, ids []int64) (list []*model.Member, err error) {
	err = m.conn.Session(ctx).Model(&model.Member{}).Where("id in (?)", ids).Find(&list).Error
	return
}

func (m *MemberDao) GetMemberById(ctx context.Context, id int64) (mem model.Member, err error) {
	err = m.conn.Session(ctx).Model(&model.Member{}).Where("id=?", id).First(&mem).Error
	return
}

func (m *MemberDao) FindMember(ctx context.Context, account, password string) (mem model.Member, err error) {
	err = m.conn.Session(ctx).Where("account=? and password=?", account, password).First(&mem).Error
	return
}

func (m *MemberDao) SaveMember(ctx context.Context, conn db.Conn, mem *model.Member) error {
	m.conn = conn.(*db.GORMConn)
	return m.conn.Tx(ctx).Create(mem).Error
}

func (m *MemberDao) ExistMemberByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := m.conn.Session(ctx).Model(&model.Member{}).Where("email=?", email).Count(&count).Error
	return count > 0, err
}

func (m *MemberDao) ExistMemberByAccount(ctx context.Context, account string) (bool, error) {
	var count int64
	err := m.conn.Session(ctx).Model(&model.Member{}).Where("account=?", account).Count(&count).Error
	return count > 0, err
}

func (m *MemberDao) ExistMemberByMobile(ctx context.Context, mobile string) (bool, error) {
	var count int64
	err := m.conn.Session(ctx).Model(&model.Member{}).Where("mobile=?", mobile).Count(&count).Error
	return count > 0, err
}

func NewMemberDao() *MemberDao {
	return &MemberDao{conn: db.NewGORM()}
}
