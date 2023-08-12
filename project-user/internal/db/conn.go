package db

type Conn interface {
	Begin()
	Rollback()
	Commit()
}
