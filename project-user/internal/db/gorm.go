package db

import (
	"context"
	"github.com/spxzx/project-user/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var _db *gorm.DB

func db() *gorm.DB {
	return _db
}

func init() {
	var err error
	if _db, err = gorm.Open(mysql.Open(config.Conf.InitMysqlOptions()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   config.Conf.Mysql.TablePrefix,
			SingularTable: true,
		},
	}); err != nil {
		panic("mysql connect error, cause by: " + err.Error())
	}
}

type GORMConn struct {
	db *gorm.DB
}

func NewGORM() *GORMConn {
	return &GORMConn{db: _db}
}

func (g *GORMConn) Session(ctx context.Context) *gorm.DB {
	return g.db.Session(&gorm.Session{Context: ctx})
}

func (g *GORMConn) Begin() {
	g.db = _db.Begin()
}

func (g *GORMConn) Rollback() {
	g.db.Rollback()
	g.db = _db
}

func (g *GORMConn) Commit() {
	g.db.Commit()
	g.db = _db
}

func (g *GORMConn) Tx(ctx context.Context) *gorm.DB {
	return g.db.WithContext(ctx)
}
