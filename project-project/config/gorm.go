package config

import (
	"fmt"
	"github.com/spxzx/project-project/internal/db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

var _db *gorm.DB

func (c *config) ReConnMysql() {
	var err error
	if c.Db.Separation {
		mst := c.Db.Master
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			mst.Username, mst.Password, mst.Host, mst.Port, mst.Db)
		if _db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   mst.TablePrefix,
				SingularTable: true,
			},
		}); err != nil {
			panic("mysql connect error, cause by: " + err.Error())
		}
		var replicas []gorm.Dialector
		for _, v := range c.Db.Slave {
			dsn_ := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
				v.Username, v.Password, v.Host, v.Port, v.Db)
			cfg := mysql.Config{DSN: dsn_}
			replicas = append(replicas, mysql.New(cfg))
		}
		_ = _db.Use(dbresolver.Register(dbresolver.Config{
			Sources: []gorm.Dialector{mysql.New(mysql.Config{
				DSN: dsn,
			})},
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{},
		}).SetMaxIdleConns(10).SetMaxOpenConns(200))
	} else {
		if _db, err = gorm.Open(mysql.Open(c.InitMysqlOptions()), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   c.Mysql.TablePrefix,
				SingularTable: true,
			},
		}); err != nil {
			panic("mysql connect error, cause by: " + err.Error())
		}
	}
	db.SetDB(_db)
}
