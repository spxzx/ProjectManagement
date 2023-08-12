package dao

import (
	"context"
	"github.com/spxzx/project-user/internal/db"
	"github.com/spxzx/project-user/pkg/model"
)

type OrganizationDao struct {
	conn *db.GORMConn
}

func (o *OrganizationDao) FindOrganizationByMemberId(ctx context.Context, id int64) (orgList []model.Organization, err error) {
	err = o.conn.Session(ctx).Where("member_id=?", id).Find(&orgList).Error
	return
}

func (o *OrganizationDao) SaveOrganization(ctx context.Context, conn db.Conn, org *model.Organization) error {
	o.conn = conn.(*db.GORMConn)
	return o.conn.Tx(ctx).Create(org).Error
}

func NewOrganizationDao() *OrganizationDao {
	return &OrganizationDao{conn: db.NewGORM()}
}
