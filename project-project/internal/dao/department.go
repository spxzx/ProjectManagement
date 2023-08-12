package dao

import (
	"context"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/pkg/model"
	"gorm.io/gorm"
)

type DepartmentDao struct {
	conn *db.GORMConn
}

func (d *DepartmentDao) Save(ctx context.Context, conn db.Conn, dpm *model.Department) error {
	d.conn = conn.(*db.GORMConn)
	return d.conn.Tx(ctx).Save(&dpm).Error
}

func (d *DepartmentDao) FindDepartment(ctx context.Context, organizationCode int64, parentDepartmentCode int64, name string) (*model.Department, error) {
	session := d.conn.Session(ctx)
	session = session.Model(&model.Department{}).Where("organization_code=? AND name=?", organizationCode, name)
	if parentDepartmentCode > 0 {
		session = session.Where("pcode=?", parentDepartmentCode)
	}
	var dp *model.Department
	err := session.Limit(1).Take(&dp).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return dp, err
}

func (d *DepartmentDao) ListDepartment(organizationCode int64, parentDepartmentCode int64, page int64, size int64) (list []*model.Department, total int64, err error) {
	session := d.conn.Session(context.Background())
	session = session.Model(&model.Department{})
	session = session.Where("organization_code=?", organizationCode)
	if parentDepartmentCode > 0 {
		session = session.Where("pcode=?", parentDepartmentCode)
	}
	err = session.Limit(int(size)).Offset(int((page - 1) * size)).Find(&list).Error
	session.Count(&total)
	return
}

func (d *DepartmentDao) FindDepartmentById(ctx context.Context, id int64) (dt *model.Department, err error) {
	err = d.conn.Session(ctx).Where("id=?", id).Find(&dt).Error
	return
}

func NewDepartmentDao() *DepartmentDao {
	return &DepartmentDao{conn: db.NewGORM()}
}
