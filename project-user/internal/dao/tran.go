package dao

import "github.com/spxzx/project-user/internal/db"

type Transaction struct {
	conn db.Conn
}

func (t *Transaction) Action(f func(conn db.Conn) error) error {
	t.conn.Begin()
	if err := f(t.conn); err != nil {
		t.conn.Rollback()
		return err
	}
	t.conn.Commit()
	return nil
}

func NewTransaction() *Transaction {
	return &Transaction{conn: db.NewGORM()}
}
