package db

import (
	"context"
	"gorm.io/gorm"
)

var _db *gorm.DB

func db() *gorm.DB {
	return _db
}

func SetDB(db__ *gorm.DB) {
	_db = db__
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
